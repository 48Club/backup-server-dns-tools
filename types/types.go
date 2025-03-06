package types

import "github.com/ethereum/go-ethereum/ethclient"

type Server struct {
	RPC string `json:"rpc"`
	IP  string `json:"ip"`
}

type Config struct {
	Server      string   `json:"server"`
	Master      Server   `json:"master"`
	Backup      string   `json:"backup"`
	RecursiveNS []string `json:"recursive-ns"`
}

type RPC struct {
	Client *ethclient.Client
	Retry  int
	Alive  bool `json:"alive"`
}
