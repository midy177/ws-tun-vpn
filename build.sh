#!/usr/bin/env bash

set -x


Version='v1.7.2'

GitHash=`git rev-parse HEAD`

BuildTime=`date +'%Y.%m.%d.%H%M%S'`


LDFlags=" \
    -X 'ws-tun-vpn/pkg.Version=${Version}' \
    -X 'ws-tun-vpn/pkg.GitHash=${GitHash}' \
    -X 'ws-tun-vpn/pkg.BuildTime=${BuildTime}' \
"

go build -ldflags "$LDFlags" -o ws-tun-vpn-server server/main.go
go build -ldflags "$LDFlags" -o ws-tun-vpn-client client/main.go
echo 'build done.'


CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o wtvs server/cmd.go

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o wtvc client/cmd.go


docker buildx build --platform linux/amd64,linux/arm64/v8 --rm --push -t 1228022817/ws-tun-vpn-server:latest -f Dockerfile .
