package audiometa

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteEmptyTagsMP3(t *testing.T) {
	path, _ := filepath.Abs("../testdata/testdata-mp3-nonEmpty.mp3")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	f.Seek(0, 0)
	tag, err := parse(f, ParseOptions{MP3})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	tag.filePath = path
	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{MP3})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "" || tag.Album() != "" || tag.Title() != "" {
		t.Fatal("Failed to remove tags for empty tag test!")
	}
}
func TestWriteTagsMP3FromEmpty(t *testing.T) {
	path, _ := filepath.Abs("../testdata/testdata-mp3-nonEmpty.mp3")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	f.Seek(0, 0)
	tag, err := parse(f, ParseOptions{MP3})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	tag.filePath = path
	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{MP3})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	tag.filePath = path
	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{MP3})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "TestArtist1" || tag.Album() != "TestAlbum1" || tag.Title() != "TestTitle1" {
		t.Fatal("Failed to validate new tags")
	}
}

func TestUpdateTagsMP3(t *testing.T) {
	TestWriteTagsMP3FromEmpty(t)
	path, _ := filepath.Abs("../testdata/testdata-mp3-nonEmpty.mp3")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	f.Seek(0, 0)
	tag, err := parse(f, ParseOptions{MP3})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	tag.filePath = path
	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{MP3})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	tag.filePath = path
	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	tag.SetArtist("TestArtist2")
	tag.filePath = path
	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{MP3})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "TestArtist2" || tag.Album() != "TestAlbum1" || tag.Title() != "TestTitle1" {
		t.Fatal("Failed to validate new tags")
	}
}

func TestWriteEmptyTagsM4A(t *testing.T) {
	path, _ := filepath.Abs("../testdata/testdata-m4a-nonEmpty.m4a")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	f.Seek(0, 0)
	tag, err := parse(f, ParseOptions{M4A})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.filePath = path
	tag.ClearAllTags()
	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{M4A})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "" || tag.Album() != "" || tag.Title() != "" {
		t.Fatal("Failed to remove tags for empty tag test!")
	}
}
func TestWriteTagsM4AFromEmpty(t *testing.T) {
	path, _ := filepath.Abs("../testdata/testdata-m4a-nonEmpty.m4a")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	f.Seek(0, 0)
	tag, err := parse(f, ParseOptions{M4A})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	tag.filePath = path
	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{M4A})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	tag.filePath = path
	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{M4A})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "TestArtist1" || tag.Album() != "TestAlbum1" || tag.Title() != "TestTitle1" {
		t.Fatal("Failed to validate new tags")
	}
}
func TestUpdateTagsM4A(t *testing.T) {
	TestWriteTagsM4AFromEmpty(t)
	path, _ := filepath.Abs("../testdata/testdata-m4a-nonEmpty.m4a")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	f.Seek(0, 0)
	tag, err := parse(f, ParseOptions{M4A})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	tag.filePath = path
	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{M4A})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	tag.filePath = path
	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	tag.SetArtist("TestArtist2")
	tag.filePath = path
	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{M4A})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "TestArtist2" || tag.Album() != "TestAlbum1" || tag.Title() != "TestTitle1" {
		t.Fatal("Failed to validate new tags")
	}
}
func TestWriteEmptyTagsFlac(t *testing.T) {
	path, _ := filepath.Abs("../testdata/testdata-flac-nonEmpty.flac")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	f.Seek(0, 0)
	tag, err := parse(f, ParseOptions{FLAC})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	tag.filePath = path
	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{FLAC})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "" || tag.Album() != "" || tag.Title() != "" {
		t.Fatal("Failed to remove tags for empty tag test!")
	}
}
func TestWriteTagsFlacFromEmpty(t *testing.T) {
	path, _ := filepath.Abs("../testdata/testdata-flac-nonEmpty.flac")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	f.Seek(0, 0)
	tag, err := parse(f, ParseOptions{FLAC})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	tag.filePath = path
	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{FLAC})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	tag.filePath = path
	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{FLAC})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "TestArtist1" || tag.Album() != "TestAlbum1" || tag.Title() != "TestTitle1" {
		t.Fatal("Failed to validate new tags")
	}
}
func TestUpdateTagsFlac(t *testing.T) {
	TestWriteTagsFlacFromEmpty(t)
	path, _ := filepath.Abs("../testdata/testdata-flac-nonEmpty.flac")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	f.Seek(0, 0)
	tag, err := parse(f, ParseOptions{FLAC})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	tag.filePath = path
	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{FLAC})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	tag.filePath = path
	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	tag.SetArtist("TestArtist2")
	tag.filePath = path
	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{FLAC})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "TestArtist2" || tag.Album() != "TestAlbum1" || tag.Title() != "TestTitle1" {
		t.Fatal("Failed to validate new tags")
	}
}
func TestWriteEmptyTagsOggVorbis(t *testing.T) {
	path, _ := filepath.Abs("../testdata/testdata-ogg-vorbis-nonEmpty.ogg")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	f.Seek(0, 0)
	tag, err := parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	tag.filePath = path
	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "" || tag.Album() != "" || tag.Title() != "" {
		t.Fatal("Failed to remove tags for empty tag test!")
	}
}
func TestWriteTagsOggVorbisFromEmpty(t *testing.T) {
	path, _ := filepath.Abs("../testdata/testdata-ogg-vorbis-nonEmpty.ogg")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	f.Seek(0, 0)
	tag, err := parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	tag.filePath = path
	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	tag.SetAlbumArtFromFilePath("../testdata/testdata-img-1.jpg")
	tag.filePath = path
	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "TestArtist1" || tag.Album() != "TestAlbum1" || tag.Title() != "TestTitle1" {
		t.Fatal("Failed to validate new tags")
	}
}
func TestWriteTagsOggVorbisFromEmptyExtended(t *testing.T) {
	path, _ := filepath.Abs("../testdata/testdata-ogg-vorbis-nonEmpty.ogg")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	f.Seek(0, 0)
	tag, err := parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	tag.filePath = path
	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	tag.SetAlbumArtFromFilePath("../testdata/testdata-img-1.jpg")
	tag.SetAdditionalTag("TEST", "TEST")
	tag.SetAdditionalTag("TEST2", "TEST2")
	tag.filePath = path
	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "TestArtist1" || tag.Album() != "TestAlbum1" || tag.Title() != "TestTitle1" {
		t.Fatal("Failed to validate new tags")
	}
	if tag.PassThrough["TEST"] != "TEST" || tag.PassThrough["TEST2"] != "TEST2" {
		t.Fatal("Extended Tags Not Found")
	}
}
func TestUpdateTagsOggVorbis(t *testing.T) {
	TestWriteTagsOggVorbisFromEmpty(t)
	path, _ := filepath.Abs("../testdata/testdata-ogg-vorbis-nonEmpty.ogg")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	f.Seek(0, 0)
	tag, err := parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	tag.filePath = path
	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	tag.filePath = path
	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	tag.SetArtist("TestArtist2")
	tag.filePath = path
	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "TestArtist2" || tag.Album() != "TestAlbum1" || tag.Title() != "TestTitle1" {
		t.Fatal("Failed to validate new tags")
	}
}
func TestUpdateTagsOggVorbisExtended(t *testing.T) {
	TestWriteTagsOggVorbisFromEmpty(t)
	path, _ := filepath.Abs("../testdata/testdata-ogg-vorbis-nonEmpty.ogg")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	f.Seek(0, 0)
	tag, err := parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	tag.filePath = path
	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	tag.SetAdditionalTag("TEST", "TEST")
	tag.SetAdditionalTag("TEST2", "TEST2")
	tag.filePath = path
	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	tag.SetArtist("TestArtist2")
	tag.SetAdditionalTag("TEST", "TEST3")
	tag.filePath = path
	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "TestArtist2" || tag.Album() != "TestAlbum1" || tag.Title() != "TestTitle1" {
		t.Fatal("Failed to validate new tags")
	}
	if tag.PassThrough["TEST"] != "TEST3" || tag.PassThrough["TEST2"] != "TEST2" {
		t.Fatal("Extended Tags Not Found")
	}
}
func TestWriteEmptyTagsOggOpus(t *testing.T) {
	path, _ := filepath.Abs("../testdata/testdata-opus-nonEmpty.ogg")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	f.Seek(0, 0)
	tag, err := parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	tag.filePath = path
	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "" || tag.Album() != "" || tag.Title() != "" {
		t.Fatal("Failed to remove tags for empty tag test!")
	}
}
func TestWriteTagsOggOpusFromEmpty(t *testing.T) {
	path, _ := filepath.Abs("../testdata/testdata-opus-nonEmpty.ogg")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	f.Seek(0, 0)
	tag, err := parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	tag.filePath = path
	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	tag.SetAlbumArtFromFilePath("../testdata/testdata-img-1.jpg")
	tag.filePath = path
	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "TestArtist1" || tag.Album() != "TestAlbum1" || tag.Title() != "TestTitle1" {
		t.Fatal("Failed to validate new tags")
	}
}
func TestWriteTagsOggOpusFromEmptyExtended(t *testing.T) {
	path, _ := filepath.Abs("../testdata/testdata-opus-nonEmpty.ogg")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	f.Seek(0, 0)
	tag, err := parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	tag.filePath = path
	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	tag.SetAlbumArtFromFilePath("../testdata/testdata-img-1.jpg")
	tag.SetAdditionalTag("TEST", "TEST")
	tag.SetAdditionalTag("TEST2", "TEST2")
	tag.filePath = path
	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "TestArtist1" || tag.Album() != "TestAlbum1" || tag.Title() != "TestTitle1" {
		t.Fatal("Failed to validate new tags")
	}
	if tag.PassThrough["TEST"] != "TEST" || tag.PassThrough["TEST2"] != "TEST2" {
		t.Fatal("Extended Tags Not Found")
	}
}
func TestUpdateTagsOggOpus(t *testing.T) {
	TestWriteTagsOggOpusFromEmpty(t)
	path, _ := filepath.Abs("../testdata/testdata-opus-nonEmpty.ogg")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	f.Seek(0, 0)
	tag, err := parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	tag.filePath = path
	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	tag.filePath = path
	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	tag.SetArtist("TestArtist2")
	tag.filePath = path
	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "TestArtist2" || tag.Album() != "TestAlbum1" || tag.Title() != "TestTitle1" {
		t.Fatal("Failed to validate new tags")
	}
}
func TestUpdateTagsOggOpusExtended(t *testing.T) {
	TestWriteTagsOggOpusFromEmpty(t)
	path, _ := filepath.Abs("../testdata/testdata-opus-nonEmpty.ogg")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	f.Seek(0, 0)
	tag, err := parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()
	tag.filePath = path
	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	tag.SetAdditionalTag("TEST", "TEST")
	tag.SetAdditionalTag("TEST2", "TEST2")
	tag.filePath = path
	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	tag.SetArtist("TestArtist2")
	tag.SetAdditionalTag("TEST", "TEST3")
	tag.filePath = path
	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	WriteFile(path, buffy.Bytes())
	f.Seek(0, 0)
	tag, err = parse(f, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "TestArtist2" || tag.Album() != "TestAlbum1" || tag.Title() != "TestTitle1" {
		t.Fatal("Failed to validate new tags")
	}
	if tag.PassThrough["TEST"] != "TEST3" || tag.PassThrough["TEST2"] != "TEST2" {
		t.Fatal("Extended Tags Not Found")
	}
}
