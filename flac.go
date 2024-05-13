package audiometa

import (
	"bytes"
	"image"
	"io"
	"os"
	"path/filepath"

	"github.com/aler9/writerseeker"
	"github.com/gcottom/audiometa/v2/flac"
)

type flacBlock struct {
	cmts   *flac.MetaDataBlockVorbisComment
	pic    *image.Image
	picIdx int
	cmtIdx int
}

func extractFLACComment(input io.Reader) (*flac.File, *flacBlock, error) {
	fb := flacBlock{}
	f, err := flac.ParseMetadata(input)
	if err != nil {
		return nil, nil, ErrFLACParse
	}
	for idx, meta := range f.Meta {
		if meta.Type == flac.VorbisComment {
			fb.cmts, err = flac.ParseFromMetaDataBlock(*meta)
			fb.cmtIdx = idx
			if err != nil {
				return nil, nil, ErrFLACCmtParse
			}
			continue
		} else if meta.Type == flac.Picture {
			if pic, err := flac.ParsePicFromMetaDataBlock(*meta); err == nil {
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

func flacSave(r io.Reader, w io.Writer, m []*flac.MetaDataBlock, needsTemp bool) error {
	if needsTemp {
		//in and out are the same file so we have to temp it
		t := &writerseeker.WriterSeeker{}
		defer t.Close()
		// Write tag in new file.
		if _, err := t.Write([]byte("fLaC")); err != nil {
			return err
		}
		for i, meta := range m {
			last := i == len(m)-1
			if _, err := t.Write(meta.Marshal(last)); err != nil {
				return err
			}
		}
		if _, err := io.Copy(t, r); err != nil {
			return err
		}
		if _, err := t.Seek(0, io.SeekStart); err != nil {
			return err
		}

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
		return nil
	}

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
