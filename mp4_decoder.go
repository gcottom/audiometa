package audiometa

import (
	"bytes"
	"encoding/binary"
	"image"
	"io"
	"strconv"
	"strings"
)

var atomTypes = map[int]string{
	0:  "implicit", // automatic based on atom name
	1:  "text",
	13: "png",
	14: "jpg",
	21: "uint8",
}

// NB: atoms does not include "----", this is handled separately
var atoms = atomNames(map[string]string{
	"\xa9alb": "album",
	"\xa9art": "artist",
	"\xa9ART": "artist",
	"aART":    "album_artist",
	"\xa9day": "year",
	"\xa9nam": "title",
	"\xa9gen": "genre",
	"trkn":    "track",
	"\xa9wrt": "composer",
	"\xa9too": "encoder",
	"cprt":    "copyright",
	"covr":    "picture",
	"\xa9grp": "grouping",
	"keyw":    "keyword",
	"\xa9lyr": "lyrics",
	"\xa9cmt": "comment",
	"tmpo":    "tempo",
	"cpil":    "compilation",
	"disk":    "disc",
})

// Detect PNG image if "implicit" class is used
var pngHeader = []byte{137, 80, 78, 71, 13, 10, 26, 10}

type atomNames map[string]string

func (f atomNames) name(n string) []string {
	res := make([]string, 1)
	for k, v := range f {
		if v == n {
			res = append(res, k)
		}
	}
	return res
}

// metadataMP4 is the implementation of Metadata for MP4 tag (atom) data.
type metadataMP4 struct {
	data map[string]interface{}
}

func readFromMP4(r io.ReadSeeker) (metadataMP4, error) {
	return readAtoms(r)
}

// ReadAtoms reads MP4 metadata atoms from the io.ReadSeeker into a Metadata, returning
// non-nil error if there was a problem.
func readAtoms(r io.ReadSeeker) (metadataMP4, error) {
	m := metadataMP4{
		data: make(map[string]interface{}),
	}
	err := m.readAtoms(r)
	return m, err
}

func (m metadataMP4) readAtoms(r io.ReadSeeker) error {
	for {
		name, size, err := readAtomHeader(r)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		switch name {
		case "meta":
			_, err := readBytes(r, 4)
			if err != nil {
				return err
			}
			fallthrough

		case "moov", "udta", "ilst":
			return m.readAtoms(r)
		}

		_, ok := atoms[name]
		var data []string

		if !ok {
			_, err := r.Seek(int64(size-8), io.SeekCurrent)
			if err != nil {
				return err
			}
			continue
		}

		err = m.readAtomData(r, name, size-8, data)
		if err != nil {
			return err
		}
	}
}

func (m metadataMP4) readAtomData(r io.ReadSeeker, name string, size uint32, processedData []string) error {
	var b []byte
	var err error
	var contentType string
	if len(processedData) > 0 {
		b = []byte(strings.Join(processedData, ";")) // add delimiter if multiple data fields
		contentType = "text"
	} else {
		// read the data
		b, err = readBytes(r, uint(size))
		if err != nil {
			return err
		}
		if len(b) < 8 {
			return ErrMP4InvalidEncoding
		}

		// "data" + size (4 bytes each)
		b = b[8:]

		if len(b) < 4 {
			return ErrMP4InvalidEncoding
		}
		class := getInt(b[1:4])
		var ok bool
		contentType, ok = atomTypes[class]
		if !ok {
			return ErrMP4InvalidCntntType
		}

		// 4: atom version (1 byte) + atom flags (3 bytes)
		// 4: NULL (usually locale indicator)
		if len(b) < 8 {
			return ErrMP4InvalidEncoding
		}
		b = b[8:]
	}

	if name == "trkn" || name == "disk" {
		if len(b) < 6 {
			return ErrMP4InvalidEncoding
		}

		m.data[name] = int(b[3])
		m.data[name+"_count"] = int(b[5])
		return nil
	}

	if contentType == "implicit" {
		if name == "covr" {
			contentType = "png"
		}
	}

	var data interface{}
	switch contentType {
	case "implicit":
		if _, ok := atoms[name]; ok {
			return ErrMP4InvalidCntntType
		}
		return nil

	case "text":
		data = string(b)

	case "uint8":
		if len(b) < 1 {
			return ErrMP4InvalidEncoding
		}
		data = getInt(b[:1])

	case "jpeg", "png":
		if img, _, err := image.Decode(bytes.NewReader(b)); err == nil {
			data = &img
		}
	}
	m.data[name] = data

	return nil
}

func readAtomHeader(r io.ReadSeeker) (name string, size uint32, err error) {
	err = binary.Read(r, binary.BigEndian, &size)
	if err != nil {
		return
	}
	name, err = readString(r, 4)
	return
}

func (m metadataMP4) getString(n []string) string {
	for _, k := range n {
		if x, ok := m.data[k]; ok {
			return x.(string)
		}
	}
	return ""
}

func (m metadataMP4) getInt(n []string) int {
	for _, k := range n {
		if x, ok := m.data[k]; ok {
			return x.(int)
		}
	}
	return 0
}

func (m metadataMP4) title() string {
	return m.getString(atoms.name("title"))
}

func (m metadataMP4) artist() string {
	return m.getString(atoms.name("artist"))
}

func (m metadataMP4) album() string {
	return m.getString(atoms.name("album"))
}

func (m metadataMP4) albumArtist() string {
	return m.getString(atoms.name("album_artist"))
}

func (m metadataMP4) composer() string {
	return m.getString(atoms.name("composer"))
}

func (m metadataMP4) genre() string {
	return m.getString(atoms.name("genre"))
}

func (m metadataMP4) year() int {
	date := m.getString(atoms.name("year"))
	if len(date) >= 4 {
		year, _ := strconv.Atoi(date[:4])
		return year
	}
	return 0
}

func (m metadataMP4) comment() string {
	t, ok := m.data["\xa9cmt"]
	if !ok {
		return ""
	}
	return t.(string)
}

func (m metadataMP4) picture() *image.Image {
	v, ok := m.data["covr"]
	if !ok {
		return nil
	}
	return v.(*image.Image)

}

func (m metadataMP4) tempo() int {
	return m.getInt(atoms.name("tempo"))
}

func (m metadataMP4) encoder() string {
	return m.getString(atoms.name("encoder"))
}

func (m metadataMP4) copyright() string {
	return m.getString(atoms.name("copyright"))
}
