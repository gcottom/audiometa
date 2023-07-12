package mp3mp4tag

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	vorbisCommentPrefix = []byte("\x03vorbis")
	opusTagsPrefix      = []byte("OpusTags")
)

var oggCRC32Poly04c11db7 = oggCRCTable(0x04c11db7)

type crc32Table [256]uint32

func oggCRCTable(poly uint32) *crc32Table {
	var t crc32Table

	for i := 0; i < 256; i++ {
		crc := uint32(i) << 24
		for j := 0; j < 8; j++ {
			if crc&0x80000000 != 0 {
				crc = (crc << 1) ^ poly
			} else {
				crc <<= 1
			}
		}
		t[i] = crc
	}

	return &t
}

func oggCRCUpdate(crc uint32, tab *crc32Table, p []byte) uint32 {
	for _, v := range p {
		crc = (crc << 8) ^ tab[byte(crc>>24)^v]
	}
	return crc
}

type oggPageHeader struct {
	Magic           [4]byte // "OggS"
	Version         uint8
	Flags           uint8
	GranulePosition uint64
	SerialNumber    uint32
	SequenceNumber  uint32
	CRC             uint32
	Segments        uint8
}

type oggDemuxer struct {
	packetBufs map[uint32]*bytes.Buffer
}

// Read ogg packets, can return empty slice of packets and nil err
// if more data is needed
func (o *oggDemuxer) Read(r io.Reader) ([][]byte, error) {
	headerBuf := &bytes.Buffer{}
	var oh oggPageHeader
	if err := binary.Read(io.TeeReader(r, headerBuf), binary.LittleEndian, &oh); err != nil {
		return nil, err
	}

	if !bytes.Equal(oh.Magic[:], []byte("OggS")) {
		return nil, errors.New("expected 'OggS'")
	}

	segmentTable := make([]byte, oh.Segments)
	if _, err := io.ReadFull(r, segmentTable); err != nil {
		return nil, err
	}
	var segmentsSize int64
	for _, s := range segmentTable {
		segmentsSize += int64(s)
	}
	segmentsData := make([]byte, segmentsSize)
	if _, err := io.ReadFull(r, segmentsData); err != nil {
		return nil, err
	}

	headerBytes := headerBuf.Bytes()
	// reset CRC to zero in header before checksum
	headerBytes[22] = 0
	headerBytes[23] = 0
	headerBytes[24] = 0
	headerBytes[25] = 0
	crc := oggCRCUpdate(0, oggCRC32Poly04c11db7, headerBytes)
	crc = oggCRCUpdate(crc, oggCRC32Poly04c11db7, segmentTable)
	crc = oggCRCUpdate(crc, oggCRC32Poly04c11db7, segmentsData)
	if crc != oh.CRC {
		return nil, fmt.Errorf("expected crc %x != %x", oh.CRC, crc)
	}

	if o.packetBufs == nil {
		o.packetBufs = map[uint32]*bytes.Buffer{}
	}

	var packetBuf *bytes.Buffer
	continued := oh.Flags&0x1 != 0
	if continued {
		if b, ok := o.packetBufs[oh.SerialNumber]; ok {
			packetBuf = b
		} else {
			return nil, fmt.Errorf("could not find continued packet %d", oh.SerialNumber)
		}
	} else {
		packetBuf = &bytes.Buffer{}
	}

	var packets [][]byte
	var p int
	for _, s := range segmentTable {
		packetBuf.Write(segmentsData[p : p+int(s)])
		if s < 255 {
			packets = append(packets, packetBuf.Bytes())
			packetBuf = &bytes.Buffer{}
		}
		p += int(s)
	}

	o.packetBufs[oh.SerialNumber] = packetBuf

	return packets, nil
}

// ReadOGGTags reads OGG metadata from the io.ReadSeeker, returning the resulting
// metadata in a Metadata implementation, or non-nil error if there was a problem.
// See http://www.xiph.org/vorbis/doc/Vorbis_I_spec.html
// and http://www.xiph.org/ogg/doc/framing.html for details.
// For Opus see https://tools.ietf.org/html/rfc7845
func ReadOGGTags(r io.Reader) (*IDTag, error) {
	od := &oggDemuxer{}
	for {
		bs, err := od.Read(r)
		if err != nil {
			return nil, err
		}

		for _, b := range bs {
			switch {
			case bytes.HasPrefix(b, vorbisCommentPrefix):
				m := &metadataOGG{
					newMetadataVorbis(),
				}
				resultTag, err := m.readVorbisComment(bytes.NewReader(b[len(vorbisCommentPrefix):]))
				return resultTag, err
			case bytes.HasPrefix(b, opusTagsPrefix):
				m := &metadataOGG{
					newMetadataVorbis(),
				}
				resultTag, err := m.readVorbisComment(bytes.NewReader(b[len(opusTagsPrefix):]))
				return resultTag, err
			}
		}
	}
}
func newMetadataVorbis() *metadataVorbis {
	return &metadataVorbis{
		c: make(map[string]string),
	}
}

type metadataOGG struct {
	*metadataVorbis
}

type metadataVorbis struct {
	c map[string]string // the vorbis comments

}

func (m *metadataVorbis) readVorbisComment(r io.Reader) (*IDTag, error) {
	var resultTag IDTag
	vendorLen, err := readUint32LittleEndian(r)
	if err != nil {
		return nil, err
	}

	vendor, err := readString(r, uint(vendorLen))
	if err != nil {
		return nil, err
	}
	m.c["vendor"] = vendor

	commentsLen, err := readUint32LittleEndian(r)
	if err != nil {
		return nil, err
	}

	for i := uint32(0); i < commentsLen; i++ {
		l, err := readUint32LittleEndian(r)
		if err != nil {
			return nil, err
		}
		cmt, err := readString(r, uint(l))
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(cmt, "album=") {
			tag := strings.Replace(cmt, "album=", "", 1)
			resultTag.album = tag
		} else if strings.HasPrefix(cmt, "ALBUM=") {
			tag := strings.Replace(cmt, "ALBUM=", "", 1)
			resultTag.album = tag
		} else if strings.HasPrefix(cmt, "artist=") {
			tag := strings.Replace(cmt, "artist=", "", 1)
			resultTag.artist = tag
		} else if strings.HasPrefix(cmt, "ARTIST=") {
			tag := strings.Replace(cmt, "ARTIST=", "", 1)
			resultTag.artist = tag
		} else if strings.HasPrefix(cmt, "date=") {
			tag := strings.Replace(cmt, "date=", "", 1)
			resultTag.id3.date = tag
		} else if strings.HasPrefix(cmt, "DATE=") {
			tag := strings.Replace(cmt, "DATE=", "", 1)
			resultTag.id3.date = tag
		} else if strings.HasPrefix(cmt, "title=") {
			tag := strings.Replace(cmt, "title=", "", 1)
			resultTag.title = tag
		} else if strings.HasPrefix(cmt, "TITLE=") {
			tag := strings.Replace(cmt, "TITLE=", "", 1)
			resultTag.title = tag
		} else if strings.HasPrefix(cmt, "genre=") {
			tag := strings.Replace(cmt, "genre=", "", 1)
			resultTag.genre = tag
		} else if strings.HasPrefix(cmt, "GENRE=") {
			tag := strings.Replace(cmt, "GENRE=", "", 1)
			resultTag.genre = tag
		}
	}

	return &resultTag, nil
}
