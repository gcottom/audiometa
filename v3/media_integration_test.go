//go:build integration
// +build integration

package audiometa_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gcottom/audiometa/v3"
)

func runTool(t *testing.T, name string, args ...string) string {
	t.Helper()
	data, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v: %v: %s", name, args, err, data)
	}
	return strings.TrimSpace(string(data))
}
func TestAudioMetaIntegration_RoundTripAndDecode(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Fatal("ffmpeg required")
	}
	for _, tc := range []struct{ ext, enc string }{
		{"mp3", "libmp3lame"}, {"m4a", "aac"}, {"m4b", "aac"}, {"ogg", "libvorbis"}, {"opus", "libopus"}, {"flac", "flac"},
	} {
		t.Run(tc.ext, func(t *testing.T) {
			dir := t.TempDir()
			in, out := filepath.Join(dir, "in."+tc.ext), filepath.Join(dir, "out."+tc.ext)
			runTool(t, "ffmpeg", "-v", "error", "-f", "lavfi", "-i", "sine=frequency=997:sample_rate=48000:duration=0.25", "-metadata", "title=Original", "-metadata", "artist=Original artist", "-c:a", tc.enc, in)
			reader, err := os.Open(in)
			if err != nil {
				t.Fatal(err)
			}
			defer reader.Close()
			tag, err := audiometa.OpenTag(reader)
			if err != nil {
				t.Fatalf("OpenTag: %v", err)
			}
			if tag.GetTitle() != "Original" {
				t.Errorf("initial title = %q", tag.GetTitle())
			}
			title := "Audit – a longer title ñ"
			artist := "Audit artist"
			tag.SetTitle(title)
			tag.SetArtist(artist)
			f, err := os.Create(out)
			if err != nil {
				t.Fatal(err)
			}
			saveErr := tag.Save(f)
			closeErr := f.Close()
			if saveErr != nil || closeErr != nil {
				t.Fatalf("save: %v; close: %v", saveErr, closeErr)
			}
			data, err := os.ReadFile(out)
			if err != nil {
				t.Fatal(err)
			}
			got, err := audiometa.OpenTag(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("reopen (%d bytes): %v", len(data), err)
			}
			if got.GetTitle() != title || got.GetArtist() != artist {
				t.Errorf("metadata mismatch: %q / %q", got.GetTitle(), got.GetArtist())
			}
			// Decode both files and hash raw PCM independently of library parsing.
			pcm := func(name string) string {
				return runTool(t, "ffmpeg", "-xerror", "-v", "error", "-i", name, "-map", "0:a:0", "-c:a", "pcm_s16le", "-f", "hash", "-hash", "sha256", "-")
			}
			if tc.ext == "mp3" {
				initial, err := os.ReadFile(in)
				if err != nil {
					t.Fatal(err)
				}
				cut := func(b []byte) int {
					if !bytes.HasPrefix(b, []byte("ID3")) || len(b) < 10 {
						t.Fatal("missing ID3 header")
					}
					return 10 + (int(b[6]) << 21) + (int(b[7]) << 14) + (int(b[8]) << 7) + int(b[9])
				}
				originalAudio, savedAudio := initial[cut(initial):], data[cut(data):]
				if !bytes.Equal(originalAudio, savedAudio) {
					t.Error("encoded MP3 bytes after ID3 changed")
				}
			}
			if want, actual := pcm(in), pcm(out); want != actual {
				t.Errorf("decoded PCM changed: before %s, after %s", want, actual)
			}
			t.Logf("%T, %d bytes, decoded audio %s", got, len(data), pcm(out))
		})
	}
}
func TestAudioMetaIntegration_NoID3Header(t *testing.T) {
	input := filepath.Join(t.TempDir(), "no-id3.mp3")
	runTool(t, "ffmpeg", "-v", "error", "-f", "lavfi", "-i", "sine=frequency=997:sample_rate=48000:duration=0.25", "-c:a", "libmp3lame", "-id3v2_version", "0", "-write_id3v1", "0", input)
	b, err := os.ReadFile(input)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.HasPrefix(b, []byte("ID3")) {
		t.Fatal("expected fixture without ID3 prefix")
	}
	runTool(t, "ffmpeg", "-xerror", "-v", "error", "-i", input, "-f", "null", "-")
	_, err = audiometa.OpenTag(bytes.NewReader(b))
	if err != nil {
		t.Fatalf("valid headerless MP3 rejected: %v", err)
	}
}

func TestAudioMetaIntegration_UntaggedMP3WritePreservesAudio(t *testing.T) {
	dir := t.TempDir()
	input, out := filepath.Join(dir, "untagged.mp3"), filepath.Join(dir, "tagged.mp3")
	runTool(t, "ffmpeg", "-v", "error", "-f", "lavfi", "-i", "sine=frequency=997:sample_rate=48000:duration=0.25", "-c:a", "libmp3lame", "-id3v2_version", "0", "-write_id3v1", "0", input)
	want, err := os.ReadFile(input)
	if err != nil {
		t.Fatal(err)
	}
	tag, err := audiometa.OpenTag(bytes.NewReader(want))
	if err != nil {
		t.Fatal(err)
	}
	tag.SetTitle("Previously untagged")
	w, err := os.Create(out)
	if err != nil {
		t.Fatal(err)
	}
	saveErr, closeErr := tag.Save(w), w.Close()
	if saveErr != nil || closeErr != nil {
		t.Fatalf("save: %v; close: %v", saveErr, closeErr)
	}
	got, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) < 10 || !bytes.HasPrefix(got, []byte("ID3")) {
		t.Fatal("output does not have a valid ID3 header")
	}
	span := 10 + (int(got[6]) << 21) + (int(got[7]) << 14) + (int(got[8]) << 7) + int(got[9])
	if span > len(got) || !bytes.Equal(want, got[span:]) {
		t.Fatal("original MPEG bytes were not preserved")
	}
	roundtrip, err := audiometa.OpenTag(bytes.NewReader(got))
	if err != nil {
		t.Fatal(err)
	}
	if roundtrip.GetTitle() != "Previously untagged" {
		t.Fatal("written title did not round-trip")
	}
	pcm := func(f string) string {
		return runTool(t, "ffmpeg", "-xerror", "-v", "error", "-i", f, "-map", "0:a:0", "-c:a", "pcm_s16le", "-f", "hash", "-hash", "sha256", "-")
	}
	if pcm(input) != pcm(out) {
		t.Fatal("decoded PCM changed while tagging the untagged file")
	}
}
