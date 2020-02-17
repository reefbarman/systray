package main

import (
	"fmt"

	"github.com/reefbarman/systray"
	"github.com/reefbarman/systray/example/icon"

	"github.com/sqweek/dialog"
)

func main() {
	systray.Run(func() {
		systray.SetIcon(icon.Data)
		systray.SetTooltip("This here is an example")
		subMenu := systray.AddSubMenuItem("Sub Menu")
		if subMenu != nil {
			subMenu.AddMenuItem("Click Me", func(item *systray.MenuItem) {
				dialog.Message("%s", "Do you want to continue?").Title("Are you sure?").YesNo()
			})

			subMenu.AddSeperator()

			subMenu.AddMenuItem("Checkable", func(checkable *systray.MenuItem) {
				checkable.ToogleChecked()
			})

			disable := subMenu.AddMenuItem("Click to Disable", func(disable *systray.MenuItem) {
				disable.ToggleDisabled()
			})

			subMenu.AddMenuItem("Click to Disable", func(toggle *systray.MenuItem) {
				disable.ToggleDisabled()

				if disable.IsDisabled() {
					toggle.SetTitle("Click to Enable")
				} else {
					toggle.SetTitle("Click to Disable")
				}
			})

			anotherSubMenu := subMenu.AddSubMenuItem("Another SubMenu")
			if anotherSubMenu != nil {
				anotherSubMenu.AddMenuItem("Click Away", func(item *systray.MenuItem) {
					fmt.Println("You clicked!")
				})
			}
		}

		systray.AddMenuItem("Checkable", func(checkable *systray.MenuItem) {
			checkable.ToogleChecked()
		})

		systray.AddMenuItem("Click to Disable", func(disable *systray.MenuItem) {
			disable.ToggleDisabled()
		})

		systray.AddSeperator()

		systray.AddMenuItem("Exit", func(item *systray.MenuItem) {
			systray.Quit()
		})

		<-systray.OnExitChan

		fmt.Println("app exiting")
	})
}
