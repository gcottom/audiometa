package mp3mp4tag

import (
	"path/filepath"
	"testing"
)

func TestWriteEmptyTagsMP3(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-mp3-nonEmpty.mp3")
	tag, err := parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	err = SaveTag(tag)
	if err != nil {
		t.Fatal("Error saving!")
	}
	tag, err = parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "" || tag.Album() != "" || tag.Title() != "" {
		t.Fatal("Failed to remove tags for empty tag test!")
	}
}
func TestWriteTagsMP3FromEmpty(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-mp3-nonEmpty.mp3")
	tag, err := parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	err = SaveTag(tag)
	if err != nil {
		t.Fatal("Error saving!")
	}
	tag, err = parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	err = SaveTag(tag)
	if err != nil {
		t.Fatal("Error saving!")
	}
	tag, err = parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "TestArtist1" || tag.Album() != "TestAlbum1" || tag.Title() != "TestTitle1" {
		t.Fatal("Failed to validate new tags")
	}
}

func TestUpdateTagsMP3(t *testing.T) {
	TestWriteTagsMP3FromEmpty(t)
	path, _ := filepath.Abs("testdata/testdata-mp3-nonEmpty.mp3")
	tag, err := parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	err = SaveTag(tag)
	if err != nil {
		t.Fatal("Error saving!")
	}
	tag, err = parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	err = SaveTag(tag)
	if err != nil {
		t.Fatal("Error saving!")
	}
	tag.SetArtist("TestArtist2")
	err = SaveTag(tag)
	if err != nil {
		t.Fatal("Error saving!")
	}
	tag, err = parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "TestArtist2" || tag.Album() != "TestAlbum1" || tag.Title() != "TestTitle1" {
		t.Fatal("Failed to validate new tags")
	}
}

func TestWriteEmptyTagsM4A(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-m4a-nonEmpty.m4a")
	tag, err := parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	err = SaveTag(tag)
	if err != nil {
		t.Fatal("Error saving!")
	}
	tag, err = parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "" || tag.Album() != "" || tag.Title() != "" {
		t.Fatal("Failed to remove tags for empty tag test!")
	}
}
func TestWriteTagsM4AFromEmpty(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-m4a-nonEmpty.m4a")
	tag, err := parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	err = SaveTag(tag)
	if err != nil {
		t.Fatal("Error saving!")
	}
	tag, err = parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	err = SaveTag(tag)
	if err != nil {
		t.Fatal("Error saving!")
	}
	tag, err = parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "TestArtist1" || tag.Album() != "TestAlbum1" || tag.Title() != "TestTitle1" {
		t.Fatal("Failed to validate new tags")
	}
}
func TestUpdateTagsM4A(t *testing.T) {
	TestWriteTagsMP3FromEmpty(t)
	path, _ := filepath.Abs("testdata/testdata-m4a-nonEmpty.m4a")
	tag, err := parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	err = SaveTag(tag)
	if err != nil {
		t.Fatal("Error saving!")
	}
	tag, err = parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	err = SaveTag(tag)
	if err != nil {
		t.Fatal("Error saving!")
	}
	tag.SetArtist("TestArtist2")
	err = SaveTag(tag)
	if err != nil {
		t.Fatal("Error saving!")
	}
	tag, err = parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "TestArtist2" || tag.Album() != "TestAlbum1" || tag.Title() != "TestTitle1" {
		t.Fatal("Failed to validate new tags")
	}
}
func TestWriteEmptyTagsFlac(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-flac-nonEmpty.flac")
	tag, err := parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	err = SaveTag(tag)
	if err != nil {
		t.Log(err)
		t.Fatal("Error saving!")
	}
	tag, err = parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "" || tag.Album() != "" || tag.Title() != "" {
		t.Fatal("Failed to remove tags for empty tag test!")
	}
}
func TestWriteTagsFlacFromEmpty(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-flac-nonEmpty.flac")
	tag, err := parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	err = SaveTag(tag)
	if err != nil {
		t.Fatal("Error saving!")
	}
	tag, err = parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	err = SaveTag(tag)
	if err != nil {
		t.Fatal("Error saving!")
	}
	tag, err = parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "TestArtist1" || tag.Album() != "TestAlbum1" || tag.Title() != "TestTitle1" {
		t.Fatal("Failed to validate new tags")
	}
}
func TestUpdateTagsFlac(t *testing.T) {
	TestWriteTagsMP3FromEmpty(t)
	path, _ := filepath.Abs("testdata/testdata-flac-nonEmpty.flac")
	tag, err := parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	err = SaveTag(tag)
	if err != nil {
		t.Fatal("Error saving!")
	}
	tag, err = parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	err = SaveTag(tag)
	if err != nil {
		t.Fatal("Error saving!")
	}
	tag.SetArtist("TestArtist2")
	err = SaveTag(tag)
	if err != nil {
		t.Fatal("Error saving!")
	}
	tag, err = parse(path)
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "TestArtist2" || tag.Album() != "TestAlbum1" || tag.Title() != "TestTitle1" {
		t.Fatal("Failed to validate new tags")
	}
}
