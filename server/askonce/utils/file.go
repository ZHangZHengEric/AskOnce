package utils

import (
	"mime/multipart"
	"os"
)

func GetFileHeader(file *os.File) (*multipart.FileHeader, error) {
	// get file size
	fileStat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// create *multipart.FileHeader
	return &multipart.FileHeader{
		Filename: fileStat.Name(),
		Size:     fileStat.Size(),
	}, nil
}
