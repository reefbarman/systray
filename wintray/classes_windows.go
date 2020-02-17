package wintray

import (
	"unsafe"

	"github.com/reefbarman/systray/win32"

	"golang.org/x/sys/windows"
)

// Contains information that the system needs to display notifications in the notification area.
// Used by Shell_NotifyIcon.
// https://msdn.microsoft.com/en-us/library/windows/desktop/bb773352(v=vs.85).aspx
// https://msdn.microsoft.com/en-us/library/windows/desktop/bb762159
type notifyIconData struct {
	Size            uint32
	Wnd             windows.Handle
	ID              uint32
	Flags           uint32
	CallbackMessage uint32
	Icon            windows.Handle
	Tip             [128]uint16
	State           uint32
	StateMask       uint32
	Info            [256]uint16
	Timeout         uint32
	Version         uint32
	InfoTitle       [64]uint16
	InfoFlags       uint32
	GUIDItem        windows.GUID
	BalloonIcon     windows.Handle
}

func (nid *notifyIconData) add() error {
	res, _, err := win32.ShellNotifyIcon.Call(
		uintptr(win32.NIM_ADD),
		uintptr(unsafe.Pointer(nid)),
	)
	if res == 0 {
		return err
	}
	return nil
}

func (nid *notifyIconData) modify() error {
	res, _, err := win32.ShellNotifyIcon.Call(
		uintptr(win32.NIM_MODIFY),
		uintptr(unsafe.Pointer(nid)),
	)
	if res == 0 {
		return err
	}
	return nil
}

func (nid *notifyIconData) delete() error {
	res, _, err := win32.ShellNotifyIcon.Call(
		uintptr(win32.NIM_DELETE),
		uintptr(unsafe.Pointer(nid)),
	)
	if res == 0 {
		return err
	}
	return nil
}

// Contains window class information.
// It is used with the RegisterClassEx and GetClassInfoEx functions.
// https://msdn.microsoft.com/en-us/library/ms633577.aspx
type wndClassEx struct {
	Size, Style uint32
	WndProc     uintptr
	ClsExtra    int32
	WndExtra    int32
	Instance    windows.Handle
	Icon        windows.Handle
	Cursor      windows.Handle
	Background  windows.Handle
	MenuName    *uint16
	ClassName   *uint16
	IconSm      windows.Handle
}

// Registers a window class for subsequent use in calls to the CreateWindow or CreateWindowEx function.
// https://msdn.microsoft.com/en-us/library/ms633587.aspx
func (w *wndClassEx) register() error {
	w.Size = uint32(unsafe.Sizeof(*w))
	res, _, err := win32.RegisterClass.Call(uintptr(unsafe.Pointer(w)))
	if res == 0 {
		return err
	}
	return nil
}

// Unregisters a window class, freeing the memory required for the class.
// https://msdn.microsoft.com/en-us/library/ms644899.aspx
func (w *wndClassEx) unregister() error {
	res, _, err := win32.UnregisterClass.Call(
		uintptr(unsafe.Pointer(w.ClassName)),
		uintptr(w.Instance),
	)
	if res == 0 {
		return err
	}
	return nil
}

// Contains information about a menu item.
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms647578(v=vs.85).aspx
type menuItemInfo struct {
	Size      uint32
	Mask      uint32
	Type      uint32
	State     uint32
	ID        uint32
	SubMenu   windows.Handle
	Checked   windows.Handle
	Unchecked windows.Handle
	ItemData  uintptr
	TypeData  *uint16
	Cch       uint32
	Item      windows.Handle
}
