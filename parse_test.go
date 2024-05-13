package audiometa

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadMP3Tags(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-mp3.mp3")
	f, err := os.Open(path)
	assert.NoError(t, err)
	tag, err := parse(f, ParseOptions{MP3})
	assert.NoError(t, err)
	assert.NotEmpty(t, tag.Artist())
	assert.NotEmpty(t, tag.Album())
	assert.NotEmpty(t, tag.Title())

}
func TestReadM4ATags(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-m4a.m4a")
	f, err := os.Open(path)
	assert.NoError(t, err)
	tag, err := parse(f, ParseOptions{M4A})
	assert.NoError(t, err)
	assert.NotEmpty(t, tag.Artist())
	assert.NotEmpty(t, tag.Album())
	assert.NotEmpty(t, tag.Title())
}
func TestReadFlacTags(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-flac.flac")
	f, err := os.Open(path)
	assert.NoError(t, err)
	tag, err := parse(f, ParseOptions{FLAC})
	assert.NoError(t, err)
	assert.NotEmpty(t, tag.Artist())
	assert.NotEmpty(t, tag.Album())
	assert.NotEmpty(t, tag.Title())
}
func TestReadOggVorbisTags(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-ogg.ogg")
	f, err := os.Open(path)
	assert.NoError(t, err)
	tag, err := parse(f, ParseOptions{OGG})
	assert.NoError(t, err)
	assert.NotEmpty(t, tag.Artist())
	assert.NotEmpty(t, tag.Album())
	assert.NotEmpty(t, tag.Title())
}
func TestReadOggOpusTags(t *testing.T) {
	path, _ := filepath.Abs("testdata/testdata-opus.ogg")
	f, err := os.Open(path)
	assert.NoError(t, err)
	tag, err := parse(f, ParseOptions{OGG})
	assert.NoError(t, err)
	assert.NotEmpty(t, tag.Artist())
	assert.NotEmpty(t, tag.Album())
	assert.NotEmpty(t, tag.Title())
}
