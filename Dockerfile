FROM golang:1.22.1-alpine3.19 as builder

# Define the project name | 定义项目名称

WORKDIR /build
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w CGO_ENABLED=0 \
    && go env -w GOOS=linux \
    && go env -w GOARCH=amd64 \
    && go env \
    && go mod tidy \
    && go build -ldflags '-s -w' \
    -gcflags="all=-trimpath=${PWD}" \
    -asmflags="all=-trimpath=${PWD}" \
    -ldflags="-s -w" \
    -o wtvs server/cmd.go

#linux/amd64,linux/arm64
FROM --platform=linux/amd64 alpine:latest

WORKDIR /app

COPY --from=builder /build/wtvs /bin/wtvs
RUN chmod +x /bin/wtvs \
    && apk add --no-cache iptables iptables-legacy \
    &&  rm /sbin/iptables \
    &&  ln -s /sbin/iptables-legacy /sbin/iptables \
    &&  apk add --no-cache ca-certificates bash iproute2 tzdata

CMD ["wtvs", "-h"]
