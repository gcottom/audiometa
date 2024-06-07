package mp3

import (
	"fmt"
	"image"
	"io"
	"strconv"
)

var mp3TextFrames = map[string]string{
	"Artist":            "TPE1",
	"Title":             "TIT2",
	"Album":             "TALB",
	"BPM":               "TBPM",
	"Genre":             "TCON",
	"Year":              "TYER",
	"AlbumArtist":       "TPE2",
	"Composer":          "TCOM",
	"Copyright":         "TCOP",
	"Date":              "TDRC",
	"EncodedBy":         "TENC",
	"Lyricist":          "TEXT",
	"Language":          "TLAN",
	"Length":            "TLEN",
	"DiscNumberString":  "TPOS",
	"Publisher":         "TPUB",
	"SubTitle":          "TIT3",
	"ISRC":              "TSRC",
	"TrackNumberString": "TRCK",
}

type MP3Tag struct {
	Album             string
	AlbumArt          *image.Image
	AlbumArtist       string
	Artist            string
	BPM               string
	Composer          string //multiple seperated by /
	Copyright         string //must begin with year yyyy and a space
	Date              string //limit 4 char DDMM
	DiscNumber        int    //TPOS
	DiscNumberString  string
	DiscTotal         int
	EncodedBy         string
	Genre             string
	ISRC              string
	Language          string
	Length            string //length in millisecs
	Lyricist          string //multiple seperated by /
	Title             string
	TrackNumber       int
	TrackNumberString string
	TrackTotal        int
	Publisher         string
	SubTitle          string
	Year              string

	reader io.ReadSeeker
}

func (t *MP3Tag) GetAlbum() string {
	return t.Album
}
func (t *MP3Tag) GetCoverArt() *image.Image {
	return t.AlbumArt
}
func (t *MP3Tag) GetAlbumArtist() string {
	return t.AlbumArtist
}
func (t *MP3Tag) GetArtist() string {
	return t.Artist
}
func (t *MP3Tag) GetBPM() int {
	bpm, err := strconv.Atoi(t.BPM)
	if err != nil {
		return 0
	}
	return bpm
}
func (t *MP3Tag) GetComposer() string {
	return t.Composer
}
func (t *MP3Tag) GetCopyright() string {
	return t.Copyright
}
func (t *MP3Tag) GetDate() string {
	return t.Date
}
func (t *MP3Tag) GetDiscNumber() int {
	return t.DiscNumber
}
func (t *MP3Tag) GetDiscTotal() int {
	return t.DiscTotal
}
func (t *MP3Tag) GetEncoder() string {
	return t.EncodedBy
}
func (t *MP3Tag) GetGenre() string {
	return t.Genre
}
func (t *MP3Tag) GetISRC() string {
	return t.ISRC
}
func (t *MP3Tag) GetLanguage() string {
	return t.Language
}
func (t *MP3Tag) GetLength() string {
	return t.Length
}
func (t *MP3Tag) GetLyricist() string {
	return t.Lyricist
}
func (t *MP3Tag) GetTitle() string {
	return t.Title
}
func (t *MP3Tag) GetTrackNumber() int {
	return t.TrackNumber
}
func (t *MP3Tag) GetTrackTotal() int {
	return t.TrackTotal
}
func (t *MP3Tag) GetPublisher() string {
	return t.Publisher
}
func (t *MP3Tag) GetSubTitle() string {
	return t.SubTitle
}
func (t *MP3Tag) GetYear() int {
	year, err := strconv.Atoi(t.Year)
	if err != nil {
		return 0
	}
	return year
}

func (t *MP3Tag) SetAlbum(album string) {
	t.Album = album
}
func (t *MP3Tag) SetCoverArt(coverArt *image.Image) {
	t.AlbumArt = coverArt
}
func (t *MP3Tag) SetAlbumArtist(albumArtist string) {
	t.AlbumArtist = albumArtist
}
func (t *MP3Tag) SetArtist(artist string) {
	t.Artist = artist
}
func (t *MP3Tag) SetBPM(bpm int) {
	t.BPM = fmt.Sprint(bpm)
}
func (t *MP3Tag) SetComposer(composer string) {
	t.Composer = composer
}
func (t *MP3Tag) SetCopyright(copyright string) {
	t.Copyright = copyright
}
func (t *MP3Tag) SetDate(date string) {
	t.Date = date
}
func (t *MP3Tag) SetDiscNumber(discNumber int) {
	t.DiscNumber = discNumber
}
func (t *MP3Tag) SetDiscTotal(discTotal int) {
	t.DiscTotal = discTotal
}
func (t *MP3Tag) SetEncoder(encoder string) {
	t.EncodedBy = encoder
}
func (t *MP3Tag) SetGenre(genre string) {
	t.Genre = genre
}
func (t *MP3Tag) SetISRC(isrc string) {
	t.ISRC = isrc
}
func (t *MP3Tag) SetLanguage(language string) {
	t.Language = language
}
func (t *MP3Tag) SetLength(length string) {
	t.Length = length
}
func (t *MP3Tag) SetLyricist(lyricist string) {
	t.Lyricist = lyricist
}
func (t *MP3Tag) SetTitle(title string) {
	t.Title = title
}
func (t *MP3Tag) SetTrackNumber(trackNumber int) {
	t.TrackNumber = trackNumber
}
func (t *MP3Tag) SetTrackTotal(trackTotal int) {
	t.TrackTotal = trackTotal
}
func (t *MP3Tag) SetPublisher(publisher string) {
	t.Publisher = publisher
}
func (t *MP3Tag) SetSubTitle(subTitle string) {
	t.SubTitle = subTitle
}
func (t *MP3Tag) SetYear(year int) {
	t.Year = fmt.Sprint(year)
}
func (t *MP3Tag) Save(w io.Writer) error {
	return SaveMP3(t, w)
}
