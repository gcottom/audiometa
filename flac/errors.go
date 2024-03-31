package flac

import "errors"

var (
	ErrorNotVorbisComment = errors.New("not a vorbis comment metadata block")
	ErrorUnexpEof         = errors.New("unexpected end of stream")
	ErrorMalformedComment = errors.New("malformed comment")
	ErrorInvalidFieldName = errors.New("malformed field Name")
	// ErrorNotPictureMetadataBlock is returned if the metadata provided is not a picture block.
	ErrorNotPictureMetadataBlock = errors.New("not a picture metadata block")
	// ErrorUnsupportedMIME is returned if the provided image MIME type is unsupported.
	ErrorUnsupportedMIME = errors.New("unsupported MIME")
	// ErrorNoFLACHeader indicates that "fLaC" marker not found at the beginning of the file
	ErrorNoFLACHeader = errors.New("fLaC head incorrect")
	// ErrorNoStreamInfo indicates that StreamInfo Metablock not present or is not the first Metablock
	ErrorNoStreamInfo = errors.New("stream info not present")
	// ErrorStreamInfoEarlyEOF indicates that an unexpected EOF is hit while reading StreamInfo Metablock
	ErrorStreamInfoEarlyEOF = errors.New("unexpected end of stream while reading stream info")
	// ErrorNoSyncCode indicates that the frames are malformed as the sync code is not present after the last Metablock
	ErrorNoSyncCode = errors.New("frames do not begin with sync code")
)
