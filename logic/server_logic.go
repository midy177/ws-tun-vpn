package logic

import (
	"bytes"
	"context"
	"errors"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"log"
	"net"
	"runtime"
	"strings"
	"time"
	"ws-tun-vpn/pkg/netutil"
	"ws-tun-vpn/pkg/nic_tool"
	"ws-tun-vpn/types"
)

// ServerLogic server logic struct.
type ServerLogic struct {
	ctx    context.Context
	config *types.ServerConfig
	iface  net.PacketConn
}

// NewServerLogic create a new server logic instance.
func NewServerLogic(ctx context.Context) (*ServerLogic, error) {
	config, ok := ctx.Value("config").(*types.ServerConfig)
	if !ok {
		return nil, errors.New("failed to get config from context")
	}
	return &ServerLogic{
		ctx:    ctx,
		config: config,
	}, nil
}

// ServerTunPacketRouteToClient 监听tun设备收到的数据，并将其转发给客户端
func (s *ServerLogic) ServerTunPacketRouteToClient() {
	nt := nic_tool.NewNicTool(s.config.IFace.Name(), s.config.BindAddress, int(s.config.MTU))
	info := nt.SetCidrAndUp()
	if s.config.Verbose && runtime.GOOS != "windows" {
		log.Printf("set tun network card(%s) cidr(%s) and up.\n", s.config.IFace.Name(), s.config.BindAddress)
		if len(info) > 0 {
			log.Println(info)
		}
	}
	info = nt.SetMtu()
	if s.config.Verbose {
		log.Printf("set tun network card(%s) mtu: %d\n", s.config.IFace.Name(), s.config.MTU)
		if len(info) > 0 {
			log.Println(info)
		}
	}
	packet := make([]byte, 0, 2048)
	packet = append(packet, packetMsg)
	for {
		n, err := s.config.IFace.Read(packet[1:])
		if err != nil {
			log.Fatalf("failed to read packet: %v", err)
			//netutil.PrintErr(err, s.config.Verbose)
		}
		b := packet[:n]
		if key := netutil.GetDstKey(b); key != "" {
			if v, ok := s.config.Cache.Get(key); ok {
				err := wsutil.WriteServerBinary(v.(net.Conn), b)
				if err != nil {
					s.config.Cache.Delete(key)
					continue
				}
			}
		}
	}
}

// Authenticate 验证客户端的authcode
func (s *ServerLogic) Authenticate(authCode string) bool {
	return s.config.AuthCode == authCode
}

// HandleConnection 处理客户端连接
func (s *ServerLogic) HandleConnection(client net.Conn) error {
	defer client.Close()
	for {
		recv, err := wsutil.ReadClientBinary(client)
		if err != nil {
			netutil.PrintErr(err, s.config.Verbose)
			break
		}
		if key := netutil.GetSrcKey(recv); key != "" {
			//counter.IncrReadBytes(len(recv))
			_, _ = s.config.IFace.Write(recv)
		}
	}
	return errors.New("")
}

// DistributeCIDR 下发cidr给客户端
func (s *ServerLogic) DistributeCIDR(client net.Conn) (string, error) {
	s.config.Cache.Set("", client, 24*time.Hour)
	addr, mask := s.config.AddressPool.GetAddressFromPool()
	if addr == "" {
		return "", errors.New("no available address in pool")
	}
	bbf := bytes.Buffer{}
	bbf.WriteRune(dhcpMsg)
	bbf.WriteString(addr + "/" + mask)
	bbf.WriteString(strings.Join(s.config.PushRoutes, ","))
	return addr, wsutil.WriteServerMessage(client, ws.OpBinary, bbf.Bytes())
}

// DistributeRote 下发cidr给客户端
func (s *ServerLogic) DistributeRote(client net.Conn) error {
	bbf := bytes.Buffer{}
	bbf.WriteRune(routeMsg)
	bbf.WriteString(strings.Join(s.config.PushRoutes, ","))
	return wsutil.WriteServerMessage(client, ws.OpBinary, bbf.Bytes())
}

// RecycleCIDR 回收下发的cidr
func (s *ServerLogic) RecycleCIDR(addr string) {
	s.config.Cache.Delete(addr)
}
