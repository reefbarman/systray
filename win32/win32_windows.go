// +build windows
package win32

import "golang.org/x/sys/windows"

// Helpful sources: https://github.com/golang/exp/blob/master/shiny/driver/internal/win32

var (
	k32 = windows.NewLazySystemDLL("Kernel32.dll")
	s32 = windows.NewLazySystemDLL("Shell32.dll")
	u32 = windows.NewLazySystemDLL("User32.dll")
)

var (
	GetModuleHandle       = k32.NewProc("GetModuleHandleW")
	ShellNotifyIcon       = s32.NewProc("Shell_NotifyIconW")
	CreatePopupMenu       = u32.NewProc("CreatePopupMenu")
	CreateWindowEx        = u32.NewProc("CreateWindowExW")
	DefWindowProc         = u32.NewProc("DefWindowProcW")
	DeleteMenu            = u32.NewProc("DeleteMenu")
	DestroyWindow         = u32.NewProc("DestroyWindow")
	DispatchMessage       = u32.NewProc("DispatchMessageW")
	GetCursorPos          = u32.NewProc("GetCursorPos")
	GetMenuItemID         = u32.NewProc("GetMenuItemID")
	GetMessage            = u32.NewProc("GetMessageW")
	InsertMenuItem        = u32.NewProc("InsertMenuItemW")
	LoadIcon              = u32.NewProc("LoadIconW")
	LoadImage             = u32.NewProc("LoadImageW")
	LoadCursor            = u32.NewProc("LoadCursorW")
	PostMessage           = u32.NewProc("PostMessageW")
	PostQuitMessage       = u32.NewProc("PostQuitMessage")
	RegisterClass         = u32.NewProc("RegisterClassExW")
	RegisterWindowMessage = u32.NewProc("RegisterWindowMessageW")
	SetForegroundWindow   = u32.NewProc("SetForegroundWindow")
	SetMenuInfo           = u32.NewProc("SetMenuInfo")
	SetMenuItemInfo       = u32.NewProc("SetMenuItemInfoW")
	ShowWindow            = u32.NewProc("ShowWindow")
	TrackPopupMenu        = u32.NewProc("TrackPopupMenu")
	TranslateMessage      = u32.NewProc("TranslateMessage")
	UnregisterClass       = u32.NewProc("UnregisterClassW")
	UpdateWindow          = u32.NewProc("UpdateWindow")
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/dd162805(v=vs.85).aspx
type Point struct {
	X int32
	Y int32
}

const IDI_APPLICATION = 32512
const IDC_ARROW = 32512 // Standard arrow
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms633548(v=vs.85).aspx
const SW_HIDE = 0
const CW_USEDEFAULT = 0x80000000

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms632600(v=vs.85).aspx
const (
	WS_CAPTION     = 0x00C00000
	WS_MAXIMIZEBOX = 0x00010000
	WS_MINIMIZEBOX = 0x00020000
	WS_OVERLAPPED  = 0x00000000
	WS_SYSMENU     = 0x00080000
	WS_THICKFRAME  = 0x00040000

	WS_OVERLAPPEDWINDOW = WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU | WS_THICKFRAME | WS_MINIMIZEBOX | WS_MAXIMIZEBOX
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/ff729176
const (
	CS_HREDRAW = 0x0002
	CS_VREDRAW = 0x0001
)
const NIF_MESSAGE = 0x00000001

const MIM_APPLYTOSUBMENUS = 0x80000000 // Settings apply to the menu and all of its submenus

const (
	WM_DESTROY    = 0x0002
	WM_CLOSE      = 0x0010
	WM_COMMAND    = 0x0111
	WM_LBUTTONUP  = 0x0202
	WM_RBUTTONUP  = 0x0205
	WM_ENDSESSION = 0x16
	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms644931(v=vs.85).aspx
	WM_USER = 0x0400
)

const (
	NIM_ADD    = 0x00000000
	NIM_MODIFY = 0x00000001
	NIM_DELETE = 0x00000002
)

const (
	TPM_BOTTOMALIGN = 0x0020
	TPM_LEFTALIGN   = 0x0000
)

const IMAGE_ICON = 1 // Loads an icon
const (
	LR_LOADFROMFILE = 0x00000010 // Loads the stand-alone image from the file
	LR_DEFAULTSIZE  = 0x00000040 // Loads default-size icon for windows(SM_CXICON x SM_CYICON) if cx, cy are set to zero
)

const (
	NIF_ICON = 0x00000002
	NIF_TIP  = 0x00000004
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms647578(v=vs.85).aspx
const (
	MIIM_STATE   = 0x00000001
	MIIM_ID      = 0x00000002
	MIIM_SUBMENU = 0x00000004
	MIIM_STRING  = 0x00000040
	MIIM_FTYPE   = 0x00000100
)

const (
	MFS_CHECKED  = 0x00000008
	MFS_DISABLED = 0x00000003
)

const (
	MFT_STRING    = 0x00000000
	MFT_SEPARATOR = 0x00000800
)
