package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func createNode(bracnh bool) *fyne.Container {
	check := widget.NewCheck("", nil)
	icon := widget.NewIcon(theme.DocumentIcon())
	if bracnh {
		icon = widget.NewIcon(theme.FolderIcon())
	}
	label := widget.NewLabel("")
	add := widget.NewButtonWithIcon("", theme.ContentAddIcon(), nil)
	add.Hidden = true
	//edit := widget.NewButtonWithIcon("", theme.SettingsIcon(), nil)
	del := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)

	return container.NewHBox(check, icon, label, layout.NewSpacer(), add, del)
}

func updateNode(win *Window, box *fyne.Container, name string) {
	hosts := win.hosts.Item(name)
	check := box.Objects[0]
	label := box.Objects[2]
	add := box.Objects[4]
	del := box.Objects[5]

	label.(*widget.Label).Text = hosts.GetName()
	if name == "System" {
		check.Hide()
		add.Hide()
		del.Hide()
		return
	}

	if name == "Backup" {
		del.Hide()
	}

	check.(*widget.Check).Checked = hosts.IsEnable()
	check.(*widget.Check).OnChanged = func(b bool) {
		hosts.SetEnable(b)
		win.save()
	}
	add.(*widget.Button).Hidden = !hosts.IsGroup()
	add.(*widget.Button).OnTapped = func() {
		group, ok := hosts.(*HostsGroup)
		if !ok {
			return
		}
		entry := widget.NewEntry()
		entry.PlaceHolder = "Please input hosts item Name"
		entry.Validator = func(s string) error {
			return validateName(s, group)
		}
		dialog.NewForm("Add Hosts Item to group ["+group.GetName()+"]", "Ok", "Cancel", []*widget.FormItem{widget.NewFormItem("GetName", entry)}, func(b bool) {
			if b {

				group.Add(NewHostsItem(entry.Text))
				win.tree.Refresh()
				win.showStatus("Create success!")
			}
		}, win.win).Show()
	}
	del.(*widget.Button).OnTapped = func() {
		info := "Confirm to delete the hosts item"
		if hosts.IsGroup() {
			info = "Confirm to delete the hosts group"
		}
		dialog.NewConfirm("Confirm", info, func(b bool) {
			if b {
				hosts.GetGroup().Remove(hosts.GetName())
				win.tree.Refresh()
				win.showStatus("Remove success")
			}
		}, win.win).Show()
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
	}, func(branch bool) fyne.CanvasObject {
		return createNode(branch)
	}, func(id widget.TreeNodeID, b bool, object fyne.CanvasObject) {
		updateNode(w, object.(*fyne.Container), id)
	})

	tree.OnSelected = func(id widget.TreeNodeID) {
		w.editor.Disable()
		w.editor.OnChanged = func(s string) {}
		if id == "System" {
			w.editor.SetText(loadSystem())
		} else {
			hosts := w.hosts.Item(id)
			w.current = hosts
			if hosts != nil {
				w.editor.SetText(hosts.GetContent())
				if !hosts.IsGroup() {
					w.editor.Enable()
					w.editor.OnChanged = func(s string) {
						hosts.SetContent(s)
					}
				}
			}
		}
		w.showStatus("Current: " + id)
	}

	return tree
}
