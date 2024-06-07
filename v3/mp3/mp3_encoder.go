package mp3

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
)

func SaveMP3(tag *MP3Tag, w io.Writer) error {
	r := tag.reader
	if _, err := r.Seek(0, 0); err != nil {
		return err
	}
	readerBytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	r = bytes.NewReader(readerBytes)
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
		textFrame := mp3TagLib.TextFrame{
			Encoding: mp3TagLib.EncodingUTF8,
			Text:     reflect.ValueOf(*tag).FieldByName(k).String(),
		}
		mp3Tag.AddFrame(v, textFrame)
	}
	mp3Tag.DeleteFrames("APIC")
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, *tag.AlbumArt, nil); err == nil {
		mp3Tag.AddAttachedPicture(mp3TagLib.PictureFrame{
			Encoding:    mp3TagLib.EncodingUTF8,
			MimeType:    "image/jpeg",
			PictureType: mp3TagLib.PTFrontCover,
			Description: "Front cover",
			Picture:     buf.Bytes(),
		})
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
		if _, err := mp3Tag.WriteTo(t); err != nil {
			return err
		}
		// Seek to a music part of original file.
		if _, err = r.Seek(originalSize, io.SeekStart); err != nil {
			return err
		}
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
