package audiometa

import (
	"image"
	"os"
	"path/filepath"
	"testing"

	"github.com/abema/go-mp4"
	"github.com/aler9/writerseeker"
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

func TestMarshalData(t *testing.T) {
	ws := &writerseeker.WriterSeeker{}
	mockWriter := mp4.NewWriter(ws)
	ctx := mp4.Context{}

	err := marshalData(mockWriter, ctx, "test string")
	assert.Error(t, err)
}

func TestWriteMeta(t *testing.T) {
	ws := &writerseeker.WriterSeeker{}
	mockWriter := mp4.NewWriter(ws)
	ctx := mp4.Context{}
	err := writeMeta(mockWriter, atomsMap["title"], ctx, "test title")
	assert.Error(t, err)
}

func TestWriteExisting(t *testing.T) {
	ws := &writerseeker.WriterSeeker{}
	mockWriter := mp4.NewWriter(ws)
	mockReadHandle := &mp4.ReadHandle{}
	tags := &IDTag{}

	done, err := writeExisting(mockReadHandle, mockWriter, tags, "title")
	assert.NoError(t, err)
	assert.True(t, done)

}

func TestContainsAtom(t *testing.T) {
	boxType := mp4.BoxType{'\251', 'n', 'a', 'm'}
	boxes := []mp4.BoxType{{'\251', 'n', 'a', 'm'}, {'\251', 'A', 'R', 'T'}}
	result := containsAtom(boxType, boxes)
	assert.Equal(t, boxType, result)

	invalidBoxType := mp4.BoxType{'\251', 'i', 'n', 'v'}
	result = containsAtom(invalidBoxType, boxes)
	assert.Equal(t, mp4.BoxType{}, result)
}

func TestContainsTag(t *testing.T) {
	delete := []string{"title", "artist"}
	result := containsTag(delete, "title")
	assert.True(t, result)

	result = containsTag(delete, "genre")
	assert.False(t, result)
}

func TestGetTag(t *testing.T) {
	boxType := mp4.BoxType{'\251', 'n', 'a', 'm'}
	result := getTag(boxType)
	assert.Equal(t, "title", result)

	invalidBoxType := mp4.BoxType{'\251', 'i', 'n', 'v'}
	result = getTag(invalidBoxType)
	assert.Equal(t, "", result)
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
