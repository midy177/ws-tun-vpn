package pages

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"log"
	"sync/atomic"
	"ws-tun-vpn/wtvc_gui/app_context"
)

// MakeMainPageContent creates the main page content
func MakeMainPageContent(ctx *app_context.AppContext) fyne.CanvasObject {
	var list *widget.List

	list = widget.NewList(
		func() int {
			return len(ctx.AppConfigs)
		},
		func() fyne.CanvasObject {
			selected := canvas.NewRectangle(color.Transparent)
			selected.SetMinSize(fyne.NewSize(10, 10))
			//indicator := widget.NewIcon(theme.ViewRefreshIcon())
			label := widget.NewLabel("")
			toolbar := widget.NewToolbar()

			return container.NewHBox(
				selected,
				//indicator,
				label,
				layout.NewSpacer(),
				toolbar,
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			cs := obj.(*fyne.Container)
			selected := cs.Objects[0].(*canvas.Rectangle)
			//indicator := container.Objects[1].(*widget.Icon)
			label := cs.Objects[1].(*widget.Label)
			toolbar := cs.Objects[3].(*widget.Toolbar)

			label.SetText(ctx.AppConfigs[id].Label)

			if int64(id) == ctx.SelectedItemID {
				// Set the selected item style
				selected.FillColor = theme.PrimaryColor()
			} else {
				// Reset the style for the unselected items
				selected.FillColor = color.Transparent
			}
			// Clear previous toolbar items and add a new delete icon
			toolbar.Items = nil
			arrowIcon := widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
				atomic.StoreInt64(&ctx.SelectedItemID, int64(id))
				list.Refresh()
				// navigate to page result for specific menu item
				ctx.NavChannel <- app_context.NavEvent{TargetPage: app_context.ModifyConfigPage}
				// Define action for the "+" icon
			})

			// Create a new delete icon action for each item
			deleteIcon := widget.NewToolbarAction(theme.DeleteIcon(), func() {
				callback := func(confirm bool) {
					if confirm {
						// Delete the item from the data slice
						ctx.AppConfigs = append(ctx.AppConfigs[:id], ctx.AppConfigs[id+1:]...)
						ctx.UpdateAppConfigs()
						// Refresh the list to update the view
						list.Refresh()
						log.Printf("Delete icon clicked - Item %d deleted", id)
					}
				}
				dg := dialog.NewConfirm("确认删除", "确定删除配置？", callback, ctx.Window)
				dg.SetConfirmText("确认✅")
				dg.SetDismissText("取消❌")
				dg.Show()
			})
			toolbar.Append(arrowIcon)
			toolbar.Append(deleteIcon)
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		//log.Printf("Selected Item ID %v", ctx.SelectedItemID)
		atomic.StoreInt64(&ctx.SelectedItemID, int64(id))
		list.Refresh()
	}
	//list.OnUnselected = func(id widget.ListItemID) {
	//	if ctx.SelectedItemID == int64(id) {
	//		ctx.SelectedItemID = -1
	//		list.Refresh()
	//		return
	//	}
	//}

	// Create the scroll container for the list
	scrollContainer := container.NewScroll(list)
	// Set a minimum size for the scroll container if needed
	scrollContainer.SetMinSize(fyne.NewSize(400, 300)) // Set width and height as needed
	// Use container.Max to allocate as much space as possible to the list
	listWithMaxHeight := container.NewStack(scrollContainer)

	header := MakePageHeader(ctx, list)

	ConnectButton := widget.NewButton("连接", func() {})
	ConnectButton.Importance = widget.HighImportance

	//statusBox := widget.NewLabel("")
	//statusBox.Wrapping = fyne.TextWrapWord

	ConnectButton.OnTapped = func() {
		if ctx.SelectedItemID >= 0 && int(ctx.SelectedItemID) < len(ctx.AppConfigs) {
			ctx.NavChannel <- app_context.NavEvent{TargetPage: app_context.ConnectTo}
		} else {
			dialog.ShowError(errors.New("未选择配置文件"), ctx.Window)
		}
	}
	// Combine the header and the scrollable list in a vertical box layout
	return container.NewVBox(
		header,
		listWithMaxHeight, // The scrollable list with enforced maximum height
		ConnectButton,
	)
}

func MakePageHeader(ctx *app_context.AppContext, list *widget.List) *fyne.Container {
	// Create the toolbar with a "+" icon
	rightToolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			ctx.NavChannel <- app_context.NavEvent{TargetPage: app_context.AddConfigPage}
		}),
	)
	// Create the header label
	headerLabel := widget.NewLabel("配置列表")
	headerLabel.TextStyle = fyne.TextStyle{Bold: true}
	headerLabel.Alignment = fyne.TextAlignCenter

	// Create the header using HBox layout
	header := container.NewHBox(
		layout.NewSpacer(),
		layout.NewSpacer(), // Spacer pushes the toolbar to the right
		rightToolbar,
	)
	return container.NewStack(header, headerLabel)
}
