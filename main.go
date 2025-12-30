package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"

	"wailts/config"
	"wailts/models"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	// Initialize configuration
	config.Init()

	// Initialize the database service
	service, err := models.NewService(config.Current.DatabasePath)
	if err != nil {
		println("Error initializing database:", err.Error())
		return
	}

	// Create an instance of the app structure, with the service
	app := NewApp(service)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "BudgetTracker",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		// BackgroundColour: &options.RGBA{R: 128, G: 38, B: 54, A: 1},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
		Linux: &linux.Options{
			Icon: icon,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
