package main

import (
	"context"
	"github.com/spf13/cobra"
	"log"
	"ws-tun-vpn/pkg/util"
	"ws-tun-vpn/service"
	"ws-tun-vpn/types"
)

var rootCmd = &cobra.Command{
	Use:   "wtvc",
	Short: "Websocket tun vpn",
	Long:  `A simple VPN based on websocket and tun devices written in Go.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Todo run service ....
		return service.NewClientService(cmd.Context())
	},
}

func main() {
	config := new(types.ClientConfig)
	rootCmd.Flags().BoolVar(&config.Verbose, "verbose", false, "Print the verbose.")
	rootCmd.Flags().BoolVar(&config.EnableTLS, "enable_tls", false, "Whether TLS is enabled on the server.")
	rootCmd.Flags().StringVar(&config.ServerUrl, "server_url", "", "Server address, for example, wtvs.com.")
	util.FlagRequiredWithFatal(rootCmd, "server_url")
	util.ValidateWithFatal(config.ServerUrl, "required,url")
	rootCmd.Flags().UintVar(&config.MTU, "mtu", 1500, "Maximum Transmission Unit.")
	rootCmd.Flags().BoolVar(&config.SkipTLSVerify, "skip_tls_verify", false,
		"Skip the validation of the server-side certificate.")
	rootCmd.Flags().StringVar(&config.CertificateFile, "certificate_file", "",
		"Use the specified certificate to verify the certificate on the server.")
	rootCmd.SetContext(context.WithValue(context.Background(), "config", config))
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
