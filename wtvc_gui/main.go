package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"ws-tun-vpn/wtvc_gui/app_context"
	"ws-tun-vpn/wtvc_gui/ico"
	"ws-tun-vpn/wtvc_gui/pages"
	"ws-tun-vpn/wtvc_gui/the"
)

func main() {
	a := app.NewWithID("com.wtv.c")
	if meta := a.Metadata(); meta.Name == "" {
		// App not packaged, probably from `go run`.
		meta.Name = "WTVC客户端"
		app.SetMetadata(meta)
	}
	a.SetIcon(ico.LoadIcon())
	a.Settings().SetTheme(&the.MyTheme{})
	var w fyne.Window
	//drv := fyne.CurrentApp().Driver()
	//if dr, ok := drv.(desktop.Driver); ok {
	//	w = dr.CreateSplashWindow()
	//} else {
	w = a.NewWindow("WTVC客户端(code by Wuly)")
	//}

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
	//w.SetContent(makeWindow(w))
	w.SetCloseIntercept(func() {
		w.Hide()
	})
	appContext := app_context.NewAppContext(w, a.Preferences())
	appContext.LoadAppConfigs()
	go pages.MakePageContentLoop(appContext)
	appContext.NavChannel <- app_context.NavEvent{
		TargetPage: app_context.MainPage,
	}
	w.Resize(fyne.NewSize(500, 80))
	w.SetFixedSize(true)
	w.ShowAndRun()
}

//func makeWindow(w fyne.Window) fyne.CanvasObject {
//	var signal = make(chan int, 1)
//	serverUrl := widget.NewEntry()
//	serverUrl.SetPlaceHolder("输入连接地址,如:yeastar.com:8080")
//	serverUrl.OnSubmitted = func(s string) {
//		log.Println("Server URL submitted:", s)
//	}
//	password := widget.NewPasswordEntry()
//	password.SetPlaceHolder("输入连接密钥")
//	password.MultiLine = true
//	password.OnSubmitted = func(s string) {
//		log.Println("Password changed:", s)
//	}
//	certFilename := widget.NewSelectEntry(nil)
//	certFilename.SetPlaceHolder("选填,服务端证书(.crt或.pem)")
//	certFilename.ActionItem = widget.NewButtonWithIcon("", theme.MenuDropDownIcon(), func() {
//		drv := fyne.CurrentApp().Driver()
//		if drv, ok := drv.(desktop.Driver); ok {
//			childWindow := drv.CreateSplashWindow()
//			childWindow.Resize(fyne.NewSize(520, 400))
//			fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
//				defer childWindow.Close()
//				if err != nil {
//					dialog.ShowError(err, w)
//					return
//				}
//				if reader == nil {
//					return
//				}
//				certFilename.SetText(reader.URI().String())
//				reader.Close()
//			}, childWindow)
//			fd.SetFilter(storage.NewExtensionFileFilter([]string{".crt", ".pem"}))
//			fd.SetConfirmText("选择")
//			fd.SetDismissText("取消")
//			fd.Resize(fyne.NewSize(720, 400))
//			childWindow.Show()
//			fd.Show()
//		}
//	})
//	certFilename.Hide()
//	skipCheckCert := widget.NewCheck("跳过TLS校验", func(value bool) {})
//	skipCheckCert.Hide()
//	check := widget.NewCheck("", func(value bool) {
//		if value {
//			signal <- enableTls
//		} else {
//			signal <- disableTls
//			skipCheckCert.Checked = false
//		}
//	})
//	form := widget.NewForm(
//		widget.NewFormItem("连接地址:", serverUrl),
//		widget.NewFormItem("连接密钥:", password),
//		widget.NewFormItem("TLS开关:", container.NewBorder(nil, nil,
//			check, skipCheckCert, certFilename)),
//	)
//
//	bottomBtn := widget.NewButtonWithIcon("连接", nil, func() {
//		signal <- connect
//	})
//	bottomBtn.Importance = widget.HighImportance
//	go func() {
//		for {
//			tp := <-signal
//			switch tp {
//			case connect:
//				popExitBtn := widget.NewButtonWithIcon("断开连接", theme.CancelIcon(), func() {})
//				popExitBtn.Importance = widget.WarningImportance
//				data := binding.NewString()
//				_ = data.Set("连接中...")
//				status := widget.NewLabelWithData(data)
//				status.TextStyle = fyne.TextStyle{Bold: true}
//				pop := widget.NewModalPopUp(
//					container.NewBorder(container.NewCenter(status), popExitBtn, layout.NewSpacer(), layout.NewSpacer()),
//					w.Canvas())
//				ctx, cancel := StartClient(&types.ClientConfig{
//					ServerUrl: serverUrl.Text,
//					BaseConfig: types.BaseConfig{
//						AuthCode:  password.Text,
//						EnableTLS: check.Checked,
//						Verbose:   true,
//						MTU:       1500,
//					},
//					CertificateFile: certFilename.Text,
//					SkipTLSVerify:   skipCheckCert.Checked,
//				}, w)
//				Counter(ctx, data)
//				pop.Resize(fyne.NewSize(w.Canvas().Size().Width, w.Canvas().Size().Height))
//				popExitBtn.OnTapped = func() {
//					pop.Hide()
//					cancel()
//					signal <- disconnect
//				}
//				pop.Show()
//				<-ctx.Done()
//				pop.Hide()
//			case disconnect:
//				log.Println("Received signal disconnect")
//			case enableTls:
//				skipCheckCert.Show()
//				certFilename.Show()
//				log.Println("TLS enabled")
//			case disableTls:
//				skipCheckCert.Hide()
//				certFilename.Hide()
//				log.Println("TLS disabled")
//			}
//		}
//	}()
//
//	return container.NewPadded(container.NewBorder(nil, bottomBtn, nil, nil, form))
//}
