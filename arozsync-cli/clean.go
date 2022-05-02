package main

import (
	"log"
	"os"
	"path/filepath"
)

//Clean all the trash files in the .trash subdirectories
func RunCleaningCycle(config *SyncConfig) {
	rootFolders := []string{}
	for _, thisFolder := range config.Folders {
		if fileExists(thisFolder.LocalFolder) {
			rootFolders = append(rootFolders, thisFolder.LocalFolder)
		}

	}

	//For each sync roots, walk all the folders and delete those named .trash
	for _, thisRoot := range rootFolders {
		markForDeleteFolders := []string{}
		filepath.Walk(thisRoot, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() && filepath.Base(path) == ".trash" {
				markForDeleteFolders = append(markForDeleteFolders, path)
			}
			return nil
		})

		for _, deletePendingFoler := range markForDeleteFolders {
			err := os.RemoveAll(deletePendingFoler)
			if err == nil {
				log.Println("[OK] Removing trash folder -> ", deletePendingFoler)
			}

			//Check if there are any more folders in parent dir. If no, remove the parent dir as well.
			itemsInparentDir, err := filepath.Glob(filepath.Clean(filepath.Dir(deletePendingFoler)) + "/*")
			if err == nil && len(itemsInparentDir) == 0 {
				os.Remove(filepath.Dir(deletePendingFoler))
			}
		}
	}
}
