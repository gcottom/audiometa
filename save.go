package audiometa

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"log"
	"strconv"

	mp4tagWriter "github.com/Sorrow446/go-mp4tag"
	mp3TagLib "github.com/bogem/id3v2"
	"github.com/go-flac/flacpicture"
	"github.com/go-flac/flacvorbis"
	"github.com/go-flac/go-flac"
)

// Save saves the corresponding IDTag to the audio file that it references and returns an error if the saving process fails
func (tag *IDTag) Save() error {
	fileType, err := GetFileType(tag.FilePath)
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
	return fmt.Errorf("no method available for filetype:%s", tag.FileType)
}

func saveMP3(tag *IDTag) error {
	mp3Tag, err := mp3TagLib.Open(tag.FilePath, mp3TagLib.Options{Parse: true})
	if err != nil {
		return err
	}
	defer mp3Tag.Close()
	mp3Tag.SetArtist(tag.Artist)
	mp3Tag.SetAlbum(tag.Album)
	mp3Tag.SetTitle(tag.Title)
	textFrame := mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.BPM,
	}
	mp3Tag.AddFrame("TBPM", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.Comments,
	}
	mp3Tag.AddFrame("COMM", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.Genre,
	}
	mp3Tag.AddFrame("TCON", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.Year,
	}
	mp3Tag.AddFrame("TYER", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.AlbumArtist,
	}
	mp3Tag.AddFrame("TPE2", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.Composer,
	}
	mp3Tag.AddFrame("TCOM", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.CopyrightMsg,
	}
	mp3Tag.AddFrame("TCOP", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.Date,
	}
	mp3Tag.AddFrame("TDRC", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.EncodedBy,
	}
	mp3Tag.AddFrame("TENC", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.Lyricist,
	}
	mp3Tag.AddFrame("TEXT", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.FileType,
	}
	mp3Tag.AddFrame("TFLT", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.Language,
	}
	mp3Tag.AddFrame("TLAN", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.Length,
	}
	mp3Tag.AddFrame("TLEN", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.PartOfSet,
	}
	mp3Tag.AddFrame("TPOS", textFrame)
	textFrame = mp3TagLib.TextFrame{
		Encoding: mp3TagLib.EncodingUTF8,
		Text:     tag.Publisher,
	}
	mp3Tag.AddFrame("TPUB", textFrame)
	if tag.AlbumArt != nil {
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
	} else {
		mp3Tag.DeleteFrames("APIC")
	}
	return mp3Tag.Save()
}

func saveMP4(tag *IDTag) error {
	var mp4tag mp4tagWriter.MP4Tags
	var delete []string
	if tag.Artist != "" {
		mp4tag.Artist = tag.Artist
	} else {
		delete = append(delete, "artist")
	}
	if tag.Album != "" {
		mp4tag.Album = tag.Album
	} else {
		delete = append(delete, "album")
	}
	if tag.AlbumArtist != "" {
		mp4tag.AlbumArtist = tag.AlbumArtist
	} else {
		delete = append(delete, "albumArtist")
	}
	if tag.Comments != "" {
		mp4tag.Comment = tag.Comments
	} else {
		delete = append(delete, "comment")
	}
	if tag.Composer != "" {
		mp4tag.Composer = tag.Composer
	} else {
		delete = append(delete, "composer")
	}
	if tag.CopyrightMsg != "" {
		mp4tag.Copyright = tag.CopyrightMsg
	} else {
		delete = append(delete, "copyright")
	}
	if tag.Genre != "" {
		mp4tag.CustomGenre = tag.Genre
	} else {
		delete = append(delete, "genre")
	}
	if tag.Title != "" {
		mp4tag.Title = tag.Title
	} else {
		delete = append(delete, "title")
	}
	if tag.Year != "" {
		y, err := strconv.Atoi(tag.Year)
		if err != nil {
			mp4tag.Year = int32(y)
		}

	} else {
		delete = append(delete, "year")
	}
	if tag.AlbumArt != nil {
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, *tag.AlbumArt, nil); err != nil {
			mp4tag.Pictures = []*mp4tagWriter.MP4Picture{{Format: mp4tagWriter.ImageTypeJPEG, Data: buf.Bytes()}}
		}
	} else {
		delete = append(delete, "cover")
	}
	mp4, err := mp4tagWriter.Open(tag.FilePath)
	if err != nil {
		return err
	}
	return mp4.Write(&mp4tag, delete)
}

func saveFLAC(tag *IDTag) error {
	f, err := flac.ParseFile(tag.FilePath)
	if err != nil {
		return err
	}
	_, idx, err := extractFLACComment(tag.FilePath)
	if err != nil {
		return err
	}
	cmts := flacvorbis.New()
	if err := cmts.Add(flacvorbis.FIELD_TITLE, tag.Title); err != nil {
		return err
	}
	if err := cmts.Add(flacvorbis.FIELD_ALBUM, tag.Album); err != nil {
		return err
	}
	if err := cmts.Add(flacvorbis.FIELD_ARTIST, tag.Artist); err != nil {
		return err
	}
	if err := cmts.Add(flacvorbis.FIELD_GENRE, tag.Genre); err != nil {
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
	if tag.AlbumArt != nil {
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, *tag.AlbumArt, nil); err == nil {
			picture, err := flacpicture.NewFromImageData(flacpicture.PictureTypeFrontCover, "Front cover", buf.Bytes(), "image/jpeg")
			if err != nil {
				return err
			}
			picturemeta := picture.Marshal()
			f.Meta = append(f.Meta, &picturemeta)
		}

	}
	return f.Save(tag.FilePath)
}

func saveOGG(tag *IDTag) error {
	if tag.Codec == "vorbis" {
		return saveVorbisTags(tag)
	} else if tag.Codec == "opus" {
		return saveOpusTags(tag)
	}
	return errors.New("codec not supported for OGG")
}
