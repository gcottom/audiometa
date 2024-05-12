package audiometa

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMP3(t *testing.T) {
	t.Run("TestWriteEmptyTagsMP3-buffers", func(t *testing.T) {
		path, _ := filepath.Abs("testdata/testdata-mp3-nonEmpty.mp3")
		f, err := os.Open(path)
		assert.NoError(t, err)
		b, err := io.ReadAll(f)
		assert.NoError(t, err)
		r := bytes.NewReader(b)
		tag, err := parse(r, ParseOptions{MP3})
		assert.NoError(t, err)
		tag.ClearAllTags()
		buffy := new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = parse(r, ParseOptions{MP3})
		assert.NoError(t, err)
		assert.Empty(t, tag.Artist())
		assert.Empty(t, tag.Album())
		assert.Empty(t, tag.Title())
	})

	t.Run("TestWriteEmptyTagsMP3-file", func(t *testing.T) {
		err := os.Mkdir("testdata/temp", 0755)
		assert.NoError(t, err)
		of, err := os.ReadFile("testdata/testdata-mp3-nonEmpty.mp3")
		assert.NoError(t, err)
		err = os.WriteFile("testdata/temp/testdata-mp3-nonEmpty.mp3", of, 0755)
		assert.NoError(t, err)
		path, _ := filepath.Abs("testdata/temp/testdata-mp3-nonEmpty.mp3")
		f, err := os.OpenFile(path, os.O_RDONLY, 0755)
		assert.NoError(t, err)
		defer f.Close()
		tag, err := Open(f, ParseOptions{MP3})
		assert.NoError(t, err)
		tag.ClearAllTags()
		err = SaveTag(tag, f)
		assert.NoError(t, err)
		f, err = os.OpenFile(path, os.O_RDONLY, 0755)
		assert.NoError(t, err)
		defer f.Close()
		tag, err = Open(f, ParseOptions{MP3})
		assert.NoError(t, err)
		f.Close()
		err = os.RemoveAll("testdata/temp")
		assert.NoError(t, err)
		assert.Empty(t, tag.Artist())
		assert.Empty(t, tag.Album())
		assert.Empty(t, tag.Title())
	})

	t.Run("TestWriteTagsMP3FromEmpty-buffers", func(t *testing.T) {
		path, _ := filepath.Abs("testdata/testdata-mp3-nonEmpty.mp3")
		f, err := os.Open(path)
		assert.NoError(t, err)
		defer f.Close()
		b, err := io.ReadAll(f)
		assert.NoError(t, err)
		r := bytes.NewReader(b)

		tag, err := Open(r, ParseOptions{MP3})
		assert.NoError(t, err)
		tag.ClearAllTags()

		buffy := new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{MP3})
		assert.NoError(t, err)
		tag.SetArtist("TestArtist1")
		tag.SetTitle("TestTitle1")
		tag.SetAlbum("TestAlbum1")

		buffy = new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{MP3})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
	})

	t.Run("TestWriteTagsMP3FromEmpty-file", func(t *testing.T) {
		err := os.Mkdir("testdata/temp", 0755)
		assert.NoError(t, err)
		of, err := os.ReadFile("testdata/testdata-mp3-nonEmpty.mp3")
		assert.NoError(t, err)
		err = os.WriteFile("testdata/temp/testdata-mp3-nonEmpty.mp3", of, 0755)
		assert.NoError(t, err)
		path, _ := filepath.Abs("testdata/temp/testdata-mp3-nonEmpty.mp3")
		f, err := os.Open(path)
		assert.NoError(t, err)
		defer f.Close()

		tag, err := Open(f, ParseOptions{MP3})
		assert.NoError(t, err)
		tag.SetArtist("TestArtist1")
		tag.SetTitle("TestTitle1")
		tag.SetAlbum("TestAlbum1")

		err = SaveTag(tag, f)
		assert.NoError(t, err)

		f, err = os.Open(path)
		assert.NoError(t, err)
		defer f.Close()
		tag, err = Open(f, ParseOptions{MP3})
		assert.NoError(t, err)
		err = os.RemoveAll("testdata/temp")
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
	})

	t.Run("TestUpdateTagsMP3-buffers", func(t *testing.T) {
		path, _ := filepath.Abs("testdata/testdata-mp3-nonEmpty.mp3")
		f, err := os.Open(path)
		assert.NoError(t, err)
		defer f.Close()
		b, err := io.ReadAll(f)
		assert.NoError(t, err)
		r := bytes.NewReader(b)

		tag, err := Open(r, ParseOptions{MP3})
		assert.NoError(t, err)
		tag.ClearAllTags()

		buffy := new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{MP3})
		assert.NoError(t, err)
		tag.SetArtist("TestArtist1")
		tag.SetTitle("TestTitle1")
		tag.SetAlbum("TestAlbum1")

		buffy = new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{MP3})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")

		tag.SetArtist("TestArtist2")

		buffy = new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)

		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{MP3})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist2")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
	})

	t.Run("TestUpdateTagsMP3-file", func(t *testing.T) {
		err := os.Mkdir("testdata/temp", 0755)
		assert.NoError(t, err)
		of, err := os.ReadFile("testdata/testdata-mp3-nonEmpty.mp3")
		assert.NoError(t, err)
		err = os.WriteFile("testdata/temp/testdata-mp3-nonEmpty.mp3", of, 0755)
		assert.NoError(t, err)
		path, _ := filepath.Abs("testdata/temp/testdata-mp3-nonEmpty.mp3")
		f, err := os.Open(path)
		assert.NoError(t, err)
		defer f.Close()

		tag, err := Open(f, ParseOptions{MP3})
		assert.NoError(t, err)
		tag.SetArtist("TestArtist1")
		tag.SetTitle("TestTitle1")
		tag.SetAlbum("TestAlbum1")
		err = SaveTag(tag, f)
		assert.NoError(t, err)

		f, err = os.Open(path)
		assert.NoError(t, err)
		defer f.Close()
		tag, err = Open(f, ParseOptions{MP3})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")

		tag.SetArtist("TestArtist2")
		err = SaveTag(tag, f)
		assert.NoError(t, err)

		f, err = os.Open(path)
		assert.NoError(t, err)
		defer f.Close()
		tag, err = Open(f, ParseOptions{MP3})
		assert.NoError(t, err)
		f.Close()
		err = os.RemoveAll("testdata/temp")
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist2")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
	})
}

func TestM4A(t *testing.T) {
	t.Run("TestWriteEmptyTagsM4A-buffers", func(t *testing.T) {
		path, _ := filepath.Abs("testdata/testdata-m4a-nonEmpty.m4a")
		f, err := os.Open(path)
		assert.NoError(t, err)
		b, err := io.ReadAll(f)
		assert.NoError(t, err)
		r := bytes.NewReader(b)
		tag, err := parse(r, ParseOptions{M4A})
		assert.NoError(t, err)
		tag.ClearAllTags()
		buffy := new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = parse(r, ParseOptions{M4A})
		assert.NoError(t, err)
		assert.Empty(t, tag.Artist())
		assert.Empty(t, tag.Album())
		assert.Empty(t, tag.Title())
	})

	t.Run("TestWriteEmptyTagsM4A-file", func(t *testing.T) {
		err := os.Mkdir("testdata/temp", 0755)
		assert.NoError(t, err)
		of, err := os.ReadFile("testdata/testdata-m4a-nonEmpty.m4a")
		assert.NoError(t, err)
		err = os.WriteFile("testdata/temp/testdata-m4a-nonEmpty.m4a", of, 0755)
		assert.NoError(t, err)
		path, _ := filepath.Abs("testdata/temp/testdata-m4a-nonEmpty.m4a")
		f, err := os.OpenFile(path, os.O_RDONLY, 0755)
		assert.NoError(t, err)
		defer f.Close()
		tag, err := Open(f, ParseOptions{M4A})
		assert.NoError(t, err)
		tag.ClearAllTags()
		err = SaveTag(tag, f)
		assert.NoError(t, err)
		f, err = os.OpenFile(path, os.O_RDONLY, 0755)
		assert.NoError(t, err)
		defer f.Close()
		tag, err = Open(f, ParseOptions{M4A})
		assert.NoError(t, err)
		f.Close()
		err = os.RemoveAll("testdata/temp")
		assert.NoError(t, err)
		assert.Empty(t, tag.Artist())
		assert.Empty(t, tag.Album())
		assert.Empty(t, tag.Title())
	})

	t.Run("TestWriteTagsM4AFromEmpty-buffers", func(t *testing.T) {
		path, _ := filepath.Abs("testdata/testdata-m4a-nonEmpty.m4a")
		f, err := os.Open(path)
		assert.NoError(t, err)
		defer f.Close()
		b, err := io.ReadAll(f)
		assert.NoError(t, err)
		r := bytes.NewReader(b)

		tag, err := Open(r, ParseOptions{M4A})
		assert.NoError(t, err)
		tag.ClearAllTags()

		buffy := new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{M4A})
		assert.NoError(t, err)
		tag.SetArtist("TestArtist1")
		tag.SetTitle("TestTitle1")
		tag.SetAlbum("TestAlbum1")

		buffy = new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{M4A})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
	})

	t.Run("TestWriteTagsM4AFromEmpty-file", func(t *testing.T) {
		err := os.Mkdir("testdata/temp", 0755)
		assert.NoError(t, err)
		of, err := os.ReadFile("testdata/testdata-m4a-nonEmpty.m4a")
		assert.NoError(t, err)
		err = os.WriteFile("testdata/temp/testdata-m4a-nonEmpty.m4a", of, 0755)
		assert.NoError(t, err)
		path, _ := filepath.Abs("testdata/temp/testdata-m4a-nonEmpty.m4a")
		f, err := os.Open(path)
		assert.NoError(t, err)
		defer f.Close()

		tag, err := Open(f, ParseOptions{M4A})
		assert.NoError(t, err)
		tag.SetArtist("TestArtist1")
		tag.SetTitle("TestTitle1")
		tag.SetAlbum("TestAlbum1")

		err = SaveTag(tag, f)
		assert.NoError(t, err)

		f, err = os.Open(path)
		assert.NoError(t, err)
		defer f.Close()
		tag, err = Open(f, ParseOptions{M4A})
		assert.NoError(t, err)
		err = os.RemoveAll("testdata/temp")
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
	})

	t.Run("TestUpdateTagsM4A-buffers", func(t *testing.T) {
		path, _ := filepath.Abs("testdata/testdata-m4a-nonEmpty.m4a")
		f, err := os.Open(path)
		assert.NoError(t, err)
		defer f.Close()
		b, err := io.ReadAll(f)
		assert.NoError(t, err)
		r := bytes.NewReader(b)

		tag, err := Open(r, ParseOptions{M4A})
		assert.NoError(t, err)
		tag.ClearAllTags()

		buffy := new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{M4A})
		assert.NoError(t, err)
		tag.SetArtist("TestArtist1")
		tag.SetTitle("TestTitle1")
		tag.SetAlbum("TestAlbum1")

		buffy = new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{M4A})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")

		tag.SetArtist("TestArtist2")

		buffy = new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)

		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{M4A})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist2")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
	})

	t.Run("TestUpdateTagsM4A-file", func(t *testing.T) {
		err := os.Mkdir("testdata/temp", 0755)
		assert.NoError(t, err)
		of, err := os.ReadFile("testdata/testdata-m4a-nonEmpty.m4a")
		assert.NoError(t, err)
		err = os.WriteFile("testdata/temp/testdata-m4a-nonEmpty.m4a", of, 0755)
		assert.NoError(t, err)
		path, _ := filepath.Abs("testdata/temp/testdata-m4a-nonEmpty.m4a")
		f, err := os.Open(path)
		assert.NoError(t, err)
		defer f.Close()

		tag, err := Open(f, ParseOptions{M4A})
		assert.NoError(t, err)
		tag.SetArtist("TestArtist1")
		tag.SetTitle("TestTitle1")
		tag.SetAlbum("TestAlbum1")
		err = SaveTag(tag, f)
		assert.NoError(t, err)

		f, err = os.Open(path)
		assert.NoError(t, err)
		defer f.Close()
		tag, err = Open(f, ParseOptions{M4A})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")

		tag.SetArtist("TestArtist2")
		err = SaveTag(tag, f)
		assert.NoError(t, err)

		f, err = os.Open(path)
		assert.NoError(t, err)
		defer f.Close()
		tag, err = Open(f, ParseOptions{M4A})
		assert.NoError(t, err)
		f.Close()
		err = os.RemoveAll("testdata/temp")
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist2")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
	})
}

func TestFLAC(t *testing.T) {
	t.Run("TestWriteEmptyTagsFLAC-buffers", func(t *testing.T) {
		path, _ := filepath.Abs("testdata/testdata-flac-nonEmpty.flac")
		f, err := os.Open(path)
		assert.NoError(t, err)
		b, err := io.ReadAll(f)
		assert.NoError(t, err)
		r := bytes.NewReader(b)
		tag, err := parse(r, ParseOptions{FLAC})
		assert.NoError(t, err)
		tag.ClearAllTags()
		buffy := new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = parse(r, ParseOptions{FLAC})
		assert.NoError(t, err)
		assert.Empty(t, tag.Artist())
		assert.Empty(t, tag.Album())
		assert.Empty(t, tag.Title())
	})

	t.Run("TestWriteEmptyTagsFLAC-file", func(t *testing.T) {
		err := os.Mkdir("testdata/temp", 0755)
		assert.NoError(t, err)
		of, err := os.ReadFile("testdata/testdata-flac-nonEmpty.flac")
		assert.NoError(t, err)
		err = os.WriteFile("testdata/temp/testdata-flac-nonEmpty.flac", of, 0755)
		assert.NoError(t, err)
		path, _ := filepath.Abs("testdata/temp/testdata-flac-nonEmpty.flac")
		f, err := os.OpenFile(path, os.O_RDONLY, 0755)
		assert.NoError(t, err)
		defer f.Close()
		tag, err := Open(f, ParseOptions{FLAC})
		assert.NoError(t, err)
		tag.ClearAllTags()
		err = SaveTag(tag, f)
		assert.NoError(t, err)
		f, err = os.OpenFile(path, os.O_RDONLY, 0755)
		assert.NoError(t, err)
		defer f.Close()
		tag, err = Open(f, ParseOptions{FLAC})
		assert.NoError(t, err)
		f.Close()
		err = os.RemoveAll("testdata/temp")
		assert.NoError(t, err)
		assert.Empty(t, tag.Artist())
		assert.Empty(t, tag.Album())
		assert.Empty(t, tag.Title())
	})

	t.Run("TestWriteTagsFLACFromEmpty-buffers", func(t *testing.T) {
		path, _ := filepath.Abs("testdata/testdata-flac-nonEmpty.flac")
		f, err := os.Open(path)
		assert.NoError(t, err)
		defer f.Close()
		b, err := io.ReadAll(f)
		assert.NoError(t, err)
		r := bytes.NewReader(b)

		tag, err := Open(r, ParseOptions{FLAC})
		assert.NoError(t, err)
		tag.ClearAllTags()

		buffy := new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{FLAC})
		assert.NoError(t, err)
		tag.SetArtist("TestArtist1")
		tag.SetTitle("TestTitle1")
		tag.SetAlbum("TestAlbum1")

		buffy = new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{FLAC})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
	})

	t.Run("TestWriteTagsFLACFromEmpty-file", func(t *testing.T) {
		err := os.Mkdir("testdata/temp", 0755)
		assert.NoError(t, err)
		of, err := os.ReadFile("testdata/testdata-flac-nonEmpty.flac")
		assert.NoError(t, err)
		err = os.WriteFile("testdata/temp/testdata-flac-nonEmpty.flac", of, 0755)
		assert.NoError(t, err)
		path, _ := filepath.Abs("testdata/temp/testdata-flac-nonEmpty.flac")
		f, err := os.Open(path)
		assert.NoError(t, err)
		defer f.Close()

		tag, err := Open(f, ParseOptions{FLAC})
		assert.NoError(t, err)
		tag.SetArtist("TestArtist1")
		tag.SetTitle("TestTitle1")
		tag.SetAlbum("TestAlbum1")

		err = SaveTag(tag, f)
		assert.NoError(t, err)

		f, err = os.Open(path)
		assert.NoError(t, err)
		defer f.Close()
		tag, err = Open(f, ParseOptions{FLAC})
		assert.NoError(t, err)
		err = os.RemoveAll("testdata/temp")
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
	})

	t.Run("TestUpdateTagsFLAC-buffers", func(t *testing.T) {
		path, _ := filepath.Abs("testdata/testdata-flac-nonEmpty.flac")
		f, err := os.Open(path)
		assert.NoError(t, err)
		defer f.Close()
		b, err := io.ReadAll(f)
		assert.NoError(t, err)
		r := bytes.NewReader(b)

		tag, err := Open(r, ParseOptions{FLAC})
		assert.NoError(t, err)
		tag.ClearAllTags()

		buffy := new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{FLAC})
		assert.NoError(t, err)
		tag.SetArtist("TestArtist1")
		tag.SetTitle("TestTitle1")
		tag.SetAlbum("TestAlbum1")

		buffy = new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{FLAC})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")

		tag.SetArtist("TestArtist2")

		buffy = new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)

		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{FLAC})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist2")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
	})

	t.Run("TestUpdateTagsFLAC-file", func(t *testing.T) {
		err := os.Mkdir("testdata/temp", 0755)
		assert.NoError(t, err)
		of, err := os.ReadFile("testdata/testdata-flac-nonEmpty.flac")
		assert.NoError(t, err)
		err = os.WriteFile("testdata/temp/testdata-flac-nonEmpty.flac", of, 0755)
		assert.NoError(t, err)
		path, _ := filepath.Abs("testdata/temp/testdata-flac-nonEmpty.flac")
		f, err := os.Open(path)
		assert.NoError(t, err)
		defer f.Close()

		tag, err := Open(f, ParseOptions{FLAC})
		assert.NoError(t, err)
		tag.SetArtist("TestArtist1")
		tag.SetTitle("TestTitle1")
		tag.SetAlbum("TestAlbum1")
		err = SaveTag(tag, f)
		assert.NoError(t, err)

		f, err = os.Open(path)
		assert.NoError(t, err)
		defer f.Close()
		tag, err = Open(f, ParseOptions{FLAC})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")

		tag.SetArtist("TestArtist2")
		err = SaveTag(tag, f)
		assert.NoError(t, err)

		f, err = os.Open(path)
		assert.NoError(t, err)
		defer f.Close()
		tag, err = Open(f, ParseOptions{FLAC})
		assert.NoError(t, err)
		f.Close()
		err = os.RemoveAll("testdata/temp")
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist2")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
	})
}

func TestWriteEmptyTagsOggVorbis(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-ogg-vorbis-nonEmpty.ogg")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal("Error reading file!")
	}
	r := bytes.NewReader(b)
	tag, err := parse(r, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()

	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	r = bytes.NewReader(buffy.Bytes())
	tag, err = parse(r, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "" || tag.Album() != "" || tag.Title() != "" {
		t.Fatal("Failed to remove tags for empty tag test!")
	}
}
func TestWriteTagsOggVorbisFromEmpty(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-ogg-vorbis-nonEmpty.ogg")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal("Error reading file!")
	}
	r := bytes.NewReader(b)
	tag, err := parse(r, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()

	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}
	r = bytes.NewReader(buffy.Bytes())
	tag, err = parse(r, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	tag.SetAlbumArtFromFilePath("testdata/testdata-img-1.jpg")

	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	r = bytes.NewReader(buffy.Bytes())
	tag, err = parse(r, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "TestArtist1" || tag.Album() != "TestAlbum1" || tag.Title() != "TestTitle1" {
		t.Fatal("Failed to validate new tags")
	}
}
func TestWriteTagsOggVorbisFromEmptyExtended(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-ogg-vorbis-nonEmpty.ogg")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal("Error reading file!")
	}
	r := bytes.NewReader(b)
	tag, err := parse(r, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()

	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	r = bytes.NewReader(buffy.Bytes())
	tag, err = parse(r, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	tag.SetAlbumArtFromFilePath("testdata/testdata-img-1.jpg")
	tag.SetAdditionalTag("TEST", "TEST")
	tag.SetAdditionalTag("TEST2", "TEST2")

	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	r = bytes.NewReader(buffy.Bytes())
	tag, err = parse(r, ParseOptions{OGG})
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
	path, _ := filepath.Abs("testdata/testdata-ogg-vorbis-nonEmpty.ogg")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal("Error reading file!")
	}
	r := bytes.NewReader(b)
	tag, err := parse(r, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()

	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	r = bytes.NewReader(buffy.Bytes())
	tag, err = parse(r, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")

	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}
	tag.SetArtist("TestArtist2")

	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	r = bytes.NewReader(buffy.Bytes())
	tag, err = parse(r, ParseOptions{OGG})
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
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal("Error reading file!")
	}
	r := bytes.NewReader(b)
	tag, err := parse(r, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()

	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	r = bytes.NewReader(buffy.Bytes())
	tag, err = parse(r, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	tag.SetAdditionalTag("TEST", "TEST")
	tag.SetAdditionalTag("TEST2", "TEST2")

	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	tag.SetArtist("TestArtist2")
	tag.SetAdditionalTag("TEST", "TEST3")

	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	r = bytes.NewReader(buffy.Bytes())
	tag, err = parse(r, ParseOptions{OGG})
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
	path, _ := filepath.Abs("testdata/testdata-opus-nonEmpty.ogg")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal("Error reading file!")
	}
	r := bytes.NewReader(b)
	tag, err := parse(r, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()

	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	r = bytes.NewReader(buffy.Bytes())
	tag, err = parse(r, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "" || tag.Album() != "" || tag.Title() != "" {
		t.Fatal("Failed to remove tags for empty tag test!")
	}
}
func TestWriteTagsOggOpusFromEmpty(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-opus-nonEmpty.ogg")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal("Error reading file!")
	}
	r := bytes.NewReader(b)
	tag, err := parse(r, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()

	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	r = bytes.NewReader(buffy.Bytes())
	tag, err = parse(r, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	tag.SetAlbumArtFromFilePath("testdata/testdata-img-1.jpg")

	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	r = bytes.NewReader(buffy.Bytes())
	tag, err = parse(r, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	if tag.Artist() != "TestArtist1" || tag.Album() != "TestAlbum1" || tag.Title() != "TestTitle1" {
		t.Fatal("Failed to validate new tags")
	}
}
func TestWriteTagsOggOpusFromEmptyExtended(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-opus-nonEmpty.ogg")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal("Error reading file!")
	}
	r := bytes.NewReader(b)
	tag, err := parse(r, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()

	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	r = bytes.NewReader(buffy.Bytes())
	tag, err = parse(r, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	tag.SetAlbumArtFromFilePath("testdata/testdata-img-1.jpg")
	tag.SetAdditionalTag("TEST", "TEST")
	tag.SetAdditionalTag("TEST2", "TEST2")

	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	r = bytes.NewReader(buffy.Bytes())
	tag, err = parse(r, ParseOptions{OGG})
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
	path, _ := filepath.Abs("testdata/testdata-opus-nonEmpty.ogg")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal("Error reading file!")
	}
	r := bytes.NewReader(b)
	tag, err := parse(r, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()

	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	r = bytes.NewReader(buffy.Bytes())
	tag, err = parse(r, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")

	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}
	tag.SetArtist("TestArtist2")

	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	r = bytes.NewReader(buffy.Bytes())
	tag, err = parse(r, ParseOptions{OGG})
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
	f, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening file!")
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal("Error reading file!")
	}
	r := bytes.NewReader(b)
	tag, err := parse(r, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.ClearAllTags()

	buffy := new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	r = bytes.NewReader(buffy.Bytes())
	tag, err = parse(r, ParseOptions{OGG})
	if err != nil {
		t.Fatal("Error parsing!")
	}
	tag.SetArtist("TestArtist1")
	tag.SetTitle("TestTitle1")
	tag.SetAlbum("TestAlbum1")
	tag.SetAdditionalTag("TEST", "TEST")
	tag.SetAdditionalTag("TEST2", "TEST2")

	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	tag.SetArtist("TestArtist2")
	tag.SetAdditionalTag("TEST", "TEST3")

	buffy = new(bytes.Buffer)
	if err = SaveTag(tag, buffy); err != nil {
		t.Fatal("error saving")
	}

	r = bytes.NewReader(buffy.Bytes())
	tag, err = parse(r, ParseOptions{OGG})
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
