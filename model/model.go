package model

func NewHostsItem(name string) *HostsItem {
	return &HostsItem{name: name}
}

func NewHostsGroup(name string) *HostsGroup {
	return &HostsGroup{name: name}
}

type Hosts interface {
	Name() string
	Content() string
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
	return "#######" + i.name + "######\n" + i.content + "\n"
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
	items  []*HostsItem
	enable bool
}

func (g *HostsGroup) Name() string {
	return g.name
}

func (g *HostsGroup) Content() string {
	content := ""
	for _, item := range g.items {
		content = content + item.Content()
	}
	return content
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

func (g *HostsGroup) Add(item *HostsItem) {
	g.items = append(g.items, item)
}

func (g *HostsGroup) RemoveIndex(index int) *HostsItem {
	item := g.items[index]
	g.items = append(g.items[:index], g.items[index:]...)
	return item
}

func (g *HostsGroup) Remove(name string) *HostsItem {
	index := -1
	for i, item := range g.items {
		if item.name == name {
			index = i
			break
		}
	}
	if index != -1 {
		return g.RemoveIndex(index)
	}
	return nil
}
