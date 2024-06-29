//go:build linux
// +build linux

package privilege

import (
	"os"
	"os/exec"
)

func (p privilege) IsAdmin() bool {
	return os.Geteuid() == 0
}

func (p privilege) Elevate() error {
	args := []string{}
	args = append(os.Args, string(os.Getpid()))
	// 使用 Polkit 工具 pkexec 来请求授权
	cmd := exec.Command("pkexec", args)
	// 执行命令
	return cmd.Run()
}
