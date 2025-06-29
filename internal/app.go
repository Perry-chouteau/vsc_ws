// internal/app.go

package internal

import (
	"log"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func Run() {
	//env
	jsonPath := getEnv("VSC_WS_JSON")
	wsPath := getEnv("VSC_WS")

	//loadfolders
	jsonData, err := LoadFolders(jsonPath)
	if err != nil {
		log.Fatalf("Failed to load JSON file: %v", err)
	}

	//app
	myApp := app.New()

	//window
	myWin := myApp.NewWindow("VS Code Workspace editor")
	ui := NewUI(myApp, myWin, jsonPath, jsonData, wsPath)

	myWin.SetContent(ui.Build())
	myWin.Resize(fyne.NewSize(400, 450))
	myWin.SetFixedSize(true)
	myWin.ShowAndRun()
}