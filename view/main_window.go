package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/tk103331/gohosts/model"
)

type MainWindow struct {
	app     fyne.App
	win     fyne.Window
	hosts   []model.Hosts
	editor  *widget.Entry
	list    *widget.List
	toolbar *widget.Toolbar
	current model.Hosts
}

func NewMainWindow() *MainWindow {
	app := app.New()
	win := app.NewWindow("Go Hosts!")
	hosts := []model.Hosts{
		model.NewHostsItem("System"),
		model.NewHostsItem("Backup"),
	}
	return &MainWindow{app: app, win: win, hosts: hosts}
}

func (w *MainWindow) ShowAndRun() {
	w.initView()
	w.win.ShowAndRun()
}

func (w *MainWindow) initView() {

	w.initToolbar()
	w.initEditor()
	w.initHostsList()

	statusText := widget.NewLabel("StatusBar")

	statusBar := container.NewHBox(statusText)

	sideBar := container.NewBorder(nil, nil, nil, nil, w.list)
	editor := container.New(layout.NewPaddedLayout(), w.editor)
	center := container.NewHSplit(sideBar, editor)
	content := container.NewBorder(w.toolbar, statusBar, nil, nil, center)

	w.win.SetContent(content)
	w.win.Resize(fyne.NewSize(800, 600))
	w.win.CenterOnScreen()

}

func (w *MainWindow) initToolbar() {
	w.toolbar = widget.NewToolbar(widget.NewToolbarAction(theme.ContentAddIcon(), func() {
		input := widget.NewEntry()
		input.PlaceHolder = "Please input hosts item name"
		dlg := dialog.NewForm("New Hosts Item", "Ok", "Cancel", []*widget.FormItem{widget.NewFormItem("Name", input)}, func(b bool) {
			if b {
				w.hosts = append(w.hosts, model.NewHostsItem(input.Text))
			}
		}, w.win)
		dlg.Resize(fyne.NewSize(300, 100))
		dlg.Show()
	}),
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			input := widget.NewEntry()
			input.PlaceHolder = "Please input hosts group name"
			dlg := dialog.NewForm("New Hosts Group", "Ok", "Cancel", []*widget.FormItem{widget.NewFormItem("Name", input)}, func(b bool) {
				if b {
					w.hosts = append(w.hosts, model.NewHostsGroup(input.Text))
				}
			}, w.win)

			dlg.Resize(fyne.NewSize(300, 100))
			dlg.Show()
		}))
}

func (w *MainWindow) initEditor() {
	w.editor = widget.NewMultiLineEntry()
	w.editor.OnChanged = func(s string) {
		item, ok := w.current.(*model.HostsItem)
		if ok && item.Name() != "System" && item.Name() != "Backup" {
			item.SetContent(w.editor.Text)
		}
	}
}

func (w *MainWindow) initHostsList() {

	hostsList := widget.NewList(
		func() int {
			return len(w.hosts)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("Loading"))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			hosts := w.hosts[i]
			if i == 0 || i == 1 {
				o.(*fyne.Container).Objects = append([]fyne.CanvasObject{}, widget.NewLabel(hosts.Name()))
			} else if hosts.IsGroup() {
				group := hosts.(*model.HostsGroup)
				o.(*fyne.Container).Objects = append([]fyne.CanvasObject{},
					widget.NewCheck("", func(b bool) {
						group.SetEnable(b)
					}),
					widget.NewLabel(hosts.Name()),
					layout.NewSpacer(),
					widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
						input := widget.NewEntry()
						dialog.NewForm("New Hosts Item", "Ok", "Cancel", []*widget.FormItem{widget.NewFormItem("Name", input)}, func(b bool) {
							if b {
								group.Add(model.NewHostsItem(input.Text))
							}
						}, w.win).Show()
					}),
					widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
						dialog.NewConfirm("Confirm", "Whether to delete the hosts group", func(b bool) {
							if b {
								w.hosts = append(w.hosts[:i], w.hosts[i:]...)
							}
						}, w.win).Show()
					}))
			} else {
				item := hosts.(*model.HostsItem)
				o.(*fyne.Container).Objects = append([]fyne.CanvasObject{},
					widget.NewCheck("", func(b bool) {
						item.SetEnable(b)
					}),
					widget.NewLabel(hosts.Name()),
					layout.NewSpacer(),
					widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
						dialog.NewConfirm("Confirm", "Whether to delete the hosts item", func(b bool) {
							if b {
								w.hosts = append(w.hosts[:i], w.hosts[i:]...)
							}
						}, w.win).Show()
					}))
			}
		})
	hostsList.OnSelected = func(id widget.ListItemID) {
		if id == 0 {
			content := loadSystem()
			w.editor.Disabled()
			w.editor.SetText(content)
		} else if id == 1 {
			content := loadBackup()
			w.editor.Disabled()
			w.editor.SetText(content)
		} else {
			w.editor.Enable()
			hosts := w.hosts[id]
			w.editor.SetText(hosts.Content())
		}
	}
	w.list = hostsList
}
