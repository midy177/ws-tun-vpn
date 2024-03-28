package nic_tool

import (
	"log"
	"os/exec"
	"strings"
)

type NicTool interface {
	SetCidrAndUp() string
	SetMtu() string
	SetRoute(cidr string) string
}

type tool struct {
	tunName string
	cidr    string
	mtu     int
}

func NewNicTool(tunName, cidr string, mtu int) NicTool {
	return &tool{tunName: tunName, cidr: cidr, mtu: mtu}
}

// execCmd executes the given command
func execCmd(c string, args ...string) string {
	//log.Printf("exec %v %v", c, args)
	cmd := exec.Command(c, args...)
	out, err := cmd.Output()
	if err != nil {
		log.Println("failed to exec cmd:", err)
	}
	if len(out) == 0 {
		return ""
	}
	s := string(out)
	return strings.ReplaceAll(s, "\n", "")
}
