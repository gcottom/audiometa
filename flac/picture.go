package flac

import (
	"bytes"
	"image/jpeg"
	"image/png"
)

// PictureType defines the type of this image
type PictureType uint32

const (
	PictureTypeOther PictureType = iota
	PictureTypeFileIcon
	PictureTypeOtherIcon
	PictureTypeFrontCover
	PictureTypeBackCover
	PictureTypeLeaflet
	PictureTypeMedia
	PictureTypeLeadArtist
	PictureTypeArtist
	PictureTypeConductor
	PictureTypeBand
	PictureTypeComposer
	PictureTypeLyricist
	PictureTypeRecordingLocation
	PictureTypeDuringRecording
	PictureTypeDuringPerformance
	PictureTypeScreenCapture
	PictureTypeBrightColouredFish
	PictureTypeIllustration
	PictureTypeBandArtistLogotype
	PictureTypePublisherStudioLogotype
)

// MetadataBlockPicture represents a picture metadata block
type MetadataBlockPicture struct {
	PictureType       PictureType
	MIME              string
	Description       string
	Width             uint32
	Height            uint32
	ColorDepth        uint32
	IndexedColorCount uint32
	ImageData         []byte
}

// Marshal encodes the PictureBlock to a standard FLAC MetaDataBloc to be accepted by go-flac
func (c *MetadataBlockPicture) Marshal() MetaDataBlock {
	res := bytes.NewBuffer([]byte{})
	res.Write(encodeUint32(uint32(c.PictureType)))
	res.Write(encodeUint32(uint32(len(c.MIME))))
	res.Write([]byte(c.MIME))
	res.Write(encodeUint32(uint32(len(c.Description))))
	res.Write([]byte(c.Description))
	res.Write(encodeUint32(c.Width))
	res.Write(encodeUint32(c.Height))
	res.Write(encodeUint32(c.ColorDepth))
	res.Write(encodeUint32(c.IndexedColorCount))
	res.Write(encodeUint32(uint32(len(c.ImageData))))
	res.Write(c.ImageData)
	return MetaDataBlock{
		Type: Picture,
		Data: res.Bytes(),
	}
}

// NewFromImageData generates a MetadataBlockPicture from image data, returns an error if the picture data connot be parsed
func NewFromImageData(pictype PictureType, description string, imgdata []byte, mime string) (*MetadataBlockPicture, error) {
	res := new(MetadataBlockPicture)
	res.PictureType = pictype
	res.Description = description
	res.MIME = mime
	res.ImageData = imgdata
	err := res.ParsePicture()
	return res, err
}

// ParseFromMetaDataBlock parses an existing picture MetaDataBlock
func ParsePicFromMetaDataBlock(meta MetaDataBlock) (*MetadataBlockPicture, error) {
	if meta.Type != Picture {
		return nil, ErrorNotPictureMetadataBlock
	}
	res := new(MetadataBlockPicture)
	data := bytes.NewBuffer(meta.Data)

	if pictype, err := readUint32(data); err != nil {
		return nil, err
	} else {
		res.PictureType = PictureType(pictype)
	}

	if mimebytes, err := readBytesWith32bitSize(data); err != nil {
		return nil, err
	} else {
		res.MIME = string(mimebytes)
	}

	if descbytes, err := readBytesWith32bitSize(data); err != nil {
		return nil, err
	} else {
		res.Description = string(descbytes)
	}

	if width, err := readUint32(data); err != nil {
		return nil, err
	} else {
		res.Width = width
	}

	if height, err := readUint32(data); err != nil {
		return nil, err
	} else {
		res.Height = height
	}

	if depth, err := readUint32(data); err != nil {
		return nil, err
	} else {
		res.ColorDepth = depth
	}

	if count, err := readUint32(data); err != nil {
		return nil, err
	} else {
		res.IndexedColorCount = count
	}

	if pic, err := readBytesWith32bitSize(data); err != nil {
		return nil, err
	} else {
		res.ImageData = pic
	}

	return res, nil
}

// ParsePicture decodes the image and inflated the Width, Height, ColorDepth, IndexedColorCount fields. This is called automatically by NewFromImageData
func (c *MetadataBlockPicture) ParsePicture() error {
	switch c.MIME {
	case "image/jpeg":
		img, err := jpeg.Decode(bytes.NewReader(c.ImageData))
		if err != nil {
			return err
		}
		c.IndexedColorCount = uint32(0)
		size := img.Bounds()
		c.Width = uint32(size.Max.X)
		c.Height = uint32(size.Max.Y)
		c.ColorDepth = uint32(24)
	case "image/png":
		img, err := png.Decode(bytes.NewReader(c.ImageData))
		if err != nil {
			return err
		}
		c.IndexedColorCount = uint32(0)
		size := img.Bounds()
		c.Width = uint32(size.Max.X)
		c.Height = uint32(size.Max.Y)
		c.ColorDepth = uint32(32)
	default:
		return ErrorUnsupportedMIME
	}
	return nil
}
