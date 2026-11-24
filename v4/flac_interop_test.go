//go:build integration
// +build integration

package audiometa_test

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/gcottom/audiometa/v4"
)

func TestFLACIndependentMetadata(t *testing.T) {
	python := os.Getenv("AUDIOMETA_MUTAGEN_PYTHON")
	if python == "" {
		t.Skip("set AUDIOMETA_MUTAGEN_PYTHON to a Python with Mutagen")
	}
	dir := t.TempDir()
	_, data := canonicalFixture(t, canonicalFormats[len(canonicalFormats)-1], dir)
	tag, err := audiometa.OpenMetadata(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if err := tag.SetTitle("International — 音"); err != nil {
		t.Fatal(err)
	}
	if err := tag.SetField(audiometa.Artist, audiometa.TextValue("One"), audiometa.TextValue("Deux")); err != nil {
		t.Fatal(err)
	}
	if err := tag.SetISRC("USABC2600001"); err != nil {
		t.Fatal(err)
	}
	if err := tag.SetField(audiometa.ReplayGainTrackGain, audiometa.DecimalValue(-3.25)); err != nil {
		t.Fatal(err)
	}
	if err := tag.SetField(audiometa.Chapters, audiometa.ChapterValue(audiometa.Chapter{ID: "intro", Title: "Opening", EndMS: 80})); err != nil {
		t.Fatal(err)
	}
	if err := tag.SetCustom("Collection", "One", "Two"); err != nil {
		t.Fatal(err)
	}

	pic := image.NewRGBA(image.Rect(0, 0, 4, 3))
	pic.Set(0, 0, color.RGBA{R: 255, A: 255})
	var pngBytes bytes.Buffer
	if err := png.Encode(&pngBytes, pic); err != nil {
		t.Fatal(err)
	}
	picture := audiometa.Picture{Role: 3, MIME: "image/png", Description: "Front", Width: 4, Height: 3, Depth: 32, Data: pngBytes.Bytes()}
	if err := tag.SetField(audiometa.Artwork, audiometa.PictureValue(picture)); err != nil {
		t.Fatal(err)
	}

	name := filepath.Join(dir, "native.flac")
	out, err := os.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	save, close := tag.Save(out), out.Close()
	if save != nil || close != nil {
		t.Fatalf("save: %v close: %v", save, close)
	}

	const check = `
import hashlib
import sys
import mutagen
from mutagen.flac import FLAC
tag = FLAC(sys.argv[1])
assert tag["title"] == ["International — 音"], tag["title"]
assert tag["artist"] == ["One", "Deux"], tag["artist"]
assert tag["isrc"] == ["USABC2600001"], tag["isrc"]
assert tag["replaygain_track_gain"] == ["-3.25 dB"], tag["replaygain_track_gain"]
assert tag["audiometa_1_custom_collection"] == ["One", "Two"]
assert len(tag["audiometa_1_chapters"]) == 1
assert len(tag.pictures) == 1
picture = tag.pictures[0]
assert (picture.type, picture.mime, picture.desc, picture.width, picture.height, picture.depth) == (3, "image/png", "Front", 4, 3, 32)
assert hashlib.sha256(picture.data).hexdigest() == sys.argv[2]
print("Mutagen", mutagen.version_string, "verified FLAC native and opaque entries and artwork SHA-256")
`
	t.Log(runTool(t, python, "-c", check, name, fmt.Sprintf("%x", sha256.Sum256(pngBytes.Bytes()))))
}
