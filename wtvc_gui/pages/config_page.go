package pages

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"os"
	"ws-tun-vpn/pkg/util"
	"ws-tun-vpn/wtvc_gui/app_context"
)

const (
	AddConfigTitle    = "添加配置"
	ModifyConfigTitle = "修改配置"
)

func MakeConfigsPage(ctx *app_context.AppContext, isAdd bool) fyne.CanvasObject {
	headerToolbarLeft := widget.NewToolbar(
		widget.NewToolbarAction(theme.NavigateBackIcon(), func() {
			ctx.NavChannel <- app_context.NavEvent{TargetPage: app_context.MainPage}
		}),
	)
	// Create the header label
	headerLabel := widget.NewLabel(ModifyConfigTitle)
	hasSelectedItem := false
	if isAdd {
		headerLabel.Text = AddConfigTitle
	} else {
		hasSelectedItem = len(ctx.AppConfigs) > int(ctx.SelectedItemID)
	}
	headerLabel.TextStyle = fyne.TextStyle{Bold: true}
	headerLabel.Alignment = fyne.TextAlignCenter

	// Create the header using HBox layout
	header := container.NewHBox(
		headerToolbarLeft,
		layout.NewSpacer(), // Spacer pushes the toolbar to the right
		layout.NewSpacer(),
	)
	cfgName := widget.NewEntry()
	cfgName.SetPlaceHolder("例如: 北京集群")

	serverUrl := widget.NewEntry()
	serverUrl.SetPlaceHolder("输入连接地址,如:yeastar.com:8080")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("输入连接密钥")
	password.MultiLine = true

	certFilename := widget.NewSelectEntry(nil)
	certFilename.SetPlaceHolder("选填,服务端证书(.crt或.pem)")
	certFilename.ActionItem = widget.NewButtonWithIcon("", theme.MenuDropDownIcon(), func() {
		drv := fyne.CurrentApp().Driver()
		if drv, ok := drv.(desktop.Driver); ok {
			childWindow := drv.CreateSplashWindow()
			childWindow.Resize(fyne.NewSize(720, 400))
			fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				defer childWindow.Close()
				if err != nil {
					dialog.ShowError(err, ctx.Window)
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
	enableTLS := widget.NewCheck("", nil)

	specialCert := widget.NewCheck("指定证书", func(value bool) {})
	specialCert.Hide()
	specialCert.OnChanged = func(value bool) {
		if value {
			skipCheckCert.Checked = false
			skipCheckCert.Hide()
			certFilename.Show()
		} else {
			skipCheckCert.Show()
			certFilename.SetText("")
			certFilename.Hide()
		}
	}
	skipCheckCert.OnChanged = func(value bool) {
		if value {
			specialCert.SetChecked(false)
			specialCert.Hide()
		} else {
			specialCert.Show()
		}
	}
	if !isAdd && hasSelectedItem {
		cfgName.SetText(ctx.AppConfigs[ctx.SelectedItemID].Label)
		serverUrl.SetText(ctx.AppConfigs[ctx.SelectedItemID].ClientConfig.ServerUrl)
		password.SetText(ctx.AppConfigs[ctx.SelectedItemID].BaseConfig.AuthCode)
		certFilename.SetText(ctx.AppConfigs[ctx.SelectedItemID].ClientConfig.CertificateFile)
		skipCheckCert.SetChecked(ctx.AppConfigs[ctx.SelectedItemID].ClientConfig.SkipTLSVerify)
		enableTLS.SetChecked(ctx.AppConfigs[ctx.SelectedItemID].BaseConfig.EnableTLS)
		if enableTLS.Checked {
			skipCheckCert.Show()
			specialCert.Show()
			if specialCert.Checked {
				certFilename.Show()
			}
		}
	}
	form := widget.NewForm(
		widget.NewFormItem("配置名称:", cfgName),
		widget.NewFormItem("连接地址:", serverUrl),
		widget.NewFormItem("连接密钥:", password),
		widget.NewFormItem("TLS开关:", container.NewHBox(
			enableTLS, specialCert, skipCheckCert)),
		widget.NewFormItem("", certFilename),
	)

	enableTLS.OnChanged = func(value bool) {
		if value {
			skipCheckCert.Show()
			specialCert.Show()
			log.Println("TLS enabled")
		} else {
			skipCheckCert.Checked = false
			skipCheckCert.Hide()
			specialCert.Hide()
			certFilename.SetText("")
			certFilename.Hide()
			log.Println("TLS disabled")
		}
	}
	bottomSaveBtn := widget.NewButtonWithIcon("保存", theme.DocumentSaveIcon(), func() {
		// 校验陪配置是否合法
		if len(cfgName.Text) == 0 {
			dialog.ShowError(errors.New("配置名称未设置"), ctx.Window)
			return
		}
		if !util.IsValidAddress(serverUrl.Text) {
			dialog.ShowError(errors.New(serverUrl.Text+" 是无效的地址"), ctx.Window)
			return
		}
		if len(password.Text) == 0 {
			dialog.ShowError(errors.New("连接密钥未设置"), ctx.Window)
			return
		}
		if enableTLS.Checked && specialCert.Checked {
			if _, err := os.Stat(certFilename.Text); err != nil {
				dialog.ShowError(err, ctx.Window)
				return
			}
		}
		acfg := app_context.AppConfigs{
			Label: cfgName.Text,
		}
		acfg.BaseConfig.EnableTLS = enableTLS.Checked
		acfg.ClientConfig.ServerUrl = serverUrl.Text
		acfg.ClientConfig.AuthCode = password.Text
		acfg.BaseConfig.MTU = 1500
		acfg.ClientConfig.SkipTLSVerify = skipCheckCert.Checked
		acfg.ClientConfig.CertificateFile = certFilename.Text
		if isAdd {
			// add
			ctx.AppConfigs = append(ctx.AppConfigs, acfg)
		} else if hasSelectedItem {
			// modify
			ctx.AppConfigs[ctx.SelectedItemID] = acfg
		}
		ctx.UpdateAppConfigs()
		ctx.NavChannel <- app_context.NavEvent{
			TargetPage: app_context.MainPage,
		}
	})
	bottomSaveBtn.Importance = widget.HighImportance

	// Combine the header and the accordions in a vertical box layout
	return container.NewBorder(nil, bottomSaveBtn, nil, nil, container.NewVBox(
		container.NewStack(header, headerLabel), form))
}
