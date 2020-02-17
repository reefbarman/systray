package wintray

import (
	"unsafe"

	"github.com/reefbarman/systray/interfaces"
	"github.com/reefbarman/systray/win32"

	"golang.org/x/sys/windows"
)

type WinTray struct {
	OnTrayMenuOpened   func()
	OnMenuItemSelected func(menuId int32)
	OnExit             func()

	instance         windows.Handle
	icon             windows.Handle
	cursor           windows.Handle
	window           windows.Handle
	loadedImages     map[string]windows.Handle
	nid              *notifyIconData
	wcex             *wndClassEx
	wmSystrayMessage uint32
	wmTaskbarCreated uint32
	visibleItems     []uint32
}

func (t *WinTray) InitInstance() error {
	const (
		className  = "SystrayClass"
		windowName = ""
	)

	t.wmSystrayMessage = win32.WM_USER + 1

	taskbarEventNamePtr, _ := windows.UTF16PtrFromString("TaskbarCreated")
	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms644947
	res, _, err := win32.RegisterWindowMessage.Call(
		uintptr(unsafe.Pointer(taskbarEventNamePtr)),
	)
	t.wmTaskbarCreated = uint32(res)

	t.loadedImages = make(map[string]windows.Handle)

	instanceHandle, _, err := win32.GetModuleHandle.Call(0)
	if instanceHandle == 0 {
		return err
	}
	t.instance = windows.Handle(instanceHandle)

	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms648072(v=vs.85).aspx
	iconHandle, _, err := win32.LoadIcon.Call(0, uintptr(win32.IDI_APPLICATION))
	if iconHandle == 0 {
		return err
	}
	t.icon = windows.Handle(iconHandle)

	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms648391(v=vs.85).aspx
	cursorHandle, _, err := win32.LoadCursor.Call(0, uintptr(win32.IDC_ARROW))
	if cursorHandle == 0 {
		return err
	}
	t.cursor = windows.Handle(cursorHandle)

	classNamePtr, err := windows.UTF16PtrFromString(className)
	if err != nil {
		return err
	}

	windowNamePtr, err := windows.UTF16PtrFromString(windowName)
	if err != nil {
		return err
	}

	t.wcex = &wndClassEx{
		Style:      win32.CS_HREDRAW | win32.CS_VREDRAW,
		WndProc:    windows.NewCallback(t.wndProc),
		Instance:   t.instance,
		Icon:       t.icon,
		Cursor:     t.cursor,
		Background: windows.Handle(6), // (COLOR_WINDOW + 1)
		ClassName:  classNamePtr,
		IconSm:     t.icon,
	}
	if err := t.wcex.register(); err != nil {
		return err
	}

	windowHandle, _, err := win32.CreateWindowEx.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(classNamePtr)),
		uintptr(unsafe.Pointer(windowNamePtr)),
		uintptr(win32.WS_OVERLAPPEDWINDOW),
		uintptr(win32.CW_USEDEFAULT),
		uintptr(win32.CW_USEDEFAULT),
		uintptr(win32.CW_USEDEFAULT),
		uintptr(win32.CW_USEDEFAULT),
		uintptr(0),
		uintptr(0),
		uintptr(t.instance),
		uintptr(0),
	)
	if windowHandle == 0 {
		return err
	}
	t.window = windows.Handle(windowHandle)

	win32.ShowWindow.Call(
		uintptr(t.window),
		uintptr(win32.SW_HIDE),
	)

	win32.UpdateWindow.Call(
		uintptr(t.window),
	)

	t.nid = &notifyIconData{
		Wnd:             windows.Handle(t.window),
		ID:              100,
		Flags:           win32.NIF_MESSAGE,
		CallbackMessage: t.wmSystrayMessage,
	}
	t.nid.Size = uint32(unsafe.Sizeof(*t.nid))

	return t.nid.add()
}

func (t *WinTray) DeInit() {
	win32.DestroyWindow.Call(uintptr(t.window))
	t.wcex.unregister()
}

func (t *WinTray) Quit() {
	win32.PostMessage.Call(uintptr(t.window), win32.WM_CLOSE, 0, 0)
}

func (t *WinTray) CreateMenu() (uintptr, error) {
	menuHandle, _, err := win32.CreatePopupMenu.Call()
	if menuHandle == 0 {
		return 0, err
	}
	menu := windows.Handle(menuHandle)

	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms647575(v=vs.85).aspx
	mi := struct {
		Size          uint32
		Mask          uint32
		Style         uint32
		Max           uint32
		Background    windows.Handle
		ContextHelpID uint32
		MenuData      uintptr
	}{
		Mask: win32.MIM_APPLYTOSUBMENUS,
	}
	mi.Size = uint32(unsafe.Sizeof(mi))

	res, _, err := win32.SetMenuInfo.Call(
		uintptr(menu),
		uintptr(unsafe.Pointer(&mi)),
	)
	if res == 0 {
		return 0, err
	}

	return uintptr(menu), nil
}

func (t *WinTray) ShowTrayMenu(menu interfaces.Menu) error {
	p := win32.Point{}
	res, _, err := win32.GetCursorPos.Call(uintptr(unsafe.Pointer(&p)))
	if res == 0 {
		return err
	}
	win32.SetForegroundWindow.Call(uintptr(t.window))

	res, _, err = win32.TrackPopupMenu.Call(
		uintptr(menu.GetHandle()),
		win32.TPM_BOTTOMALIGN|win32.TPM_LEFTALIGN,
		uintptr(p.X),
		uintptr(p.Y),
		0,
		uintptr(t.window),
		0,
	)
	if res == 0 {
		return err
	}

	return nil
}

// Loads an image from file and shows it in tray.
// LoadImage: https://msdn.microsoft.com/en-us/library/windows/desktop/ms648045(v=vs.85).aspx
// Shell_NotifyIcon: https://msdn.microsoft.com/en-us/library/windows/desktop/bb762159(v=vs.85).aspx
func (t *WinTray) SetIcon(src string) error {
	// Save and reuse handles of loaded images
	h, ok := t.loadedImages[src]
	if !ok {
		srcPtr, err := windows.UTF16PtrFromString(src)
		if err != nil {
			return err
		}
		res, _, err := win32.LoadImage.Call(
			0,
			uintptr(unsafe.Pointer(srcPtr)),
			win32.IMAGE_ICON,
			0,
			0,
			win32.LR_LOADFROMFILE|win32.LR_DEFAULTSIZE,
		)
		if res == 0 {
			return err
		}
		h = windows.Handle(res)
		t.loadedImages[src] = h
	}

	t.nid.Icon = h
	t.nid.Flags |= win32.NIF_ICON
	t.nid.Size = uint32(unsafe.Sizeof(*t.nid))

	return t.nid.modify()
}

// Sets tooltip on icon.
// Shell_NotifyIcon: https://msdn.microsoft.com/en-us/library/windows/desktop/bb762159(v=vs.85).aspx
func (t *WinTray) SetTooltip(src string) error {
	b, err := windows.UTF16FromString(src)
	if err != nil {
		return err
	}
	copy(t.nid.Tip[:], b[:])
	t.nid.Flags |= win32.NIF_TIP
	t.nid.Size = uint32(unsafe.Sizeof(*t.nid))

	return t.nid.modify()
}

func (t *WinTray) SetMenuItem(menuItem interfaces.MenuItem, parentMenu interfaces.Menu) error {
	titlePtr, err := windows.UTF16PtrFromString(menuItem.GetTitle())
	if err != nil {
		return err
	}

	mi := menuItemInfo{
		Mask:     win32.MIIM_FTYPE | win32.MIIM_STRING | win32.MIIM_ID | win32.MIIM_STATE,
		Type:     win32.MFT_STRING,
		ID:       uint32(menuItem.GetID()),
		TypeData: titlePtr,
		Cch:      uint32(len(menuItem.GetTitle())),
	}
	if menuItem.IsDisabled() {
		mi.State |= win32.MFS_DISABLED
	}
	if menuItem.IsChecked() {
		mi.State |= win32.MFS_CHECKED
	}
	mi.Size = uint32(unsafe.Sizeof(mi))

	// We set the menu item info based on the menuID
	res, _, err := win32.SetMenuItemInfo.Call(
		uintptr(parentMenu.GetHandle()),
		uintptr(menuItem.GetID()),
		0,
		uintptr(unsafe.Pointer(&mi)),
	)

	position := parentMenu.AddNewMenuItem(menuItem)

	if res == 0 {
		res, _, err = win32.InsertMenuItem.Call(
			uintptr(parentMenu.GetHandle()),
			uintptr(position),
			1,
			uintptr(unsafe.Pointer(&mi)),
		)
		if res == 0 {
			parentMenu.RemoveMenuItem(menuItem)
			return err
		}
	}

	return nil
}

func (t *WinTray) AddSeparator(menuItem interfaces.MenuItem, parentMenu interfaces.Menu) error {
	mi := menuItemInfo{
		Mask: win32.MIIM_FTYPE | win32.MIIM_ID | win32.MIIM_STATE,
		Type: win32.MFT_SEPARATOR,
		ID:   uint32(menuItem.GetID()),
	}

	mi.Size = uint32(unsafe.Sizeof(mi))

	position := parentMenu.AddNewMenuItem(menuItem)

	res, _, err := win32.InsertMenuItem.Call(
		uintptr(parentMenu.GetHandle()),
		uintptr(position),
		1,
		uintptr(unsafe.Pointer(&mi)),
	)
	if res == 0 {
		parentMenu.RemoveMenuItem(menuItem)
		return err
	}

	return nil
}

func (t *WinTray) AddSubMenuItem(menuItem interfaces.MenuItem, parentMenu interfaces.Menu) (uintptr, error) {
	titlePtr, err := windows.UTF16PtrFromString(menuItem.GetTitle())
	if err != nil {
		return 0, err
	}

	subMenuHandle, err := t.CreateMenu()
	if err != nil {
		return 0, err
	}

	mi := menuItemInfo{
		Mask:     win32.MIIM_FTYPE | win32.MIIM_STRING | win32.MIIM_ID | win32.MIIM_SUBMENU,
		Type:     win32.MFT_STRING,
		ID:       uint32(menuItem.GetID()),
		TypeData: titlePtr,
		Cch:      uint32(len(menuItem.GetTitle())),
		SubMenu:  windows.Handle(subMenuHandle),
	}
	mi.Size = uint32(unsafe.Sizeof(mi))

	// We set the menu item info based on the menuID
	res, _, err := win32.SetMenuItemInfo.Call(
		uintptr(parentMenu.GetHandle()),
		uintptr(menuItem.GetID()),
		0,
		uintptr(unsafe.Pointer(&mi)),
	)

	position := parentMenu.AddNewMenuItem(menuItem)

	if res == 0 {
		res, _, err = win32.InsertMenuItem.Call(
			uintptr(parentMenu.GetHandle()),
			uintptr(position),
			1,
			uintptr(unsafe.Pointer(&mi)),
		)
		if res == 0 {
			parentMenu.RemoveMenuItem(menuItem)
			return 0, err
		}
	}

	return subMenuHandle, nil
}

// WindowProc callback function that processes messages sent to a window.
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms633573(v=vs.85).aspx
func (t *WinTray) wndProc(hWnd windows.Handle, message uint32, wParam, lParam uintptr) (lResult uintptr) {
	switch message {
	case win32.WM_COMMAND:
		menuId := int32(wParam)
		if menuId != -1 {
			t.OnMenuItemSelected(menuId)
		}
	case win32.WM_DESTROY:
		// same as WM_ENDSESSION, but throws 0 exit code after all
		defer win32.PostQuitMessage.Call(uintptr(int32(0)))
		fallthrough
	case win32.WM_ENDSESSION:
		if t.nid != nil {
			t.nid.delete()
		}
		t.OnExit()
	case t.wmSystrayMessage:
		switch lParam {
		case win32.WM_RBUTTONUP, win32.WM_LBUTTONUP:
			t.OnTrayMenuOpened()
		}
	case t.wmTaskbarCreated: // on explorer.exe restarts
		t.nid.add()
	default:
		// Calls the default window procedure to provide default processing for any window messages that an application does not process.
		// https://msdn.microsoft.com/en-us/library/windows/desktop/ms633572(v=vs.85).aspx
		lResult, _, _ = win32.DefWindowProc.Call(
			uintptr(hWnd),
			uintptr(message),
			uintptr(wParam),
			uintptr(lParam),
		)
	}
	return
}
