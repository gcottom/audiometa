package audiometa

import "errors"

var (
	ErrFLACParse    = errors.New("error parsing flac stream")
	ErrFLACCmtParse = errors.New("error parsing flac comment")

	ErrMP4AtomOutOfBounds  = errors.New("mp4 atom out of bounds")
	ErrMP4InvalidAtomSize  = errors.New("mp4 atom has invalid size")
	ErrMP4InvalidEncoding  = errors.New("invalid encoding: got wrong number of bytes")
	ErrMP4IlstAtomMissing  = errors.New("ilst atom is missing")
	ErrMP4InvalidCntntType = errors.New("invalid content type")

	ErrOggInvalidSgmtTblSz = errors.New("invalid segment table size")
	ErrOggInvalidHeader    = errors.New("invalid ogg header")
	ErrOggInvalidCRC       = errors.New("invalid CRC")
	ErrOggMissingCOP       = errors.New("missing ogg COP packet")
	ErrOggImgConfigFail    = errors.New("failed to get image config")
	ErrOggCodecNotSpprtd   = errors.New("unsupported codec for ogg")

	ErrMP3ParseFail = errors.New("error parsing mp3")

	ErrNoMethodAvlble = errors.New("no method available for this filetype")
)
