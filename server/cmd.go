package main

import (
	"context"
	"errors"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/cobra"
	"log"
	"time"
	"ws-tun-vpn/pkg/addr_pool"
	"ws-tun-vpn/pkg/netutil"
	"ws-tun-vpn/pkg/util"
	"ws-tun-vpn/service"
	"ws-tun-vpn/types"
)

var rootCmd = &cobra.Command{
	Use:   "wtvs",
	Short: "Websocket tun vpn",
	Long:  `A simple VPN based on websocket and tun devices written in Go.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, ok := cmd.Context().Value("config").(*types.ServerConfig)
		if !ok {
			return errors.New("config not found in context")
		}
		util.ValidateWithFatal(config.AuthCode, "required", "--auth-code")
		if len(config.PushRoutes) > 0 {
			for v := range config.PushRoutes {
				util.ValidateWithFatal(v, "required,cidrv4", "--push-routes")
			}
		}
		if config.EnableTLS && !config.AutoCert {
			util.FlagRequiredWithFatal(cmd, "certificate-file")
			util.FlagRequiredWithFatal(cmd, "private-key-file")
		}
		if config.EnableTLS && config.AutoCert {
			util.ValidateWithFatal(config.Domain, "required,fqdn", "--domain")
		}
		if config.EnableTLS && !config.AutoCert {
			util.ValidateWithFatal(config.CertificateFile, "required,filepath", "--certificate-file")
			util.ValidateWithFatal(config.PrivateKeyFile, "required,filepath", "--private-key-file")
		}
		return service.NewServerService(cmd.Context())
	},
}

func main() {
	config := new(types.ServerConfig)
	rootCmd.Flags().BoolVar(&config.Verbose, "verbose", false,
		"Print the verbose.")
	rootCmd.Flags().BoolVar(&config.EnableTLS, "enable-tls", false,
		"Use TLS to start server.")
	rootCmd.Flags().StringVar(&config.ListenOn, "listen-on", ":3000",
		"Server listener address.")
	rootCmd.Flags().UintVar(&config.MTU, "mtu", 1500,
		"Maximum Transmission Unit.")
	var cidr string
	rootCmd.Flags().StringVar(&cidr, "cidr", "10.7.7.0/24",
		"Classless Inter-Domain Routing of ipv4.")
	util.ValidateWithFatal(cidr, "required,cidr", "--cidr")
	rootCmd.Flags().BoolVar(&config.AutoCert, "auto_cert", false,
		"Automatically generate a certificate that enables HTTPS server.")
	rootCmd.Flags().BoolVar(&config.AcmeCert, "acme_cert", false,
		"To use ACME to automatically issue certificates, you need to support public "+
			"network access on port 80 and configure correct DNS resolution, otherwise self-signature is used. "+
			"This configuration is valid when enable_tls is enabled.")
	rootCmd.Flags().StringVar(&config.Domain, "domain", "",
		"The domain name that is bound to the server.")
	rootCmd.Flags().StringVar(&config.CertificateFile, "certificate-file", "",
		"The certificate file that the server is bound to.")
	rootCmd.Flags().StringVar(&config.PrivateKeyFile, "private-key-file", "",
		"The private key file corresponding to the certificate bound to the server.")
	config.PushRoutes = *rootCmd.Flags().StringArray("push-routes", []string{},
		"Routes that are pushed to clients.")
	rootCmd.Flags().StringVar(&config.AuthCode, "auth-code", "",
		"The authentication code for the client to connect to the server.")

	cidrSlice := netutil.GetCidrV4SliceWithFatal(cidr)
	config.AddressPool = addr_pool.NewAddressPool(cidrSlice[1:], netutil.GetCidrV4Mask(cidr))
	config.BindAddress = cidrSlice[0] + "/" + netutil.GetCidrV4Mask(cidr)
	config.Cache = cache.New(30*time.Minute, 10*time.Minute)
	rootCmd.SetContext(context.WithValue(context.Background(), "config", config))
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
