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

func extractFLACComment(input io.Reader) (*flac.File, *flacBlock, error) {
	fb := flacBlock{}
	f, err := flac.ParseMetadata(input)
	if err != nil {
		return nil, nil, err
	}
	for idx, meta := range f.Meta {
		if meta.Type == flac.VorbisComment {
			fb.cmts, err = flacvorbis.ParseFromMetaDataBlock(*meta)
			fb.cmtIdx = idx
			if err != nil {
				return nil, nil, err
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

	return f, &fb, nil
}

func removeFLACMetaBlock(slice []*flac.MetaDataBlock, s int) []*flac.MetaDataBlock {
	return append(slice[:s], slice[s+1:]...)
}

func flacSave(r io.Reader, w io.Writer, m []*flac.MetaDataBlock) error {
	if _, err := w.Write([]byte("fLaC")); err != nil {
		return err
	}
	for i, meta := range m {
		last := i == len(m)-1
		if _, err := w.Write(meta.Marshal(last)); err != nil {
			return err
		}
	}
	if _, err := io.Copy(w, r); err != nil {
		return err
	}
	return nil

}
