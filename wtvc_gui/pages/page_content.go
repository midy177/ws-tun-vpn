package pages

import (
	"ws-tun-vpn/wtvc_gui/app_context"
)

func MakePageContentLoop(ctx *app_context.AppContext) {
	for navEvent := range ctx.NavChannel {
		switch navEvent.TargetPage {
		case app_context.MainPage:
			ctx.Window.SetContent(MakeMainPageContent(ctx))
		case app_context.AddConfigPage:
			ctx.Window.SetContent(MakeConfigsPage(ctx, true))
		case app_context.ModifyConfigPage:
			ctx.Window.SetContent(MakeConfigsPage(ctx, false))
		case app_context.ConnectTo:
			ctx.Window.SetContent(MakeConnectPage(ctx))
		}
	}
}
