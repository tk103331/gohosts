package main

func NewHostsItem(name string) *Hosts {
	return &Hosts{Name: name, IsGroup: false}
}

func NewHostsGroup(name string) *Hosts {
	group := &Hosts{Name: name, IsGroup: true}
	group.Items = make([]*Hosts,0)
	return group
}


type Hosts struct {
	Name    string
	Content string
	Enable  bool
	Items     []*Hosts
	Exclusive bool
	IsGroup bool
	parent   *Hosts
}

func (h *Hosts) Parent() *Hosts {
	return h.parent
}

func (h *Hosts) GetContent() string {
	if h.IsGroup {
		content := "#[group]" + h.Name
		for _, item := range h.Items {
			if item.Enable {
				content = content + "\n#" + item.Name + "\n" + item.GetContent()
			}
		}
		return content
	} else {
		return h.Content
	}
}

func (h *Hosts) Add(item *Hosts) {
	item.parent = h
	h.Items = append(h.Items, item)
}

func (h *Hosts) RemoveIndex(index int) *Hosts {
	if index < 0 || index >= len(h.Items) {
		return nil
	}
	item := h.Items[index]
	h.Items = append(h.Items[:index], h.Items[index+1:]...)
	return item
}

func (h *Hosts) Remove(name string) *Hosts {
	index := -1
	for i, item := range h.Items {
		if item.Name == name {
			index = i
			break
		}
	}
	if index != -1 {
		return h.RemoveIndex(index)
	}
	return nil
}

func (h *Hosts) ItemNames() []string {
	names := make([]string, len(h.Items))
	for i, it := range h.Items {
		names[i] = it.Name
	}
	return names
}

func (h *Hosts) ItemIndex(index int) *Hosts {
	item := h.Items[index]
	return item
}

func (h *Hosts) Item(name string) *Hosts {
	if name == "" {
		return h
	}
	for _, it := range h.Items {
		if it.Name == name {
			return it
		} else if it.IsGroup {
			item := it.Item(name)
			if item != nil {
				return item
			}
		}
	}
	return nil
}
