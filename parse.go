package audiometa

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"strconv"
	"strings"

	mp3TagLib "github.com/bogem/id3v2"
	"github.com/go-flac/flacpicture"
	flac "github.com/go-flac/go-flac"
)

// This operation opens the ID tag for the corresponding file that is passed in the filepath parameter regardless of the filetype as long as it is a supported file type
func parse(filepath string) (*IDTag, error) {
	fileType, err := GetFileType(filepath)
	if err != nil {
		return nil, err
	}
	if fileType == MP3 {
		return parseMP3(filepath)
	} else if fileType == FLAC {
		return parseFLAC(filepath)
	} else if fileType == OGG {
		return parseOGG(filepath)
	} else if fileType == M4P || fileType == M4A || fileType == M4B || fileType == MP4 {
		return parseMP4(filepath)
	}
	return nil, fmt.Errorf("no method available for filetype:%s", fileType)
}

func parseMP3(filepath string) (*IDTag, error) {
	resultTag := IDTag{}
	tag, err := mp3TagLib.Open(filepath, mp3TagLib.Options{Parse: true})
	if err != nil {
		return nil, fmt.Errorf("error opening mp3 [%w]", err)
	}
	defer tag.Close()
	resultTag = IDTag{Artist: tag.Artist(), Album: tag.Album(), Genre: tag.Genre(), Title: tag.Title(), Year: tag.Year()}
	if bpmFramer := tag.GetLastFrame(tag.CommonID("BPM")); bpmFramer != nil {
		if bpm, ok := bpmFramer.(mp3TagLib.TextFrame); ok {
			resultTag.BPM = bpm.Text
		}
	}
	if commentFramer := tag.GetLastFrame("COMM"); commentFramer != nil {
		if comment, ok := commentFramer.(mp3TagLib.CommentFrame); ok {
			resultTag.Comments = comment.Text
		}
	}
	if composerFramer := tag.GetLastFrame("TCOM"); composerFramer != nil {
		if composer, ok := composerFramer.(mp3TagLib.TextFrame); ok {
			resultTag.Composer = composer.Text
		}
	}
	if exFramer := tag.GetLastFrame("TPE2"); exFramer != nil {
		if ex, ok := exFramer.(mp3TagLib.TextFrame); ok {
			resultTag.AlbumArtist = ex.Text
		}
	}
	if exFramer := tag.GetLastFrame("TCOP"); exFramer != nil {
		if ex, ok := exFramer.(mp3TagLib.TextFrame); ok {
			resultTag.CopyrightMsg = ex.Text
		}
	}
	if exFramer := tag.GetLastFrame("TDRC"); exFramer != nil {
		if ex, ok := exFramer.(mp3TagLib.TextFrame); ok {
			resultTag.Date = ex.Text
		}
	}
	if exFramer := tag.GetLastFrame("TENC"); exFramer != nil {
		if ex, ok := exFramer.(mp3TagLib.TextFrame); ok {
			resultTag.EncodedBy = ex.Text
		}
	}
	if exFramer := tag.GetLastFrame("TEXT"); exFramer != nil {
		if ex, ok := exFramer.(mp3TagLib.TextFrame); ok {
			resultTag.Lyricist = ex.Text
		}
	}
	if exFramer := tag.GetLastFrame("TFLT"); exFramer != nil {
		if ex, ok := exFramer.(mp3TagLib.TextFrame); ok {
			resultTag.FileType = ex.Text
		}
	}
	if exFramer := tag.GetLastFrame("TLAN"); exFramer != nil {
		if ex, ok := exFramer.(mp3TagLib.TextFrame); ok {
			resultTag.Language = ex.Text
		}
	}
	if exFramer := tag.GetLastFrame("TLEN"); exFramer != nil {
		if ex, ok := exFramer.(mp3TagLib.TextFrame); ok {
			resultTag.Length = ex.Text
		}
	}
	if exFramer := tag.GetLastFrame("TPOS"); exFramer != nil {
		if ex, ok := exFramer.(mp3TagLib.TextFrame); ok {
			resultTag.PartOfSet = ex.Text
		}
	}
	if exFramer := tag.GetLastFrame("TPUB"); exFramer != nil {
		if ex, ok := exFramer.(mp3TagLib.TextFrame); ok {
			resultTag.Publisher = ex.Text
		}
	}
	if pictures := tag.GetFrames(tag.CommonID("Attached picture")); len(pictures) > 0 {
		pic := pictures[0].(mp3TagLib.PictureFrame)
		if img, _, err := image.Decode(bytes.NewReader(pic.Picture)); err == nil {
			resultTag.AlbumArt = &img
		}
	}
	resultTag.FilePath = filepath
	return &resultTag, nil
}

func parseFLAC(filepath string) (*IDTag, error) {
	resultTag := IDTag{FilePath: filepath}
	if cmts, _, err := extractFLACComment(filepath); cmts != nil && err == nil {
		for _, cmt := range cmts.Comments {
			if sp := strings.Split(cmt, "="); len(sp) == 2 {
				flactag := strings.ToLower(sp[0])
				if flactag == ALBUM {
					resultTag.Album = sp[1]
				} else if flactag == ARTIST {
					resultTag.Artist = sp[1]
				} else if flactag == DATE {
					resultTag.Date = sp[1]
				} else if flactag == TITLE {
					resultTag.Title = sp[1]
				} else if flactag == GENRE {
					resultTag.Genre = sp[1]
				}
			}
		}
	} else if err != nil {
		return nil, err
	}
	file, err := os.Open(filepath)
	if err != nil {
		return &resultTag, nil
	}
	f, err := flac.ParseBytes(file)
	if err != nil {
		return &resultTag, nil
	}
	var pic *flacpicture.MetadataBlockPicture
	for _, meta := range f.Meta {
		if meta.Type == flac.Picture {
			if pic, err = flacpicture.ParseFromMetaDataBlock(*meta); err == nil {
				break
			}
		}
	}
	if pic != nil {
		if img, _, err := image.Decode(bytes.NewReader(pic.ImageData)); err == nil {
			resultTag.AlbumArt = &img
		}
	}
	return &resultTag, nil
}

func parseOGG(filepath string) (*IDTag, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tag, err := readOggTags(f)
	if err != nil {
		return nil, err
	}
	tag.FilePath = filepath
	return tag, nil
}

func parseMP4(filepath string) (*IDTag, error) {
	resultTag := IDTag{}
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tag, err := readFromMP4(f)
	if err != nil {
		return nil, err
	}
	resultTag = IDTag{Artist: tag.artist(), AlbumArtist: tag.albumArtist(), Album: tag.album(),
		Comments: tag.comment(), Composer: tag.composer(), Genre: tag.genre(),
		Title: tag.title(), Year: strconv.Itoa(tag.year()), EncodedBy: tag.encoder(),
		CopyrightMsg: tag.copyright(), BPM: strconv.Itoa(tag.tempo()), FilePath: filepath}
	if tag.picture() != nil {
		if img, _, err := image.Decode(bytes.NewReader(tag.picture())); err == nil {
			resultTag.AlbumArt = &img
		}
	}
	return &resultTag, nil
}
