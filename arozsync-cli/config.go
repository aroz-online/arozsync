package main

import (
	"encoding/json"
	"io/ioutil"
)

/*
	Synchronization Setting Configuration
*/

type SyncFolder struct {
	RemoteRootID string //The remote root id for the folder, e.g. user:/
	RemoteFolder string //The relative path for the remote folder, e.g. Video/
	LocalFolder  string //The synchronization path for local folder, e.g. ~/sync/Video/

}

type SyncConfig struct {
	ServerAddr   string //Target server IPv4 address or domain, e.g. 127.0.0.1
	Port         int    //The port where ArozOS Web interface is listening
	SyncInterval int    //The interval in second where the sync routine is execute
	UseHTTPs     bool   //Use secured HTTP connection
	Folders      []*SyncFolder
}

var templateConfig SyncConfig = SyncConfig{
	ServerAddr:   "127.0.0.1",
	Port:         8080,
	SyncInterval: 300,
	UseHTTPs:     false,
	Folders: []*SyncFolder{
		{
			RemoteRootID: "user:/",
			RemoteFolder: "Desktop/",
			LocalFolder:  "./tmp",
		},
	},
}

func generateTemplateConfig(filename string) {
	js, err := json.MarshalIndent(templateConfig, "", " ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filename, js, 0755)
	if err != nil {
		panic(err)
	}
}
