package audiometa

import (
	"image"
	"io"
)

// The IDTag represents all of the metadata that can be retrieved from a file.
// The IDTag contains all tags for all audio types. Some tags may not be applicable to all types.
// Only the valid types are written to the respective data files.
// Although a tag may be set, if the function to write that tag attribute doesn't exist, the tag attribute will be ignored and the save function will not produce an error.
type IDTag struct {
	artist       string       //Artist
	albumArtist  string       //AlbumArtist
	album        string       //Album
	albumArt     *image.Image //AlbumArt for the work in image format
	comments     string       //Comments
	composer     string       //Composer
	genre        string       //Genre
	title        string       //Title
	year         string       //Year
	bpm          string       //BPM
	filePath     string       //The filepath of the file
	codec        string       //The codec of the file (ogg use only)
	copyrightMsg string       //Copyright Message
	date         string       //Date
	encodedBy    string       //Endcoded By
	lyricist     string       //Lyricist
	fileType     string       //File Type
	language     string       //Language
	length       string       //Length
	partOfSet    string       //Part of a set
	publisher    string       //Publisher

	reader io.ReadSeeker

	PassThrough map[string]string
}

// ParseOptions is a struct that is passed when parsing a tag. If included, you should set the format to one of the existing FileTypes
type ParseOptions struct {
	Format FileType
}

type MP4Delete []string

type FileType string

const (
	MP3  FileType = "mp3"
	M4P  FileType = "m4p"
	M4A  FileType = "m4a"
	M4B  FileType = "m4b"
	MP4  FileType = "mp4"
	FLAC FileType = "flac"
	OGG  FileType = "ogg"
)

const (
	ALBUM  string = "album"
	ARTIST string = "artist"
	DATE   string = "date"
	TITLE  string = "title"
	GENRE  string = "genre"
)

var mp3TextFrames = map[string]string{
	"artist":       "TPE1",
	"title":        "TIT2",
	"album":        "TALB",
	"comments":     "COMM",
	"bpm":          "TBPM",
	"genre":        "TCON",
	"year":         "TYER",
	"albumArtist":  "TPE2",
	"composer":     "TCOM",
	"copyrightMsg": "TCOP",
	"date":         "TDRC",
	"encodedBy":    "TENC",
	"lyricist":     "TEXT",
	"fileType":     "TFLT",
	"language":     "TLAN",
	"length":       "TLEN",
	"partOfSet":    "TPOS",
	"publisher":    "TPUB",
	"albumArt":     "APIC",
}

var supportedFileTypes = []FileType{MP3, M4P, M4A, M4B, MP4, FLAC, OGG}
var mp4FileTypes = []FileType{M4P, M4A, M4B, MP4}
