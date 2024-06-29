package pages

import (
	"bufio"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"io"
	"os"
	"os/exec"
)

func MakePrivilegePage(w fyne.Window) fyne.CanvasObject {
	//headerToolbarLeft := widget.NewToolbar(
	//	widget.NewToolbarAction(theme.NavigateBackIcon(), func() {
	//		ctx.NavChannel <- app_context.NavEvent{TargetPage: app_context.MainPage}
	//	}),
	//)
	// Create the header label
	//headerLabel := widget.NewLabel("提升应用权限")
	//headerLabel.TextStyle = fyne.TextStyle{Bold: true}
	//headerLabel.Alignment = fyne.TextAlignCenter

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("输入sudo提权密码")
	password.MultiLine = true

	bottomSaveBtn := widget.NewButtonWithIcon("确认", theme.ConfirmIcon(), func() {
		if len(password.Text) == 0 {
			dialog.ShowError(errors.New("连接密钥未设置"), w)
			return
		}
		err := privilege(password.Text)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		w.Close()
	})
	bottomSaveBtn.Importance = widget.HighImportance
	bottomSaveBtn.Disable()

	password.OnChanged = func(s string) {
		if len(s) != 0 {
			bottomSaveBtn.Enable()
		} else {
			bottomSaveBtn.Disable()
		}
	}

	// Combine the header and the accordions in a vertical box layout
	return container.NewBorder(nil, bottomSaveBtn, nil, nil, password)
}

func privilege(password string) error {
	// 检查程序是否以管理员权限运行
	if os.Geteuid() == 0 {
		return nil
	}
	// 提示用户输入密码
	// 创建一个 sudo 命令
	cmd := exec.Command("sudo", os.Args[0])

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer stdoutPipe.Close()
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	defer stdinPipe.Close()
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error executing sudo command:", err)
		return err
	}
	stdoutBuf := bufio.NewReader(stdoutPipe)
	for {
		readByte, err := stdoutBuf.ReadByte()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		fmt.Println(string(readByte))
	}

	return nil
}
