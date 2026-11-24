package audiometa

//go:generate go run ./internal/cmd/genmetadata

import (
	"fmt"
	"io"

	"github.com/gcottom/audiometa-core/v2"
	"github.com/gcottom/flacmeta/v2"
	"github.com/gcottom/mp3meta/v2"
	"github.com/gcottom/mp4meta/v2"
	"github.com/gcottom/oggmeta/v2"
)

// Metadata exposes typed metadata across all supported audio formats. It owns
// its edits, but the caller owns the reader and writer. Concurrent use of the
// same editor or source requires external synchronization.
type Metadata struct{ *core.Editor }

// OpenMetadata detects the format and opens an editor for canonical and custom
// fields. The reader must stay open and unmodified until the editor is unused.
func OpenMetadata(r io.ReadSeeker) (*Metadata, error) {
	kind, err := detectFormat(r)
	if err != nil {
		return nil, err
	}
	var editor *core.Editor
	switch kind {
	case mp3File:
		editor, err = mp3meta.OpenMetadata(r)
	case mp4File:
		editor, err = mp4meta.OpenMetadata(r)
	case flacFile:
		editor, err = flacmeta.OpenMetadata(r)
	case oggFile:
		editor, err = oggmeta.OpenMetadata(r)
	}
	if err != nil {
		return nil, fmt.Errorf("audiometa: open metadata: %w", err)
	}
	return &Metadata{Editor: editor}, nil
}

type (
	Field        = core.Field
	Value        = core.Value
	Kind         = core.Kind
	Date         = core.Date
	Picture      = core.Picture
	Credit       = core.Credit
	Chapter      = core.Chapter
	LyricCue     = core.LyricCue
	SyncedLyrics = core.SyncedLyrics
	TimedEvent   = core.TimedEvent
	Definition   = core.Definition
)

func Fields() []Definition                   { return core.Fields() }
func Lookup(field Field) (Definition, bool)  { return core.Lookup(field) }
func ParseDate(s string) (Date, error)       { return core.ParseDate(s) }
func TextValue(s string) Value               { return core.TextValue(s) }
func IntValue(n int64) Value                 { return core.IntValue(n) }
func BoolValue(b bool) Value                 { return core.BoolValue(b) }
func DecimalValue(n float64) Value           { return core.DecimalValue(n) }
func DateValue(d Date) Value                 { return core.DateValue(d) }
func URLValue(s string) Value                { return core.URLValue(s) }
func UUIDValue(s string) Value               { return core.UUIDValue(s) }
func BinaryValue(b []byte) Value             { return core.BinaryValue(b) }
func PictureValue(p Picture) Value           { return core.PictureValue(p) }
func CreditValue(c Credit) Value             { return core.CreditValue(c) }
func ChapterValue(c Chapter) Value           { return core.ChapterValue(c) }
func SyncedLyricsValue(l SyncedLyrics) Value { return core.SyncedLyricsValue(l) }
func TimedEventValue(e TimedEvent) Value     { return core.TimedEventValue(e) }

var (
	ErrInvalidField     = core.ErrInvalidField
	ErrInvalidValue     = core.ErrInvalidValue
	ErrReadOnly         = core.ErrReadOnly
	ErrUnsupported      = core.ErrUnsupported
	ErrInvalidCustomKey = core.ErrInvalidCustomKey
)

func (m *Metadata) text(field Field) string {
	value, _ := m.GetField(field)
	text, _ := value.Text()
	return text
}

func (m *Metadata) integer(field Field) int64 {
	value, _ := m.GetField(field)
	n, _ := value.Int()
	return n
}
