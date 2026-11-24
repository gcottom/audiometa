"""Build the candidate modules through Go's file-proxy protocol."""

import argparse
import json
import os
from pathlib import Path
import subprocess
import tempfile
import zipfile

ROOT = Path(__file__).resolve().parents[2]
MODULES = {
    "audiometa-core": ("audiometa-core", "v0.1.0"),
    "mp3meta": ("mp3meta", "v0.1.0"),
    "mp4meta": ("mp4meta", "v0.1.0"),
    "oggmeta": ("oggmeta", "v0.1.0"),
    "flacmeta": ("flacmeta", "v0.1.0"),
    "audiometa/v3": ("audiometa/v3", "v3.1.0"),
}
CONSUMER = r'''package consumer

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"

	"github.com/gcottom/audiometa/v3"
)

func TestConsumer(t *testing.T) {
	// The source only needs a stream-info block and frame sync for this SDK
	// compile and metadata smoke test. Audio decode is tested independently.
	input := make([]byte, 44)
	copy(input, "fLaC")
	binary.BigEndian.PutUint32(input[4:8], 0x80000022)
	copy(input[42:], []byte{255, 248})
	legacy, err := audiometa.OpenTag(bytes.NewReader(input))
	if err != nil { t.Fatal(err) }
	legacy.SetTitle("Legacy consumer")
	var original bytes.Buffer
	if err := legacy.Save(&original); err != nil { t.Fatal(err) }

	tag, err := audiometa.OpenMetadata(bytes.NewReader(original.Bytes()))
	if err != nil || tag.GetTitle() != "Legacy consumer" { t.Fatalf("legacy to canonical: %v", err) }
	if err := tag.SetTitle("Modern consumer"); err != nil { t.Fatal(err) }
	if err := tag.SetField(audiometa.Artist, audiometa.TextValue("First"), audiometa.TextValue("Second")); err != nil { t.Fatal(err) }
	if err := tag.SetCustom("collection", "Archive", "Favorites"); err != nil { t.Fatal(err) }
	if err := tag.SetISRC("invalid"); !errors.Is(err, audiometa.ErrInvalidValue) { t.Fatal(err) }
	var result bytes.Buffer
	if err := tag.Save(&result); err != nil { t.Fatal(err) }
	back, err := audiometa.OpenMetadata(&readSeeker{bytes.NewReader(result.Bytes())})
	if err != nil || back.GetTitle() != "Modern consumer" { t.Fatalf("canonical reopen: %v", err) }
	if a, _ := back.Get(audiometa.Artist); len(a) != 2 { t.Fatal(a) }
	if c, _ := back.GetCustom("COLLECTION"); len(c) != 2 { t.Fatal(c) }
}

type readSeeker struct { *bytes.Reader }
'''


def module_files(source):
    names = subprocess.check_output(
        ["git", "ls-files", "-z", "--cached", "--others", "--exclude-standard"], cwd=source
    ).decode("utf-8").split("\0")
    for name in sorted(set(names) - {""}):
        path = source / name
        if path.is_file() and not path.is_symlink():
            yield name, path


def package(proxy, name, source, version):
    module = f"github.com/gcottom/{name}"
    target = proxy / module / "@v"
    target.mkdir(parents=True)
    (target / f"{version}.mod").write_bytes((source / "go.mod").read_bytes())
    (target / f"{version}.info").write_text(
        json.dumps({"Version": version, "Time": "2026-11-24T00:00:00Z"}), encoding="utf-8"
    )
    (target / "list").write_text(f"{version}\n", encoding="utf-8")
    prefix = f"{module}@{version}/"
    with zipfile.ZipFile(target / f"{version}.zip", "w", zipfile.ZIP_DEFLATED) as blob:
        packaged = set()
        for name, path in module_files(source):
            blob.write(path, prefix + name)
            packaged.add(name)
        if "go.mod" not in packaged:
            raise ValueError(f"missing module manifest for {module}")
        # Go copies a parent module license into the ZIP when appropriate.
        for filename in ("LICENSE", "LICENSE.md"):
            if filename not in packaged and (source.parent / filename).is_file():
                blob.write(source.parent / filename, prefix + filename)


def run(work, go):
    proxy = work / "proxy"
    for name, (directory, version) in MODULES.items():
        package(proxy, name, ROOT / directory, version)

    consumer = work / "consumer"
    consumer.mkdir()
    (consumer / "go.mod").write_text(
        "module example.org/audiometa-consumer\n\ngo 1.18\n\n"
        "require github.com/gcottom/audiometa/v3 v3.1.0\n",
        encoding="utf-8",
    )
    (consumer / "consumer_test.go").write_text(CONSUMER, encoding="utf-8")
    env = dict(
        os.environ, GOENV="off", GOWORK="off", GOFLAGS="",
        GOTOOLCHAIN="local", GOMODCACHE=str(work / "module-cache"),
        GOPRIVATE="", GONOPROXY="", GONOSUMDB="github.com/gcottom/*",
        GOPROXY=f"{proxy.as_uri()},https://proxy.golang.org", GOSUMDB="sum.golang.org",
    )
    for args in (["version"], ["mod", "tidy"], ["list", "-m", "all"],
                 ["build", "./..."], ["test", "-count=1", "-v", "./..."]):
        subprocess.run([go, *args], cwd=consumer, env=env, check=True)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--go", default="go", help="Go executable path")
    options = parser.parse_args()
    with tempfile.TemporaryDirectory(prefix="audiometa-release-") as work:
        run(Path(work), options.go)


if __name__ == "__main__":
    main()
