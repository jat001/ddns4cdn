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

func Receiver(config *core.Config) {
	core.Store.ServiceMap = make(map[string][]*core.ServiceStats, len(config.Services))

	for {
		time.Sleep(time.Second)

		select {
		case s := <-core.Store.ServiceChan:
			if config.Store.Limit > 0 {
				if c := cap(core.Store.ServiceMap[s.ID]); c == 0 {
					core.Store.ServiceMap[s.ID] = make([]*core.ServiceStats, 0, config.Store.Limit)
				} else if l := len(core.Store.ServiceMap[s.ID]); l >= c {
					core.Store.ServiceMap[s.ID] = core.Store.ServiceMap[s.ID][l-c+1:]
				}
			}

			core.Store.ServiceMap[s.ID] = append(core.Store.ServiceMap[s.ID], s)

		default:
			continue
		}
	}
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

	go Receiver(&config)
	go Service(&config)

	api.API(&config)
}
