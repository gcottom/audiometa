package mp3mp4tag

import (
	"bytes"
	"image"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/bogem/id3v2"
	mp3TagLib "github.com/bogem/id3v2"
	mp4TagReader "github.com/dhowden/tag"
)

type IDTag struct {
	Artist      string
	AlbumArtist string
	Album       string
	AlbumArt    image.Image
	Comments    string
	Composer    string
	Genre       string
	Title       string
	Year        int
	BPM         int
	ID3         ID3Frames
}
type ID3Frames struct {
	ContentType   string //Content Type
	CopyrightMsg  string //Copyright Message
	Date          string //Date
	PlaylistDelay string //Playlist Delay
	EncodedBy     string //Endcoded By
	Lyricist      string //Lyricist
	FileType      string //File Type
	Language      string //Language
	Length        string //Length
	PartOfSet     string //Part of a set
	Publisher     string //Publisher
	TrackNumber   string //Track number ex "2/5"
}

func parse(filepath string) IDTag {
	fileTypeArr := strings.Split(filepath, ".")
	lastIndex := len(fileTypeArr) - 1
	fileType := fileTypeArr[lastIndex]
	var resultTag IDTag
	if fileType == "mp3" {
		//use the mp3TagLib
		tag, err := mp3TagLib.Open(filepath, mp3TagLib.Options{Parse: true})
		defer tag.Close()
		if err != nil {
			log.Fatal("Error while opening mp3 file: ", err)
		}
		resultTag = IDTag{Artist: tag.Artist(), AlbumArtist: tag.Artist(), Album: tag.Album(), Genre: tag.Genre(), Title: tag.Title()}
		resultTag.Year, _ = strconv.Atoi(tag.Year())
		bpmFramer := tag.GetLastFrame(tag.CommonID("BPM"))
		if bpmFramer != nil {
			bpm, ok := bpmFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert bpm frame")
			} else {
				resultTag.BPM, _ = strconv.Atoi(bpm.Text)
			}
		}
		commentFramer := tag.GetLastFrame("COMM")
		if commentFramer != nil {
			comment, ok := commentFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert comment frame")
			} else {
				resultTag.Comments = comment.Text
			}
		}
		composerFramer := tag.GetLastFrame("TCOM")
		if composerFramer != nil {
			composer, ok := composerFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert composer frame")
			} else {
				resultTag.Composer = composer.Text
			}
		}
		exFramer := tag.GetLastFrame("TCON")
		if exFramer != nil {
			ex, ok := exFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert content type frame")
			} else {
				resultTag.ID3.ContentType = ex.Text
			}
		}
		exFramer = tag.GetLastFrame("TCOP")
		if exFramer != nil {
			ex, ok := exFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert copyright frame")
			} else {
				resultTag.ID3.CopyrightMsg = ex.Text
			}
		}
		exFramer = tag.GetLastFrame("TDRC")
		if exFramer != nil {
			ex, ok := exFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert date frame")
			} else {
				resultTag.ID3.Date = ex.Text
			}
		}
		exFramer = tag.GetLastFrame("TDLY")
		if exFramer != nil {
			ex, ok := exFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert playlist delay frame")
			} else {
				resultTag.ID3.PlaylistDelay = ex.Text
			}
		}
		exFramer = tag.GetLastFrame("TENC")
		if exFramer != nil {
			ex, ok := exFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert encoded by frame")
			} else {
				resultTag.ID3.EncodedBy = ex.Text
			}
		}
		exFramer = tag.GetLastFrame("TEXT")
		if exFramer != nil {
			ex, ok := exFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert lyricist frame")
			} else {
				resultTag.ID3.Lyricist = ex.Text
			}
		}
		exFramer = tag.GetLastFrame("TFLT")
		if exFramer != nil {
			ex, ok := exFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert file type frame")
			} else {
				resultTag.ID3.FileType = ex.Text
			}
		}
		exFramer = tag.GetLastFrame("TLAN")
		if exFramer != nil {
			ex, ok := exFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert language frame")
			} else {
				resultTag.ID3.Language = ex.Text
			}
		}
		exFramer = tag.GetLastFrame("TLEN")
		if exFramer != nil {
			ex, ok := exFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert length frame")
			} else {
				resultTag.ID3.Length = ex.Text
			}
		}
		exFramer = tag.GetLastFrame("TPOS")
		if exFramer != nil {
			ex, ok := exFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert part of set frame")
			} else {
				resultTag.ID3.PartOfSet = ex.Text
			}
		}
		exFramer = tag.GetLastFrame("TPUB")
		if exFramer != nil {
			ex, ok := exFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert publisher frame")
			} else {
				resultTag.ID3.Publisher = ex.Text
			}
		}
		exFramer = tag.GetLastFrame("TRCK")
		if exFramer != nil {
			ex, ok := exFramer.(mp3TagLib.TextFrame)
			if !ok {
				log.Fatal("Couldn't assert track pos frame")
			} else {
				resultTag.ID3.TrackNumber = ex.Text
			}
		}
		pictures := tag.GetFrames(tag.CommonID("Attached picture"))
		if len(pictures) > 0 {
			pic := pictures[0].(id3v2.PictureFrame)
			img, _, err := image.Decode(bytes.NewReader(pic.Picture))
			if err != nil {
				log.Fatalln(err)
			}
			resultTag.AlbumArt = img
		}

	} else {
		f, err := os.Open(filepath)
		if err != nil {
			log.Fatal("Error while opening file: ", err)
		}
		defer f.Close()
		tag, err := mp4TagReader.ReadFrom(f)
		if err != nil {
			log.Fatal("Error while reading file: ", err)
		}
		resultTag = IDTag{Artist: tag.Artist(), AlbumArtist: tag.AlbumArtist(), Album: tag.Album(), Comments: tag.Comment(), Composer: tag.Composer(), Genre: tag.Genre(), Title: tag.Title(), Year: tag.Year(), BPM: tag.Tempo()}
		if tag.Picture() != nil {
			albumArt := tag.Picture().Data
			img, _, err := image.Decode(bytes.NewReader(albumArt))
			if err != nil {
				log.Fatal("Error opening album image")
			}
			resultTag.AlbumArt = img
		}
	}
	return resultTag
}
