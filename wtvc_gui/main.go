package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"ws-tun-vpn/pkg/privilege"
	"ws-tun-vpn/wtvc_gui/app_context"
	"ws-tun-vpn/wtvc_gui/ico"
	"ws-tun-vpn/wtvc_gui/pages"
	"ws-tun-vpn/wtvc_gui/the"
)

func main() {
	p := privilege.New()
	if p.IsAdmin() {
		RunApp()
	} else {
		err := p.Elevate()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func RunApp() {
	a := app.NewWithID("com.wtv.c")
	if meta := a.Metadata(); meta.Name == "" {
		// App not packaged, probably from `go run`.
		meta.Name = "WTVC客户端"
		app.SetMetadata(meta)
	}
	a.SetIcon(ico.LoadIcon())
	a.Settings().SetTheme(&the.MyTheme{})
	var w fyne.Window
	w = a.NewWindow("WTVC客户端(code by Wuly)")

	w.SetCloseIntercept(func() {
		w.Hide()
	})
	if desk, ok := a.(desktop.App); ok {
		quit := fyne.NewMenuItem("退出", nil)
		quit.Icon = theme.LogoutIcon()
		quit.IsQuit = true
		showMenu := fyne.NewMenuItem("显示", nil)
		showMenu.Icon = theme.VisibilityIcon()
		showMenu.Action = func() {
			w.Show()
		}
		hideMenu := fyne.NewMenuItem("隐藏", nil)
		hideMenu.Icon = theme.VisibilityOffIcon()
		hideMenu.Action = func() {
			w.Hide()
		}
		m := fyne.NewMenu("WTVC客户端(code by Wuly)", showMenu, hideMenu, fyne.NewMenuItemSeparator(), quit)
		desk.SetSystemTrayMenu(m)
	}
	w.SetCloseIntercept(func() {
		w.Hide()
	})
	appContext := app_context.NewAppContext(w, a.Preferences())
	appContext.LoadAppConfigs()
	go pages.MakePageContentLoop(appContext)
	appContext.NavChannel <- app_context.NavEvent{
		TargetPage: app_context.MainPage,
	}
	w.SetFixedSize(false)
	w.ShowAndRun()
}
