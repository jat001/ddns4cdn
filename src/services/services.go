package services

import (
	"time"

	"github.com/jat001/ddns4cdn/core"
)

type Services interface {
	Run() bool
}

func Run(ctx *Services, id string) {
	logger := core.Log.Logger.WithFields(core.LogFields{
		"module":   "services",
		"submoule": "services",
	})

	running := core.Store.RunningService

	if ok := running.CompareAndSwap(id, false, true); !ok {
		logger.Warn("Skip running service: ", id)
		return
	}

	ok := (*ctx).Run()

	c := core.GetRealStruct(ctx)
	core.Store.ServiceChan <- &core.ServiceStats{
		ID:      id,
		Type:    c.FieldByName("Type").String(),
		Success: ok,
		EndTime: time.Now(),
	}

	if ok := running.CompareAndSwap(id, true, false); !ok {
		logger.Error("This should not happen. Please report this issue with logs")
	}
}
