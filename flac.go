package audiometa

import (
	"bytes"
	"image"
	"io"

	"github.com/go-flac/flacpicture"
	"github.com/go-flac/flacvorbis"
	"github.com/go-flac/go-flac"
)

type flacBlock struct {
	cmts   *flacvorbis.MetaDataBlockVorbisComment
	pic    *image.Image
	picIdx int
	cmtIdx int
}

func extractFLACComment(input io.Reader) (*flacBlock, error) {
	fb := flacBlock{}
	f, err := flac.ParseBytes(input)
	if err != nil {
		return nil, err
	}
	for idx, meta := range f.Meta {
		if meta.Type == flac.VorbisComment {
			fb.cmts, err = flacvorbis.ParseFromMetaDataBlock(*meta)
			fb.cmtIdx = idx
			if err != nil {
				return nil, err
			}
			continue
		} else if meta.Type == flac.Picture {
			if pic, err := flacpicture.ParseFromMetaDataBlock(*meta); err == nil {
				if pic != nil {
					if img, _, err := image.Decode(bytes.NewReader(pic.ImageData)); err == nil {
						fb.pic = &img
						fb.picIdx = idx
					}
				}
				continue
			}
		}
	}

	return &fb, nil
}
func getFLACPictureIndex(metaIn []*flac.MetaDataBlock) int {
	var cmtIdx = 0
	for idx, meta := range metaIn {
		if meta.Type == flac.Picture {
			cmtIdx = idx
			break
		}
	}
	return cmtIdx
}
func removeFLACMetaBlock(slice []*flac.MetaDataBlock, s int) []*flac.MetaDataBlock {
	return append(slice[:s], slice[s+1:]...)
}
