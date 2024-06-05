package audiometa

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"image"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/aler9/writerseeker"
	"github.com/sunfish-shogi/bufseekio"
)

var (
	vorbisCommentPrefix = []byte("\x03vorbis")
	opusTagsPrefix      = []byte("OpusTags")
	crcLookup           [8][256]uint32
)

func init() {
	initCRC32Table()
}

func initCRC32Table() {
	var i, j int
	var polynomial uint32 = 0x04C11DB7
	var crc uint32

	for i = 0; i <= 0xFF; i++ {
		crc = uint32(i) << 24

		for j = 0; j < 8; j++ {
			if crc&(1<<31) != 0 {
				crc = (crc << 1) ^ polynomial
			} else {
				crc = crc << 1
			}
		}

		crcLookup[0][i] = crc
	}

	for i = 0; i <= 0xFF; i++ {
		for j = 1; j < 8; j++ {
			crcLookup[j][i] = crcLookup[0][(crcLookup[j-1][i]>>24)&0xFF] ^ (crcLookup[j-1][i] << 8)
		}
	}
}

// _osUpdateCRC updates the CRC with the given buffer and size
func _osUpdateCRC(crc uint32, buffer []byte, size int) uint32 {
	i := 0
	for size >= 8 {
		crc ^= (uint32(buffer[i]) << 24) | (uint32(buffer[i+1]) << 16) | (uint32(buffer[i+2]) << 8) | uint32(buffer[i+3])

		crc = crcLookup[7][crc>>24] ^ crcLookup[6][(crc>>16)&0xFF] ^
			crcLookup[5][(crc>>8)&0xFF] ^ crcLookup[4][crc&0xFF] ^
			crcLookup[3][buffer[i+4]] ^ crcLookup[2][buffer[i+5]] ^
			crcLookup[1][buffer[i+6]] ^ crcLookup[0][buffer[i+7]]

		i += 8
		size -= 8
	}

	for size > 0 {
		crc = (crc << 8) ^ crcLookup[0][((crc>>24)&0xFF)^uint32(buffer[i])]
		i++
		size--
	}

	return crc
}

// OggPageChecksumSet sets the checksum for the Ogg page
func OggPageChecksumSet(og *oggPage) {
	if og != nil {
		var crcReg uint32
		buf := make([]byte, 4)
		og.Header.CRC = uint32(0)
		buf[0] = 0
		buf[1] = 0
		buf[2] = 0
		buf[3] = 0

		crcReg = _osUpdateCRC(crcReg, og.Header.toBytesSlice(), len(og.Header.toBytesSlice()))
		crcReg = _osUpdateCRC(crcReg, og.Body, len(og.Body))

		buf[0] = byte(crcReg & 0xFF)
		buf[1] = byte((crcReg >> 8) & 0xFF)
		buf[2] = byte((crcReg >> 16) & 0xFF)
		buf[3] = byte((crcReg >> 24) & 0xFF)

		og.Header.CRC = binary.LittleEndian.Uint32(buf)
	}
}

type oggDemuxer struct {
	packetBufs map[uint32]*bytes.Buffer
}

// Read ogg packets, can return empty slice of packets and nil err
// if more data is needed
func (o *oggDemuxer) read(r io.Reader) ([][]byte, error) {
	var oh oggPageHeader
	if err := binary.Read(r, binary.LittleEndian, &oh); err != nil {
		return nil, err
	}

	if !bytes.Equal(oh.Magic[:], []byte("OggS")) {
		return nil, ErrOggInvalidHeader
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

	if o.packetBufs == nil {
		o.packetBufs = map[uint32]*bytes.Buffer{}
	}

	var packetBuf *bytes.Buffer
	continued := oh.Flags&0x1 != 0
	if continued {
		if b, ok := o.packetBufs[oh.SerialNumber]; ok {
			packetBuf = b
		} else {
			return nil, ErrOggMissingCOP
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

// ReadOggTags reads Ogg metadata from the io.ReadSeeker, returning the resulting
// metadata in a Metadata implementation, or non-nil error if there was a problem.
func readOggTags(r io.Reader) (*IDTag, error) {
	od := &oggDemuxer{}
	for {
		bs, err := od.read(r)
		if err != nil {
			return nil, err
		}

		for _, b := range bs {
			switch {
			case bytes.HasPrefix(b, vorbisCommentPrefix):
				m := &metadataOgg{
					newMetadataVorbis(),
				}
				resultTag, err := m.readVorbisComment(bytes.NewReader(b[len(vorbisCommentPrefix):]))
				resultTag.codec = "vorbis"
				return resultTag, err
			case bytes.HasPrefix(b, opusTagsPrefix):
				m := &metadataOgg{
					newMetadataVorbis(),
				}
				resultTag, err := m.readVorbisComment(bytes.NewReader(b[len(opusTagsPrefix):]))
				resultTag.codec = "opus"
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

type metadataOgg struct {
	*metadataVorbis
}

type metadataVorbis struct {
	c map[string]string // the vorbis comments
	p []byte
}

// Read the vorbis comments from an ogg vorbis or ogg opus file
func (m *metadataVorbis) readVorbisComment(r io.Reader) (*IDTag, error) {
	var resultTag IDTag
	resultTag.PassThrough = make(map[string]string)
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
		split := strings.Split(cmt, "=")
		if len(split) == 2 {
			temp := strings.ToUpper(split[0])
			if temp != "ALBUM" && temp != "ARTIST" && temp != "ALBUMARTIST" && temp != "DATE" && temp != "TITLE" && temp != "GENRE" && temp != "COMMENT" && temp != "COPYRIGHT" && temp != "PUBLISHER" && temp != "METADATA_BLOCK_PICTURE" && temp != "COMPOSER" {
				resultTag.PassThrough[temp] = split[1]
			} else {
				m.c[temp] = split[1]
			}
		}
	}
	resultTag.album = m.c["ALBUM"]
	resultTag.artist = m.c["ARTIST"]
	resultTag.albumArtist = m.c["ALBUMARTIST"]
	resultTag.date = m.c["DATE"]
	resultTag.title = m.c["TITLE"]
	resultTag.genre = m.c["GENRE"]
	resultTag.comments = m.c["COMMENT"]
	resultTag.copyrightMsg = m.c["COPYRIGHT"]
	resultTag.publisher = m.c["PUBLISHER"]
	resultTag.composer = m.c["COMPOSER"]

	if b64data, ok := m.c["METADATA_BLOCK_PICTURE"]; ok {
		data, err := base64.StdEncoding.DecodeString(b64data)
		if err != nil {
			return nil, err
		}
		if err = m.readPictureBlock(bytes.NewReader(data)); err != nil {
			return nil, err
		}
	}
	if len(m.p) > 0 {
		if img, _, err := image.Decode(bytes.NewReader(m.p)); err == nil {
			resultTag.albumArt = &img
		}
	}
	return &resultTag, nil
}

// Read the vorbis comment picture block
func (m *metadataVorbis) readPictureBlock(r io.Reader) error {
	//skipping picture type
	if _, err := readInt(r, 4); err != nil {
		return err
	}
	mimeLen, err := readUint(r, 4)
	if err != nil {
		return err
	}
	//skipping mime type
	if _, err := readString(r, mimeLen); err != nil {
		return err
	}
	descLen, err := readUint(r, 4)
	if err != nil {
		return err
	}
	//skipping description
	if _, err := readString(r, descLen); err != nil {
		return err
	}

	//skip width <32>, height <32>, colorDepth <32>, coloresUsed <32>

	// width
	if _, err = readInt(r, 4); err != nil {
		return err
	}
	// height
	if _, err = readInt(r, 4); err != nil {
		return err
	}
	// color depth
	if _, err = readInt(r, 4); err != nil {
		return err
	}
	// colors used
	if _, err = readInt(r, 4); err != nil {
		return err
	}

	dataLen, err := readInt(r, 4)
	if err != nil {
		return err
	}
	data := make([]byte, dataLen)
	if _, err = io.ReadFull(r, data); err != nil {
		return err
	}

	m.p = data
	return nil
}

// Saves the tags for an ogg Opus file
func saveOpusTags(tag *IDTag, w io.Writer) error {
	needsTemp := reflect.TypeOf(w) == reflect.TypeOf(new(os.File))
	var t *writerseeker.WriterSeeker
	var encoder *oggEncoder
	if needsTemp {
		//in and out are the same file so we have to temp it
		t = &writerseeker.WriterSeeker{}
		defer t.Close()
	}
	readDat, err := io.ReadAll(tag.reader)
	if err != nil {
		return err
	}
	r := bytes.NewReader(readDat)
	decoder := newOggDecoder(r)
	page, err := decoder.decodeOgg()
	if err != nil {
		return err
	}
	if needsTemp {
		encoder = newOggEncoder(page.Header.SerialNumber, t)
	} else {
		encoder = newOggEncoder(page.Header.SerialNumber, w)
	}
	if err = encoder.encodeBOS(page.Header.GranulePosition, page.Packets); err != nil {
		return err
	}
	var vorbisCommentPage *oggPage
	for {
		page, err := decoder.decodeOgg()
		if err != nil {
			if err == io.EOF {
				break // Reached the end of the input Ogg stream
			}
			return err
		}

		// Find the Vorbis comment page and store it
		if hasOpusCommentPrefix(page.Packets) {
			vorbisCommentPage = &page
			// Step 5: Prepare the new Vorbis comment packet with updated metadata and album art
			commentFields := []string{}
			if tag.album != "" {
				commentFields = append(commentFields, "ALBUM="+tag.album)
			}
			if tag.artist != "" {
				commentFields = append(commentFields, "ARTIST="+tag.artist)
			}
			if tag.genre != "" {
				commentFields = append(commentFields, "GENRE="+tag.genre)
			}
			if tag.title != "" {
				commentFields = append(commentFields, "TITLE="+tag.title)
			}
			if tag.date != "" {
				commentFields = append(commentFields, "DATE="+tag.date)
			}
			if tag.albumArtist != "" {
				commentFields = append(commentFields, "ALBUMARTIST="+tag.albumArtist)
			}
			if tag.comments != "" {
				commentFields = append(commentFields, "COMMENT="+tag.comments)
			}
			if tag.publisher != "" {
				commentFields = append(commentFields, "PUBLISHER="+tag.publisher)
			}
			if tag.copyrightMsg != "" {
				commentFields = append(commentFields, "COPYRIGHT="+tag.copyrightMsg)
			}
			if tag.composer != "" {
				commentFields = append(commentFields, "COMPOSER="+tag.composer)
			}
			for key, value := range tag.PassThrough {
				commentFields = append(commentFields, key+"="+value)
			}
			img := []byte{}
			if tag.albumArt != nil {
				// Convert album art image to JPEG format
				buf := new(bytes.Buffer)
				if err := jpeg.Encode(buf, *tag.albumArt, nil); err == nil {
					img, _ = createMetadataBlockPicture(buf.Bytes())
				}

			}

			// Create the new Vorbis comment packet
			commentPacket := createOpusCommentPacket(commentFields, img)

			// Replace the Vorbis comment packet in the original page with the new packet
			vorbisCommentPage.Packets[0] = commentPacket

			// Step 6: Write the updated Vorbis comment page to the output file
			if err = encoder.encode(vorbisCommentPage.Header.GranulePosition, vorbisCommentPage.Packets); err != nil {
				return err
			}
		} else {
			// Write non-Vorbis comment pages to the output file
			if page.Header.Flags == EOS {
				if err = encoder.encodeEOS(page.Header.GranulePosition, page.Packets); err != nil {
					return err
				}
			} else {
				if err = encoder.encode(page.Header.GranulePosition, page.Packets); err != nil {
					return err
				}
			}
		}
	}
	// Step 7: Close and rename the files to the original file
	if needsTemp {
		f := w.(*os.File)
		path, err := filepath.Abs(f.Name())
		if err != nil {
			return err
		}
		w2, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		defer w2.Close()
		if _, err := io.Copy(w2, bytes.NewReader(t.Bytes())); err != nil {
			return err
		}
		if _, err = f.Seek(0, io.SeekEnd); err != nil {
			return err
		}
	}
	return nil
}

// Saves the given tag structure to a ogg vorbis audio file
func saveVorbisTags(tag *IDTag, w io.Writer) error {
	needsTemp := reflect.TypeOf(w) == reflect.TypeOf(new(os.File))
	var t *writerseeker.WriterSeeker
	var encoder *oggEncoder
	if needsTemp {
		//in and out are the same file so we have to temp it
		t = &writerseeker.WriterSeeker{}
		defer t.Close()
	}
	r := bufseekio.NewReadSeeker(tag.reader, 128*1024, 4)
	decoder := newOggDecoder(r)
	page, err := decoder.decodeOgg()
	if err != nil {
		return err
	}
	if needsTemp {
		encoder = newOggEncoder(page.Header.SerialNumber, t)
	} else {
		encoder = newOggEncoder(page.Header.SerialNumber, w)
	}

	if err = encoder.encodeBOS(page.Header.GranulePosition, page.Packets); err != nil {
		return err
	}
	var vorbisCommentPage *oggPage
	for {
		page, err := decoder.decodeOgg()
		if err != nil {
			if err == io.EOF {
				break // Reached the end of the input Ogg stream
			}
			return err
		}

		// Find the Vorbis comment page and store it
		if hasVorbisCommentPrefix(page.Packets) {
			vorbisCommentPage = &page
			commentFields := []string{}
			if tag.album != "" {
				commentFields = append(commentFields, "ALBUM="+tag.album)
			}
			if tag.artist != "" {
				commentFields = append(commentFields, "ARTIST="+tag.artist)
			}
			if tag.genre != "" {
				commentFields = append(commentFields, "GENRE="+tag.genre)
			}
			if tag.title != "" {
				commentFields = append(commentFields, "TITLE="+tag.title)
			}
			if tag.date != "" {
				commentFields = append(commentFields, "DATE="+tag.date)
			}
			if tag.albumArtist != "" {
				commentFields = append(commentFields, "ALBUMARTIST="+tag.albumArtist)
			}
			if tag.comments != "" {
				commentFields = append(commentFields, "COMMENT="+tag.comments)
			}
			if tag.publisher != "" {
				commentFields = append(commentFields, "PUBLISHER="+tag.publisher)
			}
			if tag.composer != "" {
				commentFields = append(commentFields, "COMPOSER="+tag.composer)
			}
			if tag.copyrightMsg != "" {
				commentFields = append(commentFields, "COPYRIGHT="+tag.copyrightMsg)
			}
			for key, value := range tag.PassThrough {
				commentFields = append(commentFields, key+"="+value)
			}
			img := []byte{}
			if tag.albumArt != nil {
				// Convert album art image to JPEG format
				buf := new(bytes.Buffer)
				if err = jpeg.Encode(buf, *tag.albumArt, nil); err == nil {
					img, _ = createMetadataBlockPicture(buf.Bytes())
				}
			}

			// Create the new Vorbis comment packet
			commentPacket := createVorbisCommentPacket(commentFields, img)
			vorbisCommentPage.Packets[0] = commentPacket
			if err = encoder.encode(vorbisCommentPage.Header.GranulePosition, vorbisCommentPage.Packets); err != nil {
				return nil
			}
		} else {
			if page.Header.Flags == EOS {
				if err = encoder.encodeEOS(page.Header.GranulePosition, page.Packets); err != nil {
					return nil
				}
			} else {
				if err = encoder.encode(page.Header.GranulePosition, page.Packets); err != nil {
					return nil
				}
			}
		}
	}
	if needsTemp {
		f := w.(*os.File)
		path, err := filepath.Abs(f.Name())
		if err != nil {
			return err
		}
		w2, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		defer w2.Close()
		if _, err := io.Copy(w2, bytes.NewReader(t.Bytes())); err != nil {
			return err
		}
		if _, err = f.Seek(0, io.SeekEnd); err != nil {
			return err
		}
	}
	return nil
}

func hasOpusCommentPrefix(packets [][]byte) bool {
	return len(packets) > 0 && len(packets[0]) >= 8 && string(packets[0][:8]) == "OpusTags"
}

// Creates the comment packet for the Opus spec from the given commentFields and albumArt. The only difference between vorbis and opus is the "OpusTags" header and the framing bit
func createOpusCommentPacket(commentFields []string, albumArt []byte) []byte {
	vendorString := "audiometa"

	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(len(vendorString)))
	vorbisCommentPacket := append(buf, []byte(vendorString)...)

	if len(albumArt) > 0 {
		binary.LittleEndian.PutUint32(buf, uint32(len(commentFields)+1))
	} else {
		binary.LittleEndian.PutUint32(buf, uint32(len(commentFields)))
	}
	vorbisCommentPacket = append(vorbisCommentPacket, buf...)

	for _, field := range commentFields {
		binary.LittleEndian.PutUint32(buf, uint32(len(field)))
		vorbisCommentPacket = append(vorbisCommentPacket, buf...)
		vorbisCommentPacket = append(vorbisCommentPacket, []byte(field)...)
	}
	vorbisCommentPacket = append([]byte("OpusTags"), vorbisCommentPacket...)
	if len(albumArt) > 1 {
		albumArtBase64 := base64.StdEncoding.EncodeToString(albumArt)
		fieldLength := len("METADATA_BLOCK_PICTURE=") + len(albumArtBase64)
		binary.LittleEndian.PutUint32(buf, uint32(fieldLength))
		vorbisCommentPacket = append(vorbisCommentPacket, buf...)
		vorbisCommentPacket = append(vorbisCommentPacket, []byte("METADATA_BLOCK_PICTURE=")...)
		vorbisCommentPacket = append(vorbisCommentPacket, []byte(albumArtBase64)...)
	}
	return vorbisCommentPacket
}

// Checks if the vorbis comment header is present
func hasVorbisCommentPrefix(packets [][]byte) bool {
	return len(packets) > 0 && len(packets[0]) >= 7 && string(packets[0][:7]) == "\x03vorbis"
}

// Creates the vorbis comment packet from the given commentFields and albumArt
func createVorbisCommentPacket(commentFields []string, albumArt []byte) []byte {
	vendorString := "audiometa"

	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(len(vendorString)))
	vorbisCommentPacket := append(buf, []byte(vendorString)...)
	if len(albumArt) > 0 {
		binary.LittleEndian.PutUint32(buf, uint32(len(commentFields)+1))
	} else {
		binary.LittleEndian.PutUint32(buf, uint32(len(commentFields)))
	}
	vorbisCommentPacket = append(vorbisCommentPacket, buf...)

	for _, field := range commentFields {
		binary.LittleEndian.PutUint32(buf, uint32(len(field)))
		vorbisCommentPacket = append(vorbisCommentPacket, buf...)
		vorbisCommentPacket = append(vorbisCommentPacket, []byte(field)...)
	}
	vorbisCommentPacket = append([]byte("\x03vorbis"), vorbisCommentPacket...)
	if len(albumArt) > 1 {
		albumArtBase64 := base64.StdEncoding.EncodeToString(albumArt)
		fieldLength := len("METADATA_BLOCK_PICTURE=") + len(albumArtBase64)
		binary.LittleEndian.PutUint32(buf, uint32(fieldLength))
		vorbisCommentPacket = append(vorbisCommentPacket, buf...)
		vorbisCommentPacket = append(vorbisCommentPacket, []byte("METADATA_BLOCK_PICTURE=")...)
		vorbisCommentPacket = append(vorbisCommentPacket, []byte(albumArtBase64)...)
	}

	vorbisCommentPacket = append(vorbisCommentPacket, []byte("\x01")...)
	return vorbisCommentPacket
}

// Creates the picture block which holds the album art in the vorbis comment header
func createMetadataBlockPicture(albumArtData []byte) ([]byte, error) {
	mimeType := "image/jpeg"
	description := "Cover"
	img, _, err := image.DecodeConfig(bytes.NewReader(albumArtData))
	if err != nil {
		return nil, ErrOggImgConfigFail
	}
	res := bytes.NewBuffer([]byte{})
	res.Write(encodeUint32(uint32(3)))
	res.Write(encodeUint32(uint32(len(mimeType))))
	res.Write([]byte(mimeType))
	res.Write(encodeUint32(uint32(len(description))))
	res.Write([]byte(description))
	res.Write(encodeUint32(uint32(img.Width)))
	res.Write(encodeUint32(uint32(img.Height)))
	res.Write(encodeUint32(24))
	res.Write(encodeUint32(0))
	res.Write(encodeUint32(uint32(len(albumArtData))))
	res.Write(albumArtData)
	return res.Bytes(), nil
}
