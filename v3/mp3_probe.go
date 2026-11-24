package audiometa

import (
	"errors"
	"io"
)

var mpeg1Layer3Bitrates = [...]int{0, 32, 40, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320}
var mpeg2Layer3Bitrates = [...]int{0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160}

type mp3Frame struct {
	version byte
	rate    int
	size    int64
}

func layer3Frame(b []byte) (mp3Frame, bool) {
	if b[0] != 0xff || b[1]&0xe0 != 0xe0 || b[1]>>1&3 != 1 || b[3]&3 == 2 {
		return mp3Frame{}, false
	}
	version, bitrate, rateIndex := b[1]>>3&3, b[2]>>4, b[2]>>2&3
	if version == 1 || bitrate == 0 || bitrate == 15 || rateIndex == 3 {
		return mp3Frame{}, false
	}
	rate := [...]int{44100, 48000, 32000}[rateIndex]
	factor, kbps := 144, mpeg1Layer3Bitrates[bitrate]
	if version != 3 {
		factor, kbps = 72, mpeg2Layer3Bitrates[bitrate]
		rate /= 2
		if version == 0 {
			rate /= 2
		}
	}
	size := factor*kbps*1000/rate + int(b[2]>>1&1)
	return mp3Frame{version, rate, int64(size)}, true
}

func isUntaggedMP3(r io.ReadSeeker, header []byte) (bool, error) {
	first, valid := layer3Frame(header)
	if !valid {
		return false, nil
	}
	start, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return false, err
	}
	// Require a compatible second frame header at the calculated frame
	// boundary instead of accepting any data starting with an 11-bit sync.
	if _, err := r.Seek(first.size, io.SeekCurrent); err != nil {
		return false, err
	}
	var next [4]byte
	_, readErr := io.ReadFull(r, next[:])
	if _, err := r.Seek(start, io.SeekStart); err != nil {
		return false, err
	}
	if readErr != nil {
		if errors.Is(readErr, io.EOF) || errors.Is(readErr, io.ErrUnexpectedEOF) {
			return false, nil
		}
		return false, readErr
	}
	second, valid := layer3Frame(next[:])
	return valid && first.version == second.version && first.rate == second.rate, nil
}
