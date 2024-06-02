package flac

import (
	"bytes"
	"strings"
)

type MetaDataBlockVorbisComment struct {
	Vendor   string
	Comments []string
}

// New creates a new MetaDataBlockVorbisComment
// vendor is set to flacvorbis <version> by default
func New() *MetaDataBlockVorbisComment {
	return &MetaDataBlockVorbisComment{
		"audiometa",
		[]string{},
	}
}

// Get get all comments with field name specified by the key parameter
// If there is no match, error would still be nil
func (c *MetaDataBlockVorbisComment) Get(key string) ([]string, error) {
	res := make([]string, 0)
	for _, cmt := range c.Comments {
		p := strings.SplitN(cmt, "=", 2)
		if len(p) != 2 {
			return nil, ErrorMalformedComment
		}
		if strings.EqualFold(p[0], key) {
			res = append(res, p[1])
		}
	}
	return res, nil
}

// Add adds a key-val pair to the comments
func (c *MetaDataBlockVorbisComment) Add(key string, val string) error {
	for _, char := range key {
		if char < 0x20 || char > 0x7d || char == '=' {
			return ErrorInvalidFieldName
		}
	}
	c.Comments = append(c.Comments, key+"="+val)
	return nil
}

// Marshal marshals this block back into a flac.MetaDataBlock
func (c MetaDataBlockVorbisComment) Marshal() (MetaDataBlock, error) {
	data := bytes.NewBuffer([]byte{})
	if err := packStr(data, c.Vendor); err != nil {
		return MetaDataBlock{}, err
	}
	data.Write(encodeUint32L(uint32(len(c.Comments))))
	for _, cmt := range c.Comments {
		if err := packStr(data, cmt); err != nil {
			return MetaDataBlock{}, err
		}
	}
	return MetaDataBlock{
		Type: VorbisComment,
		Data: data.Bytes(),
	}, nil
}

// ParseFromMetaDataBlock parses an existing picture MetaDataBlock
func ParseFromMetaDataBlock(meta MetaDataBlock) (*MetaDataBlockVorbisComment, error) {
	if meta.Type != VorbisComment {
		return nil, ErrorNotVorbisComment
	}

	reader := bytes.NewReader(meta.Data)
	res := new(MetaDataBlockVorbisComment)

	vendorlen, err := readUint32L(reader)
	if err != nil {
		return nil, err
	}
	vendorbytes := make([]byte, vendorlen)
	nn, err := reader.Read(vendorbytes)
	if err != nil {
		return nil, err
	}
	if nn != int(vendorlen) {
		return nil, ErrorUnexpEof
	}
	res.Vendor = string(vendorbytes)

	cmtcount, err := readUint32L(reader)
	if err != nil {
		return nil, err
	}
	res.Comments = make([]string, cmtcount)
	for i := range res.Comments {
		cmtlen, err := readUint32L(reader)
		if err != nil {
			return nil, err
		}
		cmtbytes := make([]byte, cmtlen)
		nn, err := reader.Read(cmtbytes)
		if err != nil {
			return nil, err
		}
		if nn != int(cmtlen) {
			return nil, ErrorUnexpEof
		}
		res.Comments[i] = string(cmtbytes)
	}
	return res, nil
}
