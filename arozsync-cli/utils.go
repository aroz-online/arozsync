package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

func WebDAVFileExists(c *gowebdav.Client, name string) bool {
	_, err := c.Stat(name)
	return err == nil
}

func fileInTrash(filename string) bool {
	return strings.Contains(filepath.ToSlash(filename), "/.trash/")
}

func mtime(filename string) (int64, error) {
	file, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}

	return file.ModTime().Unix(), nil
}

func existsInLastSync(relativePath string) bool {
	_, ok := fileIndexList[relativePath]
	return ok
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
	info, err := c.Stat(remote)
	if err != nil {
		return err
	}
	err = os.Chtimes(local, info.ModTime(), info.ModTime())
	return err
}

func UploadToWebDAV(c *gowebdav.Client, remote string, local string) error {
	file, _ := os.Open(local)
	defer file.Close()

	return c.WriteStream(remote, file, 0775)
}

//Move a remote resources to trash
func WebDAVMoveToTrash(c *gowebdav.Client, name string) error {
	if !WebDAVFileExists(c, name) {
		return errors.New("remote file not exists")
	}

	trashPath := filepath.ToSlash(filepath.Join(filepath.Dir(name), ".trash", filepath.Base(name)+"."+strconv.Itoa(int(time.Now().Unix()))))
	//log.Println("Moving ", name, " to ", trashPath)
	return c.Rename(name, trashPath, true)
}
