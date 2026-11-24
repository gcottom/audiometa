package audiometa

import (
	"bytes"
	"testing"

	"github.com/gcottom/mp3meta"
)

func TestOpenUntaggedMP3(t *testing.T) {
	for _, tc := range []struct {
		name          string
		first, second []byte
		frameBytes    int
	}{
		{"MPEG1_48kHz", []byte{0xff, 0xfb, 0x94, 0}, []byte{0xff, 0xfb, 0xa4, 0}, 384},
		{"MPEG2_24kHz", []byte{0xff, 0xf3, 0x84, 0}, []byte{0xff, 0xf3, 0x94, 0}, 192},
		{"MPEG25_12kHz", []byte{0xff, 0xe3, 0x85, 0}, []byte{0xff, 0xe3, 0x85, 0}, 384},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Header-only synthetic payloads keep this format-dispatch test independent
			// of encoders. The end-to-end media check uses real audio.
			data := make([]byte, tc.frameBytes*3)
			copy(data, tc.first)
			copy(data[tc.frameBytes:], tc.second)
			tag, err := OpenTag(bytes.NewReader(data))
			if err != nil {
				t.Fatal(err)
			}
			if _, ok := tag.(*mp3meta.MP3Tag); !ok {
				t.Fatalf("got %T", tag)
			}
			if tag.GetTitle() != "" {
				t.Fatal("unexpected metadata in a tagless input")
			}
		})
	}
}
func TestRejectNonMP3Sync(t *testing.T) {
	for _, tc := range []struct {
		name          string
		first, second []byte
	}{
		{"isolated-sync", []byte{0xff, 0xfb, 0x94, 0}, nil},
		{"mismatched-sample-rate", []byte{0xff, 0xfb, 0x94, 0}, []byte{0xff, 0xfb, 0x90, 0}},
		{"aac-adts", []byte{0xff, 0xf1, 0x50, 0x80}, []byte{0xff, 0xf1, 0x50, 0x80}},
		{"reserved-layer-version", []byte{0xff, 0xeb, 0x94, 0}, []byte{0xff, 0xeb, 0x94, 0}},
		{"reserved-bitrate", []byte{0xff, 0xfb, 0xf4, 0}, []byte{0xff, 0xfb, 0xf4, 0}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data := make([]byte, 768)
			copy(data, tc.first)
			copy(data[384:], tc.second)
			if tag, err := OpenTag(bytes.NewReader(data)); err == nil {
				t.Fatalf("accepted %T as MP3", tag)
			}
		})
	}
}
