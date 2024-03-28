package types

import (
	"github.com/patrickmn/go-cache"
	"ws-tun-vpn/pkg/addr_pool"
	"ws-tun-vpn/pkg/water"
)

type BaseConfig struct {
	Verbose   bool
	EnableTLS bool
	MTU       uint
	AuthCode  string
}

type ServerConfig struct {
	BaseConfig
	ListenOn string
	//CIDRv4          string
	AutoCert        bool
	AcmeCert        bool
	Domain          string
	CertificateFile string
	PrivateKeyFile  string
	PushRoutes      []string
	BindAddress     string
	AddressPool     *addr_pool.AddressPool
	Cache           *cache.Cache
	IFace           *water.Interface
}

type ClientConfig struct {
	BaseConfig
	ServerUrl       string
	CertificateFile string
	SkipTLSVerify   bool
	GlobalMode      bool
}
