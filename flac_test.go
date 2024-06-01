package audiometa

import (
	"errors"
	"testing"

	"github.com/gcottom/audiometa/v2/flac"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExtractFLACCommentErrorCases(t *testing.T) {
	// Test case where flac.ParseMetadata returns an error
	// Mock the input io.Reader to return an error when passed to flac.ParseMetadata
	input := &mockReader{}
	_, _, err := extractFLACComment(input)
	assert.Error(t, err)
}

func TestRemoveFLACMetaBlock(t *testing.T) {
	// Test case where s is out of bounds
	// Create a slice with some elements
	slice := []*flac.MetaDataBlock{{}, {}, {}}
	// Call removeFLACMetaBlock with an index greater than the length of the slice
	result := removeFLACMetaBlock(slice, 1)
	assert.Len(t, result, 2)
	// Add more test cases as needed
}

func TestFLACSaveErrorCases(t *testing.T) {
	// Test case where needsTemp is true and creating temp file fails
	// Mock the input io.Reader and io.Writer
	input := &mockReader{}
	output := &mockWriter{}
	// Set needsTemp to true
	needsTemp := true
	err := flacSave(input, output, []*flac.MetaDataBlock{}, needsTemp)
	assert.Error(t, err)
	// Add more test cases as needed
}

func TestFLACSave(t *testing.T) {

	t.Run("Error Writing FLAC Header", func(t *testing.T) {
		mockReader := new(mockReader2)
		mockWriter := new(mockWriter2)
		metaBlocks := []*flac.MetaDataBlock{
			// Add mock MetaDataBlocks as needed
		}
		mockWriter.On("Write", mock.Anything).Return(0, errors.New("write error"))

		err := flacSave(mockReader, mockWriter, metaBlocks, false)
		assert.Error(t, err)
		assert.Equal(t, "write error", err.Error())
		mockWriter.AssertExpectations(t)
	})

	t.Run("Error Writing MetaDataBlock", func(t *testing.T) {
		mockReader := new(mockReader2)
		mockWriter := new(mockWriter2)
		metaBlocks := []*flac.MetaDataBlock{
			// Add mock MetaDataBlocks as needed
		}
		mockReader.On("Read", mock.Anything).Return(8, nil)
		mockWriter.On("Write", mock.Anything).Return(0, errors.New("write error"))

		err := flacSave(mockReader, mockWriter, metaBlocks, false)
		assert.Error(t, err)
		assert.Equal(t, "write error", err.Error())
		mockWriter.AssertExpectations(t)
	})

	t.Run("Error Copying Data", func(t *testing.T) {
		mockReader := new(mockReader2)
		mockWriter := new(mockWriter2)
		metaBlocks := []*flac.MetaDataBlock{
			// Add mock MetaDataBlocks as needed
		}
		mockWriter.On("Write", mock.Anything).Return(0, nil)
		mockReader.On("Read", mock.Anything).Return(0, errors.New("read error"))

		err := flacSave(mockReader, mockWriter, metaBlocks, false)
		assert.Error(t, err)
		assert.Equal(t, "read error", err.Error())
		mockWriter.AssertExpectations(t)
	})
}

// Helper structs for mocking io.Reader and io.Writer
type mockReader struct{}

func (m *mockReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error while mock reading")
}

type mockWriter struct{}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("mock error writing to writer")
}

type mockReader2 struct {
	mock.Mock
}

func (m *mockReader2) Read(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

type mockWriter2 struct {
	mock.Mock
}

func (m *mockWriter2) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *mockWriter2) Seek(offset int64, whence int) (int64, error) {
	args := m.Called(offset, whence)
	return args.Get(0).(int64), args.Error(1)
}
