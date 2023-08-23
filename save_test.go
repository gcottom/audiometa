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
	TestWriteTagsM4AFromEmpty(t)
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
	TestWriteTagsFlacFromEmpty(t)
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
func TestWriteEmptyTagsOggVorbis(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-ogg-vorbis-nonEmpty.ogg")
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
func TestWriteTagsOggVorbisFromEmpty(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-ogg-vorbis-nonEmpty.ogg")
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
	tag.SetAlbumArtFromFilePath("testdata/testdata-img-1.jpg")
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
func TestWriteTagsOggVorbisFromEmptyExtended(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-ogg-vorbis-nonEmpty.ogg")
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
	tag.SetAlbumArtFromFilePath("testdata/testdata-img-1.jpg")
	tag.SetAdditionalTag("TEST", "TEST")
	tag.SetAdditionalTag("TEST2", "TEST2")
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
	if tag.passThroughMap["TEST"] != "TEST" || tag.passThroughMap["TEST2"] != "TEST2" {
		t.Fatal("Extended Tags Not Found")
	}
}
func TestUpdateTagsOggVorbis(t *testing.T) {
	TestWriteTagsOggVorbisFromEmpty(t)
	path, _ := filepath.Abs("testdata/testdata-ogg-vorbis-nonEmpty.ogg")
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
func TestUpdateTagsOggVorbisExtended(t *testing.T) {
	TestWriteTagsOggVorbisFromEmpty(t)
	path, _ := filepath.Abs("testdata/testdata-ogg-vorbis-nonEmpty.ogg")
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
	tag.SetAdditionalTag("TEST", "TEST")
	tag.SetAdditionalTag("TEST2", "TEST2")
	err = SaveTag(tag)
	if err != nil {
		t.Fatal("Error saving!")
	}
	tag.SetArtist("TestArtist2")
	tag.SetAdditionalTag("TEST", "TEST3")
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
	if tag.passThroughMap["TEST"] != "TEST3" || tag.passThroughMap["TEST2"] != "TEST2" {
		t.Fatal("Extended Tags Not Found")
	}
}
func TestWriteEmptyTagsOggOpus(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-opus-nonEmpty.ogg")
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
func TestWriteTagsOggOpusFromEmpty(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-opus-nonEmpty.ogg")
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
	tag.SetAlbumArtFromFilePath("testdata/testdata-img-1.jpg")
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
func TestWriteTagsOggOpusFromEmptyExtended(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-opus-nonEmpty.ogg")
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
	tag.SetAlbumArtFromFilePath("testdata/testdata-img-1.jpg")
	tag.SetAdditionalTag("TEST", "TEST")
	tag.SetAdditionalTag("TEST2", "TEST2")
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
	if tag.passThroughMap["TEST"] != "TEST" || tag.passThroughMap["TEST2"] != "TEST2" {
		t.Fatal("Extended Tags Not Found")
	}
}
func TestUpdateTagsOggOpus(t *testing.T) {
	TestWriteTagsOggOpusFromEmpty(t)
	path, _ := filepath.Abs("testdata/testdata-opus-nonEmpty.ogg")
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
func TestUpdateTagsOggOpusExtended(t *testing.T) {
	TestWriteTagsOggOpusFromEmpty(t)
	path, _ := filepath.Abs("testdata/testdata-opus-nonEmpty.ogg")
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
	tag.SetAdditionalTag("TEST", "TEST")
	tag.SetAdditionalTag("TEST2", "TEST2")
	err = SaveTag(tag)
	if err != nil {
		t.Fatal("Error saving!")
	}
	tag.SetArtist("TestArtist2")
	tag.SetAdditionalTag("TEST", "TEST3")
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
	if tag.passThroughMap["TEST"] != "TEST3" || tag.passThroughMap["TEST2"] != "TEST2" {
		t.Fatal("Extended Tags Not Found")
	}
}
