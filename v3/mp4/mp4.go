package mp4

import (
	"fmt"
	"image"
	"io"
	"strconv"

	mp4lib "github.com/abema/go-mp4"
)

var atomsMap = map[mp4lib.BoxType]string{
	{'\251', 'a', 'l', 'b'}: "Album",
	{'a', 'A', 'R', 'T'}:    "AlbumArtist",
	{'\251', 'A', 'R', 'T'}: "Artist",
	{'\251', 'c', 'm', 't'}: "Comments",
	{'\251', 'w', 'r', 't'}: "Composer",
	{'c', 'p', 'r', 't'}:    "CopyrightMsg",
	{'c', 'o', 'v', 'r'}:    "CoverArt",
	{'\251', 'g', 'e', 'n'}: "Genre", //check for the gnre atom, can't coexis:"Genre":
	{'\251', 'n', 'a', 'm'}: "Title",
	{'\251', 'd', 'a', 'y'}: "Year",
	{'t', 'r', 'k', 'n'}:    "TrackNumber", //2uint16 (track) (totaltracks:"TrackNumber":
	{'d', 'i', 's', 'k'}:    "DiscNumber",  //2uint16 (disc) (totaldiscs:"DiscNumber":
	{'\251', 't', 'o', 'o'}: "Encoder",
	{'t', 'm', 'p', 'o'}:    "BPM", //bigEndianUin:"BPM":
}

type MP4Tag struct {
	Album       string
	AlbumArtist string
	Artist      string
	BPM         int
	Comments    string
	Composer    string
	Copyright   string
	CoverArt    *image.Image
	Encoder     string
	Genre       string
	Title       string
	TrackNumber int
	TrackTotal  int
	DiscNumber  int
	DiscTotal   int
	Year        string

	reader io.ReadSeeker
}

func (m *MP4Tag) GetAlbum() string {
	return m.Album
}

func (m *MP4Tag) GetAlbumArtist() string {
	return m.AlbumArtist
}

func (m *MP4Tag) GetArtist() string {
	return m.Artist
}

func (m *MP4Tag) GetBPM() int {
	return m.BPM
}

func (m *MP4Tag) GetComments() string {
	return m.Comments
}

func (m *MP4Tag) GetComposer() string {
	return m.Composer
}

func (m *MP4Tag) GetCopyright() string {
	return m.Copyright
}

func (m *MP4Tag) GetCoverArt() *image.Image {
	return m.CoverArt
}

func (m *MP4Tag) GetEncoder() string {
	return m.Encoder
}

func (m *MP4Tag) GetGenre() string {
	return m.Genre
}

func (m *MP4Tag) GetTitle() string {
	return m.Title
}

func (m *MP4Tag) GetTrackNumber() int {
	return m.TrackNumber
}

func (m *MP4Tag) GetTrackTotal() int {
	return m.TrackTotal
}

func (m *MP4Tag) GetDiscNumber() int {
	return m.DiscNumber
}

func (m *MP4Tag) GetDiscTotal() int {
	return m.DiscTotal
}

func (m *MP4Tag) GetYear() int {
	year, err := strconv.Atoi(m.Year)
	if err != nil {
		return 0
	}
	return year
}

func (m *MP4Tag) SetAlbum(album string) {
	m.Album = album
}
func (m *MP4Tag) SetAlbumArtist(albumArtist string) {
	m.AlbumArtist = albumArtist
}
func (m *MP4Tag) SetArtist(artist string) {
	m.Artist = artist
}
func (m *MP4Tag) SetBPM(bpm int) {
	m.BPM = bpm
}
func (m *MP4Tag) SetComments(comments string) {
	m.Comments = comments
}
func (m *MP4Tag) SetComposer(composer string) {
	m.Composer = composer
}
func (m *MP4Tag) SetCopyright(copyright string) {
	m.Copyright = copyright
}
func (m *MP4Tag) SetCoverArt(coverArt *image.Image) {
	m.CoverArt = coverArt
}
func (m *MP4Tag) SetEncoder(encoder string) {
	m.Encoder = encoder
}
func (m *MP4Tag) SetGenre(genre string) {
	m.Genre = genre
}
func (m *MP4Tag) SetTitle(title string) {
	m.Title = title
}
func (m *MP4Tag) SetTrackNumber(trackNumber int) {
	m.TrackNumber = trackNumber
}
func (m *MP4Tag) SetTrackTotal(trackTotal int) {
	m.TrackTotal = trackTotal
}
func (m *MP4Tag) SetDiscNumber(discNumber int) {
	m.DiscNumber = discNumber
}
func (m *MP4Tag) SetDiscTotal(discTotal int) {
	m.DiscTotal = discTotal
}
func (m *MP4Tag) SetYear(year int) {
	m.Year = fmt.Sprint(year)
}

func (m *MP4Tag) Save(w io.Writer) error {
	return SaveMP4(m.reader, w, m)
}
