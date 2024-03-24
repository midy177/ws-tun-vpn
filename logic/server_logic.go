package logic

import (
	"context"
	"errors"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"log"
	"net"
	"time"
	"ws-tun-vpn/pkg/cachex"
	"ws-tun-vpn/pkg/counter"
	"ws-tun-vpn/pkg/netutil"
	"ws-tun-vpn/types"
)

func ServerTunPacketRouteToClient(ctx context.Context) error {

	return nil
}

func ClientTunPacketToServer(ctx context.Context, client net.Conn) error {
	config, ok := ctx.Value("config").(*types.ServerConfig)
	if !ok {
		return errors.New("failed to get config from context")
	}
	defer client.Close()
	for {
		b, op, err := wsutil.ReadClientData(client)
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			break
		}
		if op == ws.OpText {
			if config.Verbose {
				log.Println(string(b[:]))
			}
			wsutil.WriteServerMessage(client, op, b)
		} else if op == ws.OpBinary {
			if key := netutil.GetSrcKey(b); key != "" {
				cachex.GetCache().Set(key, client, 24*time.Hour)
				counter.IncrReadBytes(len(b))
				config.IFace.Write(b)
			}
		}
	}
	return errors.New("")
}
