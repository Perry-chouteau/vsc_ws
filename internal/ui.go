// internal/ui.go

package internal

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/layout"
)

type UI struct {
	App        fyne.App
	Window     fyne.Window
	//json
	JSONPath   string
	JsonData *FolderData
	//workspace
	WSPath string

	view *fyne.Container

	// Widgets for add folder form
	emojiSelect    *widget.Select
	newEmojiEntry  *widget.Entry
	addEmojiBtn    *widget.Button
	nameEntry      *widget.Entry
	pathEntry      *widget.Entry
	addFolderBtn   *widget.Button

	selectedEmoji string
}

func NewUI(app fyne.App, window fyne.Window, jsonPath string, jsonData *FolderData, wsPath string) *UI {
	return &UI{
		App:        app,
		Window:     window,
		JsonData: jsonData,
		JSONPath:   jsonPath,
		WSPath:   wsPath,
		selectedEmoji: " ",
	}
}

func (ui *UI) Build() fyne.CanvasObject {
	ui.view = container.NewVBox()
	//update
	ui.setupEmojiSelection()
	ui.setupAddFolderForm()
	ui.refreshList()

	//view
	scroll := container.NewVScroll(ui.view)
	scroll.SetMinSize(fyne.NewSize(400, 300))

	// [tab , space , [emoji , add]]
	labelAndEmojiHandler := container.New(layout.NewGridLayout(3), 
		widget.NewLabel("Ajouter un dossier:"),
		layout.NewSpacer(),
		container.New(layout.NewGridLayout(2), ui.newEmojiEntry, ui.addEmojiBtn),
	)

	entryWrapper := container.NewMax(ui.nameEntry)
	entryWrapper.Resize(fyne.NewSize(300, 36))

	emojiAndName := container.NewHBox(
		container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 36)), ui.emojiSelect),
		container.New(layout.NewGridWrapLayout(fyne.NewSize(550, 36)), entryWrapper),
	)
	

	form := container.NewVBox(
		labelAndEmojiHandler,
		emojiAndName,
		ui.pathEntry, 
		ui.addFolderBtn,
	)

	content := container.NewVBox(
		widget.NewLabel("Folders (activate/deactivate & delete):"),
		scroll,
		form,
	)

	return content
}

func (ui *UI) setupEmojiSelection() {
	emojis := &ui.JsonData.Emojis

	if len(*emojis) == 0 {
		*emojis = []string{" "}
	}
	ui.selectedEmoji = (*emojis)[0]

	ui.emojiSelect = widget.NewSelect(*emojis, func(selected string) {
		ui.selectedEmoji = selected
	})
	ui.emojiSelect.SetSelected(ui.selectedEmoji)

	ui.newEmojiEntry = widget.NewEntry()
	ui.newEmojiEntry.SetPlaceHolder("New Emoji")

	ui.addEmojiBtn = widget.NewButton("Add Emoji", func() {
		newEmoji := strings.TrimSpace(ui.newEmojiEntry.Text)
		if newEmoji == "" {
			dialog.ShowInformation("Invalid input", "Emoji cannot be empty", ui.Window)
			return
		}
		for _, emoji := range *emojis {
			if emoji == newEmoji {
				ui.emojiSelect.SetSelected(newEmoji)
				ui.newEmojiEntry.SetText("")
				return
			}
		}

		*emojis = append(*emojis, newEmoji)
		ui.emojiSelect.Options = *emojis
		ui.emojiSelect.Refresh()
		ui.selectedEmoji = newEmoji
		ui.emojiSelect.SetSelected(newEmoji)
		ui.newEmojiEntry.SetText("")

		if err := SaveFolders(ui.JSONPath, ui.JsonData); err != nil {
			dialog.ShowError(err, ui.Window)
		}
		if err := UpdateWorkspaceFoldersOnly(ui.WSPath, ui.JsonData.Folders); err != nil {
			dialog.ShowError(err, ui.Window)
		}
	
	})
}

func (ui *UI) setupAddFolderForm() {
	ui.nameEntry = widget.NewEntry()
	ui.nameEntry.SetPlaceHolder("Folder Name (sans emoji)")

	ui.pathEntry = widget.NewEntry()
	ui.pathEntry.SetPlaceHolder("Path (ex: workspace)")

	ui.addFolderBtn = widget.NewButton("Add Folder", func() {
		nameText := strings.TrimSpace(ui.nameEntry.Text)
		path := strings.TrimSpace(ui.pathEntry.Text)
		if nameText == "" || path == "" {
			dialog.ShowInformation("Invalid input", "Name and Path cannot être vides", ui.Window)
			return
		}
		fullName := ui.selectedEmoji + " " + nameText
		for _, f := range ui.JsonData.Folders {
			if f.Path == path {
				dialog.ShowInformation("Doublon path", "Ce chemin existe déjà", ui.Window)
				return
			}
		}
		newFolder := Folder{Name: fullName, Path: path, IsActive: true}
		ui.JsonData.Folders = append(ui.JsonData.Folders, newFolder)
		if err := SaveFolders(ui.JSONPath, ui.JsonData); err != nil {
			dialog.ShowError(err, ui.Window)
			return
		}
		if err := UpdateWorkspaceFoldersOnly(ui.WSPath, ui.JsonData.Folders); err != nil {
			dialog.ShowError(err, ui.Window)
		}	
		ui.nameEntry.SetText("")
		ui.pathEntry.SetText("")
		ui.refreshList()
	})
}

func (ui *UI) refreshList() {
	ui.view.Objects = nil
	for i, folder := range ui.JsonData.Folders {
		idx := i
		f := folder

		check := widget.NewCheck(fmt.Sprintf("%s (%s)", f.Name, f.Path), func(active bool) {
			ui.JsonData.Folders[idx].IsActive = active
			if err := SaveFolders(ui.JSONPath, ui.JsonData); err != nil {
				dialog.ShowError(err, ui.Window)
			}
			if err := UpdateWorkspaceFoldersOnly(ui.WSPath, ui.JsonData.Folders); err != nil {
				dialog.ShowError(err, ui.Window)
			}		
		})
		check.SetChecked(f.IsActive)

		deleteBtn := widget.NewButton("Delete", func() {
			dialog.ShowConfirm("Confirm Delete", "Delete this folder?", func(confirm bool) {
				if confirm {
					ui.JsonData.Folders = append(ui.JsonData.Folders[:idx], ui.JsonData.Folders[idx+1:]...)
					if err := SaveFolders(ui.JSONPath, ui.JsonData); err != nil {
						dialog.ShowError(err, ui.Window)
					}
					if err := UpdateWorkspaceFoldersOnly(ui.WSPath, ui.JsonData.Folders); err != nil {
						dialog.ShowError(err, ui.Window)
					}				
					ui.refreshList()
				}
			}, ui.Window)
		})

		row := container.NewHBox(check, deleteBtn)
		ui.view.Add(row)
	}

	if err := UpdateWorkspaceFoldersOnly(ui.WSPath, ui.JsonData.Folders); err != nil {
		dialog.ShowError(err, ui.Window)
	}

	ui.view.Refresh()
}