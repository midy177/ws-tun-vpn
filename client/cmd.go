package main

import (
	"context"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
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
		config, ok := cmd.Context().Value("config").(*types.ClientConfig)
		if !ok {
			return errors.New("config not found in context")
		}
		util.ValidateWithFatal(config.ServerUrl, "required", "--server-url")
		util.ValidateWithFatal(config.AuthCode, "required", "--auth-code")
		if config.Verbose {
			configStr, _ := jsoniter.MarshalToString(config)
			fmt.Printf("push routes to client: %s\n", configStr)
		}
		return service.NewClientService(cmd.Context())
	},
}

func main() {
	config := new(types.ClientConfig)
	rootCmd.Flags().BoolVar(&config.Verbose, "verbose", false, "Print the verbose.")
	rootCmd.Flags().BoolVar(&config.EnableTLS, "enable-tls", false, "Whether TLS is enabled on the server.")
	rootCmd.Flags().StringVar(&config.ServerUrl, "server-url", "", "Server address, for example, wtvs.com.")
	util.FlagRequiredWithFatal(rootCmd, "server-url")
	rootCmd.Flags().UintVar(&config.MTU, "mtu", 1500, "Maximum Transmission Unit.")
	rootCmd.Flags().BoolVar(&config.SkipTLSVerify, "skip-tls-verify", false,
		"Skip the validation of the server-side certificate.")
	rootCmd.Flags().StringVar(&config.CertificateFile, "certificate-file", "",
		"Use the specified certificate to verify the certificate on the server.")
	rootCmd.Flags().StringVar(&config.AuthCode, "auth-code", "",
		"The authentication code for the client to connect to the server.")
	rootCmd.SetContext(context.WithValue(context.Background(), "config", config))
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
