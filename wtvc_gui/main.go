package main

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"ws-tun-vpn/service"
	"ws-tun-vpn/types"
	"ws-tun-vpn/wtvc_gui/ico"
	"ws-tun-vpn/wtvc_gui/the"
)

func main() {
	a := app.New()
	a.SetIcon(ico.LoadIcon())
	a.Settings().SetTheme(&the.MyTheme{})
	w := a.NewWindow("连接VPN服务器配置")
	w.SetCloseIntercept(func() {
		w.Hide()
	})
	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu("Websocket tun vpn client",
			fyne.NewMenuItem("显示", func() {
				w.Show()
			}),
			fyne.NewMenuItem("隐藏", func() {
				w.Hide()
			}))
		desk.SetSystemTrayMenu(m)
		//desk.SetSystemTrayIcon(ico.LoadIcon())
	}
	w.SetContent(makeWindow(w))
	w.SetCloseIntercept(func() {
		w.Hide()
	})
	w.Resize(fyne.NewSize(500, 80))
	w.ShowAndRun()
}
func makeWindow(w fyne.Window) fyne.CanvasObject {
	var signal = make(chan int, 1)
	serverUrl := widget.NewEntry()
	serverUrl.SetPlaceHolder("输入连接地址,如:yeastar.com:8080")
	serverUrl.OnSubmitted = func(s string) {
		log.Println("Server URL submitted:", s)
	}
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("输入连接密钥")
	password.MultiLine = true
	password.OnSubmitted = func(s string) {
		log.Println("Password changed:", s)
	}
	certFilename := widget.NewSelectEntry(nil)
	certFilename.SetPlaceHolder("选填,服务端证书(.crt或.pem)")
	certFilename.ActionItem = widget.NewButtonWithIcon("", theme.MenuDropDownIcon(), func() {
		drv := fyne.CurrentApp().Driver()
		if drv, ok := drv.(desktop.Driver); ok {
			childWindow := drv.CreateSplashWindow()
			childWindow.Resize(fyne.NewSize(520, 400))
			fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				defer childWindow.Close()
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				if reader == nil {
					return
				}
				certFilename.SetText(reader.URI().String())
				reader.Close()
			}, childWindow)
			fd.SetFilter(storage.NewExtensionFileFilter([]string{".crt", ".pem"}))
			fd.SetConfirmText("选择")
			fd.SetDismissText("取消")
			fd.Resize(fyne.NewSize(720, 400))
			childWindow.Show()
			fd.Show()
		}
	})
	certFilename.Hide()
	skipCheckCert := widget.NewCheck("跳过TLS校验", func(value bool) {})
	skipCheckCert.Hide()
	check := widget.NewCheck("", func(value bool) {
		if value {
			signal <- enableTls
		} else {
			signal <- disableTls
			skipCheckCert.Checked = false
		}
	})
	form := widget.NewForm(
		widget.NewFormItem("连接地址:", serverUrl),
		widget.NewFormItem("连接密钥:", password),
		widget.NewFormItem("TLS开关:", container.NewBorder(nil, nil,
			check, skipCheckCert, certFilename)),
	)

	bottomBtn := widget.NewButtonWithIcon("连接", nil, func() {
		signal <- connect
	})
	bottomBtn.Importance = widget.HighImportance
	go func() {
		for {
			tp := <-signal
			switch tp {
			case connect:
				log.Println(serverUrl.Text, password.Text, certFilename.Text, check.Checked)

				popExitBtn := widget.NewButtonWithIcon("断开连接", theme.CancelIcon(), func() {})
				popExitBtn.Importance = widget.WarningImportance
				pop := widget.NewPopUp(
					container.NewBorder(widget.NewLabel("连接成功"), popExitBtn, nil, nil),
					w.Canvas())
				ctx, cancel := StartClient(&types.ClientConfig{
					ServerUrl: serverUrl.Text,
					BaseConfig: types.BaseConfig{
						AuthCode:  password.Text,
						EnableTLS: check.Checked,
					},
					CertificateFile: certFilename.Text,
					SkipTLSVerify:   skipCheckCert.Checked,
				}, w)
				pop.Resize(fyne.NewSize(w.Canvas().Size().Width, w.Canvas().Size().Height))
				popExitBtn.OnTapped = func() {
					pop.Hide()
					cancel()
					signal <- disconnect
				}
				pop.Show()
				<-ctx.Done()
				pop.Hide()
			case disconnect:
				log.Println("Received signal disconnect")
			case enableTls:
				skipCheckCert.Show()
				certFilename.Show()
				log.Println("TLS enabled")
			case disableTls:
				skipCheckCert.Hide()
				certFilename.Hide()
				log.Println("TLS disabled")
			}
		}
	}()

	return container.NewPadded(container.NewBorder(nil, bottomBtn, nil, nil, form))
}

const (
	connect    = 0x01
	disconnect = 0x02
	enableTls  = 0x03
	disableTls = 0x04
)

func StartClient(config *types.ClientConfig, w fyne.Window) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.TODO())
	go func(ctx2 context.Context, cancel1 context.CancelFunc) {
		err := service.NewClientService(context.WithValue(ctx, "config", config))
		if err != nil {
			log.Println(err)
			dialog.ShowError(err, w)
		}
		cancel1()
	}(ctx, cancel)
	return ctx, cancel
}
