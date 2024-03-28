package logic

import (
	"context"
	"errors"
	"ws-tun-vpn/types"
)

type ClientLogic struct {
	ctx    context.Context
	config *types.ClientConfig
}

// NewClientLogic create a new client logic
func NewClientLogic(ctx context.Context) (*ClientLogic, error) {
	config, ok := ctx.Value("config").(*types.ClientConfig)
	if !ok {
		return nil, errors.New("failed to get config from context")
	}
	return &ClientLogic{
		ctx:    ctx,
		config: config,
	}, nil
}

// Start start the client logic
func (c *ClientLogic) Start() {

}
