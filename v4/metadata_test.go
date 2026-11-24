package audiometa_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/gcottom/audiometa/v4"
)

func TestNamedAccessorsShareCanonicalState(t *testing.T) {
	source := make([]byte, 44)
	copy(source, "fLaC")
	binary.BigEndian.PutUint32(source[4:], 0x80000022)
	copy(source[42:], []byte{0xff, 0xf8})
	tag, err := audiometa.OpenMetadata(bytes.NewReader(source))
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		field audiometa.Field
		get   func() string
		set   func(string) error
	}{
		{audiometa.Title, tag.GetTitle, tag.SetTitle},
		{audiometa.Subtitle, tag.GetSubtitle, tag.SetSubtitle},
		{audiometa.Album, tag.GetAlbum, tag.SetAlbum},
		{audiometa.Artist, tag.GetArtist, tag.SetArtist},
		{audiometa.AlbumArtist, tag.GetAlbumArtist, tag.SetAlbumArtist},
		{audiometa.Composer, tag.GetComposer, tag.SetComposer},
		{audiometa.Publisher, tag.GetPublisher, tag.SetPublisher},
		{audiometa.Mood, tag.GetMood, tag.SetMood},
		{audiometa.ISRC, tag.GetISRC, tag.SetISRC},
		{audiometa.Genre, tag.GetGenre, tag.SetGenre},
		{audiometa.Copyright, tag.GetCopyright, tag.SetCopyright},
		{audiometa.Encoder, tag.GetEncoder, tag.SetEncoder},
		{audiometa.Grouping, tag.GetGrouping, tag.SetGrouping},
		{audiometa.Lyricist, tag.GetLyricist, tag.SetLyricist},
		{audiometa.Conductor, tag.GetConductor, tag.SetConductor},
		{audiometa.Arranger, tag.GetArranger, tag.SetArranger},
		{audiometa.RecordLabel, tag.GetRecordLabel, tag.SetRecordLabel},
		{audiometa.Work, tag.GetWork, tag.SetWork},
		{audiometa.Movement, tag.GetMovement, tag.SetMovement},
	} {
		want := "Metadata"
		if c.field == audiometa.ISRC {
			want = "USABC2600001"
		}
		if err := c.set(want); err != nil {
			t.Fatal(err)
		}
		if value, ok := tag.GetField(c.field); !ok || !value.Equal(audiometa.TextValue(want)) {
			t.Fatalf("named setter did not set %s", c.field)
		}
		if got := c.get(); got != want {
			t.Errorf("named getter %s = %q", c.field, got)
		}
	}
	for _, c := range []struct {
		field audiometa.Field
		get   func() int64
		set   func(int64) error
	}{
		{audiometa.BPM, tag.GetBPM, tag.SetBPM},
		{audiometa.TrackNumber, tag.GetTrackNumber, tag.SetTrackNumber},
		{audiometa.TrackTotal, tag.GetTrackTotal, tag.SetTrackTotal},
		{audiometa.DiscNumber, tag.GetDiscNumber, tag.SetDiscNumber},
		{audiometa.DiscTotal, tag.GetDiscTotal, tag.SetDiscTotal},
	} {
		const want = 7
		if err := c.set(want); err != nil {
			t.Fatal(err)
		}
		if value, ok := tag.GetField(c.field); !ok || !value.Equal(audiometa.IntValue(want)) {
			t.Fatalf("named setter did not set %s", c.field)
		}
		if got := c.get(); got != want {
			t.Errorf("named getter %s = %d", c.field, got)
		}
	}
}
