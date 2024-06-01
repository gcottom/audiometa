package audiometa

import (
	"bytes"
	"encoding/binary"
)

// The MIME type as defined in RFC 3534.
const MIMEType = "application/ogg"

const headerSize = 27

// max segment size
const maxSegSize = 255

// max sequence-of-segments size in a page
const mps = maxSegSize * 255

// == 65307, per the RFC
const maxPageSize = headerSize + maxSegSize + mps

// The byte order of integers in ogg page headers.
var byteOrder = binary.LittleEndian

type oggPageHeader struct {
	Magic           [4]byte // 0-3, always == "OggS"
	Version         byte    // 4, always == 0
	Flags           byte    // 5 Flags is a bitmask of COP, BOS, and/or EOS.
	GranulePosition int64   // 6-13, codec-specific, GranulePosition represents the granule position, its interpretation depends on the encapsulated codec.
	SerialNumber    uint32  // 14-17, associated with a logical stream, SerialNumber represents the bitstream serial number.
	SequenceNumber  uint32  // 18-21, sequence number of page in packet
	CRC             uint32  // 22-25
	Segments        byte    // 26
}

const (
	// Continuation of packet
	COP byte = 1 << iota
	// Beginning of stream
	BOS = 1 << iota
	// End of stream
	EOS = 1 << iota
)

func (o oggPageHeader) toBytesSlice() []byte {
	b := new(bytes.Buffer)
	_ = binary.Write(b, byteOrder, o.Magic)
	_ = binary.Write(b, byteOrder, o.Version)
	_ = binary.Write(b, byteOrder, o.Flags)
	_ = binary.Write(b, byteOrder, o.GranulePosition)
	_ = binary.Write(b, byteOrder, o.SerialNumber)
	_ = binary.Write(b, byteOrder, o.SequenceNumber)
	_ = binary.Write(b, byteOrder, o.CRC)
	_ = binary.Write(b, byteOrder, o.Segments)
	return b.Bytes()
}
func (o oggPageHeader) toBytesBuffer() *bytes.Buffer {
	b := new(bytes.Buffer)
	_ = binary.Write(b, byteOrder, o.Magic)
	_ = binary.Write(b, byteOrder, o.Version)
	_ = binary.Write(b, byteOrder, o.Flags)
	_ = binary.Write(b, byteOrder, o.GranulePosition)
	_ = binary.Write(b, byteOrder, o.SerialNumber)
	_ = binary.Write(b, byteOrder, o.SequenceNumber)
	_ = binary.Write(b, byteOrder, o.CRC)
	_ = binary.Write(b, byteOrder, o.Segments)
	return b
}
