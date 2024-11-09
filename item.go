package postman

type ItemList []*Item

func (b ItemList) FindItem(name string) *Item {
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
	Item        ItemList `json:"item,omitempty"`
}

func NewItem(name string) *Item {
	return &Item{
		Name: name,
		Item: ItemList{},
	}
}

func CreateFolder(itemList *ItemList, path string, folderOpts *FolderOptions) *Item {
	slicedPath := SlicePath(path)
	if len(slicedPath) == 0 {
		return nil
	}

	root := itemList.FindItem(slicedPath[0])
	if root == nil {
		root = NewItem(slicedPath[0])
		if folderOpts != nil {
			root.Description = folderOpts.Description
		}
		*itemList = append(*itemList, root)
	}

	for _, value := range slicedPath[1:] {
		child := root.Item.FindItem(value)
		if child == nil {
			child = NewItem(value)
			if folderOpts != nil {
				child.Description = folderOpts.Description
			}
			root.Item = append(root.Item, child)
		}
		root = child
	}
	return root
}
