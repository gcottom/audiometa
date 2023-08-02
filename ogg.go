package mp3mp4tag

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"os"
	"path/filepath"
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

type oggDemuxer struct {
	packetBufs map[uint32]*bytes.Buffer
}

// Read ogg packets, can return empty slice of packets and nil err
// if more data is needed
func (o *oggDemuxer) Read(r io.Reader) ([][]byte, error) {
	headerBuf := &bytes.Buffer{}
	var oh oggPageHeader
	if err := binary.Read(io.TeeReader(r, headerBuf), binary.LittleEndian, &oh); err != nil {
		fmt.Println("Error in binary read")
		return nil, err
	}

	if !bytes.Equal(oh.Magic[:], []byte("OggS")) {
		return nil, errors.New("expected 'OggS'")
	}

	segmentTable := make([]byte, oh.Segments)
	if _, err := io.ReadFull(r, segmentTable); err != nil {
		fmt.Println("Error in segment table")
		return nil, err
	}
	var segmentsSize int64
	for _, s := range segmentTable {
		segmentsSize += int64(s)
	}
	segmentsData := make([]byte, segmentsSize)
	if _, err := io.ReadFull(r, segmentsData); err != nil {
		fmt.Println("Error in segments data")
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
func ReadOGGTags(r io.Reader) (*IDTag, error) {
	od := &oggDemuxer{}
	for {
		bs, err := od.Read(r)
		if err != nil {
			fmt.Println("Error in read function")
			return nil, err
		}

		for _, b := range bs {
			switch {
			case bytes.HasPrefix(b, vorbisCommentPrefix):
				m := &metadataOGG{
					newMetadataVorbis(),
				}
				resultTag, err := m.readVorbisComment(bytes.NewReader(b[len(vorbisCommentPrefix):]))
				resultTag.codec = "vorbis"
				return resultTag, err
			case bytes.HasPrefix(b, opusTagsPrefix):
				m := &metadataOGG{
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

type metadataOGG struct {
	*metadataVorbis
}

type metadataVorbis struct {
	c map[string]string // the vorbis comments
	p *Picture
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
		fmt.Println(cmt)
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
		} else if strings.HasPrefix(cmt, "albumartist=") {
			tag := strings.Replace(cmt, "albumartist=", "", 1)
			resultTag.albumArtist = tag
		} else if strings.HasPrefix(cmt, "ALBUMARTIST=") {
			tag := strings.Replace(cmt, "ALBUMARTIST=", "", 1)
			resultTag.albumArtist = tag
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
		} else if strings.HasPrefix(cmt, "comment=") {
			tag := strings.Replace(cmt, "comment=", "", 1)
			resultTag.genre = tag
		} else if strings.HasPrefix(cmt, "COMMENT=") {
			tag := strings.Replace(cmt, "COMMENT=", "", 1)
			resultTag.genre = tag
		}
	}
	if b64data, ok := m.c["metadata_block_picture"]; ok {
		data, err := base64.StdEncoding.DecodeString(b64data)
		if err != nil {
			return nil, err
		}
		m.readPictureBlock(bytes.NewReader(data))
	}
	albumArt := m.p
	if albumArt != nil {
		img, _, err := image.Decode(bytes.NewReader(albumArt.Data))
		if err != nil {
			log.Fatal("Error opening album image")
		}
		resultTag.albumArt = &img
	}
	return &resultTag, nil
}
func (m *metadataVorbis) readPictureBlock(r io.Reader) error {
	b, err := readInt(r, 4)
	if err != nil {
		return err
	}
	pictureType, ok := vorbisPictureTypes[byte(b)]
	if !ok {
		return fmt.Errorf("invalid picture type: %v", b)
	}
	mimeLen, err := readUint(r, 4)
	if err != nil {
		return err
	}
	mime, err := readString(r, mimeLen)
	if err != nil {
		return err
	}

	ext := ""
	switch mime {
	case "image/jpeg":
		ext = "jpg"
	case "image/png":
		ext = "png"
	case "image/gif":
		ext = "gif"
	}

	descLen, err := readUint(r, 4)
	if err != nil {
		return err
	}
	desc, err := readString(r, descLen)
	if err != nil {
		return err
	}

	// We skip width <32>, height <32>, colorDepth <32>, coloresUsed <32>
	_, err = readInt(r, 4) // width
	if err != nil {
		return err
	}
	_, err = readInt(r, 4) // height
	if err != nil {
		return err
	}
	_, err = readInt(r, 4) // color depth
	if err != nil {
		return err
	}
	_, err = readInt(r, 4) // colors used
	if err != nil {
		return err
	}

	dataLen, err := readInt(r, 4)
	if err != nil {
		return err
	}
	data := make([]byte, dataLen)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return err
	}

	m.p = &Picture{
		Ext:         ext,
		MIMEType:    mime,
		Type:        pictureType,
		Description: desc,
		Data:        data,
	}
	return nil
}

func clearTagsOpus(path string) error {
	inputFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer inputFile.Close()
	decoder := NewOGGDecoder(inputFile)
	tempOut, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	tempOut += "/output_file.ogg"
	outputFile, err := os.Create(tempOut)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	page, err := decoder.DecodeOGG()
	if err != nil {
		return err
	}
	encoder := NewOGGEncoder(page.Serial, outputFile)
	err = encoder.EncodeBOS(page.Granule, page.Packets)
	if err != nil {
		return err
	}
	var vorbisCommentPage *OGGPage
	for {
		page, err := decoder.DecodeOGG()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if hasOpusCommentPrefix(page.Packets) {
			vorbisCommentPage = &page
			emptyImage := []byte{}
			emptyComments := []string{}
			commentPacket := createOpusCommentPacket(emptyComments, emptyImage)

			vorbisCommentPage.Packets[0] = commentPacket
			err = encoder.Encode(vorbisCommentPage.Granule, vorbisCommentPage.Packets)
			if err != nil {
				return err
			}
			if len(page.Packets) == 1 {
				page, err := decoder.DecodeOGG()
				if err != nil {
					if err == io.EOF {
						break
					}
					return err
				}
				if page.Type == COP {
					if len(page.Packets) > 1 {
						err = encoder.Encode(page.Granule, page.Packets[1:])
						if err != nil {
							return err
						}
					}
				} else {
					err = encoder.Encode(page.Granule, page.Packets)
					if err != nil {
						return err
					}
				}
			}
		} else {
			// Write non-Vorbis comment pages to the output file
			if page.Type == EOS {
				err = encoder.EncodeEOS(page.Granule, page.Packets)
				if err != nil {
					return err
				}
			} else {
				err = encoder.Encode(page.Granule, page.Packets)
				if err != nil {
					return err
				}
			}
		}
	}
	inputFile.Close()
	outputFile.Close()
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	os.Rename(tempOut, abs)
	return nil
}

func saveOpusTags(tag *IDTag) error {
	// Step 1: Clear existing tags from the file
	err := clearTagsOpus(tag.fileUrl)
	if err != nil {
		return err
	}

	// Step 2: Open the input file and create an Ogg decoder
	inputFile, err := os.Open(tag.fileUrl)
	if err != nil {
		return err
	}
	defer inputFile.Close()
	decoder := NewOGGDecoder(inputFile)

	// Step 3: Create a temporary output file
	tempOut, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	tempOut += "/output_file.ogg"
	outputFile, err := os.Create(tempOut)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	page, err := decoder.DecodeOGG()
	if err != nil {
		return err
	}
	encoder := NewOGGEncoder(page.Serial, outputFile)
	err = encoder.EncodeBOS(page.Granule, page.Packets)
	if err != nil {
		return err
	}
	var vorbisCommentPage *OGGPage
	for {
		page, err := decoder.DecodeOGG()
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
			if tag.id3.date != "" {
				commentFields = append(commentFields, "DATE="+tag.title)
			}
			if tag.albumArtist != "" {
				commentFields = append(commentFields, "ALBUMARTIST="+tag.albumArtist)
			}
			img := []byte{}
			if tag.albumArt != nil {
				// Convert album art image to JPEG format
				buf := new(bytes.Buffer)
				err = jpeg.Encode(buf, *tag.albumArt, nil)
				if err != nil {
					return err
				}
				img = createMetadataBlockPicture(buf.Bytes())
			}

			// Create the new Vorbis comment packet
			commentPacket := createOpusCommentPacket(commentFields, img)

			// Replace the Vorbis comment packet in the original page with the new packet
			vorbisCommentPage.Packets[0] = commentPacket

			// Step 6: Write the updated Vorbis comment page to the output file
			err = encoder.Encode(vorbisCommentPage.Granule, vorbisCommentPage.Packets)
			if err != nil {
				return err
			}
		} else {
			// Write non-Vorbis comment pages to the output file
			if page.Type == EOS {
				err = encoder.EncodeEOS(page.Granule, page.Packets)
				if err != nil {
					return err
				}
			} else {
				err = encoder.Encode(page.Granule, page.Packets)
				if err != nil {
					return err
				}
			}
		}
	}
	// Step 7: Close and rename the files to the original file
	inputFile.Close()
	outputFile.Close()
	err = os.Rename(tempOut, tag.fileUrl)
	if err != nil {
		return err
	}

	return nil
}

func hasOpusCommentPrefix(packets [][]byte) bool {
	return len(packets) > 0 && len(packets[0]) >= 8 && string(packets[0][:8]) == "OpusTags"
}
func createOpusCommentPacket(commentFields []string, albumArt []byte) []byte {
	vendorString := "mp3mp4tag"

	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(len(vendorString)))
	vorbisCommentPacket := append(buf, []byte(vendorString)...)

	binary.LittleEndian.PutUint32(buf, uint32(len(commentFields)))
	vorbisCommentPacket = append(vorbisCommentPacket, buf...)

	for _, field := range commentFields {
		binary.LittleEndian.PutUint32(buf, uint32(len(field)))
		vorbisCommentPacket = append(vorbisCommentPacket, buf...)
		vorbisCommentPacket = append(vorbisCommentPacket, []byte(field)...)
	}
	vorbisCommentPacket = append([]byte("OpusTags"), vorbisCommentPacket...)
	vorbisCommentPacket = append(vorbisCommentPacket, albumArt...)
	return vorbisCommentPacket
}

func createMetadataBlockPicture(albumArtData []byte) []byte {
	// Replace these values with the appropriate information about the album art
	description := "Cover" // Description of the album art image
	width := uint32(100)   // Replace with the width of the image in pixels
	height := uint32(100)  // Replace with the height of the image in pixels
	depth := uint32(24)    // Replace with the color depth of the image in bits per pixel
	colors := uint32(0)    // Replace with the number of colors in the image
	mimeType := "image/jpeg"
	// Create the METADATA_BLOCK_PICTURE field with the image data
	imageType := []byte("\x00\x00\x00\x03") // 3 indicates the MIME type is URL

	// Append the uint32 values as big endian byte representations
	tempBuf := make([]byte, 4)

	// Create the METADATA_BLOCK_PICTURE field with the image data
	blockData := []byte{}
	blockData = append(blockData, imageType...)
	blockData = append(blockData, []byte(mimeType)...)
	blockData = append(blockData, byte(0)) // Null-terminated string
	blockData = append(blockData, byte(0)) // Optional field for the image index (converted to byte)
	blockData = append(blockData, byte(0)) // Optional field for the image count (converted to byte)
	blockData = append(blockData, []byte(description)...)
	blockData = append(blockData, byte(0)) // Null-terminated string

	// Append the uint32 values as big endian byte representations
	binary.BigEndian.PutUint32(tempBuf, width)
	blockData = append(blockData, tempBuf...)

	binary.BigEndian.PutUint32(tempBuf, height)
	blockData = append(blockData, tempBuf...)

	binary.BigEndian.PutUint32(tempBuf, depth)
	blockData = append(blockData, tempBuf...)

	binary.BigEndian.PutUint32(tempBuf, colors)
	blockData = append(blockData, tempBuf...)

	// Append the length of albumArtData as a big endian uint32
	binary.BigEndian.PutUint32(tempBuf, uint32(len(albumArtData)))
	blockData = append(blockData, tempBuf...)

	blockData = append(blockData, albumArtData...)

	// Calculate the length of the METADATA_BLOCK_PICTURE field
	blockLen := uint32(len(blockData))

	// Create the METADATA_BLOCK_PICTURE field with the calculated length and data
	pictureBlock := make([]byte, 4)
	binary.BigEndian.PutUint32(pictureBlock, blockLen)

	pictureBlock = append(pictureBlock, blockData...)

	return pictureBlock
}
func clearTagsVorbis(path string) error {
	inputFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer inputFile.Close()
	decoder := NewOGGDecoder(inputFile)
	tempOut, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	tempOut += "/output_file.ogg"
	outputFile, err := os.Create(tempOut)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	page, err := decoder.DecodeOGG()
	if err != nil {
		return err
	}
	encoder := NewOGGEncoder(page.Serial, outputFile)
	err = encoder.EncodeBOS(page.Granule, page.Packets)
	if err != nil {
		return err
	}
	var vorbisCommentPage *OGGPage
	for {
		page, err := decoder.DecodeOGG()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if hasVorbisCommentPrefix(page.Packets) {
			vorbisCommentPage = &page
			emptyImage := []byte{}
			emptyComments := []string{}
			commentPacket := createVorbisCommentPacket(emptyComments, emptyImage)

			vorbisCommentPage.Packets[0] = commentPacket
			err = encoder.Encode(vorbisCommentPage.Granule, vorbisCommentPage.Packets)
			if err != nil {
				return err
			}
			if len(page.Packets) == 1 {
				page, err := decoder.DecodeOGG()
				if err != nil {
					if err == io.EOF {
						break
					}
					return err
				}
				if page.Type == COP {
					if len(page.Packets) > 1 {
						err = encoder.Encode(page.Granule, page.Packets[1:])
						if err != nil {
							return err
						}
					}
				} else {
					err = encoder.Encode(page.Granule, page.Packets)
					if err != nil {
						return err
					}
				}
			}
		} else {
			// Write non-Vorbis comment pages to the output file
			if page.Type == EOS {
				err = encoder.EncodeEOS(page.Granule, page.Packets)
				if err != nil {
					return err
				}
			} else {
				err = encoder.Encode(page.Granule, page.Packets)
				if err != nil {
					return err
				}
			}
		}
	}
	inputFile.Close()
	outputFile.Close()
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	os.Rename(tempOut, abs)
	return nil
}

func saveVorbisTags(tag *IDTag) error {
	// Step 1: Clear existing tags from the file
	err := clearTagsVorbis(tag.fileUrl)
	if err != nil {
		return err
	}

	// Step 2: Open the input file and create an Ogg decoder
	inputFile, err := os.Open(tag.fileUrl)
	if err != nil {
		return err
	}
	defer inputFile.Close()
	decoder := NewOGGDecoder(inputFile)

	// Step 3: Create a temporary output file
	tempOut, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	tempOut += "/output_file.ogg"
	outputFile, err := os.Create(tempOut)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	page, err := decoder.DecodeOGG()
	if err != nil {
		return err
	}
	encoder := NewOGGEncoder(page.Serial, outputFile)
	err = encoder.EncodeBOS(page.Granule, page.Packets)
	if err != nil {
		return err
	}
	var vorbisCommentPage *OGGPage
	for {
		page, err := decoder.DecodeOGG()
		if err != nil {
			if err == io.EOF {
				break // Reached the end of the input Ogg stream
			}
			return err
		}

		// Find the Vorbis comment page and store it
		if hasVorbisCommentPrefix(page.Packets) {
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
			if tag.id3.date != "" {
				commentFields = append(commentFields, "DATE="+tag.title)
			}
			if tag.albumArtist != "" {
				commentFields = append(commentFields, "ALBUMARTIST="+tag.albumArtist)
			}
			if tag.comments != "" {
				commentFields = append(commentFields, "COMMENT="+tag.comments)
			}
			img := []byte{}
			if tag.albumArt != nil {
				// Convert album art image to JPEG format
				buf := new(bytes.Buffer)
				err = jpeg.Encode(buf, *tag.albumArt, nil)
				if err != nil {
					return err
				}
				img = createMetadataBlockPicture(buf.Bytes())
			}

			// Create the new Vorbis comment packet
			commentPacket := createVorbisCommentPacket(commentFields, img)

			// Replace the Vorbis comment packet in the original page with the new packet
			vorbisCommentPage.Packets[0] = commentPacket

			// Step 6: Write the updated Vorbis comment page to the output file
			err = encoder.Encode(vorbisCommentPage.Granule, vorbisCommentPage.Packets)
			if err != nil {
				return err
			}
		} else {
			// Write non-Vorbis comment pages to the output file
			if page.Type == EOS {
				err = encoder.EncodeEOS(page.Granule, page.Packets)
				if err != nil {
					return err
				}
			} else {
				err = encoder.Encode(page.Granule, page.Packets)
				if err != nil {
					return err
				}
			}
		}
	}
	// Step 7: Close and rename the files to the original file
	inputFile.Close()
	outputFile.Close()
	err = os.Rename(tempOut, tag.fileUrl)
	if err != nil {
		return err
	}

	return nil
}

func hasVorbisCommentPrefix(packets [][]byte) bool {
	return len(packets) > 0 && len(packets[0]) >= 7 && string(packets[0][:7]) == "\x03vorbis"
}
func createVorbisCommentPacket(commentFields []string, albumArt []byte) []byte {
	vendorString := "mp3mp4tag"

	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(len(vendorString)))
	vorbisCommentPacket := append(buf, []byte(vendorString)...)

	binary.LittleEndian.PutUint32(buf, uint32(len(commentFields)))
	vorbisCommentPacket = append(vorbisCommentPacket, buf...)

	for _, field := range commentFields {
		binary.LittleEndian.PutUint32(buf, uint32(len(field)))
		vorbisCommentPacket = append(vorbisCommentPacket, buf...)
		vorbisCommentPacket = append(vorbisCommentPacket, []byte(field)...)
	}
	vorbisCommentPacket = append([]byte("\x03vorbis"), vorbisCommentPacket...)
	vorbisCommentPacket = append(vorbisCommentPacket, albumArt...)
	vorbisCommentPacket = append(vorbisCommentPacket, []byte("\x01")...)
	return vorbisCommentPacket
}
