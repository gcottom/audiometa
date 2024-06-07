package mp4

import (
	"bytes"
	"encoding/binary"
	"image"
	"io"
	"reflect"

	mp4lib "github.com/abema/go-mp4"
)

func ReadMP4(r io.ReadSeeker) (*MP4Tag, error) {
	tag := new(MP4Tag)
	tag.reader = r
	tptr := reflect.ValueOf(tag).Elem()

	// I know this looks like O(n^3), it's not. The middle loop realisticly only
	// has 1 iteration, and the inner loop is only executed twice
ol:
	for atom, field := range atomsMap {
		boxes, err := mp4lib.ExtractBoxWithPayload(r, nil, mp4lib.BoxPath{mp4lib.BoxTypeMoov(), mp4lib.BoxTypeUdta(), mp4lib.BoxTypeMeta(), mp4lib.BoxTypeIlst(), atom, mp4lib.BoxTypeData()})
		if err != nil {
			return nil, err
		}
		for _, box := range boxes {
			data := box.Payload.(*mp4lib.Data)
			switch atom {
			case mp4lib.BoxType{'t', 'r', 'k', 'n'}, mp4lib.BoxType{'d', 'i', 's', 'k'}:
				var num uint16
				if err := binary.Read(bytes.NewReader(data.Data[2:4]), binary.BigEndian, &num); err != nil {
					return nil, err
				}
				tptr.FieldByName(field).SetInt(int64(num))
				typ := reflect.TypeOf(*tag)
				fNum := 0
			strL:
				for i := 0; i < typ.NumField(); i++ {
					if typ.Field(i).Name == field {
						fNum = i + 1
						break strL
					}
				}
				if err = binary.Read(bytes.NewReader(data.Data[4:6]), binary.BigEndian, &num); err != nil {
					return nil, err
				}
				tptr.Field(fNum).SetInt(int64(num))
				continue ol
			case mp4lib.BoxType{'t', 'm', 'p', 'o'}:
				tag.BPM = getInt(data.Data[:1])
				continue ol
			case mp4lib.BoxType{'c', 'o', 'v', 'r'}:
				img, _, err := image.Decode(bytes.NewReader(data.Data))
				if err != nil {
					return nil, err
				}
				tag.CoverArt = &img
				continue ol
			default:
				tptr.FieldByName(field).SetString(string(data.Data))
				continue ol
			}
		}
	}
	return tag, nil
}
