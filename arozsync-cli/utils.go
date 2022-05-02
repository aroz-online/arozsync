package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/studio-b12/gowebdav"
)

func fileExists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

func mtime(filename string) (int64, error) {
	file, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}

	return file.ModTime().Unix(), nil
}

//Download file from WebDAV endpoint, overwrite it with remote modtime
func downloadFromWebDAV(c *gowebdav.Client, remote string, local string) error {
	reader, err := c.ReadStream(remote)
	if err != nil {
		return err
	}

	if !fileExists(filepath.Dir(local)) {
		os.MkdirAll(filepath.Dir(local), 0775)
	}

	file, err := os.Create(local)
	if err != nil {
		return err
	}

	io.Copy(file, reader)
	file.Close()

	//Overwrite the modtime
	info, _ := c.Stat(remote)
	err = os.Chtimes(local, info.ModTime(), info.ModTime())

	return nil
}

func UploadToWebDAV(c *gowebdav.Client, remote string, local string) error {
	file, _ := os.Open(local)
	defer file.Close()

	return c.WriteStream(remote, file, 0775)
}
