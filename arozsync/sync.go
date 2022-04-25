package main

import (
	"errors"

	"github.com/studio-b12/gowebdav"
)

/*
	WebDAV File System Helpers
*/

func wd_fileExists(c *gowebdav.Client, webdavFilePath string) bool {
	_, err := c.Stat(webdavFilePath)
	if err != nil {
		return false
	}
	return true
}

func wd_isDir(c *gowebdav.Client, webdavFilePath string) bool {
	info, err := c.Stat(webdavFilePath)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func wd_getMtime(c *gowebdav.Client, webdavFilePath string) (int64, error) {
	info, err := c.Stat(webdavFilePath)
	if err != nil {
		return 0, err
	}

	return info.ModTime().Unix(), nil
}

/*
	Recursive logic for Syncing a remote folder to local folder
*/
func SyncFolder(c *gowebdav.Client, webdavRelPath string, localPath string) error {
	//Check if the webdav relative path exists
	if !wd_fileExists(c, webdavRelPath) {
		return errors.New("target folder not found")
	}

	//Check if the remote filepath is a directory
	if !wd_isDir(c, webdavRelPath) {
		return errors.New("target path is not a valid folder")
	}

	//Folder exists on server side
	//Download files that is newer than local versions

	//Upload files that is older than local versions

	return nil
}
func syncFolder(c *gowebdav.Client, webdavRelPath string, localPath string) error {

	return nil
}
