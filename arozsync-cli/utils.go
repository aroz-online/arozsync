package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/studio-b12/gowebdav"
	"gopkg.in/toast.v1"
)

func notification(title string, message string) error {
	iconPath, _ := filepath.Abs("./icon.png")
	notification := toast.Notification{
		AppID:   "Arozsync",
		Title:   title,
		Message: message,
		Icon:    iconPath,
		Actions: []toast.Action{
			{"protocol", "Open Web-desktop", "http://localhost:8080/"},
			{"protocol", "Okay!", ""},
		},
	}

	return notification.Push()
}

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

func fileInLocalVerBuffer(filename string) bool {
	return strings.Contains(filepath.ToSlash(filename), "/.localver/")
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

	err := c.WriteStream(remote, file, 0775)
	if err != nil {
		return err
	}
	//Update the local file modtime so it wont get pull down again
	uptime := time.Now().Local()

	//Set both access time and modified time of the file to the current time
	err = os.Chtimes(local, uptime, uptime)
	if err != nil {
		return err
	}

	return nil
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

//Get file MD5 Sum
func GetFileMD5Sum(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer file.Close()

	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}
