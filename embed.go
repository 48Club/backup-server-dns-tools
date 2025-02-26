package main

import (
	_ "embed"

	"github.com/48Club/backup-server-dns-tools/types"
)

//go:embed config.json
var _config []byte

var config = types.Config{}
