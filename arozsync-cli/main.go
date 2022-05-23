package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/studio-b12/gowebdav"
)

/*
	ArozOS WebDAV Sync CLI Interface
*/

//Connection Related
var congifPath = flag.String("c", "config.json", "The configuration file path")
var username = flag.String("user", "", "Username for authentication")
var password = flag.String("pass", "", "Password for authentication")

//Synchronization Related
var enableRemoteWrite = flag.Bool("rw", true, "Enable upload local change to remote file system. Set this to false for DOWNLOAD FROM REMOTE ONLY.")
var enableRemoteDelete = flag.Bool("rd", false, "Enable remove remote file by deleting in local file system")
var enableLocalWrite = flag.Bool("lw", true, "Enable remote changed to overwrite local changes")
var enableLocalRemove = flag.Bool("ld", true, "Enable remove local file by deleteing in remote file system")
var keepOverwriteVersions = flag.Bool("keepver", true, "Keep older version locally when overwrite local from remote")

//Command Related
var cleanMode = flag.Bool("clean", false, "[DANGER] Execute system cleaning to remove deleted file backups")

//Runtime Global Variables
var executingSyncConfig *SyncConfig
var WebDAVEndpoint string = ""
var syncRunning bool = false
var fileIndexList map[string]int64
var lastSyncTime int64 = 0

func main() {
	flag.Parse()

	//Generate a notification agent

	//Generate a template config if not exists
	if !fileExists(*congifPath) {
		log.Println("Configuration not found. A template has been generated for you. Please modify the template and restart this application.")
		generateTemplateConfig("config.json")
		os.Exit(0)
	}

	if *username == "" || *password == "" {
		log.Println("Missing username or password. Usage: arozsync-cli -user={username} -pass={password}")
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

	if *cleanMode {
		RunCleaningCycle(executingSyncConfig)
		log.Println("[INFO] All sync folder trash cleared")
		os.Exit(0)
	}

	//Start the connection to server
	serverConnEndpt := ""
	if executingSyncConfig.UseHTTPs {
		serverConnEndpt += "https://"
	} else {
		serverConnEndpt += "http://"
	}

	serverConnEndpt += executingSyncConfig.ServerAddr + ":" + strconv.Itoa(executingSyncConfig.Port) + "/webdav"
	WebDAVEndpoint = serverConnEndpt

	//Test Connections
	c := gowebdav.NewClient(serverConnEndpt+"/user", *username, *password)
	_, err = c.ReadDir("/")
	if err != nil {
		notification("Arozsync Start Failed", "Arozsync is unable to start. Please make sure your configuration file is correct and you have enabled WebDAV on your account.")
		log.Println("Sync test failed. Please make sure your configuration file is correct and you have enabled WebDAV on your account.")
		log.Fatal(err)
	}

	notification("Arozsync Started", "Arozsync started. You will receive notification if there are any errors.")

	//Setup Ready! Sync Now
	lastSyncTime = time.Now().Unix()
	SyncFoldersFromConfig(executingSyncConfig)

	//Start Sync Progress Ticker
	ticker := time.NewTicker(time.Duration(executingSyncConfig.SyncInterval) * time.Second)
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			err = SyncFoldersFromConfig(executingSyncConfig)
			if err != nil {
				notification("Sync Failed!", "Failed to execute file synchronization sequence: "+err.Error()+"\n\n "+time.Now().Format("2006.01.02 15:04:05"))
			}
		}
	}

}
