package worker

import (
	"sync"

	"github.com/jat001/ddns4cdn/api"
	"github.com/jat001/ddns4cdn/core"
	"github.com/pelletier/go-toml/v2"
)

type worker struct {
	Logger core.LogEntry
}

func Start(raw []byte) {
	ctx := worker{
		Logger: core.Logger.WithFields(core.LogFields{
			"module":   "worker",
			"submoule": "start",
		}),
	}

	ctx.Logger.Info("Run worker")

	config := core.Config{}
	err := toml.Unmarshal(raw, &config)
	if err != nil {
		ctx.Logger.Fatal(err)
		return
	}

	core.Log.SetLevel(config.Log.Level)

	m := sync.Map{}

	go Service(&config, &m)

	api.API(&config, &m)
}
