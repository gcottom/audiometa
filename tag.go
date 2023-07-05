package mp3mp4tag

import (
	"bytes"
	"image"
	"os"
)

func OpenTag(filepath string) (*IDTag, error) {
	tag, err := parse(filepath)
	if err != nil {
		return nil, err
	} else {
		return tag, nil
	}
}
func SaveTag(tag *IDTag) error {
	err := tag.Save()
	if err != nil {
		return err
	} else {
		return nil
	}
}
func (tag *IDTag) ClearAllTags() {
	//clears all tags except the fileUrl tag which is used to reference the file
	tag.artist = ""
	tag.albumArtist = ""
	tag.album = ""
	tag.albumArt = nil
	tag.comments = ""
	tag.composer = ""
	tag.genre = ""
	tag.title = ""
	tag.year = 0
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

func (tag *IDTag) Artist() string {
	return tag.artist
}
func (tag *IDTag) SetArtist(artist string) {
	tag.artist = artist
}
func (tag *IDTag) AlbumArtist() string {
	return tag.albumArtist
}
func (tag *IDTag) SetAlbumArtist(albumArtist string) {
	tag.albumArtist = albumArtist
}
func (tag *IDTag) Album() string {
	return tag.album
}
func (tag *IDTag) SetAlbum(album string) {
	tag.album = album
}
func (tag *IDTag) Comments() string {
	return tag.comments
}
func (tag *IDTag) SetComments(comments string) {
	tag.comments = comments
}
func (tag *IDTag) Composer() string {
	return tag.composer
}
func (tag *IDTag) SetComposer(composer string) {
	tag.composer = composer
}
func (tag *IDTag) Genre() string {
	return tag.genre
}
func (tag *IDTag) SetGenre(genre string) {
	tag.genre = genre
}
func (tag *IDTag) Title() string {
	return tag.title
}
func (tag *IDTag) SetTitle(title string) {
	tag.title = title
}
func (tag *IDTag) Year() int {
	return tag.year
}
func (tag *IDTag) SetYear(year int) {
	tag.year = year
}
func (tag *IDTag) BPM() string {
	return tag.bpm
}
func (tag *IDTag) SetBPM(bpm string) {
	tag.bpm = bpm
}
func (tag *IDTag) ContentType() string {
	return tag.id3.contentType
}
func (tag *IDTag) SetContentType(contentType string) {
	tag.id3.contentType = contentType
}
func (tag *IDTag) CopyrightMsg() string {
	return tag.id3.copyrightMsg
}
func (tag *IDTag) SetCopyrightMsg(copyrightMsg string) {
	tag.id3.copyrightMsg = copyrightMsg
}
func (tag *IDTag) Date() string {
	return tag.id3.date
}
func (tag *IDTag) SetDate(date string) {
	tag.id3.date = date
}
func (tag *IDTag) EncodedBy() string {
	return tag.id3.encodedBy
}
func (tag *IDTag) SetEncodedBy(encodedBy string) {
	tag.id3.encodedBy = encodedBy
}
func (tag *IDTag) Lyricist() string {
	return tag.id3.lyricist
}
func (tag *IDTag) SetLyricist(lyricist string) {
	tag.id3.lyricist = lyricist
}
func (tag *IDTag) FileType() string {
	return tag.id3.fileType
}
func (tag *IDTag) SetFileType(fileType string) {
	tag.id3.fileType = fileType
}
func (tag *IDTag) Language() string {
	return tag.id3.language
}
func (tag *IDTag) SetLanguage(language string) {
	tag.id3.language = language
}
func (tag *IDTag) Length() string {
	return tag.id3.length
}
func (tag *IDTag) SetLength(length string) {
	tag.id3.length = length
}
func (tag *IDTag) PartOfSet() string {
	return tag.id3.partOfSet
}
func (tag *IDTag) SetPartOfSet(partOfSet string) {
	tag.id3.partOfSet = partOfSet
}
func (tag *IDTag) Publisher() string {
	return tag.id3.publisher
}
func (tag *IDTag) SetPublisher(publisher string) {
	tag.id3.publisher = publisher
}
func (tag *IDTag) AlbumArt() *image.Image {
	return tag.albumArt
}
func (tag *IDTag) SetAlbumArtFromByteArray(albumArt []byte) error {
	img, _, err := image.Decode(bytes.NewReader(albumArt))
	if err != nil {
		return err
	} else {
		tag.albumArt = &img
		return nil
	}
}
func (tag *IDTag) SetAlbumArtFromImage(albumArt *image.Image) {
	tag.albumArt = albumArt
}
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
