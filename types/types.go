package types

import (
	"github.com/patrickmn/go-cache"
	"ws-tun-vpn/pkg/addr_pool"
	"ws-tun-vpn/pkg/water"
)

type BaseConfig struct {
	Verbose   bool   `json:"verbose,omitempty"`
	EnableTLS bool   `json:"enable_tls,omitempty"`
	MTU       uint   `json:"mtu,omitempty"`
	AuthCode  string `json:"auth_code,omitempty"`
}

type ServerConfig struct {
	BaseConfig
	ListenOn string `json:"listen_on"`
	//CIDRv4          string
	AutoCert        bool                   `json:"auto_cert"`
	AcmeCert        bool                   `json:"acme_cert"`
	Domain          string                 `json:"domain"`
	CertificateFile string                 `json:"certificate_file"`
	PrivateKeyFile  string                 `json:"private_key_file"`
	PushRoutes      []string               `json:"push_routes"`
	PushDns         string                 `json:"push_dns"`
	BindAddress     string                 `json:"bind_address"`
	AddressPool     *addr_pool.AddressPool `json:"-"`
	Cache           *cache.Cache           `json:"-"`
	IFace           *water.Interface       `json:"-"`
}

type ClientConfig struct {
	BaseConfig
	ServerUrl       string `json:"server_url,omitempty"`
	CertificateFile string `json:"certificate_file,omitempty"`
	SkipTLSVerify   bool   `json:"skip_tls_verify,omitempty"`
	GlobalMode      bool   `json:"global_mode,omitempty"`
}
