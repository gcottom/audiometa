package audiometa

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOggPageHeaderToBytesSlice(t *testing.T) {
	header := oggPageHeader{
		Magic:           [4]byte{'O', 'g', 'g', 'S'},
		Version:         0,
		Flags:           BOS,
		GranulePosition: 12345,
		SerialNumber:    30,
		SequenceNumber:  1,
		CRC:             0,
		Segments:        2,
	}

	expected := []byte{
		'O', 'g', 'g', 'S', // Magic
		0,                        // Version
		BOS,                      // Flags
		57, 48, 0, 0, 0, 0, 0, 0, // GranulePosition
		30, 0, 0, 0, // SerialNumber
		1, 0, 0, 0, // SequenceNumber
		0, 0, 0, 0, // CRC
		2, // Segments
	}

	result := header.toBytesSlice()
	assert.Equal(t, expected, result, "The byte slice representation of the header should match the expected value.")
}

func TestOggPageHeaderToBytesBuffer(t *testing.T) {
	header := oggPageHeader{
		Magic:           [4]byte{'O', 'g', 'g', 'S'},
		Version:         0,
		Flags:           BOS,
		GranulePosition: 12345,
		SerialNumber:    15,
		SequenceNumber:  1,
		CRC:             0,
		Segments:        2,
	}

	expected := []byte{
		'O', 'g', 'g', 'S', // Magic
		0,                        // Version
		BOS,                      // Flags
		57, 48, 0, 0, 0, 0, 0, 0, // GranulePosition
		15, 0, 0, 0, // SerialNumber
		1, 0, 0, 0, // SequenceNumber
		0, 0, 0, 0, // CRC
		2, // Segments
	}

	result := header.toBytesBuffer()
	assert.Equal(t, expected, result.Bytes(), "The byte buffer representation of the header should match the expected value.")
}

func TestCRCFunctions(t *testing.T) {

	buffer := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	crc := uint32(0)

	updatedCRC := _osUpdateCRC(crc, buffer, len(buffer))
	expectedCRC := uint32(0x7d0f3681) // Replace with the actual expected CRC value based on the test buffer

	assert.Equal(t, expectedCRC, updatedCRC, "The CRC should match the expected value.")
}

func TestOggPageChecksumSet(t *testing.T) {
	og := &oggPage{
		Header: &oggPageHeader{
			Magic:           [4]byte{'O', 'g', 'g', 'S'},
			Version:         0,
			Flags:           BOS,
			GranulePosition: 12345,
			SerialNumber:    67890,
			SequenceNumber:  1,
			CRC:             0,
			Segments:        2,
		},
		Body: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
	}

	OggPageChecksumSet(og)

	// Verify CRC in the header
	crcReg := og.Header.CRC
	expectedCRC := uint32(0x9b4396e8)

	assert.Equal(t, expectedCRC, crcReg, "The CRC in the header should match the expected value.")
}

func TestOggRead(t *testing.T) {

	og := &oggPage{
		Header: &oggPageHeader{
			Magic:           [4]byte{'O', 'g', 'g', 'Z'},
			Version:         0,
			Flags:           BOS,
			GranulePosition: 12345,
			SerialNumber:    67890,
			SequenceNumber:  1,
			CRC:             0,
			Segments:        2,
		},
		Body: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
	}
	t.Run("invalid magic", func(t *testing.T) {
		r := bytes.NewReader(og.Header.toBytesSlice())
		dem := oggDemuxer{}
		res, err := dem.read(r)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
	t.Run("empty dat", func(t *testing.T) {
		r := bytes.NewReader((&oggPage{Header: &oggPageHeader{}}).Header.toBytesSlice())
		dem := oggDemuxer{}
		res, err := dem.read(r)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}
