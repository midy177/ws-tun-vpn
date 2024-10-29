package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
	"ws-tun-vpn/pkg/logview"
	"ws-tun-vpn/pkg/util"
	"ws-tun-vpn/service"
	"ws-tun-vpn/types"
)

var rootCmd = &cobra.Command{
	Use:   "wtvc_gui",
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
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
		lev, err := zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
		if err != nil || lev == zerolog.NoLevel {
			err = nil
			lev = zerolog.InfoLevel
		}
		logger := zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.StampMilli,
		}).With().Timestamp().Logger().Level(lev)
		return service.NewClientService(cmd.Context(), &logV{&logger})
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

type logV struct {
	*zerolog.Logger
}

func (l *logV) Print(level logview.LogLevel, log string) {
	switch level {
	case logview.LogError:
		l.Error().Msg(log)
	case logview.LogInfo:
		l.Info().Msg(log)
	case logview.LogDebug:
		l.Debug().Msg(log)
	case logview.LogWarm:
		l.Warn().Msg(log)
	case logview.LogVerbose:
		l.Debug().Msg(log)
	default:
		l.Info().Msg(log)
	}
}

func (l *logV) GetView() fyne.CanvasObject {
	return nil
}

func (l *logV) SetLogLineSize(maxSize int) {
	return
}

func (l *logV) GetText() *bytes.Buffer {
	return nil
}

func (l *logV) Clear() {
}
