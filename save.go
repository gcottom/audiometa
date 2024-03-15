package audiometa

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"reflect"

	"github.com/aler9/writerseeker"
	mp3TagLib "github.com/bogem/id3v2/v2"
	"github.com/go-flac/flacpicture"
	"github.com/go-flac/flacvorbis"
	"github.com/go-flac/go-flac"
	"github.com/sunfish-shogi/bufseekio"
)

// Save saves the corresponding IDTag to the audio file that it references and returns an error if the saving process fails
func (tag *IDTag) Save() error {
	fileType, err := GetFileType(tag.filePath)
	if err != nil {
		return err
	}
	if _, err := tag.reader.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if fileType == MP3 {
		return saveMP3(tag)
	} else if fileType == M4A || fileType == M4B || fileType == M4P || fileType == MP4 {
		return saveMP4(tag)
	} else if fileType == FLAC {
		return saveFLAC(tag)
	} else if fileType == OGG {
		return saveOGG(tag)
	}
	return fmt.Errorf("no method available for filetype:%s", tag.fileType)
}

func saveMP3(tag *IDTag) error {
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
	ws := &writerseeker.WriterSeeker{}
	// Write tag in new file.
	if _, err = mp3Tag.WriteTo(ws); err != nil {
		return err
	}
	// Seek to a music part of original file.
	if _, err = r.Seek(originalSize, io.SeekStart); err != nil {
		return err
	}

	// Write to new file the music part.
	buf := getByteSlice(128 * 1024)
	defer putByteSlice(buf)
	if _, err = io.CopyBuffer(ws, r, buf); err != nil {
		return err
	}
	buffy := ws.Bytes()
	return WriteFile(tag.filePath, buffy)
}

func saveMP4(tag *IDTag) error {
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
	buffy, err := writeMP4(r, tag, delete)
	if err != nil {
		return err
	}
	return WriteFile(tag.filePath, buffy)

}

func saveFLAC(tag *IDTag) error {
	r := bufseekio.NewReadSeeker(tag.reader, 128*1024, 4)
	f, err := flac.ParseBytes(r)
	if err != nil {
		return err
	}
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return err
	}
	fb, err := extractFLACComment(r)
	if err != nil {
		return err
	}
	cmts := flacvorbis.New()
	if err := cmts.Add(flacvorbis.FIELD_TITLE, tag.title); err != nil {
		return err
	}
	if err := cmts.Add(flacvorbis.FIELD_ALBUM, tag.album); err != nil {
		return err
	}
	if err := cmts.Add(flacvorbis.FIELD_ARTIST, tag.artist); err != nil {
		return err
	}
	if err := cmts.Add(flacvorbis.FIELD_GENRE, tag.genre); err != nil {
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
			picture, err := flacpicture.NewFromImageData(flacpicture.PictureTypeFrontCover, "Front cover", buf.Bytes(), "image/jpeg")
			if err != nil {
				return err
			}
			picturemeta := picture.Marshal()
			f.Meta = append(f.Meta, &picturemeta)
		}

	}

	//[]byte for future use
	buffy := f.Marshal()
	return WriteFile(tag.filePath, buffy)
}

func saveOGG(tag *IDTag) error {
	if tag.codec == "vorbis" {
		buffy, err := saveVorbisTags(tag)
		if err != nil {
			return err
		}
		return WriteFile(tag.filePath, buffy)
	} else if tag.codec == "opus" {
		buffy, err := saveOpusTags(tag)
		if err != nil {
			return err
		}
		return WriteFile(tag.filePath, buffy)
	}
	return errors.New("codec not supported for OGG")
}
