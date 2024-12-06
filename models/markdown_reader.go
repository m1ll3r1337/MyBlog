package models

import (
	"io"
	"os"
)

type MarkdownReader interface {
	Read(path string) (string, error)
}

type FileReader struct {}

func (fr FileReader) Read (path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
