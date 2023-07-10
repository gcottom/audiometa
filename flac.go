package mp3mp4tag

import (
	"github.com/go-flac/flacvorbis"
	"github.com/go-flac/go-flac"
)

func extractFLACComment(fileName string) (*flacvorbis.MetaDataBlockVorbisComment, int) {
	f, err := flac.ParseFile(fileName)
	if err != nil {
		panic(err)
	}

	var cmt *flacvorbis.MetaDataBlockVorbisComment
	var cmtIdx int
	for idx, meta := range f.Meta {
		if meta.Type == flac.VorbisComment {
			cmt, err = flacvorbis.ParseFromMetaDataBlock(*meta)
			cmtIdx = idx
			if err != nil {
				panic(err)
			}
		}
	}
	return cmt, cmtIdx
}
func remove(slice []*flac.MetaDataBlock, s int)[]*flac.MetaDataBlock {
    return append(slice[:s], slice[s+1:]...)
}