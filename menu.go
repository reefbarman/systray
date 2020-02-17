package winsystray

import (
	"sync/atomic"
	"winsystray/interfaces"
)

// Menu represents the top level or sub menus of a tray application
type Menu struct {
	handle uintptr
	items  []interfaces.MenuItem
}

// AddSeperator will add a seperator to the menu
func (m *Menu) AddSeperator() {
	id := atomic.AddInt32(&currentID, 1)
	menuItem := &MenuItem{
		id:     id,
		parent: m,
	}

	addSeperator(menuItem)
}

// AddMenuItem will add an item to the menu
func (m *Menu) AddMenuItem(title string, onClick func(*MenuItem)) *MenuItem {
	menuItem := createMenuItem(title, m)
	menuItem.onClick = onClick

	setMenuItem(menuItem)
	return menuItem
}

// AddSubMenuItem will add a sub menu to the menu
func (m *Menu) AddSubMenuItem(title string) *Menu {
	menuItem := createMenuItem(title, m)
	return addSubMenuItem(menuItem)
}

// GetHandle will return the platform specific pointer to the raw menu resource
func (m Menu) GetHandle() uintptr {
	return m.handle
}

// AddNewMenuItem implements the interfaces.Menu interface allowing platform specific code to keep track of menu positions when adding or updating a menu
func (m *Menu) AddNewMenuItem(menuItem interfaces.MenuItem) int32 {
	for i, v := range m.items {
		if v == menuItem {
			return int32(i)
		}
	}

	m.items = append(m.items, menuItem)
	return int32(len(m.items) - 1)
}

// RemoveMenuItem implements the interfaces.Menu interface allowing platform specific code to remove items from the menu in their current positions
func (m *Menu) RemoveMenuItem(menuItem interfaces.MenuItem) {
	indexToRemove := -1

	for i, v := range m.items {
		if v == menuItem {
			indexToRemove = i
			break
		}
	}

	if indexToRemove >= 0 {
		m.items = append(m.items[:indexToRemove], m.items[indexToRemove+1:]...)
	}
}
