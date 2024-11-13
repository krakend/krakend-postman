package postman

type itemList []*Item

func (b itemList) findItem(name string) *Item {
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
	Item        itemList `json:"item,omitempty"`
}

func newItem(name string) *Item {
	return &Item{
		Name: name,
		Item: itemList{},
	}
}

func createFolder(itemList *itemList, path string, folderOpts *folderOptions) *Item {
	slicedPath := slicePath(path)
	if len(slicedPath) == 0 {
		return nil
	}

	root := itemList.findItem(slicedPath[0])
	if root == nil {
		root = newItem(slicedPath[0])
		if folderOpts != nil {
			root.Description = folderOpts.Description
		}
		*itemList = append(*itemList, root)
	}

	for _, value := range slicedPath[1:] {
		child := root.Item.findItem(value)
		if child == nil {
			child = newItem(value)
			if folderOpts != nil {
				child.Description = folderOpts.Description
			}
			root.Item = append(root.Item, child)
		}
		root = child
	}
	return root
}
