package logic

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"net"
	"net/http"
	"runtime"
	"strings"
	"time"
	"ws-tun-vpn/pkg/counter"
	"ws-tun-vpn/pkg/logview"
	"ws-tun-vpn/pkg/nic_tool"
	"ws-tun-vpn/pkg/util"
	"ws-tun-vpn/pkg/water"
	"ws-tun-vpn/types"
)

type ClientLogic struct {
	ctx      context.Context
	config   *types.ClientConfig
	conn     net.Conn
	iFace    *water.Interface
	eg       *errgroup.Group
	nicTool  nic_tool.NicTool
	nicReady bool
	logView  logview.LogView
}

// NewClientLogic create a new client logic
func NewClientLogic(ctx context.Context, logView logview.LogView) (*ClientLogic, error) {
	config, ok := ctx.Value("config").(*types.ClientConfig)
	if !ok {
		return nil, errors.New("failed to get config from context")
	}
	errg, childCtx := errgroup.WithContext(ctx)
	return &ClientLogic{
		ctx:     childCtx,
		eg:      errg,
		config:  config,
		logView: logView,
	}, nil
}

// Start the client logic
func (c *ClientLogic) Start() error {
	if err := c.connectServer(); err != nil {
		return err
	}
	c.eg.Go(c.directLoop)
	c.eg.Go(c.receiveLoop)
	c.eg.Go(func() error {
		<-c.ctx.Done()
		_ = c.conn.Close()
		if c.iFace != nil {
			_ = c.iFace.Close()
		}
		return errors.New("主动关闭连接！")
	})
	err := c.eg.Wait()
	return err
}

func (c *ClientLogic) connectServer() error {
	header := make(http.Header)
	header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36")
	header.Set("Authorization", c.config.AuthCode)
	dialer := ws.Dialer{
		Header:  ws.HandshakeHeaderHTTP(header),
		Timeout: 30 * time.Second,
	}

	url := "ws://" + c.config.ServerUrl + "/connect"
	if c.config.EnableTLS {
		url = "wss://" + c.config.ServerUrl + "/connect"
		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: c.config.SkipTLSVerify,
		}
	}
	conn, _, _, err := dialer.Dial(context.Background(), url)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

// direct tun data to server
func (c *ClientLogic) directLoop() error {
	if c.conn == nil {
		return errors.New("connection not established")
	}
	defer c.conn.Close()
	// Reuse the same buffer to reduce memory allocation
	packet := make([]byte, 2048)
	for {
		select {
		case <-c.ctx.Done():
		default:
			if !c.nicReady {
				continue
			}
			n, err := c.iFace.Read(packet)
			// when read err of io.EOF, should continue to the next loop
			if err != nil && err != io.EOF {
				c.logView.Print(logview.LogError, err.Error())
				return err
			}
			if n == 0 {
				continue
			}
			if c.config.Verbose {
				// 打印接收到的包大小
				log.Printf("send packet size: %d", n)
			}
			counter.IncrWrittenBytes(n)
			if err := wsutil.WriteClientBinary(c.conn, packet[:n]); err != nil {
				return err
			}
		}
	}
}

// receive data from server and forward to tun interface
func (c *ClientLogic) receiveLoop() error {
	if c.conn == nil {
		c.logView.Print(logview.LogError, "connection not established")
		return errors.New("connection not established")
	}
	defer c.conn.Close()
	for {
		select {
		case <-c.ctx.Done():
		default:
			data, err := wsutil.ReadServerBinary(c.conn)
			if err != nil && err != io.EOF {
				c.logView.Print(logview.LogError, err.Error())
				return err
			}
			lenData := len(data)
			if lenData == 0 {
				continue
			}
			counter.IncrReadBytes(lenData)
			switch data[0] {
			case dhcpMsg:
				if err := c.handleDhcpMsg(data[1:]); err != nil {
					c.logView.Print(logview.LogError, err.Error())
					return err
				}
			case routeMsg:
				c.handleRouteMsg(data[1:])
			case packetMsg:
				if _, err := c.iFace.Write(data[1:]); err != nil {
					c.logView.Print(logview.LogError, err.Error())
					return err
				}
				if c.config.Verbose {
					// 打印接收到的包大小
					log.Printf("received packet size: %d", len(data[1:]))
				}
			case dnsMsg:
				c.handleDnsMsg(data[1:])
			default:
				c.logView.Print(logview.LogError, fmt.Sprintf("unknown message type: %v", data[0]))
				//return errors.New(fmt.Sprintf("unknown message type: %v", data[0]))
			}
		}
	}
}

func (c *ClientLogic) handleDhcpMsg(cidr []byte) error {
	if c.conn == nil {
		c.logView.Print(logview.LogError, "connection not established")
		return errors.New("connection not established")
	}
	cidrS := string(cidr)
	if c.config.Verbose {
		log.Printf("received dhcp message: %v", cidrS)
	}
	_, _, err := net.ParseCIDR(cidrS)
	if err != nil {
		return err
	}

	wc := water.Config{DeviceType: water.TUN}
	wc.PlatformSpecificParams = water.PlatformSpecificParams{}
	os := runtime.GOOS
	wc.PlatformSpecificParams.Name = util.GenerateTunName(4)
	if os == "windows" {
		wc.PlatformSpecificParams.Network = []string{cidrS}
	}
	iFace, err := water.New(wc)
	if err != nil {
		return err
	}
	c.iFace = iFace
	log.Printf("interface created successfully: %v", iFace.Name())
	c.logView.Print(logview.LogInfo, "interface created successfully: "+iFace.Name())
	c.nicTool = nic_tool.NewNicTool(iFace.Name(), cidrS, int(c.config.MTU))
	info := c.nicTool.SetCidrAndUp()

	c.logView.Print(logview.LogInfo,
		fmt.Sprintf("set tun network card(%s) cidr(%s) and up.\n", iFace.Name(), cidrS))

	if c.config.Verbose && os != "windows" {
		log.Printf("set tun network card(%s) cidr(%s) and up.\n", iFace.Name(), cidrS)
		if len(info) > 0 {
			log.Println(info)
		}
	}
	info = c.nicTool.SetMtu()
	c.logView.Print(logview.LogInfo, fmt.Sprintf("set tun network card(%s) mtu: %d\n", iFace.Name(), c.config.MTU))
	if c.config.Verbose {
		log.Printf("set tun network card(%s) mtu: %d\n", iFace.Name(), c.config.MTU)
		if len(info) > 0 {
			log.Println(info)
		}
	}
	//c.nicTool.SetRoute(cidrS)
	c.nicReady = true
	return nil
}

func (c *ClientLogic) handleRouteMsg(list []byte) {
	routes := strings.Split(string(list), ",")
	for _, route := range routes {
		c.nicTool.SetRoute(route)
		c.logView.Print(logview.LogInfo, fmt.Sprintf("set tun network card(%s) route: %s\n", c.iFace.Name(), route))

	}
}

func (c *ClientLogic) handleDnsMsg(dns []byte) {
	dnsIp := string(dns)
	if net.ParseIP(dnsIp) != nil {
		c.nicTool.SetPrimaryDnsServer(dnsIp)
		c.logView.Print(logview.LogInfo, fmt.Sprintf("set tun network card(%s) dns: %s\n", c.iFace.Name(), string(dns)))
	}
}
