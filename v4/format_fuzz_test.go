package audiometa

import (
	"bytes"
	"io"
	"testing"
)

func FuzzDetectFormat(f *testing.F) {
	for _, seed := range []string{"", "fLaC\x80\x00\x00\x22", "ID3\x04\x00\x00\x00\x00", "\x00\x00\x01\x00ftyp", "OggS\x00\x02\x00\x00", "\xff\xfb\x90\x00\x00\x00\x00\x00"} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		reader := bytes.NewReader(data)
		kind, err := detectFormat(reader)
		if err != nil {
			return
		}
		if kind < mp3File || kind > oggFile {
			t.Fatalf("unexpected file kind %d", kind)
		}
		if pos, err := reader.Seek(0, io.SeekCurrent); err != nil || pos != 0 {
			t.Fatalf("successful detection consumed input: %d / %v", pos, err)
		}
	})
}
