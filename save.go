package audiometa

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"reflect"

	"github.com/aler9/writerseeker"
	mp3TagLib "github.com/bogem/id3v2/v2"
	"github.com/gcottom/audiometa/v2/flac"
	"github.com/sunfish-shogi/bufseekio"
)

// Save writes the full ID Tag and audio to the io.Writer w.
// If w is of type *os.File, Save overwrites the existing file
// and when complete, w points to the end of the file.
func (tag *IDTag) Save(w io.Writer) error {
	fileType := FileType(tag.fileType)
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
	return ErrNoMethodAvlble
}

func saveMP3(tag *IDTag, w io.Writer) error {
	r := tag.reader

	readerBytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	r = bytes.NewReader(readerBytes)
	fmt.Println(len(readerBytes))
	if _, err = r.Seek(0, 0); err != nil {
		return err
	}

	mp3Tag, err := mp3TagLib.ParseReader(r, mp3TagLib.Options{Parse: true})
	if err != nil {
		return err
	}
	originalSize := int64(mp3Tag.Size())
	fmt.Println(originalSize)
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
	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	if reflect.TypeOf(w) == reflect.TypeOf(new(os.File)) {
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
		//in and out are the same file so we have to temp it
		t := &writerseeker.WriterSeeker{}
		defer t.Close()
		// Write tag in new file.
		if tagbytes, err := mp3Tag.WriteTo(t); err != nil {
			return err
		} else {
			fmt.Println(tagbytes)
		}

		// Seek to a music part of original file.
		if _, err = r.Seek(originalSize, io.SeekStart); err != nil {
			return err
		}
		// Write to new file the music part.
		//musicData, err := io.ReadAll(r)
		//if err != nil {
		//	return err
		//}
		if _, err := io.Copy(t, r); err != nil {
			return err
		}
		if _, err = t.Seek(0, io.SeekStart); err != nil {
			return err
		}
		if _, err := io.Copy(w2, bytes.NewReader(t.Bytes())); err != nil {
			return err
		}
		if _, err = f.Seek(0, io.SeekEnd); err != nil {
			return err
		}
		return nil
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
	if _, err = io.Copy(w, r); err != nil {
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
	needsTemp := reflect.TypeOf(w) == reflect.TypeOf(new(os.File))
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
	cmtsmeta, err := cmts.Marshal()
	if err != nil {
		return err
	}
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
	if _, err = r.Seek(0, io.SeekStart); err != nil {
		return err
	}
	return flacSave(r, w, f.Meta, needsTemp)
}

func saveOGG(tag *IDTag, w io.Writer) error {
	if tag.codec == "vorbis" {
		return saveVorbisTags(tag, w)
	} else if tag.codec == "opus" {
		return saveOpusTags(tag, w)
	}
	return ErrOggCodecNotSpprtd
}
