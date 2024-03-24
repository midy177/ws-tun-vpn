package service

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
	"ws-tun-vpn/types"
)

func NewServerService(ctx context.Context) error {
	if cfg, ok := ctx.Value("config").(*types.ServerConfig); ok {
		fmt.Printf("%v\n", cfg.Verbose)
		// Todo run service ....
		fmt.Println("run hugo...")
	}
	go func(ctx context.Context) {
		<-ctx.Done()
		os.Exit(0)
	}(ctx)
	//checkCertificateStatusLogic := logic.NewCheckCertificateStatusLogic(svcCtx)
	//logrus.Infof("Start running the scheduled check certificate and private key task.")
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		logrus.Infof("Start task.")
		//checkCertificateStatusLogic.StartCheck()
		ticker.Reset(time.Hour * 6)
		logrus.Infof("The current inspection task is completed.")
	}
	return nil
}
func NewClientService(ctx context.Context) error {
	if cfg, ok := ctx.Value("config").(*types.ClientConfig); ok {
		fmt.Printf("%v\n", cfg.Verbose)
		// Todo run service ....
		fmt.Println("run hugo...")
	}
	go func(ctx context.Context) {
		<-ctx.Done()
		os.Exit(0)
	}(ctx)
	//checkCertificateStatusLogic := logic.NewCheckCertificateStatusLogic(svcCtx)
	//logrus.Infof("Start running the scheduled check certificate and private key task.")
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		logrus.Infof("Start task.")
		//checkCertificateStatusLogic.StartCheck()
		ticker.Reset(time.Hour * 6)
		logrus.Infof("The current inspection task is completed.")
	}
	return nil
}
