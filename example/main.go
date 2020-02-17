package main

import (
	"fmt"
	"winsystray"
	"winsystray/example/icon"

	"github.com/sqweek/dialog"
)

func main() {
	winsystray.Run(func() {
		winsystray.SetIcon(icon.Data)
		winsystray.SetTooltip("This here is an example")
		subMenu := winsystray.AddSubMenuItem("Sub Menu")
		if subMenu != nil {
			subMenu.AddMenuItem("Click Me", func(item *winsystray.MenuItem) {
				dialog.Message("%s", "Do you want to continue?").Title("Are you sure?").YesNo()
			})

			subMenu.AddSeperator()

			subMenu.AddMenuItem("Checkable", func(checkable *winsystray.MenuItem) {
				checkable.ToogleChecked()
			})

			disable := subMenu.AddMenuItem("Click to Disable", func(disable *winsystray.MenuItem) {
				disable.ToggleDisabled()
			})

			subMenu.AddMenuItem("Click to Disable", func(toggle *winsystray.MenuItem) {
				disable.ToggleDisabled()

				if disable.IsDisabled() {
					toggle.SetTitle("Click to Enable")
				} else {
					toggle.SetTitle("Click to Disable")
				}
			})

			anotherSubMenu := subMenu.AddSubMenuItem("Another SubMenu")
			if anotherSubMenu != nil {
				anotherSubMenu.AddMenuItem("Click Away", func(item *winsystray.MenuItem) {
					fmt.Println("You clicked!")
				})
			}
		}

		winsystray.AddMenuItem("Checkable", func(checkable *winsystray.MenuItem) {
			checkable.ToogleChecked()
		})

		winsystray.AddMenuItem("Click to Disable", func(disable *winsystray.MenuItem) {
			disable.ToggleDisabled()
		})

		winsystray.AddSeperator()

		winsystray.AddMenuItem("Exit", func(item *winsystray.MenuItem) {
			winsystray.Quit()
		})

		<-winsystray.OnExitChan

		fmt.Println("app exiting")
	})
}
