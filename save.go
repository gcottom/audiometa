package audiometa

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"reflect"

	mp3TagLib "github.com/bogem/id3v2/v2"
	"github.com/gcottom/audiometa/v2/flac"
	"github.com/sunfish-shogi/bufseekio"
)

// Save saves the corresponding IDTag to the audio file that it references and returns an error if the saving process fails
func (tag *IDTag) Save(w io.Writer) error {
	fileType, err := GetFileType(tag.filePath)
	if err != nil {
		return err
	}
	if _, err := tag.reader.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if fileType == MP3 {
		return saveMP3(tag, w)
	} else if fileType == M4A || fileType == M4B || fileType == M4P || fileType == MP4 {
		return saveMP4(tag, w)
	} else if fileType == FLAC {
		return saveFLAC(tag, w)
	} else if fileType == OGG {
		return saveOGG(tag, w)
	}
	return fmt.Errorf("no method available for filetype:%s", tag.fileType)
}

func saveMP3(tag *IDTag, w io.Writer) error {
	r := tag.reader
	mp3Tag, err := mp3TagLib.ParseReader(r, mp3TagLib.Options{Parse: true})
	if err != nil {
		return err
	}
	for k, v := range mp3TextFrames {
		if reflect.ValueOf(*tag).FieldByName(k).IsZero() {
			mp3Tag.DeleteFrames(v)
			continue
		}
		if k == "albumArt" {
			buf := new(bytes.Buffer)
			if err := jpeg.Encode(buf, *tag.albumArt, nil); err == nil {
				mp3Tag.AddAttachedPicture(mp3TagLib.PictureFrame{
					Encoding:    mp3TagLib.EncodingUTF8,
					MimeType:    "image/jpeg",
					PictureType: mp3TagLib.PTFrontCover,
					Description: "Front cover",
					Picture:     buf.Bytes(),
				})
			}
			continue
		}
		textFrame := mp3TagLib.TextFrame{
			Encoding: mp3TagLib.EncodingUTF8,
			Text:     reflect.ValueOf(*tag).FieldByName(k).String(),
		}
		mp3Tag.AddFrame(v, textFrame)
	}
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return err
	}
	originalSize, err := parseHeader(r)
	if err != nil {
		return err
	}
	// Write tag in new file.
	if _, err = mp3Tag.WriteTo(w); err != nil {
		return err
	}
	// Seek to a music part of original file.
	if _, err = r.Seek(originalSize, io.SeekStart); err != nil {
		return err
	}

	// Write to new file the music part.
	buf := getByteSlice(128 * 1024)
	defer putByteSlice(buf)
	if _, err = io.CopyBuffer(w, r, buf); err != nil {
		return err
	}
	return nil
}

func saveMP4(tag *IDTag, w io.Writer) error {
	var delete MP4Delete
	fields := reflect.VisibleFields(reflect.TypeOf(*tag))
	for _, field := range fields {
		fieldName := field.Name
		if fieldName == "data" || fieldName == "reader" {
			continue
		}
		if fieldName == "albumArt" && reflect.ValueOf(*tag).FieldByName(fieldName).IsNil() {
			delete = append(delete, fieldName)
			continue
		}
		if reflect.ValueOf(*tag).FieldByName(fieldName).String() == "" {
			delete = append(delete, fieldName)
		}
	}
	r := bufseekio.NewReadSeeker(tag.reader, 128*1024, 4)
	return writeMP4(r, w, tag, delete)

}

func saveFLAC(tag *IDTag, w io.Writer) error {
	r := bufseekio.NewReadSeeker(tag.reader, 128*1024, 4)
	f, fb, err := extractFLACComment(r)
	if err != nil {
		return err
	}
	cmts := flac.New()
	if err := cmts.Add(flac.FIELD_TITLE, tag.title); err != nil {
		return err
	}
	if err := cmts.Add(flac.FIELD_ALBUM, tag.album); err != nil {
		return err
	}
	if err := cmts.Add(flac.FIELD_ARTIST, tag.artist); err != nil {
		return err
	}
	if err := cmts.Add(flac.FIELD_GENRE, tag.genre); err != nil {
		return err
	}
	cmtsmeta := cmts.Marshal()
	if fb.cmtIdx > 0 {
		f.Meta = removeFLACMetaBlock(f.Meta, fb.cmtIdx)
		f.Meta = append(f.Meta, &cmtsmeta)
	} else {
		f.Meta = append(f.Meta, &cmtsmeta)
	}
	if fb.picIdx > 0 {
		f.Meta = removeFLACMetaBlock(f.Meta, fb.picIdx)
	}
	if tag.albumArt != nil {
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, *tag.albumArt, nil); err == nil {
			picture, err := flac.NewFromImageData(flac.PictureTypeFrontCover, "Front cover", buf.Bytes(), "image/jpeg")
			if err != nil {
				return err
			}
			picturemeta := picture.Marshal()
			f.Meta = append(f.Meta, &picturemeta)
		}

	}
	return flacSave(r, w, f.Meta)
}

func saveOGG(tag *IDTag, w io.Writer) error {
	if tag.codec == "vorbis" {
		return saveVorbisTags(tag, w)
	} else if tag.codec == "opus" {
		return saveOpusTags(tag, w)
	}
	return errors.New("codec not supported for OGG")
}
