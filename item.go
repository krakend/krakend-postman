package postman

type Branch []*Item

func (b Branch) FindByPath(path string) *Item {
	slicedPath := SlicePath(path)

	current := b.FindItem(slicedPath[0])
	for _, value := range slicedPath[1:] {
		child := current.FindChild(value)
		if child == nil {
			return nil
		}
		current = child
	}

	return current
}

func (b Branch) FindItem(name string) *Item {
	for _, current := range b {
		if current.Name == name {
			return current
		}
	}
	return nil
}

type Item struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Request     *Request `json:"request,omitempty"`
	Item        Branch   `json:"item,omitempty"`
}

func NewItem(name string) *Item {
	return &Item{
		Name: name,
		Item: Branch{},
	}
}

func (i *Item) FindChild(name string) *Item {
	for _, child := range i.Item {
		if child.Name == name {
			return child
		}
	}
	return nil
}
