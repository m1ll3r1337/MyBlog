package models

import (
	"io"
	"os"
)

type MarkdownReader interface {
	Read(path string) (string, error)
}

type MarkdownWriter interface {
	Write(path, content string) error
}

type FileReader struct {}

type FileWriter struct {}

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

func (fw FileWriter) Write (path, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
