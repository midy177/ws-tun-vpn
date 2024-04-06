package service

import (
	"context"
	"io"
	"log"
	"net/http"
	"runtime"
	"ws-tun-vpn/handler"
	"ws-tun-vpn/logic"
	"ws-tun-vpn/pkg/loadlib"
	"ws-tun-vpn/pkg/logview"
	"ws-tun-vpn/pkg/util"
	"ws-tun-vpn/pkg/water"
	"ws-tun-vpn/types"
)

func NewServerService(ctx context.Context) error {
	config, ok := ctx.Value("config").(*types.ServerConfig)
	if !ok {
		log.Fatalln("config not found in context")
	}
	wc := water.Config{DeviceType: water.TUN}
	wc.PlatformSpecificParams = water.PlatformSpecificParams{}
	os := runtime.GOOS
	wc.PlatformSpecificParams.Name = util.GenerateTunName(4)
	if os == "windows" {
		wc.PlatformSpecificParams.Network = []string{config.BindAddress}
		wc.PlatformSpecificParams.Mtu = int(config.MTU)
	}
	iFace, err := water.New(wc)
	if err != nil {
		return err
	}
	config.IFace = iFace
	handler.RegisterHandlers(ctx)
	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		_, _ = io.WriteString(w, "ok")
	})
	//eg := errgroup.Group{}
	log.Printf("vtun websocket server started on: %v", config.ListenOn)
	//if config.EnableTLS && !config.AutoCert && !config.AcmeCert {
	//	if config.EnableTLS && config.AutoCert && !config.AcmeCert {
	//	}
	//	if config.EnableTLS && config.AutoCert && config.AcmeCert {
	//		// 设置 autocert.Manager，用于管理证书
	//		m := &autocert.Manager{
	//			Prompt:     autocert.AcceptTOS,
	//			HostPolicy: autocert.HostWhitelist(config.Domain),
	//			Cache:      autocert.DirCache("certs"), // 存储证书的缓存目录
	//		}
	//		m.TLSConfig()
	//		eg.Go(func() error {
	//			log.Println("Make sure that port 80 can be accessed through each bound IP address to obtain the automatically issued certificate.")
	//			return http.ListenAndServe(":80", nil)
	//		})
	//	}
	//	eg.Go(func() error {
	//		return http.ListenAndServeTLS(config.ListenOn, config.CertificateFile, config.PrivateKeyFile, nil)
	//	})
	//} else {
	//	eg.Go(func() error {
	//		return http.ListenAndServe(config.ListenOn, nil)
	//	})
	//}
	return http.ListenAndServe(config.ListenOn, nil)
}

var once uint32

func NewClientService(ctx context.Context, logView logview.LogView) error {
	err := loadlib.LoadTunLib()
	if err != nil {
		return err
	}
	clientLogic, err := logic.NewClientLogic(ctx, logView)
	if err != nil {
		return err
	}
	return clientLogic.Start()
}
