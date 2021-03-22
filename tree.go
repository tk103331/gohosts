package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func createNode() *fyne.Container {
	check := widget.NewCheck("", nil)
	label := widget.NewLabel("")
	add := widget.NewButtonWithIcon("", theme.ContentAddIcon(), nil)
	add.Hidden = true
	//edit := widget.NewButtonWithIcon("", theme.SettingsIcon(), nil)
	del := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
	return container.NewHBox(check, label, layout.NewSpacer(), add, del)
}

func updateNode(win *Window, box *fyne.Container, name string) {
	hosts := win.hosts.Item(name)
	check := box.Objects[0]
	label := box.Objects[1]
	add := box.Objects[3]
	del := box.Objects[4]

	if name == "System" || name == "Backup" {
		check.Hide()
		add.Hide()
		del.Hide()
	}

	check.(*widget.Check).Checked = hosts.IsEnable()
	check.(*widget.Check).OnChanged = func(b bool) {
		hosts.SetEnable(b)
		win.save()
	}
	label.(*widget.Label).Text = hosts.Name()
	add.(*widget.Button).Hidden = !hosts.IsGroup()
	add.(*widget.Button).OnTapped = func() {
		group, ok := hosts.(*HostsGroup)
		if !ok {
			return
		}
		entry := widget.NewEntry()
		dialog.NewForm("Add Hosts Item to group ["+group.Name()+"]", "Ok", "Cancel", []*widget.FormItem{widget.NewFormItem("Name", entry)}, func(b bool) {
			if b {
				group.Add(NewHostsItem(entry.Text))
				win.tree.Refresh()
			}
		}, win.win).Show()
	}
	del.(*widget.Button).OnTapped = func() {
		dialog.NewConfirm("Confirm", "Whether to delete the hosts group", func(b bool) {
			if b {
				win.hosts.Remove(hosts.Name())
			}
		}, win.win).Show()
	}

	if name == "System" || name == "Backup" {
		check.Hide()
		add.Hide()
		del.Hide()
	}
}

func (w *Window) createTree() *widget.Tree {

	tree := widget.NewTree(func(id widget.TreeNodeID) []widget.TreeNodeID {
		h := w.hosts.Item(id)
		if h != nil {
			if h.IsGroup() {
				return h.(*HostsGroup).ItemNames()
			}
		}
		return nil
	}, func(id widget.TreeNodeID) bool {
		h := w.hosts.Item(id)
		if h != nil {
			return h.IsGroup()
		}
		return false
	}, func(b bool) fyne.CanvasObject {
		return createNode()
	}, func(id widget.TreeNodeID, b bool, object fyne.CanvasObject) {
		updateNode(w, object.(*fyne.Container), id)
	})

	tree.OnSelected = func(id widget.TreeNodeID) {
		if id == "System" {
			w.editor.SetText(loadSystem())
			w.editor.Disable()
		} else if id == "Backup" {
			w.editor.SetText(loadBackup())
			w.editor.Disable()
		} else {
			hosts := w.hosts.Item(id)
			w.current = hosts
			if hosts != nil {
				w.editor.SetText(hosts.Content())
				w.editor.Enable()
			}
		}
	}
	return tree
}
