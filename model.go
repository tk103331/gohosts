package main

import "strings"

func NewHostsItem(name string) *HostsItem {
	return &HostsItem{name: name}
}

func NewHostsGroup(name string) *HostsGroup {
	return &HostsGroup{name: name}
}

type Hosts interface {
	Name() string
	Content() string
	SetContent(string)
	IsEnable() bool
	SetEnable(b bool)
	IsGroup() bool
}

type HostsItem struct {
	name    string
	content string
	enable  bool
}

func (i *HostsItem) Name() string {
	return i.name
}

func (i *HostsItem) Content() string {
	return i.content
}

func (i *HostsItem) SetContent(content string) {
	i.content = content
}

func (i *HostsItem) IsEnable() bool {
	return i.enable
}

func (i *HostsItem) SetEnable(b bool) {
	i.enable = b
}

func (i *HostsItem) IsGroup() bool {
	return false
}

type HostsGroup struct {
	name   string
	items  []Hosts
	enable bool
}

func (g *HostsGroup) Name() string {
	return g.name
}

func (g *HostsGroup) Content() string {
	content := ""
	for _, item := range g.items {
		if item.IsEnable() {
			content = content + "\n#" + item.Name() + "\n" + item.Content()
		}
	}
	return content
}

func (g *HostsGroup) SetContent(string) {

}

func (g *HostsGroup) IsEnable() bool {
	return g.enable
}

func (i *HostsGroup) SetEnable(b bool) {
	i.enable = b
}

func (g *HostsGroup) IsGroup() bool {
	return true
}

func (g *HostsGroup) Add(item Hosts) {
	g.items = append(g.items, item)
}

func (g *HostsGroup) RemoveIndex(index int) Hosts {
	item := g.items[index]
	g.items = append(g.items[:index], g.items[index:]...)
	return item
}

func (g *HostsGroup) Remove(name string) Hosts {
	index := -1
	for i, item := range g.items {
		if item.Name() == name {
			index = i
			break
		}
	}
	if index != -1 {
		return g.RemoveIndex(index)
	}
	return nil
}

func (g *HostsGroup) ItemNames() []string {
	names := make([]string, len(g.items))
	for i, it := range g.items {
		if g.name == "" {
			names[i] = it.Name()
		} else {
			names[i] = g.name + "." + it.Name()
		}
	}
	return names
}

func (g *HostsGroup) ItemIndex(index int) Hosts {
	item := g.items[index]
	return item
}

func (g *HostsGroup) Item(name string) Hosts {
	if name == "" {
		return g
	}
	strs := strings.SplitN(name, ".", 2)
	for _, it := range g.items {
		if it.Name() == strs[0] {
			if len(strs) == 1 {
				return it
			} else if group, ok := it.(*HostsGroup); it.IsGroup() && ok {
				return group.Item(strs[1])
			}
		}
	}
	return nil
}
