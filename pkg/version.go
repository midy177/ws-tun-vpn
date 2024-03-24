package pkg

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"runtime"
)

var (
	Version   = "v1.7.2"
	GitHash   = ""
	BuildTime = ""
	Title     = "A simple VPN based on websocket and tun devices written in Go."
)

func DisplayVersionInfo() {
	fmt.Println("-")
	figure.NewFigure("ws tun vpn", "doom", true).Print()
	fmt.Println(Title)
	fmt.Printf("Version -> %s\n", Version)
	if GitHash != "" {
		fmt.Printf("Git hash -> %s\n", GitHash)
	}
	if BuildTime != "" {
		fmt.Printf("Build time -> %s\n", BuildTime)
	}

	fmt.Printf("Go version -> %s\n", runtime.Version())

}
