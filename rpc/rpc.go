package rpc

import (
	"context"
	"time"

	"github.com/48Club/backup-server-dns-tools/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func NewAliveCheck(s types.Server) *types.RPC {
	ec, err := ethclient.Dial(s.RPC)
	r := &types.RPC{
		Client: ec,
		Retry:  0,
		Alive:  true,
	}
	if err != nil {
		r.Retry++
	}
	return r
}

func LoopCheckAlive(r *types.RPC) {
	tc := time.NewTicker(500 * time.Millisecond)
	for {
		<-tc.C
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		_, err := r.Client.NetworkID(ctx)
		cancel()
		if err != nil {
			r.Retry++
		} else {
			r.Alive = true
			r.Retry = 0
		}
		if r.Retry > 3 {
			r.Alive = false
		}
	}
}
