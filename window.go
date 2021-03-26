package main

import (
	"encoding/json"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Window struct {
	app fyne.App
	win fyne.Window

	tree   *widget.Tree
	editor *widget.Entry
	status *widget.Label

	current Hosts
	hosts   *HostsGroup
}

func Run() {
	myApp := app.New()
	win := myApp.NewWindow("Go Hosts!")

	root := NewHostsGroup("")
	data := myApp.Preferences().String("hosts")

	_ = json.Unmarshal([]byte(data), &root.Items)
	root.Add(NewHostsItem("System"))

	backup := NewHostsItem("Backup")
	root.Add(backup)
	(&Window{app: myApp, win: win, hosts: root}).Run()


	system := loadSystem()
	backup.SetContent(system)
}

func (w *Window) Run() {

	w.init()

	system := loadSystem()
	err := saveBackup(system)
	if err != nil {
		log.Println(err)
		dialog.NewInformation("Error", "Saving backup file error!\n"+err.Error(), w.win).Show()
	}

	w.win.ShowAndRun()
}

func (w *Window) save() {
	data, _ := json.Marshal(w.hosts)
	w.app.Preferences().SetString("hosts", string(data))

	system := loadSystem()
	err := saveBackup(system)
	if err != nil {
		log.Println(err)
		dialog.NewInformation("Error", "Saving backup file error!\n"+err.Error(), w.win).Show()
		return
	}
	content := w.hosts.GetContent()
	if content == "" {
		return
	}
	err = saveSystem(content)
	if err != nil {
		log.Println(err)
		dialog.NewInformation("Error", "Saving system hosts file error!\n"+err.Error(), w.win)
		return
	}
	w.showStatus("Save success!")
}

func (w *Window) init() {

	toolbar := w.createToolbar()
	editor := w.createEditor()
	tree := w.createTree()

	statusLabel := widget.NewLabel("Ready")

	w.tree = tree
	w.editor = editor
	w.status = statusLabel

	statusBar := container.NewHBox(statusLabel)
	center := container.NewHSplit(container.NewBorder(nil, nil, nil, nil, tree),
		container.New(layout.NewPaddedLayout(), editor))
	content := container.NewBorder(toolbar, statusBar, nil, nil, center)

	w.win.SetContent(content)
	w.win.Resize(fyne.NewSize(800, 600))
	w.win.CenterOnScreen()

}

func (w *Window) createToolbar() *widget.Toolbar {
	toolbar := widget.NewToolbar(widget.NewToolbarAction(theme.DocumentIcon(), func() {
		input := widget.NewEntry()
		input.PlaceHolder = "Please input hosts item Name"
		input.Validator = func(s string) error {
			return validateName(s, w.hosts)
		}
		dlg := dialog.NewForm("New Hosts Item", "Ok", "Cancel", []*widget.FormItem{widget.NewFormItem("GetName", input)}, func(b bool) {
			if b {
				w.hosts.Add(NewHostsItem(input.Text))
				w.tree.Refresh()
				w.showStatus("Create success!")
			}
		}, w.win)
		dlg.Resize(fyne.NewSize(300, 100))
		dlg.Show()
	}),
		widget.NewToolbarAction(theme.FolderIcon(), func() {
			input := widget.NewEntry()
			input.PlaceHolder = "Please input hosts group Name"
			input.Validator = func(s string) error {
				return validateName(s, w.hosts)
			}
			dlg := dialog.NewForm("New Hosts GetGroup", "Ok", "Cancel", []*widget.FormItem{widget.NewFormItem("GetName", input)}, func(b bool) {
				if b {
					w.hosts.Add(NewHostsGroup(input.Text))
					w.tree.Refresh()
					w.showStatus("Create success!")
				}
			}, w.win)

			dlg.Resize(fyne.NewSize(300, 100))
			dlg.Show()
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			w.save()
		}),
	)
	return toolbar
}

func (w *Window) createEditor() *widget.Entry {
	editor := widget.NewMultiLineEntry()

	return editor
}

func (w *Window) showStatus(status string) {
	w.status.SetText(status)
}
