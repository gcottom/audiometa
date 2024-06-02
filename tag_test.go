package audiometa

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTag(t *testing.T) {
	artist := "the talented guy"
	albumArtist := "also the talented guy"
	album := "I couldn't come up with a name EP"
	comments := "some comments that I wrote"
	composer := "bob the composer man"
	genre := "Heavy Metal"
	title := "the title for thou I am"
	year := "2024"
	bpm := "107"
	copyrightMsg := "don't steal things"
	date := "05-31-2024"
	encodedBy := "me: the encoder"
	lyricist := "the lyrics guy"
	fileType := "mp3"
	language := "english"
	length := "3:08"
	partOfSet := "false"
	publisher := "someone rich"

	path, _ := filepath.Abs("testdata/testdata-mp3.mp3")
	f, err := os.Open(path)
	assert.NoError(t, err)
	tag, err := parse(f, ParseOptions{MP3})
	assert.NoError(t, err)

	t.Run("test set album", func(t *testing.T) {
		tag.SetAlbum(album)
		alb := tag.album
		assert.Equal(t, album, alb)
	})
	t.Run("test set/get album", func(t *testing.T) {
		tag.SetAlbum(album)
		alb := tag.Album()
		assert.Equal(t, album, alb)
	})
	t.Run("test set/get album artist", func(t *testing.T) {
		tag.SetAlbumArtist(albumArtist)
		aart := tag.AlbumArtist()
		assert.Equal(t, albumArtist, aart)
	})
	t.Run("test set album artist", func(t *testing.T) {
		tag.SetAlbumArtist(albumArtist)
		aart := tag.albumArtist
		assert.Equal(t, albumArtist, aart)
	})
	t.Run("test set/get artist", func(t *testing.T) {
		tag.SetArtist(artist)
		art := tag.Artist()
		assert.Equal(t, artist, art)
	})
	t.Run("test set artist", func(t *testing.T) {
		tag.SetArtist(artist)
		art := tag.artist
		assert.Equal(t, artist, art)
	})
	t.Run("test set bpm", func(t *testing.T) {
		tag.SetBPM(bpm)
		bp := tag.bpm
		assert.Equal(t, bpm, bp)
	})
	t.Run("test set/get bpm", func(t *testing.T) {
		tag.SetBPM(bpm)
		bp := tag.BPM()
		assert.Equal(t, bpm, bp)
	})
	t.Run("test set comments", func(t *testing.T) {
		tag.SetComments(comments)
		cmts := tag.comments
		assert.Equal(t, comments, cmts)
	})
	t.Run("test set/get comments", func(t *testing.T) {
		tag.SetComments(comments)
		cmts := tag.Comments()
		assert.Equal(t, comments, cmts)
	})
	t.Run("test set composer", func(t *testing.T) {
		tag.SetComposer(composer)
		cmpsr := tag.composer
		assert.Equal(t, composer, cmpsr)
	})
	t.Run("test set/get composer", func(t *testing.T) {
		tag.SetComposer(composer)
		cmpsr := tag.Composer()
		assert.Equal(t, composer, cmpsr)
	})
	t.Run("test set copyright", func(t *testing.T) {
		tag.SetCopyrightMsg(copyrightMsg)
		cprt := tag.copyrightMsg
		assert.Equal(t, copyrightMsg, cprt)
	})
	t.Run("test set/get copyright", func(t *testing.T) {
		tag.SetCopyrightMsg(copyrightMsg)
		cprt := tag.CopyrightMsg()
		assert.Equal(t, copyrightMsg, cprt)
	})
	t.Run("test set date", func(t *testing.T) {
		tag.SetDate(date)
		dte := tag.date
		assert.Equal(t, date, dte)
	})
	t.Run("test set/get date", func(t *testing.T) {
		tag.SetDate(date)
		dte := tag.Date()
		assert.Equal(t, date, dte)
	})
	t.Run("test set/get encoded by", func(t *testing.T) {
		tag.SetEncodedBy(encodedBy)
		enc := tag.EncodedBy()
		assert.Equal(t, encodedBy, enc)
	})
	t.Run("test set encoded by", func(t *testing.T) {
		tag.SetEncodedBy(encodedBy)
		enc := tag.encodedBy
		assert.Equal(t, encodedBy, enc)
	})
	t.Run("test set filetype", func(t *testing.T) {
		tag.SetFileType(fileType)
		ft := tag.fileType
		assert.Equal(t, fileType, ft)
	})
	t.Run("test set/get filetype", func(t *testing.T) {
		tag.SetFileType(fileType)
		ft := tag.FileType()
		assert.Equal(t, fileType, ft)
	})
	t.Run("test set language", func(t *testing.T) {
		tag.SetLanguage(language)
		lang := tag.language
		assert.Equal(t, language, lang)
	})
	t.Run("test set/get language", func(t *testing.T) {
		tag.SetLanguage(language)
		lang := tag.Language()
		assert.Equal(t, language, lang)
	})
	t.Run("test set length", func(t *testing.T) {
		tag.SetLength(length)
		lgth := tag.length
		assert.Equal(t, length, lgth)
	})
	t.Run("test set/get length", func(t *testing.T) {
		tag.SetLength(length)
		lgth := tag.Length()
		assert.Equal(t, length, lgth)
	})
	t.Run("test set lyricist", func(t *testing.T) {
		tag.SetLyricist(lyricist)
		lcst := tag.lyricist
		assert.Equal(t, lyricist, lcst)
	})
	t.Run("test set/get lyricist", func(t *testing.T) {
		tag.SetLyricist(lyricist)
		lcst := tag.Lyricist()
		assert.Equal(t, lyricist, lcst)
	})
	t.Run("test set part of set", func(t *testing.T) {
		tag.SetPartOfSet(partOfSet)
		pos := tag.partOfSet
		assert.Equal(t, partOfSet, pos)
	})
	t.Run("test set/get part of set", func(t *testing.T) {
		tag.SetPartOfSet(partOfSet)
		pos := tag.PartOfSet()
		assert.Equal(t, partOfSet, pos)
	})
	t.Run("test set publisher", func(t *testing.T) {
		tag.SetPublisher(publisher)
		pub := tag.publisher
		assert.Equal(t, publisher, pub)
	})
	t.Run("test set/get publisher", func(t *testing.T) {
		tag.SetPublisher(publisher)
		pub := tag.Publisher()
		assert.Equal(t, publisher, pub)
	})
	t.Run("test set/get title", func(t *testing.T) {
		tag.SetTitle(title)
		titl := tag.Title()
		assert.Equal(t, title, titl)
	})
	t.Run("test set title", func(t *testing.T) {
		tag.SetTitle(title)
		titl := tag.title
		assert.Equal(t, title, titl)
	})
	t.Run("test set/get year", func(t *testing.T) {
		tag.SetYear(year)
		yr := tag.Year()
		assert.Equal(t, year, yr)
	})
	t.Run("test set year", func(t *testing.T) {
		tag.SetYear(year)
		yr := tag.year
		assert.Equal(t, year, yr)
	})
	t.Run("test set genre", func(t *testing.T) {
		tag.SetGenre(genre)
		gnr := tag.genre
		assert.Equal(t, genre, gnr)
	})
	t.Run("test set/get genre", func(t *testing.T) {
		tag.SetGenre(genre)
		gnr := tag.Genre()
		assert.Equal(t, genre, gnr)
	})
}
func TestOpenTagFromPath_NoParseOptions(t *testing.T) {
	filepath := "testdata-mp3.mp3"
	file, err := os.Create(filepath)
	assert.NoError(t, err, "Expected no error creating test file")
	defer os.Remove(filepath)
	file.Close()

	idTag, err := OpenTagFromPath(filepath)
	assert.NoError(t, err, "Expected no error opening tag from path")
	assert.NotNil(t, idTag, "Expected non-nil IDTag")
}

func TestOpenTagFromPath_WithParseOptions(t *testing.T) {
	filepath := "testdata-mp3.mp3"
	file, err := os.Create(filepath)
	assert.NoError(t, err, "Expected no error creating test file")
	defer os.Remove(filepath)
	file.Close()

	parseOptions := ParseOptions{Format: MP3}
	idTag, err := OpenTagFromPath(filepath, parseOptions)
	assert.NoError(t, err, "Expected no error opening tag from path")
	assert.NotNil(t, idTag, "Expected non-nil IDTag")
}

func TestOpenTagFromPath_WithParseOptionsInvalidFileType(t *testing.T) {
	filepath := "testdata-mp3.mp3"
	file, err := os.Create(filepath)
	assert.NoError(t, err, "Expected no error creating test file")
	defer os.Remove(filepath)
	file.Close()

	parseOptions := ParseOptions{Format: "noooo"}
	idTag, err := OpenTagFromPath(filepath, parseOptions)
	assert.Error(t, err)
	assert.Nil(t, idTag)
}

func TestOpenTagFromPath_NoParseOptionsInvalidFileType(t *testing.T) {
	filepath := "testdata-text.txt"
	file, err := os.Create(filepath)
	assert.NoError(t, err, "Expected no error creating test file")
	defer os.Remove(filepath)
	file.Close()

	idTag, err := OpenTagFromPath(filepath)
	assert.Error(t, err)
	assert.Nil(t, idTag)
}

func TestAlbumArtTags(t *testing.T) {
	t.Run("TestSetImageFromByteSlice", func(t *testing.T) {
		path, _ := filepath.Abs("testdata/testdata-mp3-nonEmpty.mp3")
		f, err := os.Open(path)
		assert.NoError(t, err)
		b, err := io.ReadAll(f)
		assert.NoError(t, err)
		r := bytes.NewReader(b)
		tag, err := parse(r, ParseOptions{MP3})
		assert.NoError(t, err)
		p, err := filepath.Abs("./testdata/withAlbumArt/testdata-img-1.jpg")
		assert.NoError(t, err)
		img, err := os.Open(p)
		assert.NoError(t, err)

		decImg, _, err := image.Decode(img)
		assert.NoError(t, err)
		buf := new(bytes.Buffer)
		err = jpeg.Encode(buf, decImg, nil)
		assert.NoError(t, err)

		err = tag.SetAlbumArtFromByteArray(buf.Bytes())
		assert.NoError(t, err)
		picFile, err := os.Open(p)
		assert.NoError(t, err)
		picData, _, err := image.Decode(picFile)
		assert.NoError(t, err)
		img1data := image_2_array_at(picData)
		img2data := image_2_array_at(*tag.albumArt)

		assert.True(t, compareImages(img1data, img2data))
	})

	t.Run("TestSetImageFromImage", func(t *testing.T) {
		path, _ := filepath.Abs("testdata/testdata-mp3-nonEmpty.mp3")
		f, err := os.Open(path)
		assert.NoError(t, err)
		b, err := io.ReadAll(f)
		assert.NoError(t, err)
		r := bytes.NewReader(b)
		tag, err := parse(r, ParseOptions{MP3})
		assert.NoError(t, err)
		p, err := filepath.Abs("./testdata/withAlbumArt/testdata-img-1.jpg")
		assert.NoError(t, err)
		img, err := os.Open(p)
		assert.NoError(t, err)
		decImg, _, err := image.Decode(img)
		assert.NoError(t, err)

		tag.SetAlbumArtFromImage(&decImg)
		picFile, err := os.Open(p)
		assert.NoError(t, err)
		picData, _, err := image.Decode(picFile)
		assert.NoError(t, err)
		img1data := image_2_array_at(picData)
		img2data := image_2_array_at(*tag.albumArt)

		assert.True(t, compareImages(img1data, img2data))
	})
}

func TestPassthroughOGG(t *testing.T) {
	t.Run("TestPassThroughMap-OGG", func(t *testing.T) {
		path, _ := filepath.Abs("testdata/testdata-ogg-vorbis-nonEmpty.ogg")
		f, err := os.Open(path)
		assert.NoError(t, err)
		b, err := io.ReadAll(f)
		assert.NoError(t, err)
		r := bytes.NewReader(b)
		tag, err := parse(r, ParseOptions{OGG})
		assert.NoError(t, err)

		tag.SetAdditionalTag("testtag", "testvalue")
		buffy := new(bytes.Buffer)
		err = tag.Save(buffy)
		assert.NoError(t, err)
		reader := bytes.NewReader(buffy.Bytes())
		tag, err = parse(reader, ParseOptions{OGG})
		assert.NoError(t, err)

		assert.Equal(t, (tag.AdditionalTags())["TESTTAG"], "testvalue")

	})
}
