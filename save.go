package audiometa

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"log"

	mp3TagLib "github.com/bogem/id3v2"
	"github.com/go-flac/flacpicture"
	"github.com/go-flac/flacvorbis"
	"github.com/go-flac/go-flac"
)

// Save saves the corresponding IDTag to the audio file that it references and returns an error if the saving process fails
func (tag *IDTag) Save() error {
	fileType, err := GetFileType(tag.filePath)
	if err != nil {
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
	mp3Tag, err := mp3TagLib.Open(tag.filePath, mp3TagLib.Options{Parse: true})
	if err != nil {
		return err
	}
	defer mp3Tag.Close()
	mp3Tag.SetArtist(tag.artist)
	mp3Tag.SetAlbum(tag.album)
	mp3Tag.SetTitle(tag.title)
	textFrame := mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.bpm,
	}
	mp3Tag.AddFrame("TBPM", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.comments,
	}
	mp3Tag.AddFrame("COMM", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.genre,
	}
	mp3Tag.AddFrame("TCON", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.year,
	}
	mp3Tag.AddFrame("TYER", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.albumArtist,
	}
	mp3Tag.AddFrame("TPE2", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.composer,
	}
	mp3Tag.AddFrame("TCOM", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.copyrightMsg,
	}
	mp3Tag.AddFrame("TCOP", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.date,
	}
	mp3Tag.AddFrame("TDRC", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.encodedBy,
	}
	mp3Tag.AddFrame("TENC", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.lyricist,
	}
	mp3Tag.AddFrame("TEXT", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.fileType,
	}
	mp3Tag.AddFrame("TFLT", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.language,
	}
	mp3Tag.AddFrame("TLAN", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.length,
	}
	mp3Tag.AddFrame("TLEN", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.partOfSet,
	}
	mp3Tag.AddFrame("TPOS", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.publisher,
	}
	mp3Tag.AddFrame("TPUB", textFrame)
	if tag.albumArt != nil {
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
	} else {
		mp3Tag.DeleteFrames("APIC")
	}
	return mp3Tag.Save()
}

func saveMP4(tag *IDTag) error {
	var mp4tag Tags
	var delete []string
	if tag.artist != "" {
		mp4tag.Artist = tag.artist
	} else {
		delete = append(delete, "artist")
	}
	if tag.album != "" {
		mp4tag.Album = tag.album
	} else {
		delete = append(delete, "album")
	}
	if tag.albumArtist != "" {
		mp4tag.AlbumArtist = tag.albumArtist
	} else {
		delete = append(delete, "albumArtist")
	}
	if tag.comments != "" {
		mp4tag.Comment = tag.comments
	} else {
		delete = append(delete, "comment")
	}
	if tag.composer != "" {
		mp4tag.Composer = tag.composer
	} else {
		delete = append(delete, "composer")
	}
	if tag.copyrightMsg != "" {
		mp4tag.Copyright = tag.copyrightMsg
	} else {
		delete = append(delete, "copyright")
	}
	if tag.genre != "" {
		mp4tag.Genre = tag.genre
	} else {
		delete = append(delete, "genre")
	}
	if tag.title != "" {
		mp4tag.Title = tag.title
	} else {
		delete = append(delete, "title")
	}
	if tag.year != "" {
		mp4tag.Year = tag.year
	} else {
		delete = append(delete, "year")
	}
	if tag.albumArt != nil {
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, *tag.albumArt, nil); err != nil {
			mp4tag.Cover = buf.Bytes()
		}
	} else {
		delete = append(delete, "cover")
	}
	mp4tag.Delete = delete
	return WriteMP4(tag.filePath, &mp4tag)
}

func saveFLAC(tag *IDTag) error {
	f, err := flac.ParseFile(tag.filePath)
	if err != nil {
		return err
	}
	_, idx, err := extractFLACComment(tag.filePath)
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
	if idx > 0 {
		f.Meta = removeFLACMetaBlock(f.Meta, idx)
		f.Meta = append(f.Meta, &cmtsmeta)
	} else {
		f.Meta = append(f.Meta, &cmtsmeta)
		log.Printf("length %d", len(f.Meta))
	}
	idx = getFLACPictureIndex(f.Meta)
	if idx > 0 {
		f.Meta = removeFLACMetaBlock(f.Meta, idx)
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
	return f.Save(tag.filePath)
}

func saveOGG(tag *IDTag) error {
	if tag.codec == "vorbis" {
		return saveVorbisTags(tag)
	} else if tag.codec == "opus" {
		return saveOpusTags(tag)
	}
	return errors.New("codec not supported for OGG")
}
