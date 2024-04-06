package pages

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"time"
	"ws-tun-vpn/pkg/counter"
	"ws-tun-vpn/pkg/logview"
	"ws-tun-vpn/service"
	"ws-tun-vpn/wtvc_gui/app_context"
)

func MakeConnectPage(ctx *app_context.AppContext) fyne.CanvasObject {
	download := widget.NewButtonWithIcon("下载: 0kb", theme.DownloadIcon(), nil)
	download.Importance = widget.SuccessImportance
	upload := widget.NewButtonWithIcon("上传: 0kb", theme.UploadIcon(), nil)
	upload.Importance = widget.WarningImportance
	headerLabel := container.NewCenter(container.NewHBox(upload, download))

	mask := widget.NewMultiLineEntry()
	mask.Disable()
	msg := logview.NewLogView(ctx.Window)
	saveBtn := widget.NewButtonWithIcon("断开连接", theme.ContentClearIcon(), nil)
	saveBtn.Importance = widget.DangerImportance
	childCtx, cancel := StartClient(ctx, msg, saveBtn)
	saveBtn.OnTapped = func() {
		if saveBtn.Text == "断开连接" {
			cancel()
			saveBtn.SetIcon(theme.NavigateBackIcon())
			saveBtn.SetText("返回主页")
			saveBtn.Refresh()
		} else if saveBtn.Text == "返回主页" {
			ctx.NavChannel <- app_context.NavEvent{
				TargetPage: app_context.MainPage,
			}
		}
	}
	Counter(childCtx, upload, download)
	// Combine the header and the accordions in a vertical box layout
	return container.NewBorder(headerLabel, saveBtn, nil, nil, container.NewStack(mask, container.NewPadded(msg.GetView())))
}

func StartClient(ctx *app_context.AppContext, logView logview.LogView, saveBtn *widget.Button) (context.Context, context.CancelFunc) {
	childCtx, cancel := context.WithCancel(context.TODO())
	go func() {
		cfg := ctx.AppConfigs[ctx.SelectedItemID].ClientConfig
		err := service.NewClientService(context.WithValue(childCtx, "config", &cfg), logView)
		if err != nil {
			logView.Print(logview.LogError, err.Error())
		}
		cancel()
		if saveBtn.Text == "断开连接" {
			saveBtn.SetIcon(theme.NavigateBackIcon())
			saveBtn.SetText("返回主页")
			saveBtn.Refresh()
		}
	}()
	return childCtx, cancel
}

func Counter(ctx context.Context, upload, download *widget.Button) {
	go func() {
		t := time.NewTicker(time.Second)
		defer counter.ResetBytes()
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				upload.SetText("上传: 0kb")
				upload.Refresh()
				download.SetText("下载: 0kb")
				download.Refresh()
				return
			case <-t.C:
				d, u := counter.PrintBytes(false)
				upload.SetText("上传: " + u)
				upload.Refresh()
				download.SetText("下载: " + d)
				download.Refresh()
			}
		}
	}()
}
