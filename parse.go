package audiometa

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"os"
	"strconv"
	"strings"

	mp3TagLib "github.com/bogem/id3v2/v2"
	"github.com/sunfish-shogi/bufseekio"
)

// This operation opens the ID tag for the corresponding file that is passed in the filepath parameter regardless of the filetype as long as it is a supported file type

// parse will assume that the file is in memory
func parse(input io.Reader, opts ParseOptions) (*IDTag, error) {
	format := opts.Format
	b := new(bytes.Buffer)
	b.ReadFrom(input)
	data := b.Bytes()
	switch {
	case format == MP3:
		tag, err := parseMP3(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		tag.data = data
		return tag, nil
	case fileTypesContains(format, mp4FileTypes):
		tag, err := parseMP4(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		tag.data = data
		return tag, nil
	case format == FLAC:
		tag, err := parseFLAC(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		tag.data = data
		return tag, nil
	case format == OGG:
		tag, err := parseOGG(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		tag.data = data
		return tag, nil
	}
	return nil, fmt.Errorf("no method available for filetype:%s", format)
}

// parseFile assumes that the file is to be read from the path
func parseFile(filepath string) (*IDTag, error) {
	fileType, err := GetFileType(filepath)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	switch {
	case fileType == MP3:
		tag, err := parseMP3(file)
		if err != nil {
			return nil, err
		}
		tag.filePath = filepath
		return tag, nil
	case fileTypesContains(fileType, mp4FileTypes):
		tag, err := parseMP4(file)
		if err != nil {
			return nil, err
		}
		tag.filePath = filepath
		return tag, nil
	case fileType == FLAC:
		tag, err := parseFLAC(file)
		if err != nil {
			return nil, err
		}
		tag.filePath = filepath
		return tag, nil

	case fileType == OGG:
		tag, err := parseOGG(file)
		if err != nil {
			return nil, err
		}
		tag.filePath = filepath
		return tag, nil
	}
	return nil, fmt.Errorf("no method available for filetype:%s", fileType)
}
func parseMP3(input io.Reader) (*IDTag, error) {
	resultTag := IDTag{}
	tag, err := mp3TagLib.ParseReader(input, mp3TagLib.Options{Parse: true})
	if err != nil {
		return nil, fmt.Errorf("error opening mp3 [%w]", err)
	}
	defer tag.Close()
	resultTag = IDTag{artist: tag.Artist(), album: tag.Album(), genre: tag.Genre(), title: tag.Title(), year: tag.Year()}
	if bpmFramer := tag.GetLastFrame(tag.CommonID("BPM")); bpmFramer != nil {
		if bpm, ok := bpmFramer.(mp3TagLib.TextFrame); ok {
			resultTag.bpm = bpm.Text
		}
	}
	if commentFramer := tag.GetLastFrame("COMM"); commentFramer != nil {
		if comment, ok := commentFramer.(mp3TagLib.CommentFrame); ok {
			resultTag.comments = comment.Text
		}
	}
	if composerFramer := tag.GetLastFrame("TCOM"); composerFramer != nil {
		if composer, ok := composerFramer.(mp3TagLib.TextFrame); ok {
			resultTag.composer = composer.Text
		}
	}
	if exFramer := tag.GetLastFrame("TPE2"); exFramer != nil {
		if ex, ok := exFramer.(mp3TagLib.TextFrame); ok {
			resultTag.albumArtist = ex.Text
		}
	}
	if exFramer := tag.GetLastFrame("TCOP"); exFramer != nil {
		if ex, ok := exFramer.(mp3TagLib.TextFrame); ok {
			resultTag.copyrightMsg = ex.Text
		}
	}
	if exFramer := tag.GetLastFrame("TDRC"); exFramer != nil {
		if ex, ok := exFramer.(mp3TagLib.TextFrame); ok {
			resultTag.date = ex.Text
		}
	}
	if exFramer := tag.GetLastFrame("TENC"); exFramer != nil {
		if ex, ok := exFramer.(mp3TagLib.TextFrame); ok {
			resultTag.encodedBy = ex.Text
		}
	}
	if exFramer := tag.GetLastFrame("TEXT"); exFramer != nil {
		if ex, ok := exFramer.(mp3TagLib.TextFrame); ok {
			resultTag.lyricist = ex.Text
		}
	}
	if exFramer := tag.GetLastFrame("TFLT"); exFramer != nil {
		if ex, ok := exFramer.(mp3TagLib.TextFrame); ok {
			resultTag.fileType = ex.Text
		}
	}
	if exFramer := tag.GetLastFrame("TLAN"); exFramer != nil {
		if ex, ok := exFramer.(mp3TagLib.TextFrame); ok {
			resultTag.language = ex.Text
		}
	}
	if exFramer := tag.GetLastFrame("TLEN"); exFramer != nil {
		if ex, ok := exFramer.(mp3TagLib.TextFrame); ok {
			resultTag.length = ex.Text
		}
	}
	if exFramer := tag.GetLastFrame("TPOS"); exFramer != nil {
		if ex, ok := exFramer.(mp3TagLib.TextFrame); ok {
			resultTag.partOfSet = ex.Text
		}
	}
	if exFramer := tag.GetLastFrame("TPUB"); exFramer != nil {
		if ex, ok := exFramer.(mp3TagLib.TextFrame); ok {
			resultTag.publisher = ex.Text
		}
	}
	if pictures := tag.GetFrames(tag.CommonID("Attached picture")); len(pictures) > 0 {
		pic := pictures[0].(mp3TagLib.PictureFrame)
		if img, _, err := image.Decode(bytes.NewReader(pic.Picture)); err == nil {
			resultTag.albumArt = &img
		}
	}
	return &resultTag, nil
}

func parseFLAC(input io.Reader) (*IDTag, error) {
	resultTag := IDTag{}
	fb, err := extractFLACComment(input)
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
