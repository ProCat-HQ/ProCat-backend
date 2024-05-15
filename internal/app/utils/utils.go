package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

var (
	ImagesFileExtensions = map[string]string{
		"pngExt":  "89504e47",
		"gifExt":  "47494638",
		"jpg1Ext": "ffd8ffe1",
		"jpg2Ext": "ffd8ffdb",
		"jpg3Ext": "ffd8ffe0",
		"jpg4Ext": "ffd8ffee",
	}
)

func SaveUploadedFileAndCheckExtension(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	fileExt := make([]byte, 4)
	n, err := src.Read(fileExt)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("can't read from file")
	}
	fileExtHex := fmt.Sprintf("%x", fileExt)
	flag := false
	for _, v := range ImagesFileExtensions {
		if fileExtHex == v {
			flag = true
			break
		}
	}
	if !flag {
		return errors.New("bad file extension")
	}

	_, err = src.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func RemoveFiles(filenames []string, dir string) error {
	for _, v := range filenames {
		err := os.Remove(dir + v)
		if err != nil {
			return err
		}
	}
	return nil
}
