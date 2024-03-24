package ws

import (
	"context"
	"github.com/golang/snappy"
	"log"
	"net"
	"time"
	"ws-tun-vpn/config"
	"ws-tun-vpn/pkg/cachex"
	"ws-tun-vpn/pkg/cipher"
	"ws-tun-vpn/pkg/counter"
	"ws-tun-vpn/pkg/netutil"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/net-byte/water"
)

const ConnTag = "conn"

var _ctx context.Context
var _cancel context.CancelFunc

func StartClientForApi(config config.Config, outputStream <-chan []byte, inputStream chan<- []byte, writeCallback, readCallback func(int), _ctx context.Context) {
	go tunToWs(config, outputStream, _ctx, writeCallback)
	for xtun.ContextOpened(_ctx) {
		ctx, cancel := context.WithCancel(_ctx)
		conn := netutil.ConnectServer(config.ServerAddr, "", false)
		if conn == nil {
			cancel()
			time.Sleep(3 * time.Second)
			continue
		}
		cachex.GetCache().Set(ConnTag, conn, 24*time.Hour)
		go wsToTun(config, conn, inputStream, ctx, cancel, readCallback)
		ping(conn, config, ctx, cancel)
		cachex.GetCache().Delete(ConnTag)
		conn.Close()
	}
}

// StartClient starts the ws client
func StartClient(iFace *water.Interface, config config.Config) {
	log.Println("vtun websocket client started")
	_ctx, _cancel = context.WithCancel(context.Background())
	outputStream := make(chan []byte)
	go xtun.ReadFromTun(iFace, config, outputStream, _ctx)
	inputStream := make(chan []byte)
	go xtun.WriteToTun(iFace, config, inputStream, _ctx)
	StartClientForApi(
		config, outputStream, inputStream,
		func(n int) { counter.IncrWrittenBytes(n) },
		func(n int) { counter.IncrReadBytes(n) },
		_ctx,
	)
}

func ping(conn net.Conn, config config.Config, _ctx context.Context, _cancel context.CancelFunc) {
	defer _cancel()
	for xtun.ContextOpened(_ctx) {
		err := wsutil.WriteClientMessage(conn, ws.OpText, []byte("ping"))
		if err != nil {
			break
		}
		time.Sleep(3 * time.Second)
	}
}

// wsToTun sends packets from ws to tun
func wsToTun(config config.Config, conn net.Conn, inputStream chan<- []byte, _ctx context.Context, _cancel context.CancelFunc, callback func(int)) {
	defer _cancel()
	for xtun.ContextOpened(_ctx) {
		packet, err := wsutil.ReadServerBinary(conn)
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			break
		}
		n := len(packet)
		if config.Compress {
			packet, _ = snappy.Decode(nil, packet)
		}
		if config.Obfs {
			packet = cipher.XOR(packet)
		}
		inputStream <- packet[:]
		callback(n)
	}
}

// tunToWs sends packets from tun to ws
func tunToWs(config config.Config, outputStream <-chan []byte, _ctx context.Context, callback func(int)) {
	for xtun.ContextOpened(_ctx) {
		b := <-outputStream
		n := len(b)
		if v, ok := cachex.GetCache().Get(ConnTag); ok {
			if config.Obfs {
				b = cipher.XOR(b)
			}
			if config.Compress {
				b = snappy.Encode(nil, b)
			}
			conn := v.(net.Conn)
			if err := wsutil.WriteClientBinary(conn, b); err != nil {
				netutil.PrintErr(err, config.Verbose)
				continue
			}
			callback(n)
		}
	}
}

func Close() {
	_cancel()
}
