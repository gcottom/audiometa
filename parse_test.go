package mp3mp4tag

import (
	"path/filepath"
	"testing"
)

func TestReadMP3Tags(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-mp3.mp3")
	tag, err := parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.artist == "" || tag.album == "" || tag.title == "" {
		t.Fatal("Data parsed was blank!")
	}

}
func TestReadM4ATags(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-m4a.m4a")
	tag, err := parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.artist == "" || tag.album == "" || tag.title == "" {
		t.Fatal("Data parsed was blank!")
	}
}
func TestReadFlacTags(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-flac.flac")
	tag, err := parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.artist == "" || tag.album == "" || tag.title == "" {
		t.Fatal("Data parsed was blank!")
	}
}
func TestReadOggTags(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-ogg.ogg")
	tag, err := parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.artist == "" || tag.album == "" || tag.title == "" {
		t.Fatal("Data parsed was blank!")
	}
}
