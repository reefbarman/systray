package systray

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/getlantern/golog"
)

var (
	// OnExitChan can be optionally waited on for detecting and handling the shutdown of the tray application
	OnExitChan = make(chan bool)

	currentID = int32(-1)
	menuItems = make(map[int32]*MenuItem)
	log       = golog.LoggerFor("systray")

	trayMenu      *Menu
	menuItemsLock sync.RWMutex
	onTrayRun     func()
)

// Run is called to start the tray application and the callback is triggered when it is up and running
func Run(onRun func()) {
	runtime.LockOSThread()

	onTrayRun = func() {
		menu, err := createMenu()
		if err != nil {
			log.Errorf("Unable to create root menu: %v", err)
			return
		}
		trayMenu = menu

		if onRun != nil {
			go onRun()
		}
	}

	nativeLoop()
}

// Quit will close the tray application, the OnExitChan will be triggered after the application has shut down
func Quit() {
	quit()
}

// SetIcon will set the icon for the tray application in the system tray
func SetIcon(iconBytes []byte) {
	bh := md5.Sum(iconBytes)
	dataHash := hex.EncodeToString(bh[:])
	iconFilePath := filepath.Join(os.TempDir(), "systray_temp_icon_"+dataHash)

	if _, err := os.Stat(iconFilePath); os.IsNotExist(err) {
		if err := ioutil.WriteFile(iconFilePath, iconBytes, 0644); err != nil {
			log.Errorf("Unable to write icon data to temp file: %v", err)
			return
		}
	}

	setIcon(iconFilePath)
}

// SetTooltip will set a tooltip on hover over the system tray icon
func SetTooltip(tooltip string) {
	setTooltip(tooltip)
}

// AddSeperator will add a seperator between items in the tray menu
func AddSeperator() {
	id := atomic.AddInt32(&currentID, 1)
	menuItem := &MenuItem{
		id:     id,
		parent: trayMenu,
	}

	addSeperator(menuItem)
}

// AddMenuItem will add a new item to the tray menu with an on click callback
func AddMenuItem(title string, onClick func(*MenuItem)) *MenuItem {
	menuItem := createMenuItem(title, trayMenu)
	menuItem.onClick = onClick

	setMenuItem(menuItem)

	return menuItem
}

// AddSubMenuItem will add a new sub menu to the tray menu. The sub menu is returned, allowing the adding of items to it
func AddSubMenuItem(title string) *Menu {
	menuItem := createMenuItem(title, trayMenu)

	return addSubMenuItem(menuItem)
}

func createMenuItem(title string, parent *Menu) *MenuItem {
	id := atomic.AddInt32(&currentID, 1)

	menuItem := &MenuItem{
		id:     id,
		title:  title,
		parent: parent,
	}

	menuItemsLock.Lock()
	defer menuItemsLock.Unlock()
	menuItems[id] = menuItem

	return menuItem
}

func onMenuItemSelected(menuID int32) {
	menuItemsLock.RLock()
	item := menuItems[menuID]
	menuItemsLock.RUnlock()

	item.onClick(item)
}

func onExit() {
	select {
	case OnExitChan <- true:
		break
	default:
		break
	}
}
