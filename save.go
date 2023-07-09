package mp3mp4tag

import (
	"bytes"
	"image/jpeg"
	"log"

	mp4tagWriter "github.com/Sorrow446/go-mp4tag"
	mp3TagLib "github.com/bogem/id3v2"
)

// This operation saves the corresponding IDTag to the mp3/mp4 file that it references and returns an error if the saving process fails
func (tag *IDTag) Save() error {
	fileType, err := getFileType(tag.fileUrl)
	if err != nil {
		return err
	}
	if *fileType == "mp3" {
		mp3Tag, err := mp3TagLib.Open(tag.fileUrl, mp3TagLib.Options{Parse: true})
		if err != nil {
			log.Fatal("Error while opening mp3 file for writing: ", err)
		}
		mp3Tag.SetArtist(tag.artist)
		mp3Tag.SetAlbum(tag.album)

		mp3Tag.SetTitle(tag.title)

		//if tag.bpm != "" {
		textFrame := mp3TagLib.TextFrame{
			Encoding: mp3TagLib.EncodingUTF8,
			Text:     tag.bpm,
		}
		mp3Tag.AddFrame("TBPM", textFrame)
		//}
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
			Text:     tag.id3.copyrightMsg,
		}
		mp3Tag.AddFrame("TCOP", textFrame)
		textFrame = mp3TagLib.TextFrame{
			Encoding: mp3TagLib.EncodingUTF8,
			Text:     tag.id3.date,
		}
		mp3Tag.AddFrame("TDRC", textFrame)
		textFrame = mp3TagLib.TextFrame{
			Encoding: mp3TagLib.EncodingUTF8,
			Text:     tag.id3.encodedBy,
		}
		mp3Tag.AddFrame("TENC", textFrame)
		textFrame = mp3TagLib.TextFrame{
			Encoding: mp3TagLib.EncodingUTF8,
			Text:     tag.id3.lyricist,
		}
		mp3Tag.AddFrame("TEXT", textFrame)
		textFrame = mp3TagLib.TextFrame{
			Encoding: mp3TagLib.EncodingUTF8,
			Text:     tag.id3.fileType,
		}
		mp3Tag.AddFrame("TFLT", textFrame)
		textFrame = mp3TagLib.TextFrame{
			Encoding: mp3TagLib.EncodingUTF8,
			Text:     tag.id3.language,
		}
		mp3Tag.AddFrame("TLAN", textFrame)
		textFrame = mp3TagLib.TextFrame{
			Encoding: mp3TagLib.EncodingUTF8,
			Text:     tag.id3.length,
		}
		mp3Tag.AddFrame("TLEN", textFrame)
		textFrame = mp3TagLib.TextFrame{
			Encoding: mp3TagLib.EncodingUTF8,
			Text:     tag.id3.partOfSet,
		}
		mp3Tag.AddFrame("TPOS", textFrame)
		textFrame = mp3TagLib.TextFrame{
			Encoding: mp3TagLib.EncodingUTF8,
			Text:     tag.id3.publisher,
		}
		mp3Tag.AddFrame("TPUB", textFrame)
		if tag.albumArt != nil {
			buf := new(bytes.Buffer)
			jpeg.Encode(buf, *tag.albumArt, nil)
			bytes := buf.Bytes()
			pic := mp3TagLib.PictureFrame{
				Encoding:    mp3TagLib.EncodingUTF8,
				MimeType:    "image/jpeg",
				PictureType: mp3TagLib.PTFrontCover,
				Description: "Front cover",
				Picture:     bytes,
			}
			mp3Tag.AddAttachedPicture(pic)
		} else {
			mp3Tag.DeleteFrames("APIC")
		}
		err = mp3Tag.Save()
		if err != nil {
			return err
		}
		mp3Tag.Close()
	} else {
		var mp4tag mp4tagWriter.Tags
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
		if tag.id3.copyrightMsg != "" {
			mp4tag.Copyright = tag.id3.copyrightMsg
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
			jpeg.Encode(buf, *tag.albumArt, nil)
			bytes := buf.Bytes()
			mp4tag.Cover = bytes
		} else {
			delete = append(delete, "cover")
		}

		mp4tag.Delete = delete

		err := mp4tagWriter.Write(tag.fileUrl, &mp4tag)
		if err != nil {
			return err
		}
	}
	return nil
}
