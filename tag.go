package mp3mp4tag

import (
	"bytes"
	"image"
	"os"
)

// Opens the ID tag for the corresponding file as long as it is a supported filetype
// Use the OpenTag command and you will be able to access all metadata associated with the file
func OpenTag(filepath string) (*IDTag, error) {
	tag, err := parse(filepath)
	if err != nil {
		return nil, err
	} else {
		return tag, nil
	}
}

// This operation saves the corresponding IDTag to the mp3/mp4 file that it references and returns an error if the saving process fails
func SaveTag(tag *IDTag) error {
	err := tag.Save()
	if err != nil {
		return err
	} else {
		return nil
	}
}

// clears all tags except the fileUrl tag which is used to reference the file
func (tag *IDTag) ClearAllTags() {
	tag.artist = ""
	tag.albumArtist = ""
	tag.album = ""
	tag.albumArt = nil
	tag.comments = ""
	tag.composer = ""
	tag.genre = ""
	tag.title = ""
	tag.year = ""
	tag.bpm = ""

	tag.id3.contentType = ""
	tag.id3.copyrightMsg = ""
	tag.id3.date = ""
	tag.id3.encodedBy = ""
	tag.id3.lyricist = ""
	tag.id3.fileType = ""
	tag.id3.language = ""
	tag.id3.length = ""
	tag.id3.partOfSet = ""
	tag.id3.publisher = ""
}

// Get the artist for a tag
func (tag *IDTag) Artist() string {
	return tag.artist
}

// Set the artist for a tag
func (tag *IDTag) SetArtist(artist string) {
	tag.artist = artist
}

// Get the album artist for a tag
func (tag *IDTag) AlbumArtist() string {
	return tag.albumArtist
}

// Set teh album artist for a tag
func (tag *IDTag) SetAlbumArtist(albumArtist string) {
	tag.albumArtist = albumArtist
}

// Get the album for a tag
func (tag *IDTag) Album() string {
	return tag.album
}

// Set the album for a tag
func (tag *IDTag) SetAlbum(album string) {
	tag.album = album
}

// Get the commnets for a tag
func (tag *IDTag) Comments() string {
	return tag.comments
}

// Set the comments for a tag
func (tag *IDTag) SetComments(comments string) {
	tag.comments = comments
}

// Get the composer for a tag
func (tag *IDTag) Composer() string {
	return tag.composer
}

// Set the composer for a tag
func (tag *IDTag) SetComposer(composer string) {
	tag.composer = composer
}

// Get the genre for a tag
func (tag *IDTag) Genre() string {
	return tag.genre
}

// Set the genre for a tag
func (tag *IDTag) SetGenre(genre string) {
	tag.genre = genre
}

// Get the title for a tag
func (tag *IDTag) Title() string {
	return tag.title
}

// Set the title for a tag
func (tag *IDTag) SetTitle(title string) {
	tag.title = title
}

// Get the year for a tag
func (tag *IDTag) Year() string {
	return tag.year
}

// Set the year for a tag
func (tag *IDTag) SetYear(year string) {
	tag.year = year
}

// Get the BPM for a tag
func (tag *IDTag) BPM() string {
	return tag.bpm
}

// Set the BPM for a tag
func (tag *IDTag) SetBPM(bpm string) {
	tag.bpm = bpm
}

// Get the Copyright Messgae for a tag
func (tag *IDTag) CopyrightMsg() string {
	return tag.id3.copyrightMsg
}

// Set the Copyright Message for a tag
func (tag *IDTag) SetCopyrightMsg(copyrightMsg string) {
	tag.id3.copyrightMsg = copyrightMsg
}

// Get the date for a tag
func (tag *IDTag) Date() string {
	return tag.id3.date
}

// Set the date for a tag
func (tag *IDTag) SetDate(date string) {
	tag.id3.date = date
}

// Get who encoded the tag
func (tag *IDTag) EncodedBy() string {
	return tag.id3.encodedBy
}

// Set who encoded the tag
func (tag *IDTag) SetEncodedBy(encodedBy string) {
	tag.id3.encodedBy = encodedBy
}

// Get the lyricist for the tag
func (tag *IDTag) Lyricist() string {
	return tag.id3.lyricist
}

// Set the lyricist for the tag
func (tag *IDTag) SetLyricist(lyricist string) {
	tag.id3.lyricist = lyricist
}

// Get the filetype of the tag
func (tag *IDTag) FileType() string {
	return tag.id3.fileType
}

// Set the filtype of the tag
func (tag *IDTag) SetFileType(fileType string) {
	tag.id3.fileType = fileType
}

// Get the language of the tag
func (tag *IDTag) Language() string {
	return tag.id3.language
}

// Set the lanuguage of the tag
func (tag *IDTag) SetLanguage(language string) {
	tag.id3.language = language
}

// Get the langth of the tag
func (tag *IDTag) Length() string {
	return tag.id3.length
}

// Set the length of the tag
func (tag *IDTag) SetLength(length string) {
	tag.id3.length = length
}

// Get if tag is part of a set
func (tag *IDTag) PartOfSet() string {
	return tag.id3.partOfSet
}

// Set if the tag is part of a set
func (tag *IDTag) SetPartOfSet(partOfSet string) {
	tag.id3.partOfSet = partOfSet
}

// Get publisher for the tag
func (tag *IDTag) Publisher() string {
	return tag.id3.publisher
}

// Set publihser for the tag
func (tag *IDTag) SetPublisher(publisher string) {
	tag.id3.publisher = publisher
}

// Get the album art for the tag as an *image.Image
func (tag *IDTag) AlbumArt() *image.Image {
	return tag.albumArt
}

// Set the album art by passing a byte array for the album art
func (tag *IDTag) SetAlbumArtFromByteArray(albumArt []byte) error {
	img, _, err := image.Decode(bytes.NewReader(albumArt))
	if err != nil {
		return err
	} else {
		tag.albumArt = &img
		return nil
	}
}

// Set the album art by passing an *image.Image as the album art
func (tag *IDTag) SetAlbumArtFromImage(albumArt *image.Image) {
	tag.albumArt = albumArt
}

// Set the album art by passing a filepath as a string
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
	tag.albumArt = &img
	return nil
}
