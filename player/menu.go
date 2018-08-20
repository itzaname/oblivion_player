package player

type SubItem struct {
	Name   string
	Object interface{}
	Selected bool
}

type Menu struct {
	Items    []SubItem
	Name     string
	Callback func(interface{})
	Active   bool
}

func (m *Menu) GetActiveItemName() string {
	for _,item := range m.Items {
		if item.Selected {
			return item.Name
		}
	}
	return "Kek"
}

func (m *Menu) GetActiveItem() SubItem {
	for _,item := range m.Items {
		if item.Selected {
			return item
		}
	}
	return SubItem{"Error", nil, false}
}

func (m *Menu) SetActiveItem(index int) {
	for i := 0; i < len(m.Items); i++ {
		if i == index {
			m.Items[i].Selected = true
		} else {
			m.Items[i].Selected = false
		}
	}
}