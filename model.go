package main

import "strings"

func NewHostsItem(name string) *HostsItem {
	return &HostsItem{Name: name}
}

func NewHostsGroup(name string) *HostsGroup {
	group := &HostsGroup{}
	group.HostsItem = NewHostsItem(name)
	return group
}

type Hosts interface {
	GetName() string
	GetContent() string
	SetContent(string)
	IsEnable() bool
	SetEnable(bool)
	IsGroup() bool
	GetGroup() *HostsGroup
	SetGroup(*HostsGroup)
}

type HostsItem struct {
	Name    string
	Content string
	Enable  bool
	group   *HostsGroup
}

func (i *HostsItem) GetName() string {
	return i.Name
}

func (i *HostsItem) GetContent() string {
	return i.Content
}

func (i *HostsItem) SetContent(content string) {
	i.Content = content
}

func (i *HostsItem) IsEnable() bool {
	return i.Enable
}

func (i *HostsItem) SetEnable(b bool) {
	i.Enable = b
}

func (i *HostsItem) IsGroup() bool {
	return false
}

func (i *HostsItem) GetGroup() *HostsGroup {
	return i.group
}

func (i *HostsItem) SetGroup(group *HostsGroup) {
	i.group = group
}

type HostsGroup struct {
	*HostsItem
	Items     []Hosts
	Exclusive bool
}

func (g *HostsGroup) GetContent() string {
	content := ""
	for _, item := range g.Items {
		if item.IsEnable() {
			content = content + "\n#" + item.GetName() + "\n" + item.GetContent()
		}
	}
	return content
}

func (g *HostsGroup) SetContent(string) {

}

func (g *HostsGroup) IsGroup() bool {
	return true
}

func (g *HostsGroup) Add(item Hosts) {
	item.SetGroup(g)
	g.Items = append(g.Items, item)
}

func (g *HostsGroup) RemoveIndex(index int) Hosts {
	if index < 0 || index >= len(g.Items) {
		return nil
	}
	item := g.Items[index]
	g.Items = append(g.Items[:index], g.Items[index+1:]...)
	return item
}

func (g *HostsGroup) Remove(name string) Hosts {
	index := -1
	for i, item := range g.Items {
		if item.GetName() == name {
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
	names := make([]string, len(g.Items))
	for i, it := range g.Items {
		if g.Name == "" {
			names[i] = it.GetName()
		} else {
			names[i] = g.Name + "." + it.GetName()
		}
	}
	return names
}

func (g *HostsGroup) ItemIndex(index int) Hosts {
	item := g.Items[index]
	return item
}

func (g *HostsGroup) Item(name string) Hosts {
	if name == "" {
		return g
	}
	strs := strings.SplitN(name, ".", 2)
	for _, it := range g.Items {
		if it.GetName() == strs[0] {
			if len(strs) == 1 {
				return it
			} else if group, ok := it.(*HostsGroup); it.IsGroup() && ok {
				return group.Item(strs[1])
			}
		}
	}
	return nil
}
