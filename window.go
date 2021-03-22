package main

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
)

type Window struct {
	app fyne.App
	win fyne.Window

	tree   *widget.Tree
	editor *widget.Entry

	current Hosts
	hosts   HostsGroup
}

func Run() {
	myApp := app.New()
	win := myApp.NewWindow("Go Hosts!")

	root := HostsGroup{items: make([]Hosts, 0)}
	data := myApp.Preferences().String("hosts")

	_ = json.Unmarshal([]byte(data), &root.items)
	root.Add(NewHostsItem("System"))
	root.Add(NewHostsItem("Backup"))
	(&Window{app: myApp, win: win, hosts: root}).Run()
}

func (w *Window) Run() {

	w.init()
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
	content := w.hosts.Content()
	err = saveSystem(content)
	if err != nil {
		log.Println(err)
		dialog.NewInformation("Error", "Saving system hosts file error!\n"+err.Error(), w.win)
		return
	}
}

func (w *Window) init() {

	toolbar := w.createToolbar()
	editor := w.createEditor()
	tree := w.createTree()

	w.tree = tree
	w.editor = editor

	statusBar := container.NewHBox(widget.NewLabel(" "))
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
		input.PlaceHolder = "Please input hosts item name"
		dlg := dialog.NewForm("New Hosts Item", "Ok", "Cancel", []*widget.FormItem{widget.NewFormItem("Name", input)}, func(b bool) {
			if b {
				w.hosts.Add(NewHostsItem(input.Text))
				w.tree.Refresh()
			}
		}, w.win)
		dlg.Resize(fyne.NewSize(300, 100))
		dlg.Show()
	}),
		widget.NewToolbarAction(theme.FolderIcon(), func() {
			input := widget.NewEntry()
			input.PlaceHolder = "Please input hosts group name"
			dlg := dialog.NewForm("New Hosts Group", "Ok", "Cancel", []*widget.FormItem{widget.NewFormItem("Name", input)}, func(b bool) {
				if b {
					w.hosts.Add(NewHostsGroup(input.Text))
					w.tree.Refresh()
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
	editor.OnChanged = func(s string) {
		if w.current != nil {
			item, ok := w.current.(*HostsItem)
			if ok && !item.IsGroup() && item.Name() != "System" && item.Name() != "Backup" {
				item.SetContent(w.editor.Text)
			}
		}
	}
	return editor
}
