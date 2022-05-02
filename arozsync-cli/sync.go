package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/studio-b12/gowebdav"
)

func SyncFoldersFromConfig(config *SyncConfig) {
	if syncRunning {
		log.Println("[FAILED] Another sync routine running in progress. Skipping.")
		return
	}
	syncRunning = true
	log.Println("[INFO] Starting folder sync routine")
	for _, thisFolder := range config.Folders {
		folderRootname := thisFolder.RemoteRootID
		folderRootname = strings.ReplaceAll(folderRootname, ":/", "")

		//Establish connection to target endpoint
		cleanedRemoteFolder := filepath.ToSlash(filepath.Clean(thisFolder.RemoteFolder))
		connectPath := WebDAVEndpoint + filepath.ToSlash(filepath.Join("/", folderRootname, cleanedRemoteFolder))
		c := gowebdav.NewClient(connectPath, *username, *password)

		//Check if the target local folder exists. If not, create it
		if !fileExists(thisFolder.LocalFolder) {
			os.MkdirAll(thisFolder.LocalFolder, 0775)
		}

		//Do a recursive folder update
		syncRemoteFolderProcedure(c, thisFolder.LocalFolder, "/")
		SyncLocalFolderprocedure(c, thisFolder.LocalFolder, "/")
		lastSyncTime = time.Now().Unix()
	}
	syncRunning = false
	log.Println("[INFO] Folder sync routine completed. Next sync in ", config.SyncInterval, " seconds.")
}

//Perform sync on the given filepath
func syncRemoteFolderProcedure(c *gowebdav.Client, localBase string, remoteBase string) {
	files, err := c.ReadDir(remoteBase)
	if err != nil {
		log.Println("[FAILED] Unable to sync folder to ", localBase)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			//Is folder.
			syncRemoteFolderProcedure(c, localBase, remoteBase+file.Name()+"/")
		} else {
			//Is File.
			thisRelativePath := filepath.ToSlash(filepath.Join(remoteBase, file.Name()))

			//Check for sync actions
			expectedLocalPath := filepath.Join(localBase, thisRelativePath)
			if !fileExists(expectedLocalPath) {
				//This file not exists locally.

				if len(fileIndexList) > 0 && existsInLastSync(thisRelativePath) {
					if *enableRemoteDelete {
						//This file exists in last sync. This file is recently deleted
						err := WebDAVMoveToTrash(c, thisRelativePath)
						if err != nil {
							log.Println("[FAILED] Unale to delete remote file ", thisRelativePath, err.Error())
							continue
						} else {
							log.Println("[OK] Remote Deleted file ", thisRelativePath)
						}
					} else {
						log.Println("[ERROR] Remote Delete (-rd) flag is set to false. Enable this flag in order to delete file from local sync folder.")
						//Re-download the missing file to keep local file system structured based on startup rules
						downloadFromWebDAV(c, thisRelativePath, expectedLocalPath)
					}

				} else {
					//Download file from server
					err := downloadFromWebDAV(c, thisRelativePath, expectedLocalPath)
					if err != nil {
						log.Println("[FAILED] Unable to sync ", thisRelativePath, err.Error())
						continue
					} else {
						log.Println("[OK] Downloaded ", thisRelativePath)
					}
				}

			} else {
				//File already exists. Compare mtime and synchronize it
				localFileModTime, err := mtime(expectedLocalPath)
				if err != nil {
					continue
				}
				if file.ModTime().Unix() < localFileModTime && *enableRemoteWrite {
					//The local file is newer. Uplaod it
					err := UploadToWebDAV(c, thisRelativePath, expectedLocalPath)
					if err != nil {
						log.Println("[FAILED] Unable to sync ", thisRelativePath, err.Error())
						continue
					} else {
						log.Println("[OK] Updated Remote Copy of ", thisRelativePath)
					}

				} else if file.ModTime().Unix() > localFileModTime && *enableLocalWrite {
					//The remote file is newer. Download it
					err := downloadFromWebDAV(c, thisRelativePath, expectedLocalPath)
					if err != nil {
						log.Println("[FAILED] Unable to sync ", thisRelativePath, err.Error())
						continue
					} else {
						log.Println("[OK] Updated Local Copy of ", thisRelativePath)
					}

				}
			}

		}
	}
}

//Sync local files
func SyncLocalFolderprocedure(c *gowebdav.Client, localBase string, remoteBase string) {
	//Upload Newly Created Files

	if len(fileIndexList) > 0 && *enableRemoteWrite {
		//Upload newly created files
		filepath.Walk(localBase, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() && !fileInTrash(path) {
				//Is File. Check if it is created after last sync time and not exists in fileIndexList
				relPath, err := filepath.Rel(localBase, path)
				if err != nil {
					return nil
				}
				relPath = "/" + filepath.ToSlash(relPath)

				if !existsInLastSync(relPath) {
					//This is a newly created file. Upload it to server
					err = UploadToWebDAV(c, relPath, path)
					if err != nil {
						log.Println("[FAILED] Unable to upload ", relPath, err.Error())
					} else {
						log.Println("[OK] Uploaded ", relPath)
					}

				}
			}
			return nil
		})
	}

	//Delete files that no longer exists
	thisScanFilelist := map[string]int64{}
	filepath.Walk(localBase, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && !fileInTrash(path) {
			//Is File (that is not trash)
			relPath, err := filepath.Rel(localBase, path)
			if err != nil {
				return nil
			}

			relPath = "/" + filepath.ToSlash(relPath)

			if len(fileIndexList) > 0 && !WebDAVFileExists(c, relPath) && *enableLocalRemove {
				//This file has been removed from server side. Remove it locally
				log.Println("[INFO] Moving " + path + " to trash")
				os.MkdirAll(filepath.Join(filepath.Dir(path), ".trash"), 0775)
				os.Rename(path, filepath.Join(filepath.Dir(path), ".trash", filepath.Base(path)))
			} else if !WebDAVFileExists(c, relPath) && *enableLocalRemove {
				//Just started up and find some new files that do not exists on remote side. Upload this file
				err = UploadToWebDAV(c, relPath, path)
				if err != nil {
					log.Println("[FAILED] Unable to upload ", relPath, err.Error())
				} else {
					log.Println("[OK] Uploaded ", relPath)
					thisScanFilelist[relPath] = info.ModTime().Unix()
				}
			} else {
				thisScanFilelist[relPath] = info.ModTime().Unix()
			}
		}
		return nil
	})

	fileIndexList = thisScanFilelist
}
