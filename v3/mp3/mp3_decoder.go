package mp3

import (
	"bytes"
	"errors"
	"image"
	"io"
	"reflect"

	mp3TagLib "github.com/bogem/id3v2/v2"
)

func ParseMP3(r io.ReadSeeker) (*MP3Tag, error) {
	resultTag := MP3Tag{}
	resultTag.reader = r
	tag, err := mp3TagLib.ParseReader(r, mp3TagLib.Options{Parse: true})
	if err != nil {
		return nil, errors.New("error parsing mp3")
	}
	rtPtr := reflect.ValueOf(&resultTag).Elem()
	for k, v := range mp3TextFrames {
		framer := tag.GetTextFrame(v)
		if framer.Text == "" {
			continue
		}
		rtPtr.FieldByName(k)
	}
	if pictures := tag.GetFrames("APIC"); len(pictures) > 0 {
		pic := pictures[0].(mp3TagLib.PictureFrame)
		if img, _, err := image.Decode(bytes.NewReader(pic.Picture)); err == nil {
			resultTag.AlbumArt = &img
		}
	}
	return &resultTag, nil
}
