package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"ws-tun-vpn/wtvc_gui/app_context"
)

func MakePageContentLoop(ctx *app_context.AppContext) {
	for navEvent := range ctx.NavChannel {
		switch navEvent.TargetPage {
		case app_context.MainPage:
			ctx.Window.SetContent(container.NewPadded(MakeMainPageContent(ctx)))
			//ctx.Window.Resize(fyne.NewSize(500, 400))
		case app_context.AddConfigPage:
			ctx.Window.SetContent(container.NewPadded(MakeConfigsPage(ctx, true)))
			ctx.Window.Resize(fyne.NewSize(500, 400))
		case app_context.ModifyConfigPage:
			ctx.Window.SetContent(container.NewPadded(MakeConfigsPage(ctx, false)))
			ctx.Window.Resize(fyne.NewSize(500, 400))
		case app_context.ConnectTo:
			ctx.Window.SetContent(container.NewPadded(MakeConnectPage(ctx)))
			ctx.Window.Resize(fyne.NewSize(500, 400))
		}
	}
}
