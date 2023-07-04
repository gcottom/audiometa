package mp3mp4tag

import (
	"strings"
	"errors"
)

func getFileType(filepath string) (*string, error) {
	fileTypeArr := strings.Split(filepath, ".")
	lastIndex := len(fileTypeArr) - 1
	fileType := fileTypeArr[lastIndex]
	fileType = strings.ToLower(fileType)
	if(fileType == "mp3" || fileType == "m4p" || fileType == "m4a" || fileType == "m4b"){
		return &fileType, nil
	}else{
		return nil, errors.New("Format: Unsupported Format: "+ fileType)
	}
}
