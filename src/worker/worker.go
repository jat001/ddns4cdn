package worker

import (
	"time"

	"github.com/jat001/ddns4cdn/api"
	"github.com/jat001/ddns4cdn/core"
	"github.com/pelletier/go-toml/v2"
)

type worker struct {
	Logger *core.LogEntry
}

func Worker(raw []byte) {
	ctx := worker{
		Logger: core.Log.Logger.WithFields(core.LogFields{
			"module":   "worker",
			"submoule": "worker",
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

	go func() {
		for {
			time.Sleep(time.Second)
			select {
			case s := <-core.Store.ServiceStats:
				ctx.Logger.Debug(s)
				core.Store.ServiceStats2 = append(core.Store.ServiceStats2, s)
			default:
				continue
			}
		}
	}()

	go Service(&config)

	api.API(&config)
}
