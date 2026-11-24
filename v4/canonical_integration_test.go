//go:build integration
// +build integration

package audiometa_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gcottom/audiometa-core/v2"
	"github.com/gcottom/audiometa-core/v2/contracttest"
	"github.com/gcottom/audiometa/v4"
)

type integrationFormat struct {
	name      string
	extension string
	codec     string
	args      []string
}

var canonicalFormats = []integrationFormat{
	{"ID3v2.4", "mp3", "libmp3lame", []string{"-id3v2_version", "4"}},
	{"ID3v2.3", "mp3", "libmp3lame", []string{"-id3v2_version", "3"}},
	{"MP4", "mp4", "aac", []string{"-movflags", "+faststart"}},
	{"M4A", "m4a", "aac", []string{"-movflags", "+faststart"}},
	{"M4B", "m4b", "aac", nil},
	{"Vorbis", "ogg", "libvorbis", nil},
	{"Opus", "opus", "libopus", nil},
	{"FLAC", "flac", "flac", nil},
}

func canonicalFixture(t *testing.T, c integrationFormat, dir string) (string, []byte) {
	t.Helper()
	filename := filepath.Join(dir, "source."+c.extension)
	args := []string{"-v", "error", "-f", "lavfi", "-i", "sine=frequency=997:sample_rate=48000:duration=0.20", "-c:a", c.codec, "-metadata", "title=Initial", "-metadata", "artist=Original Artist"}
	args = append(args, c.args...)
	runTool(t, "ffmpeg", append(args, filename)...)
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return filename, data
}

func openCanonical(r io.ReadSeeker) (*core.Editor, error) {
	tag, err := audiometa.OpenMetadata(r)
	if err != nil {
		return nil, err
	}
	return tag.Editor, nil
}

func TestCanonicalContract(t *testing.T) {
	for _, c := range canonicalFormats {
		t.Run(c.name, func(t *testing.T) {
			_, data := canonicalFixture(t, c, t.TempDir())
			contracttest.Run(t, openCanonical, data)
		})
	}
}

func TestCanonicalAudioIntegrityAndTransfer(t *testing.T) {
	dir := t.TempDir()
	flacDir := filepath.Join(dir, "flac")
	if err := os.Mkdir(flacDir, 0700); err != nil {
		t.Fatal(err)
	}
	_, source := canonicalFixture(t, canonicalFormats[len(canonicalFormats)-1], flacDir)
	from := contracttest.OpenFixture(t, openCanonical, source)
	set := map[core.Field][]core.Value{
		core.Title:               {core.TextValue("Transfer — 金属")},
		core.Artist:              {core.TextValue("One"), core.TextValue("Dos")},
		core.Mood:                {core.TextValue("Restless")},
		core.RecordingDate:       {core.DateValue("1984-02")},
		core.ReleaseDate:         {core.DateValue("2025-07-06")},
		core.TrackNumber:         {core.IntValue(2)},
		core.TrackTotal:          {core.IntValue(9)},
		core.RecordingID:         {core.UUIDValue("27b771fd-f244-4341-9ff4-79c7dff56001")},
		core.ReplayGainTrackGain: {core.DecimalValue(-4.5)},
		core.ReplayGainTrackPeak: {core.DecimalValue(1.025)},
		core.Chapters:            {core.ChapterValue(core.Chapter{ID: "a", Title: "Intro", EndMS: 80})},
	}
	for field, values := range set {
		if err := from.SetField(field, values...); err != nil {
			t.Fatal(err)
		}
	}
	if err := from.SetCustom("Collection", "", "Archive 日本語"); err != nil {
		t.Fatal(err)
	}
	from = contracttest.OpenFixture(t, openCanonical, contracttest.SaveFixture(t, from))

	for _, c := range canonicalFormats {
		t.Run(c.name, func(t *testing.T) {
			subdir := t.TempDir()
			input, data := canonicalFixture(t, c, subdir)
			target, err := audiometa.OpenMetadata(bytes.NewReader(data))
			if err != nil {
				t.Fatal(err)
			}
			for field, values := range from.Populated() {
				d, _ := core.Lookup(field)
				if !d.Mutable {
					continue
				}
				if err := target.SetField(field, values...); err != nil {
					t.Fatal(err)
				}
			}
			for key, values := range from.Custom() {
				if err := target.SetCustom(key, values...); err != nil {
					t.Fatal(err)
				}
			}
			out := filepath.Join(subdir, "changed."+c.extension)
			file, err := os.Create(out)
			if err != nil {
				t.Fatal(err)
			}
			save, close := target.Save(file), file.Close()
			if save != nil || close != nil {
				t.Fatalf("save: %v; close: %v", save, close)
			}
			data, err = os.ReadFile(out)
			if err != nil {
				t.Fatal(err)
			}
			back := contracttest.OpenFixture(t, openCanonical, data)
			for field, want := range set {
				actual, _ := back.Get(field)
				t.Run(string(field), func(t *testing.T) { contracttest.AssertValues(t, actual, want) })
			}
			if got, _ := back.GetCustom("COLLECTION"); len(got) != 2 || got[0] != "" || got[1] != "Archive 日本語" {
				t.Fatal(got)
			}
			pcm := func(p string) string {
				return runTool(t, "ffmpeg", "-xerror", "-v", "error", "-i", p, "-map", "0:a:0", "-c:a", "pcm_s32le", "-f", "hash", "-hash", "sha256", "-")
			}
			if initial, result := pcm(input), pcm(out); initial != result {
				t.Fatalf("decoded PCM changed from %s to %s", initial, result)
			} else {
				t.Logf("PCM %s", result)
			}
			type audioStream struct {
				CodecName  string `json:"codec_name"`
				Profile    string `json:"profile"`
				SampleRate string `json:"sample_rate"`
				Channels   int    `json:"channels"`
				TimeBase   string `json:"time_base"`
				Duration   string `json:"duration"`
			}
			stream := func(p string) audioStream {
				var probe struct{ Streams []audioStream }
				s := runTool(t, "ffprobe", "-v", "error", "-select_streams", "a:0", "-show_entries", "stream=codec_name,profile,sample_rate,channels,time_base,duration", "-of", "json", p)
				if err := json.Unmarshal([]byte(s), &probe); err != nil || len(probe.Streams) != 1 {
					t.Fatalf("ffprobe audio stream: %v %s", err, s)
				}
				return probe.Streams[0]
			}
			if a, b := stream(input), stream(out); a != b {
				t.Fatalf("audio stream properties changed: %+v / %+v", a, b)
			}
			packets := func(p string) []string {
				s := runTool(t, "ffprobe", "-v", "error", "-select_streams", "a:0", "-show_packets", "-show_entries", "packet=data_hash", "-show_data_hash", "sha256", "-of", "json", p)
				var probe struct {
					Packets []struct {
						Hash string `json:"data_hash"`
					}
				}
				if err := json.Unmarshal([]byte(s), &probe); err != nil || len(probe.Packets) == 0 {
					t.Fatalf("ffprobe audio packets: %v %s", err, s)
				}
				hashes := make([]string, len(probe.Packets))
				for i, packet := range probe.Packets {
					if len(packet.Hash) != 71 {
						t.Fatalf("packet %d lacks SHA-256: %q", i, packet.Hash)
					}
					hashes[i] = packet.Hash
				}
				return hashes
			}
			if a, b := packets(input), packets(out); !reflect.DeepEqual(a, b) {
				t.Fatal("compressed audio packet hashes changed")
			}
		})
	}
}
