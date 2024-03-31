package audiometa

import (
	"bytes"
	"image"
	"io"
	"reflect"
	"strconv"
	"strings"

	mp3TagLib "github.com/bogem/id3v2/v2"
	"github.com/sunfish-shogi/bufseekio"
)

// This operation opens the ID tag for the corresponding file that is passed in the filepath parameter regardless of the filetype as long as it is a supported file type
func parse(input io.ReadSeeker, opts ParseOptions) (*IDTag, error) {
	if _, err := input.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	format := opts.Format
	switch {
	case format == MP3:
		tag, err := parseMP3(input)
		if err != nil {
			return nil, err
		}
		tag.fileType = "mp3"
		tag.reader = input
		return tag, nil
	case fileTypesContains(format, mp4FileTypes):
		tag, err := parseMP4(input)
		if err != nil {
			return nil, err
		}
		tag.reader = input
		tag.fileType = "mp4"
		return tag, nil
	case format == FLAC:
		tag, err := parseFLAC(input)
		if err != nil {
			return nil, err
		}
		tag.reader = input
		tag.fileType = "flac"
		return tag, nil
	case format == OGG:
		tag, err := parseOGG(input)
		if err != nil {
			return nil, err
		}
		tag.fileType = "ogg"
		tag.reader = input
		return tag, nil
	}
	return nil, ErrNoMethodAvlble
}

func parseMP3(input io.Reader) (*IDTag, error) {
	resultTag := IDTag{}
	tag, err := mp3TagLib.ParseReader(input, mp3TagLib.Options{Parse: true})
	if err != nil {
		return nil, ErrMP3ParseFail
	}
	resultTag = IDTag{artist: tag.Artist(), album: tag.Album(), genre: tag.Genre(), title: tag.Title(), year: tag.Year()}
	rtPtr := reflect.ValueOf(&resultTag)
	for k, v := range mp3TextFrames {
		field := k
		if k == "albumArt" {
			continue
		}
		if k == "bpm" {
			field = "BPM"
		}
		if framer := tag.GetLastFrame(v); framer != nil {
			if t, ok := framer.(mp3TagLib.TextFrame); ok {
				if t.Text == "" {
					continue
				}
				rtPtr.MethodByName("Set" + strings.ToUpper(field[:1]) + field[1:]).Call([]reflect.Value{reflect.ValueOf(t.Text)})
			}
		}
	}
	if pictures := tag.GetFrames("APIC"); len(pictures) > 0 {
		pic := pictures[0].(mp3TagLib.PictureFrame)
		if img, _, err := image.Decode(bytes.NewReader(pic.Picture)); err == nil {
			resultTag.albumArt = &img
		}
	}
	return &resultTag, nil
}

func parseFLAC(input io.Reader) (*IDTag, error) {
	resultTag := IDTag{}
	_, fb, err := extractFLACComment(input)
	if err != nil {
		return nil, err
	}
	if fb.cmts != nil {
		for _, cmt := range fb.cmts.Comments {
			if sp := strings.Split(cmt, "="); len(sp) == 2 {
				flactag := strings.ToLower(sp[0])
				if flactag == ALBUM {
					resultTag.album = sp[1]
				} else if flactag == ARTIST {
					resultTag.artist = sp[1]
				} else if flactag == DATE {
					resultTag.date = sp[1]
				} else if flactag == TITLE {
					resultTag.title = sp[1]
				} else if flactag == GENRE {
					resultTag.genre = sp[1]
				}
			}
		}
	}
	resultTag.albumArt = fb.pic
	return &resultTag, nil
}

func parseOGG(input io.Reader) (*IDTag, error) {
	return readOggTags(input)
}

func parseMP4(input io.ReadSeeker) (*IDTag, error) {
	resultTag := IDTag{}
	r := bufseekio.NewReadSeeker(input, 128*1024, 4)
	tag, err := readFromMP4(r)
	if err != nil {
		return nil, err
	}
	resultTag = IDTag{artist: tag.artist(), albumArtist: tag.albumArtist(), album: tag.album(),
		albumArt: tag.picture(), comments: tag.comment(), composer: tag.composer(), genre: tag.genre(),
		title: tag.title(), year: strconv.Itoa(tag.year()), encodedBy: tag.encoder(),
		copyrightMsg: tag.copyright(), bpm: strconv.Itoa(tag.tempo())}

	return &resultTag, nil
}
