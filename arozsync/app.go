package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/studio-b12/gowebdav"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Scan return a list of scanned nodes
func (a *App) ScanNearbyNodes() string {
	hosts := MdnsScanner.Scan(10, "arozos.com")
	js, _ := json.Marshal(hosts)
	return string(js)
}

func (a *App) OpenLinkInLocalBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Println("Unable to open browser with given link: ", err)
	}
}

//Try all the ips and see if it is connectable
func (a *App) TryConnect(ips []string, username string, password string, remember bool) []string {
	succConn := []string{}
	for _, ip := range ips {
		thisRoot := "http://" + ip + "/webdav/user"
		c := gowebdav.NewClient(thisRoot, username, password)
		_, err := c.ReadDir("/")
		if err != nil {
			log.Println(err)
			continue
		} else {
			//ok!
			succConn = append(succConn, ip)
		}
	}

	if remember && len(succConn) > 0 {
		SaveCred(&LoginCred{
			Username: username,
			Password: password,
			IPs:      ips,
		})
	}
	return succConn
}
