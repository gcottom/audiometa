package audiometa

import (
	"image"
	"os"
	"path/filepath"
	"testing"

	"github.com/abema/go-mp4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWriter struct {
	mock.Mock
}

func (m *MockWriter) StartBox(boxInfo *mp4.BoxInfo) (n int, err error) {
	args := m.Called(boxInfo)
	return args.Int(0), args.Error(1)
}

func (m *MockWriter) EndBox() (n int, err error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockWriter) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func TestGetAtomsList(t *testing.T) {
	result := getAtomsList()
	expected := []mp4.BoxType{
		{'\251', 'a', 'l', 'b'},
		{'a', 'A', 'R', 'T'},
		{'\251', 'A', 'R', 'T'},
		{'\251', 'c', 'm', 't'},
		{'\251', 'w', 'r', 't'},
		{'c', 'p', 'r', 't'},
		{'c', 'o', 'v', 'r'},
		{'\251', 'g', 'e', 'n'},
		{'\251', 'n', 'a', 'm'},
		{'\251', 'd', 'a', 'y'},
	}
	assert.ElementsMatch(t, expected, result)
}

func TestMetadataGetFunctions(t *testing.T) {
	imgPath, err := filepath.Abs("./testdata/withAlbumArt/testdata-img-1.jpg")
	assert.NoError(t, err)
	f, err := os.Open(imgPath)
	assert.NoError(t, err)

	imgData, _, err := image.Decode(f)
	assert.NoError(t, err)

	m := metadataMP4{
		data: map[string]interface{}{"\xa9day": "2008-08-09", "\xa9cmt": "a comment about comments", "covr": &imgData, "intData": 6},
	}
	t.Run("test getInt", func(t *testing.T) {
		intData := m.getInt([]string{"intData"})
		assert.Equal(t, 6, intData)
	})
	t.Run("test getYear", func(t *testing.T) {
		year := m.year()
		assert.Equal(t, 2008, year)
	})
	t.Run("test getComment", func(t *testing.T) {
		cmnt := m.comment()
		assert.Equal(t, "a comment about comments", cmnt)
	})
	t.Run("test getPicture", func(t *testing.T) {
		pic := m.picture()
		picFile, err := os.Open(imgPath)
		assert.NoError(t, err)
		picData, _, err := image.Decode(picFile)
		assert.NoError(t, err)
		img1data := image_2_array_at(picData)
		img2data := image_2_array_at(*pic)

		assert.True(t, compareImages(img1data, img2data))
	})
}
