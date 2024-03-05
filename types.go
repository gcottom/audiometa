package audiometa

import (
	"image"
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

	PassThrough map[string]string
}

// The Picture type contains a byte representation of an image
type Picture struct {
	Ext         string // Extension of the picture file.
	MIMEType    string // MIMEType of the picture.
	Type        string // Type of the picture (see pictureTypes).
	Description string // Description.
	Data        []byte // Raw picture data.
}

const (
	MP3  string = "mp3"
	M4P  string = "m4p"
	M4A  string = "m4a"
	M4B  string = "m4b"
	MP4  string = "mp4"
	FLAC string = "flac"
	OGG  string = "ogg"
)

const (
	ALBUM  string = "album"
	ARTIST string = "artist"
	DATE   string = "date"
	TITLE  string = "title"
	GENRE  string = "genre"
)

var supportedFileTypes = []string{MP3, M4P, M4A, M4B, MP4, FLAC, OGG}

var vorbisPictureTypes = map[byte]string{
	0x00: "Other",
	0x01: "32x32 pixels 'file icon' (PNG only)",
	0x02: "Other file icon",
	0x03: "Cover (front)",
	0x04: "Cover (back)",
	0x05: "Leaflet page",
	0x06: "Media (e.g. lable side of CD)",
	0x07: "Lead artist/lead performer/soloist",
	0x08: "Artist/performer",
	0x09: "Conductor",
	0x0A: "Band/Orchestra",
	0x0B: "Composer",
	0x0C: "Lyricist/text writer",
	0x0D: "Recording Location",
	0x0E: "During recording",
	0x0F: "During performance",
	0x10: "Movie/video screen capture",
	0x11: "A bright coloured fish",
	0x12: "Illustration",
	0x13: "Band/artist logotype",
	0x14: "Publisher/Studio logotype",
}

type Tags struct {
	Album       string
	AlbumArtist string
	Artist      string
	Comment     string
	Composer    string
	Copyright   string
	Cover       []byte
	Custom      map[string]string
	Delete      []string
	DiskNumber  int
	DiskTotal   int
	Genre       string
	Label       string
	Title       string
	TrackNumber int
	TrackTotal  int
	Year        string
}
