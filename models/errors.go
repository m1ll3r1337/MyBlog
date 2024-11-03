package models

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

var (
	ErrNotFound   = errors.New("recourse could not be found")
	ErrEmailTaken = errors.New("models: email already taken")
)

type FileError struct {
	Issue string
}

func (fe FileError) Error() string {
	return fmt.Sprintf("invalid file: %v", fe.Issue)
}

func checkContentType(r io.ReadSeeker, allowedTypes []string) error {
	testBytes := make([]byte, 512)
	_, err := r.Read(testBytes)
	if err != nil {
		return fmt.Errorf("checking content type: %v", err)
	}
	_, err = r.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("checking content type: %v", err)
	}

	contentType := http.DetectContentType(testBytes)
	for _, v := range allowedTypes {
		if contentType == v {
			return nil
		}
	}
	return FileError{
		Issue: fmt.Sprintf("content-type %s not allowed", contentType),
	}
}

func checkExtension(filename string, allowedExtensions []string) error {
	if hasExtension(filename, allowedExtensions) {
		return nil
	}
	return FileError{
		Issue: fmt.Sprintf("extension %s not allowed", filepath.Ext(filename)),
	}
}
