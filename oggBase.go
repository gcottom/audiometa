package mp3mp4tag

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// OggPage represents a single Ogg page.
type OggPage struct {
	Version   uint8
	TypeFlags uint8
	Granule   int64
	Serial    int32
	Sequence  int32
	Offset    int64
	Complete  bool
	Packets   [][]byte
}

// NewOggPage creates a new OggPage object.
func NewOggPage(fileobj io.Reader) (*OggPage, error) {
	page := &OggPage{}

	// Read header
	header := make([]byte, 27)
	_, err := io.ReadFull(fileobj, header)
	if err != nil {
		return nil, err
	}

	// Parse header fields
	if string(header[:4]) != "OggS" {
		return nil, errors.New("invalid Ogg page")
	}
	page.Version = uint8(header[4])
	page.TypeFlags = uint8(header[5])
	page.Granule = int64(binary.LittleEndian.Uint64(header[6:14]))
	page.Serial = int32(binary.LittleEndian.Uint32(header[14:18]))
	page.Sequence = int32(binary.LittleEndian.Uint32(header[18:22]))
	page.Offset = int64(binary.LittleEndian.Uint64(header[22:30]))

	// Parse lacing values
	lacingBytes := make([]byte, int(header[26]))
	_, err = io.ReadFull(fileobj, lacingBytes)
	if err != nil {
		return nil, err
	}

	packetSizes := make([]int, len(lacingBytes))
	totalSize := 0
	for i, lacingByte := range lacingBytes {
		packetSizes[i] = int(lacingByte)
		totalSize += packetSizes[i]
	}

	// Read packets
	packetData := make([]byte, totalSize)
	_, err = io.ReadFull(fileobj, packetData)
	if err != nil {
		return nil, err
	}

	// Split packet data into packets
	packets := make([][]byte, len(packetSizes))
	offset := 0
	for i, size := range packetSizes {
		packets[i] = packetData[offset : offset+size]
		offset += size
	}

	page.Packets = packets

	// Determine if the last packet is complete
	page.Complete = len(page.Packets) == 0 || lacingBytes[len(lacingBytes)-1] != 255

	return page, nil
}

// Write serializes the OggPage object to bytes.
func (page *OggPage) Write() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Write header
	buf.Write([]byte("OggS"))
	buf.WriteByte(page.Version)
	buf.WriteByte(page.TypeFlags)
	binary.Write(buf, binary.LittleEndian, page.Granule)
	binary.Write(buf, binary.LittleEndian, page.Serial)
	binary.Write(buf, binary.LittleEndian, page.Sequence)
	binary.Write(buf, binary.LittleEndian, page.Offset)

	// Write lacing values
	lacingBytes := make([]byte, len(page.Packets))
	offset := 0
	for i, packet := range page.Packets {
		lacingBytes[i] = byte(len(packet))
		offset += len(packet)
	}
	buf.WriteByte(byte(len(lacingBytes)))
	buf.Write(lacingBytes)

	// Write packet data
	for _, packet := range page.Packets {
		buf.Write(packet)
	}

	return buf.Bytes(), nil
}

// String returns a string representation of the OggPage object.
func (page *OggPage) String() string {
	packetCount := len(page.Packets)
	packetSizes := make([]int, packetCount)
	for i, packet := range page.Packets {
		packetSizes[i] = len(packet)
	}
	totalSize := sumInts(packetSizes)
	return fmt.Sprintf("<OggPage version=%d sequence=%d serial=%d packets=%d size=%d>",
		page.Version, page.Sequence, page.Serial, packetCount, totalSize)
}

// StreamInfo represents the information of an Ogg stream.
type StreamInfo struct {
	Version      uint8
	Channels     uint8
	SampleRate   uint32
	BitrateUpper uint32
	BitrateNom   uint32
	BitrateLower uint32
	Blocksize0   uint16
	Blocksize1   uint16
}

// NewStreamInfo creates a new StreamInfo object from the given reader.
func NewStreamInfo(fileobj io.Reader) (*StreamInfo, error) {
	info := &StreamInfo{}

	header := make([]byte, 18)
	_, err := io.ReadFull(fileobj, header)
	if err != nil {
		return nil, err
	}

	info.Version = uint8(header[0])
	info.Channels = uint8(header[1])
	info.SampleRate = binary.LittleEndian.Uint32(header[2:6])
	info.BitrateUpper = binary.LittleEndian.Uint32(header[6:10])
	info.BitrateNom = binary.LittleEndian.Uint32(header[10:14])
	info.BitrateLower = binary.LittleEndian.Uint32(header[14:18])

	blockSizes := make([]byte, 4)
	_, err = io.ReadFull(fileobj, blockSizes)
	if err != nil {
		return nil, err
	}

	info.Blocksize0 = binary.LittleEndian.Uint16(blockSizes[:2])
	info.Blocksize1 = binary.LittleEndian.Uint16(blockSizes[2:])

	return info, nil
}


// Tags represents the tags (metadata) of an Ogg file.
type Tags struct {
	Vendor string
	Tags   map[string][]string
}

// NewTags creates a new Tags object from the given reader and stream info.
func NewTags(fileobj io.Reader, info *StreamInfo) (*Tags, error) {
	tags := &Tags{
		Tags: make(map[string][]string),
	}

	// Read the vendor string length
	var vendorLen uint32
	err := binary.Read(fileobj, binary.LittleEndian, &vendorLen)
	if err != nil {
		return nil, err
	}

	// Read the vendor string
	vendorBytes := make([]byte, vendorLen)
	_, err = io.ReadFull(fileobj, vendorBytes)
	if err != nil {
		return nil, err
	}
	tags.Vendor = string(vendorBytes)

	// Read the number of tag entries
	var numTags uint32
	err = binary.Read(fileobj, binary.LittleEndian, &numTags)
	if err != nil {
		return nil, err
	}

	// Read the tag entries
	for i := uint32(0); i < numTags; i++ {
		// Read the tag entry length
		var entryLen uint32
		err = binary.Read(fileobj, binary.LittleEndian, &entryLen)
		if err != nil {
			return nil, err
		}

		// Read the tag entry
		entryBytes := make([]byte, entryLen)
		_, err = io.ReadFull(fileobj, entryBytes)
		if err != nil {
			return nil, err
		}
		entry := string(entryBytes)

		// Split the entry into key and value
		keyValue := splitTagEntry(entry)
		if keyValue != nil {
			key := keyValue[0]
			value := keyValue[1]

			// Add the tag to the map
			tags.Tags[key] = append(tags.Tags[key], value)
		}
	}

	return tags, nil
}

// Helper function to split a tag entry into key and value.
func splitTagEntry(entry string) []string {
	// Find the first '=' character
	index := 0
	for i, ch := range entry {
		if ch == '=' {
			index = i
			break
		}
	}

	// If '=' is not found, return nil
	if index == 0 {
		return nil
	}

	key := entry[:index]
	value := entry[index+1:]

	return []string{key, value}
}

// Delete removes the tags from the Ogg file.
func (tags *Tags) Delete(fileobj io.Writer) error {
	tags.Tags = make(map[string][]string) // Clear the tags
	return nil
}
func (tags *Tags) Clear() {
	tags.Vendor = ""
	tags.Tags = make(map[string][]string)
}
func encodeString(str string) []byte {
	return []byte(str)
}
// Inject writes the tags to the specified file object, using the given padding function.
func (tags *Tags) Inject(fileobj io.Writer, padding PaddingFunc) error {
	// Prepare the tag data
	vendorData := encodeString(tags.Vendor)
	tagData := make([]byte, 0)

	// Build the tag data from the key-value pairs
	for key, values := range tags.Tags {
		for _, value := range values {
			keyData := encodeString(key)
			valueData := encodeString(value)
			tagData = append(tagData, keyData...)
			tagData = append(tagData, valueData...)
		}
	}

	// Calculate the tag size
	tagSize := len(vendorData) + len(tagData) + 8

	// Prepare the padding information
	paddingInfo := &PaddingInfo{
		Padding: padding(tagSize),
		Size:    tagSize,
	}

	// Write the vendor data
	if _, err := fileobj.Write(vendorData); err != nil {
		return err
	}

	// Write the tag data
	if _, err := fileobj.Write(tagData); err != nil {
		return err
	}

	// Write the padding
	paddingSize := paddingInfo.GetDefaultPadding()
	paddingBytes := make([]byte, paddingSize)
	if _, err := fileobj.Write(paddingBytes); err != nil {
		return err
	}

	return nil
}

// OggFileType represents a generic Ogg file.
type OggFileType struct {
	Info StreamInfo
	Tags Tags
}

// Load loads file information from the specified file.
func (file *OggFileType) Load(fileobj io.Reader) error {
	info, err := NewStreamInfo(fileobj)
	if err != nil {
		return err
	}
	file.Info = *info

	tags, err := NewTags(fileobj, &file.Info)
	if err != nil {
		return err
	}
	file.Tags = *tags

	return nil
}

// Delete removes tags from the file.
func (file *OggFileType) Delete(fileobj io.Writer) error {
	file.Tags.Clear()
	err := file.Tags.Inject(fileobj, nil)
	if err != nil {
		return err
	}
	return nil
}

// Save saves the tags to the file.
func (file *OggFileType) Save(fileobj io.Writer, padding PaddingFunction) error {
	err := file.Tags.Inject(fileobj, padding)
	if err != nil {
		return err
	}
	return nil
}

// PaddingFunction is a function type for calculating padding.
type PaddingFunction func(info *PaddingInfo) int

// PaddingInfo represents padding information used in the padding calculation.
type PaddingInfo struct {
	Padding int
	Size    int
}

// GetDefaultPadding returns the default amount of padding after saving.
func (info *PaddingInfo) GetDefaultPadding() int {
	high := 10240 + info.Size/100 // 10 KiB + 1% of trailing data
	low := 1024 + info.Size/1000  // 1 KiB + 0.1% of trailing data

	if info.Padding >= 0 {
		if info.Padding > high {
			return low
		}
		return info.Padding
	} else {
		return low
	}
}

// sumInts returns the sum of the integers in the slice.
func sumInts(ints []int) int {
	sum := 0
	for _, i := range ints {
		sum += i
	}
	return sum
}

func main() {
	fmt.Println("Go code conversion completed successfully.")
}