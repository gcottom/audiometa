package audiometa

import (
	"bytes"
	"image"
	"os"
)

// OpenTag Opens the ID tag for the corresponding file as long as it is a supported filetype
// Use the OpenTag command and you will be able to access all metadata associated with the file
func OpenTag(filepath string) (*IDTag, error) {
	return parse(filepath)
}

// SaveTag saves the corresponding IDTag to the audio file that it references and returns an error if the saving process fails
func SaveTag(tag *IDTag) error {
	return tag.Save()
}

// ClearAllTags clears all tags except the fileUrl tag which is used to reference the file, takes an optional parameter "preserveUnkown": when this is true passThroughMap is not cleared and unknown tags are preserved
func (tag *IDTag) ClearAllTags(preserveUnknown ...bool) {
	tag.Artist = ""
	tag.AlbumArtist = ""
	tag.Album = ""
	tag.AlbumArt = nil
	tag.Comments = ""
	tag.Composer = ""
	tag.Genre = ""
	tag.Title = ""
	tag.Year = ""
	tag.BPM = ""

	tag.CopyrightMsg = ""
	tag.Date = ""
	tag.EncodedBy = ""
	tag.Lyricist = ""
	tag.FileType = ""
	tag.Language = ""
	tag.Length = ""
	tag.PartOfSet = ""
	tag.Publisher = ""

	preserve := false
	if len(preserveUnknown) != 0 {
		preserve = preserveUnknown[0]
	}
	if !preserve {
		tag.PassThrough = make(map[string]string)
	}

}

// SetAdditionalTag sets an additional (unmapped) tag taking an id and value (id,value) (ogg only)
func (tag *IDTag) SetAdditionalTag(id string, value string) {
	tag.PassThrough[id] = value
}

// SetAlbumArtFromByteArray sets the album art by passing a byte array for the album art
func (tag *IDTag) SetAlbumArtFromByteArray(albumArt []byte) error {
	img, _, err := image.Decode(bytes.NewReader(albumArt))
	if err != nil {
		return err
	} else {
		tag.AlbumArt = &img
		return nil
	}
}

// SetAlbumArtFromImage sets the album art by passing an *image.Image as the album art
func (tag *IDTag) SetAlbumArtFromImage(albumArt *image.Image) {
	tag.AlbumArt = albumArt
}

// SetAlbumArtFromFilePath sets the album art by passing a filepath as a string
func (tag *IDTag) SetAlbumArtFromFilePath(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return err
	}
	tag.AlbumArt = &img
	return nil
}
