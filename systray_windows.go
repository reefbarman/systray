// +build windows

package winsystray

import (
	"unsafe"
	"winsystray/win32"
	"winsystray/wintray"

	"golang.org/x/sys/windows"
)

var wt wintray.WinTray

func init() {
	wt.OnTrayMenuOpened = func() {
		wt.ShowTrayMenu(trayMenu)
	}

	wt.OnMenuItemSelected = onMenuItemSelected
	wt.OnExit = onExit
}

func quit() {
	wt.Quit()
}

func setTooltip(tooltip string) {
	if err := wt.SetTooltip(tooltip); err != nil {
		log.Errorf("Unable to set tooltip: %v", err)
	}
}

func addSeperator(menuItem *MenuItem) {
	if err := wt.AddSeparator(menuItem, menuItem.parent); err != nil {
		log.Errorf("Unable to add seperator: %v", err)
	}
}

func setMenuItem(menuItem *MenuItem) {
	if err := wt.SetMenuItem(menuItem, menuItem.parent); err != nil {
		log.Errorf("Unable to add menu item: %v", err)
	}
}

func addSubMenuItem(menuItem *MenuItem) *Menu {
	subMenuHandle, err := wt.AddSubMenuItem(menuItem, menuItem.parent)
	if err != nil {
		log.Errorf("Unable to add menu item: %v", err)
		return nil
	}

	return &Menu{handle: subMenuHandle}
}

func createMenu() (*Menu, error) {
	menuHandle, err := wt.CreateMenu()
	if err != nil {
		return nil, err
	}

	return &Menu{handle: menuHandle}, nil
}

func setIcon(iconFilePath string) {
	if err := wt.SetIcon(iconFilePath); err != nil {
		log.Errorf("Unable to set icon: %v", err)
	}
}

func nativeLoop() {
	if err := wt.InitInstance(); err != nil {
		log.Errorf("Unable to init instance: %v", err)
		return
	}

	defer func() {
		wt.DeInit()
	}()

	onTrayRun()

	// Main message pump.
	m := &struct {
		WindowHandle windows.Handle
		Message      uint32
		Wparam       uintptr
		Lparam       uintptr
		Time         uint32
		Pt           win32.Point
	}{}
	for {
		ret, _, err := win32.GetMessage.Call(uintptr(unsafe.Pointer(m)), 0, 0, 0)

		// If the function retrieves a message other than WM_QUIT, the return value is nonzero.
		// If the function retrieves the WM_QUIT message, the return value is zero.
		// If there is an error, the return value is -1
		// https://msdn.microsoft.com/en-us/library/windows/desktop/ms644936(v=vs.85).aspx
		switch int32(ret) {
		case -1:
			log.Errorf("Error at message loop: %v", err)
			return
		case 0:
			return
		default:
			win32.TranslateMessage.Call(uintptr(unsafe.Pointer(m)))
			win32.DispatchMessage.Call(uintptr(unsafe.Pointer(m)))
		}
	}
}
