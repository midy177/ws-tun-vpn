package logic

import (
	"bytes"
	"context"
	"errors"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"log"
	"net"
	"time"
	"ws-tun-vpn/pkg/netutil"
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
	packet := make([]byte, 0, 2048)
	for {
		n, err := s.config.IFace.Read(packet)
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
		recv, op, err := wsutil.ReadClientData(client)
		if err != nil {
			netutil.PrintErr(err, s.config.Verbose)
			break
		}
		if op == ws.OpText {
			if s.config.Verbose {
				log.Println(string(recv[:]))
			}
			wsutil.WriteServerMessage(client, op, recv)
		} else if op == ws.OpBinary {
			if key := netutil.GetSrcKey(recv); key != "" {
				//counter.IncrReadBytes(len(recv))
				s.config.IFace.Write(recv)
			}
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
	err := wsutil.WriteServerMessage(client, ws.OpBinary, []byte(addr+"/"+mask))
	if err != nil {
		return "", err
	}
	recv, op, err := wsutil.ReadClientData(client)
	if err != nil {
		return "", err
	}
	if op == ws.OpBinary && bytes.Equal(recv, []byte("ok")) {
		return addr, nil
	}
	return "", errors.New("client response not ok")
}

// RecycleCIDR 回收下发的cidr
func (s *ServerLogic) RecycleCIDR(addr string) {
	s.config.Cache.Delete(addr)
}
