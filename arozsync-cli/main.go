package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/studio-b12/gowebdav"
)

/*
	ArozOS WebDAV Sync CLI Interface
*/

var congifPath = flag.String("c", "config.json", "The configuration file path")
var username = flag.String("user", "", "Username for authentication")
var password = flag.String("pass", "", "Password for authentication")

var executingSyncConfig *SyncConfig
var WebDAVEndpoint string = ""

func main() {
	flag.Parse()

	//Generate a template config if not exists
	if !fileExists(*congifPath) {
		log.Println("Configuration not found. A template has been generated for you. Please modify the template and restart this application.")
		generateTemplateConfig("config.json")
		os.Exit(0)
	}

	//Read the configuration file
	configuration, err := ioutil.ReadFile(*congifPath)
	if err != nil {
		log.Println("Unable to read config file. Terminating.")
		panic(err)
	}

	//Parse the configuration file
	err = json.Unmarshal(configuration, &executingSyncConfig)
	if err != nil {
		log.Println("Unable to parse config file. Terminating.")
		panic(err)
	}

	//Start the connection to server
	serverConnEndpt := ""
	if executingSyncConfig.UseHTTPs {
		serverConnEndpt += "https://"
	} else {
		serverConnEndpt += "http://"
	}

	serverConnEndpt += executingSyncConfig.ServerIpv4 + ":" + strconv.Itoa(executingSyncConfig.Port) + "/webdav"
	WebDAVEndpoint = serverConnEndpt

	//Test Connections
	c := gowebdav.NewClient(serverConnEndpt+"/user", *username, *password)
	_, err = c.ReadDir("/")
	if err != nil {
		log.Println("Sync test failed. Please make sure your configuration file is correct.")
		log.Fatal(err)
	}

	//Setup Ready! Sync Now
	SyncFoldersFromConfig(executingSyncConfig)

	//Start Sync Progress Ticker
	ticker := time.NewTicker(time.Duration(executingSyncConfig.SyncInterval) * time.Second)
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			SyncFoldersFromConfig(executingSyncConfig)
		}
	}

}

func SyncFoldersFromConfig(config *SyncConfig) {
	log.Println("Starting folder sync routine")
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
		syncFolderProcedure(c, thisFolder.LocalFolder, "/")
	}
	//If file not exists locally, download it from server side

}

func syncFolderProcedure(c *gowebdav.Client, localBase string, remoteBase string) {
	fmt.Println(remoteBase)
	files, _ := c.ReadDir(remoteBase)
	for _, file := range files {
		if file.IsDir() {
			//Is folder.
			syncFolderProcedure(c, localBase, remoteBase+file.Name()+"/")
		} else {
			//Is File
			thisRelativePath := filepath.ToSlash(filepath.Join(remoteBase, file.Name()))
			log.Println(thisRelativePath)
		}

	}
}
