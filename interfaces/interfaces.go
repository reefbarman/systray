package interfaces

type MenuItem interface {
	GetID() int32
	GetTitle() string
	IsChecked() bool
	IsDisabled() bool
}

type Menu interface {
	AddNewMenuItem(menuItem MenuItem) int32
	RemoveMenuItem(menuItem MenuItem)
	GetHandle() uintptr
}
