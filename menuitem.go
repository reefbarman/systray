package winsystray

// MenuItem represents an item displayed in the root or a sub menu of the tray application
// It can be disabled, checked or have the title updated
type MenuItem struct {
	id       int32
	title    string
	checked  bool
	disabled bool
	onClick  func(*MenuItem)
	parent   *Menu
}

// GetID will return the unique id of this menu item
func (m MenuItem) GetID() int32 {
	return m.id
}

// SetTitle allows the updating of the items title
func (m *MenuItem) SetTitle(title string) {
	m.title = title
	setMenuItem(m)
}

// GetTitle allows retrieving the current title
func (m MenuItem) GetTitle() string {
	return m.title
}

// ToogleChecked will switch the checked state on the item
func (m *MenuItem) ToogleChecked() {
	m.checked = !m.checked
	setMenuItem(m)
}

// IsChecked allows checking the checked state of the item
func (m MenuItem) IsChecked() bool {
	return m.checked
}

// ToggleDisabled will switch the disabled state on the item
func (m *MenuItem) ToggleDisabled() {
	m.disabled = !m.disabled
	setMenuItem(m)
}

// IsDisabled will allow the checking of the disabled state of the item
func (m MenuItem) IsDisabled() bool {
	return m.disabled
}
