package mp3mp4tag

import (
	"bytes"
	"image"
	"log"
	"os"
	"strconv"
	"strings"

	mp3TagLib "github.com/bogem/id3v2"
	"github.com/go-flac/flacpicture"
	flac "github.com/go-flac/go-flac"
)

// This operation opens the ID tag for the corresponding file that is passed in the filepath parameter regardless of the filetype as long as it is a supported file type
func parse(filepath string) (*IDTag, error) {
	fileType, err := getFileType(filepath)
	if err != nil {
		return nil, err
	}
	var resultTag IDTag
	if *fileType == "mp3" {
		//use the mp3TagLib
		tag, err := mp3TagLib.Open(filepath, mp3TagLib.Options{Parse: true})
		if err != nil {
			log.Fatal("Error while opening mp3 file: ", err)
			return nil, err
		}
		defer tag.Close()
		resultTag = IDTag{artist: tag.Artist(), album: tag.Album(), genre: tag.Genre(), title: tag.Title()}
		resultTag.year = tag.Year()
		bpmFramer := tag.GetLastFrame(tag.CommonID("BPM"))
		if bpmFramer != nil {
			bpm, ok := bpmFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert bpm frame")
			} else {
				resultTag.bpm = bpm.Text
			}
		}
		commentFramer := tag.GetLastFrame("COMM")
		if commentFramer != nil {
			comment, ok := commentFramer.(mp3TagLib.CommentFrame)
			if !ok {
				log.Fatal("Couldn't assert comment frame")
			} else {
				resultTag.comments = comment.Text
			}
		}
		composerFramer := tag.GetLastFrame("TCOM")
		if composerFramer != nil {
			composer, ok := composerFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert composer frame")
			} else {
				resultTag.composer = composer.Text
			}
		}
		exFramer := tag.GetLastFrame("TPE2")
		if exFramer != nil {
			ex, ok := exFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert album artist frame")
			} else {
				resultTag.albumArtist = ex.Text
			}
		}
		exFramer = tag.GetLastFrame("TCOP")
		if exFramer != nil {
			ex, ok := exFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert copyright frame")
			} else {
				resultTag.id3.copyrightMsg = ex.Text
			}
		}
		exFramer = tag.GetLastFrame("TDRC")
		if exFramer != nil {
			ex, ok := exFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert date frame")
			} else {
				resultTag.id3.date = ex.Text
			}
		}
		exFramer = tag.GetLastFrame("TENC")
		if exFramer != nil {
			ex, ok := exFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert encoded by frame")
			} else {
				resultTag.id3.encodedBy = ex.Text
			}
		}
		exFramer = tag.GetLastFrame("TEXT")
		if exFramer != nil {
			ex, ok := exFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert lyricist frame")
			} else {
				resultTag.id3.lyricist = ex.Text
			}
		}
		exFramer = tag.GetLastFrame("TFLT")
		if exFramer != nil {
			ex, ok := exFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert file type frame")
			} else {
				resultTag.id3.fileType = ex.Text
			}
		}
		exFramer = tag.GetLastFrame("TLAN")
		if exFramer != nil {
			ex, ok := exFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert language frame")
			} else {
				resultTag.id3.language = ex.Text
			}
		}
		exFramer = tag.GetLastFrame("TLEN")
		if exFramer != nil {
			ex, ok := exFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert length frame")
			} else {
				resultTag.id3.length = ex.Text
			}
		}
		exFramer = tag.GetLastFrame("TPOS")
		if exFramer != nil {
			ex, ok := exFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert part of set frame")
			} else {
				resultTag.id3.partOfSet = ex.Text
			}
		}
		exFramer = tag.GetLastFrame("TPUB")
		if exFramer != nil {
			ex, ok := exFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert publisher frame")
			} else {
				resultTag.id3.publisher = ex.Text
			}
		}
		pictures := tag.GetFrames(tag.CommonID("Attached picture"))
		if len(pictures) > 0 {
			pic := pictures[0].(mp3TagLib.PictureFrame)
			img, _, err := image.Decode(bytes.NewReader(pic.Picture))
			if err != nil {
				resultTag.albumArt = nil
			}
			resultTag.albumArt = &img
		} else {
			resultTag.albumArt = nil
		}

	} else if *fileType == "flac" {
		cmts, _ := extractFLACComment(filepath)
		resultTag = IDTag{}
		if cmts != nil {
			for _, cmt := range cmts.Comments {
				log.Println(cmt)
				if strings.HasPrefix(cmt, "album=") {
					tag := strings.Replace(cmt, "album=", "", 1)
					resultTag.album = tag
				} else if strings.HasPrefix(cmt, "ALBUM=") {
					tag := strings.Replace(cmt, "ALBUM=", "", 1)
					resultTag.album = tag
				} else if strings.HasPrefix(cmt, "artist=") {
					tag := strings.Replace(cmt, "artist=", "", 1)
					resultTag.artist = tag
				} else if strings.HasPrefix(cmt, "ARTIST=") {
					tag := strings.Replace(cmt, "ARTIST=", "", 1)
					resultTag.artist = tag
				} else if strings.HasPrefix(cmt, "date=") {
					tag := strings.Replace(cmt, "date=", "", 1)
					resultTag.id3.date = tag
				} else if strings.HasPrefix(cmt, "DATE=") {
					tag := strings.Replace(cmt, "DATE=", "", 1)
					resultTag.id3.date = tag
				} else if strings.HasPrefix(cmt, "title=") {
					tag := strings.Replace(cmt, "title=", "", 1)
					resultTag.title = tag
				} else if strings.HasPrefix(cmt, "TITLE=") {
					tag := strings.Replace(cmt, "TITLE=", "", 1)
					resultTag.title = tag
				} else if strings.HasPrefix(cmt, "genre=") {
					tag := strings.Replace(cmt, "genre=", "", 1)
					resultTag.genre = tag
				} else if strings.HasPrefix(cmt, "GENRE=") {
					tag := strings.Replace(cmt, "GENRE=", "", 1)
					resultTag.genre = tag
				}
			}
		}
		log.Println("DONE PARSING FILE")
		file, err := os.Open(filepath)
		if err != nil {
			log.Println("Error while opening file: ", err)
			return nil, err
		}
		f, err := flac.ParseBytes(file)
		if err != nil {
			log.Fatal("Error opening album image")
		}
		var pic *flacpicture.MetadataBlockPicture
		for _, meta := range f.Meta {
			if meta.Type == flac.Picture {
				pic, err = flacpicture.ParseFromMetaDataBlock(*meta)
				if err != nil {
					log.Println("Error opening album art")
				}
			}
		}
		if pic != nil {
			img, _, err := image.Decode(bytes.NewReader(pic.ImageData))
			if err != nil {
				log.Fatal("Error opening album image")
			}
			resultTag.albumArt = &img
		} else {
			resultTag.albumArt = nil
		}
	} else {
		f, err := os.Open(filepath)
		if err != nil {
			log.Fatal("Error while opening file: ", err)
			return nil, err
		}
		defer f.Close()
		tag, err := ReadFromMP4(f)
		if err != nil {
			log.Fatal("Error while reading file: ", err)
			return nil, err
		}
		resultTag = IDTag{artist: tag.Artist(), albumArtist: tag.AlbumArtist(), album: tag.Album(), comments: tag.Comment(), composer: tag.Composer(), genre: tag.Genre(), title: tag.Title(), year: strconv.Itoa(tag.Year())}
		resultTag.id3.encodedBy = tag.Encoder()
		resultTag.id3.copyrightMsg = tag.Copyright()
		if tag.Picture() != nil {
			albumArt := tag.Picture()
			img, _, err := image.Decode(bytes.NewReader(albumArt))
			if err != nil {
				log.Fatal("Error opening album image")
			}
			resultTag.albumArt = &img
		} else {
			resultTag.albumArt = nil
		}
	}
	resultTag.fileUrl = filepath
	return &resultTag, nil
}
