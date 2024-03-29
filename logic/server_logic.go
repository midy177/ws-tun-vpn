package logic

import (
	"bytes"
	"context"
	"errors"
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
	info = nt.EnableIpForward()
	if s.config.Verbose {
		log.Println("set ip forward enable.")
		if len(info) > 0 {
			log.Println(info)
		}
	}
	info = nt.EnableNat()
	if s.config.Verbose {
		log.Println("set ip nat enable.")
		if len(info) > 0 {
			log.Println(info)
		}
	}
	packet := make([]byte, 2048)
	var buf bytes.Buffer
	for {
		buf.Reset()
		buf.WriteRune(packetMsg)
		n, err := s.config.IFace.Read(packet)
		if err != nil {
			log.Fatalf("failed to read packet: %v", err)
			//netutil.PrintErr(err, s.config.Verbose)
		}
		b := packet[:n]
		if key := netutil.GetDstKey(b); key != "" {
			if v, ok := s.config.Cache.Get(key); ok {
				if s.config.Verbose {
					log.Printf(" start send data to client: %v\n", key)
				}
				buf.Write(b)
				err := wsutil.WriteServerBinary(v.(net.Conn), buf.Bytes())
				if err != nil {
					log.Printf("failed to write data to client: %v,release key: %s\n", err, key)
					s.config.Cache.Delete(key)
					continue
				}
			} else {
				if s.config.Verbose {
					log.Printf("client not found: %s\n", key)
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
			if s.config.Verbose {
				log.Printf("failed to read client data,err: %v\n", err)
			}
			return err
		}
		if key := netutil.GetSrcKey(recv); key != "" {
			//counter.IncrReadBytes(len(recv))
			_, _ = s.config.IFace.Write(recv)
			if s.config.Verbose {
				log.Printf("recv data from client: %s, len: %d\n", key, len(recv))
			}
		}
	}
}

// DistributeCIDR 下发cidr给客户端
func (s *ServerLogic) DistributeCIDR(client net.Conn) (string, error) {
	addr, mask := s.config.AddressPool.GetAddressFromPool()
	if addr == "" {
		return "", errors.New("no available address in pool")
	}
	// store client connection in cache
	s.config.Cache.Set(addr, client, 24*time.Hour)
	var buf bytes.Buffer
	buf.WriteRune(dhcpMsg)
	buf.WriteString(addr + "/" + mask)
	return addr, wsutil.WriteServerBinary(client, buf.Bytes())
}

// DistributeRote 下发cidr给客户端
func (s *ServerLogic) DistributeRote(client net.Conn) error {
	var buf bytes.Buffer
	buf.WriteRune(routeMsg)
	if len(s.config.PushRoutes) == 0 {
		return nil
	}
	buf.WriteString(strings.Join(s.config.PushRoutes, ","))
	return wsutil.WriteServerBinary(client, buf.Bytes())
}

// RecycleCIDR 回收下发的cidr
func (s *ServerLogic) RecycleCIDR(addr string) {
	s.config.Cache.Delete(addr)
}
