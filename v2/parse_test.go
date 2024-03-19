package audiometa

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadMP3Tags(t *testing.T) {
	path, _ := filepath.Abs("../testdata/testdata-mp3.mp3")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	tag, err := parse(f, ParseOptions{MP3})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.artist == "" || tag.album == "" || tag.title == "" {
		t.Fatal("Data parsed was blank!")
	}

}
func TestReadM4ATags(t *testing.T) {
	path, _ := filepath.Abs("../testdata/testdata-m4a.m4a")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	tag, err := parse(f, ParseOptions{M4A})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.artist == "" || tag.album == "" || tag.title == "" {
		t.Fatal("Data parsed was blank!")
	}
}
func TestReadFlacTags(t *testing.T) {
	path, _ := filepath.Abs("../testdata/testdata-flac.flac")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	tag, err := parse(f, ParseOptions{FLAC})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.artist == "" || tag.album == "" || tag.title == "" {
		t.Fatal("Data parsed was blank!")
	}
}
func TestReadOggVorbisTags(t *testing.T) {
	path, _ := filepath.Abs("../testdata/testdata-ogg.ogg")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	tag, err := parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.artist == "" || tag.album == "" || tag.title == "" {
		t.Fatal("Data parsed was blank!")
	}
}
func TestReadOggOpusTags(t *testing.T) {
	path, _ := filepath.Abs("../testdata/testdata-opus.ogg")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	tag, err := parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.artist == "" || tag.album == "" || tag.title == "" {
		t.Fatal("Data parsed was blank!")
	}
}
