package audiometa

import (
	"fmt"
	"image"
	"io"

	"github.com/gcottom/audiometa/mp3/v3"
)

type Tag interface {
	Save(io.Writer) error

	GetAlbum() string
	GetAlbumArtist() string
	GetArtist() string
	GetBPM() int
	GetComposer() string
	GetCopyright() string
	GetCoverArt() *image.Image
	GetDiscNumber() int
	GetDiscTotal() int
	GetEncoder() string
	GetGenre() string
	GetTitle() string
	GetTrackNumber() int
	GetTrackTotal() int

	SetAlbum(string)
	SetAlbumArtist(string)
	SetArtist(string)
	SetBPM(int)
	SetComposer(string)
	SetCoverArt(*image.Image)
	SetDiscNumber(int)
	SetDiscTotal(int)
	SetEncoder(string)
	SetGenre(string)
	SetTitle(string)
	SetTrackNumber(int)
	SetTrackTotal(int)
}

func OpenTag(r io.ReadSeeker) (Tag, error) {
	b, err := readBytes(r, 8)
	if err != nil {
		return nil, err
	}

	if _, err = r.Seek(-8, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("error seeking back to original position: %v", err)
	}

	switch {
	case string(b[0:3]) == "ID3":
		return mp3.ParseMP3(r)
	case string(b[4:8]) == "ftyp":
		return mp4.parseMP4(r)
	case string(b[0:4]) == "fLaC":
		return //flac reader
	case string(b[0:4]) == "OggS":
		return //ogg reader
	}

}
