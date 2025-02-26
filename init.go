package main

import (
	"encoding/json"

	"github.com/48Club/backup-server-dns-tools/rpc"
)

func init() {
	err := json.Unmarshal(_config, &config)
	if err != nil {
		panic(err)
	}
	master = rpc.NewAliveCheck(config.Master)
	go rpc.LoopCheckAlive(master)
}
