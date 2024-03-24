package ws

import (
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/golang/snappy"
	"log"
	"net"
	"net/http"
	"time"
	"ws-tun-vpn/pkg/cachex"
	"ws-tun-vpn/pkg/cipher"
	"ws-tun-vpn/pkg/counter"
	"ws-tun-vpn/pkg/netutil"
	"ws-tun-vpn/pkg/water"
	"ws-tun-vpn/types"
)

// StartServer starts the ws server
func StartServer(iFace *water.Interface, config types.ServerConfig) {
	// server -> client
	go toClient(config, iFace)
	// client -> server
	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		if !checkPermission(w, r, config) {
			return
		}
		wsconn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			log.Printf("[server] failed to upgrade http %v", err)
			return
		}
		toServer(config, wsconn, iFace)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", "11")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Cache-Control", "no-cachex")
		w.Header().Set("CF-Cache-Status", "DYNAMIC")
		w.Header().Set("Server", "ws-tun-vpn")
		w.Write([]byte(`I'm health!`))
	})

	log.Printf("tun websocket vpn server started on -> %v", config.ListenOn)

	if config.EnableTLS && !config.AutoCert {
		err := http.ListenAndServeTLS(config.ListenOn, config.CertificateFile, config.PrivateKeyFile, nil)
		if err != nil {
			log.Fatalln(err)
		}

	} else if config.EnableTLS && config.AutoCert {
		if config.AcmeCert {

		} else {

		}
	} else {
		err := http.ListenAndServe(config.ListenOn, nil)
		if err != nil {
			log.Fatalln(err)
		}
	}

}

// checkPermission checks the permission of the request
func checkPermission(w http.ResponseWriter, req *http.Request, config config.Config) bool {
	if config.Key == "" {
		return true
	}
	key := req.Header.Get("key")
	if key != config.Key {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("No permission"))
		return false
	}
	return true
}

// toClient sends data to client
func toClient(config config.Config, iFace *water.Interface) {
	packet := make([]byte, config.BufferSize)
	for {
		n, err := iFace.Read(packet)
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			break
		}
		b := packet[:n]
		if key := netutil.GetDstKey(b); key != "" {
			if v, ok := cachex.GetCache().Get(key); ok {
				//if config.Obfs {
				//	b = cipher.XOR(b)
				//}
				//if config.Compress {
				//	b = snappy.Encode(nil, b)
				//}
				err := wsutil.WriteServerBinary(v.(net.Conn), b)
				if err != nil {
					cachex.GetCache().Delete(key)
					continue
				}
				counter.IncrWrittenBytes(n)
			}
		}
	}
}

// toServer sends data to server
func toServer(config config.Config, wsconn net.Conn, iFace *water.Interface) {
	defer wsconn.Close()
	for {
		b, op, err := wsutil.ReadClientData(wsconn)
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			break
		}
		if op == ws.OpText {
			if config.Verbose {
				log.Println(string(b[:]))
			}
			wsutil.WriteServerMessage(wsconn, op, b)
		} else if op == ws.OpBinary {
			if config.Compress {
				b, _ = snappy.Decode(nil, b)
			}
			if config.Obfs {
				b = cipher.XOR(b)
			}
			if key := netutil.GetSrcKey(b); key != "" {
				cachex.GetCache().Set(key, wsconn, 24*time.Hour)
				counter.IncrReadBytes(len(b))
				iFace.Write(b)
			}
		}
	}
}
