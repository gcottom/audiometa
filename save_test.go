package audiometa

import (
	"bytes"
	"errors"
	"image"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func compareImages(src1 [][][3]float32, src2 [][][3]float32) bool {
	dif := 0
	for i, dat1 := range src1 {
		for j := range dat1 {
			if len(src1[i][j]) != len(src2[i][j]) {
				dif++
			}
		}
	}
	return dif == 0
}

func image_2_array_at(src image.Image) [][][3]float32 {
	bounds := src.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	iaa := make([][][3]float32, height)

	for y := 0; y < height; y++ {
		row := make([][3]float32, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			// Shifting by 8 reduces this to the range [0, 255].
			row[x] = [3]float32{float32(r >> 8), float32(g >> 8), float32(b >> 8)}
		}
		iaa[y] = row
	}

	return iaa
}

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
		_, err = f.Seek(0, io.SeekStart)
		assert.NoError(t, err)
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
		p, err := filepath.Abs("./testdata/withAlbumArt/testdata-img-1.jpg")
		assert.NoError(t, err)
		tag.SetAlbumArtFromFilePath(p)

		buffy = new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{MP3})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")

		picFile, err := os.Open(p)
		assert.NoError(t, err)
		picData, _, err := image.Decode(picFile)
		assert.NoError(t, err)
		img1data := image_2_array_at(picData)
		img2data := image_2_array_at(*tag.albumArt)

		assert.True(t, compareImages(img1data, img2data))
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
		p, err := filepath.Abs("./testdata/withAlbumArt/testdata-img-1.jpg")
		assert.NoError(t, err)
		tag.SetAlbumArtFromFilePath(p)

		err = SaveTag(tag, f)
		assert.NoError(t, err)

		_, err = f.Seek(0, io.SeekStart)
		assert.NoError(t, err)
		tag, err = Open(f, ParseOptions{MP3})
		assert.NoError(t, err)
		err = os.RemoveAll("testdata/temp")
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")

		picFile, err := os.Open(p)
		assert.NoError(t, err)
		picData, _, err := image.Decode(picFile)
		assert.NoError(t, err)
		img1data := image_2_array_at(picData)
		img2data := image_2_array_at(*tag.albumArt)

		assert.True(t, compareImages(img1data, img2data))

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
		p, err := filepath.Abs("./testdata/withAlbumArt/testdata-img-1.jpg")
		assert.NoError(t, err)
		tag.SetAlbumArtFromFilePath(p)

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
		picFile, err := os.Open(p)
		assert.NoError(t, err)
		picData, _, err := image.Decode(picFile)
		assert.NoError(t, err)
		img1data := image_2_array_at(picData)
		img2data := image_2_array_at(*tag.albumArt)

		assert.True(t, compareImages(img1data, img2data))
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
		p, err := filepath.Abs("./testdata/withAlbumArt/testdata-img-1.jpg")
		assert.NoError(t, err)
		tag.SetAlbumArtFromFilePath(p)
		err = SaveTag(tag, f)
		assert.NoError(t, err)

		_, err = f.Seek(0, io.SeekStart)
		assert.NoError(t, err)
		tag, err = Open(f, ParseOptions{MP3})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")

		tag.SetArtist("TestArtist2")
		err = SaveTag(tag, f)
		assert.NoError(t, err)

		_, err = f.Seek(0, io.SeekStart)
		assert.NoError(t, err)
		tag, err = Open(f, ParseOptions{MP3})
		assert.NoError(t, err)
		f.Close()
		err = os.RemoveAll("testdata/temp")
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist2")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
		picFile, err := os.Open(p)
		assert.NoError(t, err)
		picData, _, err := image.Decode(picFile)
		assert.NoError(t, err)
		img1data := image_2_array_at(picData)
		img2data := image_2_array_at(*tag.albumArt)

		assert.True(t, compareImages(img1data, img2data))
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
		_, err = f.Seek(0, io.SeekStart)
		assert.NoError(t, err)
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

		_, err = f.Seek(0, io.SeekStart)
		assert.NoError(t, err)
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

		_, err = f.Seek(0, io.SeekStart)
		assert.NoError(t, err)
		tag, err = Open(f, ParseOptions{M4A})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")

		tag.SetArtist("TestArtist2")
		err = SaveTag(tag, f)
		assert.NoError(t, err)

		_, err = f.Seek(0, io.SeekStart)
		assert.NoError(t, err)
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
		_, err = f.Seek(0, io.SeekStart)
		assert.NoError(t, err)
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
		p, err := filepath.Abs("./testdata/withAlbumArt/testdata-img-1.jpg")
		assert.NoError(t, err)
		tag.SetAlbumArtFromFilePath(p)

		buffy = new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{FLAC})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
		picFile, err := os.Open(p)
		assert.NoError(t, err)
		picData, _, err := image.Decode(picFile)
		assert.NoError(t, err)
		img1data := image_2_array_at(picData)
		img2data := image_2_array_at(*tag.albumArt)

		assert.True(t, compareImages(img1data, img2data))
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
		p, err := filepath.Abs("./testdata/withAlbumArt/testdata-img-1.jpg")
		assert.NoError(t, err)
		tag.SetAlbumArtFromFilePath(p)

		err = SaveTag(tag, f)
		assert.NoError(t, err)

		_, err = f.Seek(0, io.SeekStart)
		assert.NoError(t, err)
		tag, err = Open(f, ParseOptions{FLAC})
		assert.NoError(t, err)
		err = os.RemoveAll("testdata/temp")
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
		picFile, err := os.Open(p)
		assert.NoError(t, err)
		picData, _, err := image.Decode(picFile)
		assert.NoError(t, err)
		img1data := image_2_array_at(picData)
		img2data := image_2_array_at(*tag.albumArt)

		assert.True(t, compareImages(img1data, img2data))
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

		_, err = f.Seek(0, io.SeekStart)
		assert.NoError(t, err)
		tag, err = Open(f, ParseOptions{FLAC})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")

		tag.SetArtist("TestArtist2")
		err = SaveTag(tag, f)
		assert.NoError(t, err)

		_, err = f.Seek(0, io.SeekStart)
		assert.NoError(t, err)
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

func TestOggVorbis(t *testing.T) {
	t.Run("TestWriteEmptyTagsOggVorbis-buffers", func(t *testing.T) {
		path, _ := filepath.Abs("testdata/testdata-ogg-vorbis-nonEmpty.ogg")
		f, err := os.Open(path)
		assert.NoError(t, err)
		b, err := io.ReadAll(f)
		assert.NoError(t, err)
		r := bytes.NewReader(b)
		tag, err := parse(r, ParseOptions{OGG})
		assert.NoError(t, err)
		tag.ClearAllTags()
		buffy := new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = parse(r, ParseOptions{OGG})
		assert.NoError(t, err)
		assert.Empty(t, tag.Artist())
		assert.Empty(t, tag.Album())
		assert.Empty(t, tag.Title())
	})

	t.Run("TestWriteEmptyTagsOggVorbis-file", func(t *testing.T) {
		err := os.Mkdir("testdata/temp", 0755)
		assert.NoError(t, err)
		of, err := os.ReadFile("testdata/testdata-ogg-vorbis-nonEmpty.ogg")
		assert.NoError(t, err)
		err = os.WriteFile("testdata/temp/testdata-ogg-vorbis-nonEmpty.ogg", of, 0755)
		assert.NoError(t, err)
		path, _ := filepath.Abs("testdata/temp/testdata-ogg-vorbis-nonEmpty.ogg")
		f, err := os.OpenFile(path, os.O_RDONLY, 0755)
		assert.NoError(t, err)
		defer f.Close()
		tag, err := Open(f, ParseOptions{OGG})
		assert.NoError(t, err)
		tag.ClearAllTags()
		err = SaveTag(tag, f)
		assert.NoError(t, err)
		_, err = f.Seek(0, io.SeekStart)
		assert.NoError(t, err)
		tag, err = Open(f, ParseOptions{OGG})
		assert.NoError(t, err)
		f.Close()
		err = os.RemoveAll("testdata/temp")
		assert.NoError(t, err)
		assert.Empty(t, tag.Artist())
		assert.Empty(t, tag.Album())
		assert.Empty(t, tag.Title())
	})

	t.Run("TestWriteTagsOggVorbisFromEmpty-buffers", func(t *testing.T) {
		path, _ := filepath.Abs("testdata/testdata-ogg-vorbis-nonEmpty.ogg")
		f, err := os.Open(path)
		assert.NoError(t, err)
		defer f.Close()
		b, err := io.ReadAll(f)
		assert.NoError(t, err)
		r := bytes.NewReader(b)

		tag, err := Open(r, ParseOptions{OGG})
		assert.NoError(t, err)
		tag.ClearAllTags()

		buffy := new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{OGG})
		assert.NoError(t, err)
		tag.SetArtist("TestArtist1")
		tag.SetTitle("TestTitle1")
		tag.SetAlbum("TestAlbum1")
		tag.SetGenre("Trap")
		tag.SetDate("2024-06-01")
		tag.SetAlbumArtist("a talented guy")
		tag.SetComments("I wrote some comments about your song")
		tag.SetPublisher("I am the publisher")
		tag.SetCopyrightMsg("hey please don't steal")
		tag.SetComposer("someone composed I suppose")
		p, err := filepath.Abs("./testdata/withAlbumArt/testdata-img-1.jpg")
		assert.NoError(t, err)
		tag.SetAlbumArtFromFilePath(p)

		buffy = new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{OGG})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
		assert.Equal(t, tag.Genre(), "Trap")
		assert.Equal(t, tag.Date(), "2024-06-01")
		assert.Equal(t, tag.AlbumArtist(), "a talented guy")
		assert.Equal(t, tag.Comments(), "I wrote some comments about your song")
		assert.Equal(t, tag.Publisher(), "I am the publisher")
		assert.Equal(t, tag.CopyrightMsg(), "hey please don't steal")
		assert.Equal(t, tag.Composer(), "someone composed I suppose")
		picFile, err := os.Open(p)
		assert.NoError(t, err)
		picData, _, err := image.Decode(picFile)
		assert.NoError(t, err)
		img1data := image_2_array_at(picData)
		img2data := image_2_array_at(*tag.albumArt)

		assert.True(t, compareImages(img1data, img2data))
	})

	t.Run("TestWriteTagsOggVorbisFromEmpty-file", func(t *testing.T) {
		err := os.Mkdir("testdata/temp", 0755)
		assert.NoError(t, err)
		of, err := os.ReadFile("testdata/testdata-ogg-vorbis-nonEmpty.ogg")
		assert.NoError(t, err)
		err = os.WriteFile("testdata/temp/testdata-ogg-vorbis-nonEmpty.ogg", of, 0755)
		assert.NoError(t, err)
		path, _ := filepath.Abs("testdata/temp/testdata-ogg-vorbis-nonEmpty.ogg")
		f, err := os.Open(path)
		assert.NoError(t, err)
		defer f.Close()

		tag, err := Open(f, ParseOptions{OGG})
		assert.NoError(t, err)
		tag.SetArtist("TestArtist1")
		tag.SetTitle("TestTitle1")
		tag.SetAlbum("TestAlbum1")
		p, err := filepath.Abs("./testdata/withAlbumArt/testdata-img-1.jpg")
		assert.NoError(t, err)
		tag.SetAlbumArtFromFilePath(p)

		err = SaveTag(tag, f)
		assert.NoError(t, err)

		_, err = f.Seek(0, io.SeekStart)
		assert.NoError(t, err)
		tag, err = Open(f, ParseOptions{OGG})
		assert.NoError(t, err)
		err = os.RemoveAll("testdata/temp")
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
		picFile, err := os.Open(p)
		assert.NoError(t, err)
		picData, _, err := image.Decode(picFile)
		assert.NoError(t, err)
		img1data := image_2_array_at(picData)
		img2data := image_2_array_at(*tag.albumArt)

		assert.True(t, compareImages(img1data, img2data))
	})

	t.Run("TestUpdateTagsOggVorbis-buffers", func(t *testing.T) {
		path, _ := filepath.Abs("testdata/testdata-ogg-vorbis-nonEmpty.ogg")
		f, err := os.Open(path)
		assert.NoError(t, err)
		defer f.Close()
		b, err := io.ReadAll(f)
		assert.NoError(t, err)
		r := bytes.NewReader(b)

		tag, err := Open(r, ParseOptions{OGG})
		assert.NoError(t, err)
		tag.ClearAllTags()

		buffy := new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{OGG})
		assert.NoError(t, err)
		tag.SetArtist("TestArtist1")
		tag.SetTitle("TestTitle1")
		tag.SetAlbum("TestAlbum1")
		p, err := filepath.Abs("./testdata/withAlbumArt/testdata-img-1.jpg")
		assert.NoError(t, err)
		tag.SetAlbumArtFromFilePath(p)

		buffy = new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{OGG})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")

		tag.SetArtist("TestArtist2")

		buffy = new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)

		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{OGG})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist2")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
		picFile, err := os.Open(p)
		assert.NoError(t, err)
		picData, _, err := image.Decode(picFile)
		assert.NoError(t, err)
		img1data := image_2_array_at(picData)
		img2data := image_2_array_at(*tag.albumArt)

		assert.True(t, compareImages(img1data, img2data))
	})

	t.Run("TestUpdateTagsOggVorbis-file", func(t *testing.T) {
		err := os.Mkdir("testdata/temp", 0755)
		assert.NoError(t, err)
		of, err := os.ReadFile("testdata/testdata-ogg-vorbis-nonEmpty.ogg")
		assert.NoError(t, err)
		err = os.WriteFile("testdata/temp/testdata-ogg-vorbis-nonEmpty.ogg", of, 0755)
		assert.NoError(t, err)
		path, _ := filepath.Abs("testdata/temp/testdata-ogg-vorbis-nonEmpty.ogg")
		f, err := os.Open(path)
		assert.NoError(t, err)
		defer f.Close()

		tag, err := Open(f, ParseOptions{OGG})
		assert.NoError(t, err)
		tag.SetArtist("TestArtist1")
		tag.SetTitle("TestTitle1")
		tag.SetAlbum("TestAlbum1")
		p, err := filepath.Abs("./testdata/withAlbumArt/testdata-img-1.jpg")
		assert.NoError(t, err)
		tag.SetAlbumArtFromFilePath(p)
		err = SaveTag(tag, f)
		assert.NoError(t, err)

		_, err = f.Seek(0, io.SeekStart)
		assert.NoError(t, err)
		tag, err = Open(f, ParseOptions{OGG})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")

		tag.SetArtist("TestArtist2")
		err = SaveTag(tag, f)
		assert.NoError(t, err)

		_, err = f.Seek(0, io.SeekStart)
		assert.NoError(t, err)
		tag, err = Open(f, ParseOptions{OGG})
		assert.NoError(t, err)
		f.Close()
		err = os.RemoveAll("testdata/temp")
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist2")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
		picFile, err := os.Open(p)
		assert.NoError(t, err)
		picData, _, err := image.Decode(picFile)
		assert.NoError(t, err)
		img1data := image_2_array_at(picData)
		img2data := image_2_array_at(*tag.albumArt)

		assert.True(t, compareImages(img1data, img2data))
	})
}

func TestOggOpus(t *testing.T) {
	t.Run("TestWriteEmptyTagsOggOpus-buffers", func(t *testing.T) {
		path, _ := filepath.Abs("testdata/testdata-opus-nonEmpty.ogg")
		f, err := os.Open(path)
		assert.NoError(t, err)
		b, err := io.ReadAll(f)
		assert.NoError(t, err)
		r := bytes.NewReader(b)
		tag, err := parse(r, ParseOptions{OGG})
		assert.NoError(t, err)
		tag.ClearAllTags()
		buffy := new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = parse(r, ParseOptions{OGG})
		assert.NoError(t, err)
		assert.Empty(t, tag.Artist())
		assert.Empty(t, tag.Album())
		assert.Empty(t, tag.Title())
	})

	t.Run("TestWriteEmptyTagsOggOpus-file", func(t *testing.T) {
		err := os.Mkdir("testdata/temp", 0755)
		assert.NoError(t, err)
		of, err := os.ReadFile("testdata/testdata-opus-nonEmpty.ogg")
		assert.NoError(t, err)
		err = os.WriteFile("testdata/temp/testdata-opus-nonEmpty.ogg", of, 0755)
		assert.NoError(t, err)
		path, _ := filepath.Abs("testdata/temp/testdata-opus-nonEmpty.ogg")
		f, err := os.OpenFile(path, os.O_RDONLY, 0755)
		assert.NoError(t, err)
		defer f.Close()
		tag, err := Open(f, ParseOptions{OGG})
		assert.NoError(t, err)
		tag.ClearAllTags()
		err = SaveTag(tag, f)
		assert.NoError(t, err)
		_, err = f.Seek(0, io.SeekStart)
		assert.NoError(t, err)
		tag, err = Open(f, ParseOptions{OGG})
		assert.NoError(t, err)
		f.Close()
		err = os.RemoveAll("testdata/temp")
		assert.NoError(t, err)
		assert.Empty(t, tag.Artist())
		assert.Empty(t, tag.Album())
		assert.Empty(t, tag.Title())
	})

	t.Run("TestWriteTagsOggOpusFromEmpty-buffers", func(t *testing.T) {
		path, _ := filepath.Abs("testdata/testdata-opus-nonEmpty.ogg")
		f, err := os.Open(path)
		assert.NoError(t, err)
		defer f.Close()
		b, err := io.ReadAll(f)
		assert.NoError(t, err)
		r := bytes.NewReader(b)

		tag, err := Open(r, ParseOptions{OGG})
		assert.NoError(t, err)
		tag.ClearAllTags()

		buffy := new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{OGG})
		assert.NoError(t, err)
		tag.SetArtist("TestArtist1")
		tag.SetTitle("TestTitle1")
		tag.SetAlbum("TestAlbum1")
		tag.SetGenre("Trap")
		tag.SetDate("2024-06-01")
		tag.SetAlbumArtist("a talented guy")
		tag.SetComments("I wrote some comments about your song")
		tag.SetPublisher("I am the publisher")
		tag.SetCopyrightMsg("hey please don't steal")
		tag.SetComposer("someone composed I suppose")
		p, err := filepath.Abs("./testdata/withAlbumArt/testdata-img-1.jpg")
		assert.NoError(t, err)
		tag.SetAlbumArtFromFilePath(p)
		buffy = new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)

		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{OGG})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
		assert.Equal(t, tag.Genre(), "Trap")
		assert.Equal(t, tag.Date(), "2024-06-01")
		assert.Equal(t, tag.AlbumArtist(), "a talented guy")
		assert.Equal(t, tag.Comments(), "I wrote some comments about your song")
		assert.Equal(t, tag.Publisher(), "I am the publisher")
		assert.Equal(t, tag.CopyrightMsg(), "hey please don't steal")
		assert.Equal(t, tag.Composer(), "someone composed I suppose")
		picFile, err := os.Open(p)
		assert.NoError(t, err)
		picData, _, err := image.Decode(picFile)
		assert.NoError(t, err)
		img1data := image_2_array_at(picData)
		img2data := image_2_array_at(*tag.albumArt)

		assert.True(t, compareImages(img1data, img2data))
	})

	t.Run("TestWriteTagsOggOpusFromEmpty-file", func(t *testing.T) {
		err := os.Mkdir("testdata/temp", 0755)
		assert.NoError(t, err)
		of, err := os.ReadFile("testdata/testdata-opus-nonEmpty.ogg")
		assert.NoError(t, err)
		err = os.WriteFile("testdata/temp/testdata-opus-nonEmpty.ogg", of, 0755)
		assert.NoError(t, err)
		path, _ := filepath.Abs("testdata/temp/testdata-opus-nonEmpty.ogg")
		f, err := os.Open(path)
		assert.NoError(t, err)
		defer f.Close()

		tag, err := Open(f, ParseOptions{OGG})
		assert.NoError(t, err)
		tag.SetArtist("TestArtist1")
		tag.SetTitle("TestTitle1")
		tag.SetAlbum("TestAlbum1")
		p, err := filepath.Abs("./testdata/withAlbumArt/testdata-img-1.jpg")
		assert.NoError(t, err)
		tag.SetAlbumArtFromFilePath(p)

		err = SaveTag(tag, f)
		assert.NoError(t, err)

		_, err = f.Seek(0, io.SeekStart)
		assert.NoError(t, err)
		tag, err = Open(f, ParseOptions{OGG})
		assert.NoError(t, err)
		err = os.RemoveAll("testdata/temp")
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
		picFile, err := os.Open(p)
		assert.NoError(t, err)
		picData, _, err := image.Decode(picFile)
		assert.NoError(t, err)
		img1data := image_2_array_at(picData)
		img2data := image_2_array_at(*tag.albumArt)

		assert.True(t, compareImages(img1data, img2data))
	})

	t.Run("TestUpdateTagsOggOpus-buffers", func(t *testing.T) {
		path, _ := filepath.Abs("testdata/testdata-opus-nonEmpty.ogg")
		f, err := os.Open(path)
		assert.NoError(t, err)
		defer f.Close()
		b, err := io.ReadAll(f)
		assert.NoError(t, err)
		r := bytes.NewReader(b)

		tag, err := Open(r, ParseOptions{OGG})
		assert.NoError(t, err)
		tag.ClearAllTags()

		buffy := new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{OGG})
		assert.NoError(t, err)
		tag.SetArtist("TestArtist1")
		tag.SetTitle("TestTitle1")
		tag.SetAlbum("TestAlbum1")
		p, err := filepath.Abs("./testdata/withAlbumArt/testdata-img-1.jpg")
		assert.NoError(t, err)
		tag.SetAlbumArtFromFilePath(p)

		buffy = new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)
		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{OGG})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")

		tag.SetArtist("TestArtist2")

		buffy = new(bytes.Buffer)
		err = SaveTag(tag, buffy)
		assert.NoError(t, err)

		r = bytes.NewReader(buffy.Bytes())
		tag, err = Open(r, ParseOptions{OGG})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist2")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
		picFile, err := os.Open(p)
		assert.NoError(t, err)
		picData, _, err := image.Decode(picFile)
		assert.NoError(t, err)
		img1data := image_2_array_at(picData)
		img2data := image_2_array_at(*tag.albumArt)

		assert.True(t, compareImages(img1data, img2data))
	})

	t.Run("TestUpdateTagsOggOpus-file", func(t *testing.T) {
		err := os.Mkdir("testdata/temp", 0755)
		assert.NoError(t, err)
		of, err := os.ReadFile("testdata/testdata-opus-nonEmpty.ogg")
		assert.NoError(t, err)
		err = os.WriteFile("testdata/temp/testdata-opus-nonEmpty.ogg", of, 0755)
		assert.NoError(t, err)
		path, _ := filepath.Abs("testdata/temp/testdata-opus-nonEmpty.ogg")
		f, err := os.Open(path)
		assert.NoError(t, err)
		defer f.Close()

		tag, err := Open(f, ParseOptions{OGG})
		assert.NoError(t, err)
		tag.SetArtist("TestArtist1")
		tag.SetTitle("TestTitle1")
		tag.SetAlbum("TestAlbum1")
		p, err := filepath.Abs("./testdata/withAlbumArt/testdata-img-1.jpg")
		assert.NoError(t, err)
		tag.SetAlbumArtFromFilePath(p)
		err = SaveTag(tag, f)
		assert.NoError(t, err)

		_, err = f.Seek(0, io.SeekStart)
		assert.NoError(t, err)
		tag, err = Open(f, ParseOptions{OGG})
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist1")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")

		tag.SetArtist("TestArtist2")
		err = SaveTag(tag, f)
		assert.NoError(t, err)

		_, err = f.Seek(0, io.SeekStart)
		assert.NoError(t, err)
		tag, err = Open(f, ParseOptions{OGG})
		assert.NoError(t, err)
		f.Close()
		err = os.RemoveAll("testdata/temp")
		assert.NoError(t, err)
		assert.Equal(t, tag.Artist(), "TestArtist2")
		assert.Equal(t, tag.Album(), "TestAlbum1")
		assert.Equal(t, tag.Title(), "TestTitle1")
		picFile, err := os.Open(p)
		assert.NoError(t, err)
		picData, _, err := image.Decode(picFile)
		assert.NoError(t, err)
		img1data := image_2_array_at(picData)
		img2data := image_2_array_at(*tag.albumArt)

		assert.True(t, compareImages(img1data, img2data))
	})
}

type saveMockReader struct {
	mock.Mock
}

func (m *saveMockReader) Seek(offset int64, whence int) (int64, error) {
	args := m.Called(offset, whence)
	return args.Get(0).(int64), args.Error(1)
}

func (m *saveMockReader) Read(p []byte) (int, error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

type saveMockWriter struct {
	mock.Mock
}

func (m *saveMockWriter) Write(p []byte) (int, error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func TestSave_SeekError(t *testing.T) {
	mockReader := new(saveMockReader)
	mockReader.On("Seek", int64(0), io.SeekStart).Return(int64(0), errors.New("seek error"))

	tag := &IDTag{
		reader: mockReader,
	}

	mockWriter := new(saveMockWriter)

	err := tag.Save(mockWriter)
	assert.EqualError(t, err, "seek error")
	mockReader.AssertExpectations(t)
}

func TestSave_MP3Error(t *testing.T) {
	mockReader := new(saveMockReader)
	mockReader.On("Read", mock.Anything).Return(100, io.EOF)
	mockReader.On("Seek", int64(0), io.SeekStart).Return(int64(0), nil)

	tag := &IDTag{
		fileType: "mp3",
		reader:   mockReader,
	}

	mockWriter := new(saveMockWriter)
	mockWriter.On("Write", mock.Anything).Return(0, errors.New("mp3 write error"))

	err := tag.Save(mockWriter)
	assert.EqualError(t, err, "mp3 write error")
	mockWriter.AssertExpectations(t)
}
func TestSave_MP3ErrorRead(t *testing.T) {
	mockReader := new(saveMockReader)
	mockReader.On("Seek", int64(0), io.SeekStart).Return(int64(0), nil)
	mockReader.On("Read", mock.Anything).Return(0, errors.New("read error"))

	tag := &IDTag{
		fileType: "mp3",
		reader:   mockReader,
	}

	mockWriter := new(saveMockWriter)

	err := tag.Save(mockWriter)
	assert.EqualError(t, err, "read error")
	mockWriter.AssertExpectations(t)
}

func TestSave_MP4Error(t *testing.T) {
	mockReader := new(saveMockReader)
	mockReader.On("Read", mock.Anything).Return(400, io.EOF)
	mockReader.On("Seek", int64(0), mock.Anything).Return(int64(0), nil)

	tag := &IDTag{
		fileType: "m4a",
		reader:   mockReader,
	}

	mockWriter := new(saveMockWriter)

	err := tag.Save(mockWriter)
	assert.Error(t, err, "mp4 write error")
	mockWriter.AssertExpectations(t)
}

func TestSave_FLACError(t *testing.T) {
	mockReader := new(saveMockReader)
	mockReader.On("Read", mock.Anything).Return(400, io.EOF)
	mockReader.On("Seek", int64(0), io.SeekStart).Return(int64(0), nil)

	tag := &IDTag{
		fileType: "flac",
		reader:   mockReader,
	}

	mockWriter := new(saveMockWriter)

	err := tag.Save(mockWriter)
	assert.EqualError(t, err, "error parsing flac stream")
	mockReader.AssertExpectations(t)
}

func TestSave_OGGError(t *testing.T) {
	mockReader := new(saveMockReader)
	mockReader.On("Read", mock.Anything).Return(156, io.EOF)
	mockReader.On("Seek", int64(0), io.SeekStart).Return(int64(0), nil)

	tag := &IDTag{
		fileType: "ogg",
		codec:    "vorbis",
		reader:   mockReader,
	}

	mockWriter := new(saveMockWriter)

	err := tag.Save(mockWriter)
	assert.Error(t, err)
	mockReader.AssertExpectations(t)
}

func TestSave_UnsupportedFileType(t *testing.T) {
	mockReader := new(saveMockReader)
	mockReader.On("Seek", int64(0), io.SeekStart).Return(int64(0), nil)

	tag := &IDTag{
		fileType: "UNKNOWN",
		reader:   mockReader,
	}

	mockWriter := new(saveMockWriter)

	err := tag.Save(mockWriter)
	assert.EqualError(t, err, ErrNoMethodAvlble.Error())
	mockReader.AssertExpectations(t)
}
