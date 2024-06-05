package audiometa

import (
	"bytes"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"reflect"

	"github.com/abema/go-mp4"
	"github.com/aler9/writerseeker"
	"github.com/sunfish-shogi/bufseekio"
)

var atomsMap = map[string]mp4.BoxType{
	"album":        {'\251', 'a', 'l', 'b'},
	"albumArtist":  {'a', 'A', 'R', 'T'},
	"artist":       {'\251', 'A', 'R', 'T'},
	"comments":     {'\251', 'c', 'm', 't'},
	"composer":     {'\251', 'w', 'r', 't'},
	"copyrightMsg": {'c', 'p', 'r', 't'},
	"albumArt":     {'c', 'o', 'v', 'r'},
	"genre":        {'\251', 'g', 'e', 'n'},
	"title":        {'\251', 'n', 'a', 'm'},
	"year":         {'\251', 'd', 'a', 'y'},
}

// Make new atoms and write to.
func createAndWrite(w *mp4.Writer, ctx mp4.Context, _tags *IDTag) error {
	for tagName, boxType := range atomsMap {
		if tagName == "albumArt" {
			if _tags.albumArt != nil {
				buf := new(bytes.Buffer)
				if err := png.Encode(buf, *_tags.albumArt); err == nil {
					if _, err := w.StartBox(&mp4.BoxInfo{Type: boxType}); err != nil {
						return err
					}
					if _, err := w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeData()}); err != nil {
						return err
					}
					var boxData = &mp4.Data{
						DataType: mp4.DataTypeBinary,
						Data:     buf.Bytes(),
					}
					dataCtx := ctx
					dataCtx.UnderIlstMeta = true
					if _, err := mp4.Marshal(w, boxData, dataCtx); err != nil {
						return err
					}
					if _, err := w.EndBox(); err != nil {
						return err
					}
					if _, err := w.EndBox(); err != nil {
						return err
					}
				}
			}
			continue
		}
		val := reflect.ValueOf(*_tags).FieldByName(tagName).String()
		if val == "" {
			continue
		}

		if _, err := w.StartBox(&mp4.BoxInfo{Type: boxType}); err != nil {
			return err
		}
		if _, err := w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeData()}); err != nil {
			return err
		}
		var boxData = &mp4.Data{
			DataType: mp4.DataTypeStringUTF8,
			Data:     []byte(val),
		}
		dataCtx := ctx
		dataCtx.UnderIlstMeta = true
		if _, err := mp4.Marshal(w, boxData, dataCtx); err != nil {
			return err
		}
		if _, err := w.EndBox(); err != nil {
			return err
		}
		if _, err := w.EndBox(); err != nil {
			return err
		}
	}
	return nil

}

func containsAtom(boxType mp4.BoxType, boxes []mp4.BoxType) mp4.BoxType {
	for _, _boxType := range boxes {
		if boxType == _boxType {
			return boxType
		}
	}
	return mp4.BoxType{}
}

func getAtomsList() []mp4.BoxType {
	var atomsList []mp4.BoxType
	for _, atom := range atomsMap {
		atomsList = append(atomsList, atom)
	}
	return atomsList
}

func writeMP4(r *bufseekio.ReadSeeker, wo io.Writer, _tags *IDTag, delete MP4Delete) error {
	atomsList := getAtomsList()

	ws := &writerseeker.WriterSeeker{}
	defer ws.Close()

	w := mp4.NewWriter(ws)
	var mdatOffsetDiff int64
	var stcoOffsets []int64
	closedTags := false

	_, err := mp4.ReadBoxStructure(r, func(h *mp4.ReadHandle) (interface{}, error) {
		switch h.BoxInfo.Type {

		case containsAtom(h.BoxInfo.Type, atomsList):
			return nil, nil

		case mp4.BoxTypeFree():
			if !closedTags {
				_, err := w.EndBox()
				if err != nil {
					return nil, err
				}
				if err := w.CopyBox(r, &h.BoxInfo); err != nil {
					return nil, err
				}
				_, err = w.EndBox()
				if err != nil {
					return nil, err
				}
				_, err = w.EndBox()
				if err != nil {
					return nil, err
				}
				_, err = w.EndBox()
				if err != nil {
					return nil, err
				}
				closedTags = true
				return nil, nil
			}
			if err := w.CopyBox(r, &h.BoxInfo); err != nil {
				return nil, err
			}
			return nil, nil

		case mp4.BoxTypeMeta():
			w.StartBox(&h.BoxInfo)
			box, _, err := h.ReadPayload()
			if err != nil {
				return nil, err
			}
			if _, err = mp4.Marshal(w, box, h.BoxInfo.Context); err != nil {
				return nil, err
			}
			return h.Expand()

		case mp4.BoxTypeMoov(),
			mp4.BoxTypeUdta():
			w.StartBox(&h.BoxInfo)
			box, _, err := h.ReadPayload()
			if err != nil {
				return nil, err
			}
			if _, err = mp4.Marshal(w, box, h.BoxInfo.Context); err != nil {
				return nil, err
			}
			return h.Expand()

		case mp4.BoxTypeIlst():
			_, err := w.StartBox(&h.BoxInfo)
			if err != nil {
				return nil, err
			}
			ctx := h.BoxInfo.Context
			if err = createAndWrite(w, ctx, _tags); err != nil {
				return nil, err
			}
			return h.Expand()

		default:
			if h.BoxInfo.Type == mp4.BoxTypeStco() {
				offset, _ := w.Seek(0, io.SeekCurrent)
				stcoOffsets = append(stcoOffsets, offset)
			}
			if h.BoxInfo.Type == mp4.BoxTypeMdat() {
				iOffset := int64(h.BoxInfo.Offset)
				oOffset, _ := w.Seek(0, io.SeekCurrent)
				mdatOffsetDiff = oOffset - iOffset
			}
			if err := w.CopyBox(r, &h.BoxInfo); err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	ts := bufseekio.NewReadSeeker(bytes.NewReader(ws.Bytes()), 1024*1024, 3)

	_, err = mp4.ReadBoxStructure(ts, func(h *mp4.ReadHandle) (any, error) {
		switch h.BoxInfo.Type {
		case mp4.BoxTypeStco():
			stcoOffsets = append(stcoOffsets, int64(h.BoxInfo.Offset))
		default:
			return h.Expand()
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	if _, err = ws.Seek(0, io.SeekStart); err != nil {
		return err
	}
	// if mdat box is moved, update stco box
	if mdatOffsetDiff != 0 {
		for _, stcoOffset := range stcoOffsets {
			// seek to stco box header
			if _, err := ts.Seek(stcoOffset, io.SeekStart); err != nil {
				return err
			}
			// read box header
			bi, err := mp4.ReadBoxInfo(ts)
			if err != nil {
				return err
			}
			// read stco box payload
			var stco mp4.Stco
			_, err = mp4.Unmarshal(ts, bi.Size-bi.HeaderSize, &stco, bi.Context)
			if err != nil {
				return err
			}
			// update chunk offsets
			for i := range stco.ChunkOffset {
				stco.ChunkOffset[i] += uint32(mdatOffsetDiff)
			}
			// seek to stco box payload
			_, err = bi.SeekToPayload(ws)
			if err != nil {
				return err
			}
			// write stco box payload
			if _, err := mp4.Marshal(ws, &stco, bi.Context); err != nil {
				return err
			}
		}
	}
	if _, err = ws.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if reflect.TypeOf(wo) == reflect.TypeOf(new(os.File)) {
		f := wo.(*os.File)
		path, err := filepath.Abs(f.Name())
		if err != nil {
			return err
		}
		w2, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		defer w2.Close()
		if _, err = w2.Write(ws.Bytes()); err != nil {
			return err
		}
		if _, err = f.Seek(0, io.SeekEnd); err != nil {
			return err
		}
		return nil
	}
	if _, err := wo.Write(ws.Bytes()); err != nil {
		return err
	}
	return nil

}
