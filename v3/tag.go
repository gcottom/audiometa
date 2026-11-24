package audiometa

import (
	"errors"
	"fmt"
	"image"
	"io"

	"github.com/gcottom/flacmeta"
	"github.com/gcottom/mp3meta"
	"github.com/gcottom/mp4meta"
	"github.com/gcottom/oggmeta"
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
	kind, err := detectFormat(r)
	if err != nil {
		return nil, err
	}
	switch kind {
	case mp3File:
		return mp3meta.ParseMP3(r)
	case mp4File:
		return mp4meta.ReadMP4(r)
	case flacFile:
		return flacmeta.ReadFLAC(r)
	case oggFile:
		return oggmeta.ReadOGG(r)
	default:
		return nil, errors.New("unsupported file type")
	}
}

type fileKind uint8

const (
	mp3File fileKind = iota + 1
	mp4File
	flacFile
	oggFile
)

func detectFormat(r io.ReadSeeker) (fileKind, error) {
	if r == nil {
		return 0, errors.New("audiometa: nil source")
	}
	b, err := readBytes(r, 8)
	if err != nil {
		return 0, err
	}

	if _, err = r.Seek(-8, io.SeekCurrent); err != nil {
		return 0, fmt.Errorf("error seeking back to original position: %w", err)
	}

	switch {
	case string(b[0:3]) == "ID3":
		return mp3File, nil
	case string(b[4:8]) == "ftyp":
		return mp4File, nil
	case string(b[0:4]) == "fLaC":
		return flacFile, nil
	case string(b[0:4]) == "OggS":
		return oggFile, nil
	default:
		ok, err := isUntaggedMP3(r, b[:4])
		if err != nil {
			return 0, err
		}
		if ok {
			return mp3File, nil
		}
		return 0, errors.New("unsupported file type")
	}
}
