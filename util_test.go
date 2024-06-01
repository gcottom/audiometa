package audiometa

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileType(t *testing.T) {
	tests := []struct {
		filepath string
		expected FileType
		err      error
	}{
		{"file.mp3", "mp3", nil},
		{"document.m4a", "m4a", nil},
		{"image.ogg", "ogg", nil},
		{"noextensionfile", "", errors.New("unsupported file extension or no extension")},
		{"unsupportedfile.txt", "", errors.New("unsupported file extension or no extension")},
	}

	for _, test := range tests {
		ft, err := GetFileType(test.filepath)
		if test.err != nil {
			assert.EqualError(t, err, test.err.Error(), "GetFileType(%s) should return error %v", test.filepath, test.err)
		} else {
			assert.NoError(t, err, "GetFileType(%s) should not return error", test.filepath)
			assert.Equal(t, test.expected, ft, "GetFileType(%s) should return %v", test.filepath, test.expected)
		}
	}
}

func TestGetInt(t *testing.T) {
	tests := []struct {
		b        []byte
		expected int
	}{
		{[]byte{0x00}, 0},
		{[]byte{0x01}, 1},
		{[]byte{0x01, 0x00}, 256},
		{[]byte{0x01, 0x00, 0x01}, 65537},
	}

	for _, test := range tests {
		result := getInt(test.b)
		assert.Equal(t, test.expected, result, "getInt(%v) should return %d", test.b, test.expected)
	}
}

func TestReadInt(t *testing.T) {
	tests := []struct {
		input    []byte
		n        uint
		expected int
		err      error
	}{
		{[]byte{0x01}, 1, 1, nil},
		{[]byte{0x01, 0x00}, 2, 256, nil},
		{[]byte{0x01, 0x00, 0x01}, 3, 65537, nil},
		{[]byte{}, 1, 0, io.EOF},
	}

	for _, test := range tests {
		r := bytes.NewReader(test.input)
		result, err := readInt(r, test.n)
		if test.err != nil {
			assert.EqualError(t, err, test.err.Error(), "readInt(%v, %d) should return error %v", test.input, test.n, test.err)
		} else {
			assert.NoError(t, err, "readInt(%v, %d) should not return error", test.input, test.n)
			assert.Equal(t, test.expected, result, "readInt(%v, %d) should return %d", test.input, test.n, test.expected)
		}
	}
}

func TestReadUint(t *testing.T) {
	tests := []struct {
		input    []byte
		n        uint
		expected uint
		err      error
	}{
		{[]byte{0x01}, 1, 1, nil},
		{[]byte{0x01, 0x00}, 2, 256, nil},
		{[]byte{0x01, 0x00, 0x01}, 3, 65537, nil},
		{[]byte{}, 1, 0, io.EOF},
	}

	for _, test := range tests {
		r := bytes.NewReader(test.input)
		result, err := readUint(r, test.n)
		if test.err != nil {
			assert.EqualError(t, err, test.err.Error(), "readUint(%v, %d) should return error %v", test.input, test.n, test.err)
		} else {
			assert.NoError(t, err, "readUint(%v, %d) should not return error", test.input, test.n)
			assert.Equal(t, test.expected, result, "readUint(%v, %d) should return %d", test.input, test.n, test.expected)
		}
	}
}

func TestReadBytes(t *testing.T) {
	tests := []struct {
		input    []byte
		n        uint
		expected []byte
		err      error
	}{
		{[]byte("hello"), 5, []byte("hello"), nil},
		{[]byte("hello"), 10, nil, errors.New("unexpected EOF")},
		{make([]byte, readBytesMaxUpfront+1), uint(readBytesMaxUpfront + 1), make([]byte, readBytesMaxUpfront+1), nil},
	}

	for _, test := range tests {
		r := bytes.NewReader(test.input)
		result, err := readBytes(r, test.n)
		if test.err != nil {
			assert.EqualError(t, err, test.err.Error(), "readBytes(%v, %d) should return error %v", test.input, test.n, test.err)
		} else {
			assert.NoError(t, err, "readBytes(%v, %d) should not return error", test.input, test.n)
			assert.Equal(t, test.expected, result, "readBytes(%v, %d) should return %v", test.input, test.n, test.expected)
		}
	}
}

func TestReadString(t *testing.T) {
	tests := []struct {
		input    []byte
		n        uint
		expected string
		err      error
	}{
		{[]byte("hello"), 5, "hello", nil},
		{[]byte("hello"), 10, "", errors.New("unexpected EOF")},
	}

	for _, test := range tests {
		r := bytes.NewReader(test.input)
		result, err := readString(r, test.n)
		if test.err != nil {
			assert.EqualError(t, err, test.err.Error(), "readString(%v, %d) should return error %v", test.input, test.n, test.err)
		} else {
			assert.NoError(t, err, "readString(%v, %d) should not return error", test.input, test.n)
			assert.Equal(t, test.expected, result, "readString(%v, %d) should return %s", test.input, test.n, test.expected)
		}
	}
}

func TestReadUint32LittleEndian(t *testing.T) {
	tests := []struct {
		input    []byte
		expected uint32
		err      error
	}{
		{[]byte{0x01, 0x00, 0x00, 0x00}, 1, nil},
		{[]byte{0xff, 0x00, 0x00, 0x00}, 255, nil},
		{[]byte{0x01, 0x00, 0x00}, 0, io.ErrUnexpectedEOF},
	}

	for _, test := range tests {
		r := bytes.NewReader(test.input)
		result, err := readUint32LittleEndian(r)
		if test.err != nil {
			assert.EqualError(t, err, test.err.Error(), "readUint32LittleEndian(%v) should return error %v", test.input, test.err)
		} else {
			assert.NoError(t, err, "readUint32LittleEndian(%v) should not return error", test.input)
			assert.Equal(t, test.expected, result, "readUint32LittleEndian(%v) should return %d", test.input, test.expected)
		}
	}
}

func TestEncodeUint32(t *testing.T) {
	tests := []struct {
		input    uint32
		expected []byte
	}{
		{1, []byte{0x00, 0x00, 0x00, 0x01}},
		{255, []byte{0x00, 0x00, 0x00, 0xff}},
	}

	for _, test := range tests {
		result := encodeUint32(test.input)
		assert.Equal(t, test.expected, result, "encodeUint32(%d) should return %v", test.input, test.expected)
	}
}

func TestFileTypesContains(t *testing.T) {
	tests := []struct {
		input    FileType
		expected bool
	}{
		{"txt", false},
		{"pdf", false},
		{"jpg", false},
		{"png", false},
		{"mp3", true},
		{"m4a", true},
		{"m4b", true},
		{"m4p", true},
		{"mp4", true},
		{"flac", true},
		{"ogg", true},
	}

	for _, test := range tests {
		result := fileTypesContains(test.input, supportedFileTypes)
		assert.Equal(t, test.expected, result, "fileTypesContains(%s, %v) should return %t", test.input, supportedFileTypes, test.expected)
	}
}
