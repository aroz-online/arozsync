package main

/*

	arozsync
	Author: tobychui

	This is a simple desktop application for synchronizing your files
	on arozos to your local PC, best suited when you love to work with
	multiple PCs and laptops
*/
import (
	"arozsync/mod/mdns"
	"embed"
	"log"
	"os"

	"github.com/wailsapp/wails/v2/pkg/options/mac"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed frontend/src
var assets embed.FS
var MdnsScanner *mdns.MDNSHost

//go:embed build/appicon.png
var icon []byte

func main() {
	/*
		Create all required managers
	*/

	name, err := os.Hostname()
	if err != nil {
		name = "Unknown Computer"
	}
	s, err := mdns.NewMDNS(mdns.NetworkHost{
		HostName: name,
		Domain:   "scan.aroz.sync",
		Model:    "computer",
	})

	if err != nil {
		//Show error page

	} else {
		MdnsScanner = s
	}

	/*
		Create an App Instance
	*/
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err = wails.Run(&options.App{
		Title:             "arozsync",
		Width:             400,
		Height:            600,
		MinWidth:          400,
		MinHeight:         600,
		MaxWidth:          400,
		MaxHeight:         600,
		DisableResize:     true,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		RGBA:              &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		Assets:            assets,
		Menu:              nil,
		Logger:            nil,
		LogLevel:          logger.DEBUG,
		OnStartup:         app.startup,
		OnDomReady:        app.domReady,
		OnBeforeClose:     app.beforeClose,
		OnShutdown:        app.shutdown,
		WindowStartState:  options.Normal,
		Bind: []interface{}{
			app,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			// DisableFramelessWindowDecorations: false,
			WebviewUserDataPath: "",
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "arozsync",
				Message: "Sync your arozos files to local disk",
				Icon:    icon,
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
