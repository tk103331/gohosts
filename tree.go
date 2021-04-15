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
	holder := widget.NewLabel("")
	icon := widget.NewIcon(theme.DocumentIcon())
	if bracnh {
		icon = widget.NewIcon(theme.FolderIcon())
	}
	label := widget.NewLabel("")
	add := widget.NewButtonWithIcon("", theme.ContentAddIcon(), nil)
	add.Hidden = true
	//edit := widget.NewButtonWithIcon("", theme.SettingsIcon(), nil)
	del := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)

	return container.NewHBox(holder, icon, label, layout.NewSpacer(), add, del)
}

func updateNode(win *Window, box *fyne.Container, name string) {
	hosts := win.root.Item(name)
	if hosts == nil {
		return
	}
	label := box.Objects[2]
	add := box.Objects[4]
	del := box.Objects[5]


	if hosts.Parent().Exclusive {
		radio := widget.NewRadioGroup([]string{""}, nil)
		radio.Horizontal = true
		box.Objects[0] = radio
		radio.OnChanged = func(s string) {
			for _, it := range hosts.Parent().Items {
				it.Enable = false
				it.radio.Selected = "nil"
				it.radio.Refresh()
			}
			hosts.Enable = s == ""
			radio.Selected = s
			win.refreshEditor()
		}
		if hosts.Enable {
			radio.Selected = ""
		} else {
			radio.Selected = "nil"
		}

		hosts.radio = radio

	} else {
		check := widget.NewCheck("", nil)
		box.Objects[0] = check
		check.Checked = hosts.Enable
		check.OnChanged = func(b bool) {
			hosts.Enable = b
			win.refreshEditor()
		}

	}
	label.(*widget.Label).Text = hosts.Name
	if name == "System" {
		box.Objects[0].Hide()
		add.Hide()
		del.Hide()
		return
	}

	if name == "Backup" {
		del.Hide()
	}
	add.(*widget.Button).Hidden = !hosts.IsGroup
	add.(*widget.Button).OnTapped = func() {
		if !hosts.IsGroup {
			return
		}
		entry := widget.NewEntry()
		entry.PlaceHolder = "Please input hosts item Name"
		entry.Validator = func(s string) error {
			return validateName(s, hosts)
		}
		dialog.NewForm("Add Hosts Item to group ["+hosts.Name+"]", "Ok", "Cancel", []*widget.FormItem{widget.NewFormItem("Name", entry)}, func(b bool) {
			if b {
				hosts.Add(NewHostsItem(entry.Text))
				win.tree.Refresh()
				win.showStatus("Create success!")
			}
		}, win.win).Show()
	}
	del.(*widget.Button).OnTapped = func() {
		info := "Confirm to delete the hosts item"
		if hosts.IsGroup {
			info = "Confirm to delete the hosts group"
		}
		dialog.NewConfirm("Confirm", info, func(b bool) {
			if b {
				hosts.Parent().Remove(hosts.Name)
				win.tree.Refresh()
				win.showStatus("Remove success")
			}
		}, win.win).Show()
	}
}

func (w *Window) createTree() *widget.Tree {

	tree := widget.NewTree(func(id widget.TreeNodeID) []widget.TreeNodeID {
		h := w.root.Item(id)
		if h != nil {
			if h.IsGroup {
				return h.ItemNames()
			}
		}
		return nil
	}, func(id widget.TreeNodeID) bool {
		h := w.root.Item(id)
		if h != nil {
			return h.IsGroup
		}
		return false
	}, func(branch bool) fyne.CanvasObject {
		return createNode(branch)
	}, func(id widget.TreeNodeID, b bool, object fyne.CanvasObject) {
		updateNode(w, object.(*fyne.Container), id)
	})

	tree.OnSelected = func(id widget.TreeNodeID) {
		w.current = id
		w.editor.Disable()
		w.editor.OnChanged = func(s string) {}
		if id == "System" {
			w.editor.SetText(loadSystem())
		} else {
			hosts := w.root.Item(id)
			if hosts != nil {
				w.editor.SetText(hosts.GetContent())
				if !hosts.IsGroup {
					w.editor.Enable()
					w.editor.OnChanged = func(s string) {
						hosts.Content = s
					}
				}
			}
		}
		w.showStatus("Current: " + id)
	}

	return tree
}
