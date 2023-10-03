package services

import (
	"reflect"
	"time"

	"github.com/jat001/ddns4cdn/core"
)

type Services interface {
	Run()
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

	defer func() {
		c := reflect.ValueOf(ctx).Elem().Elem().Elem()
		core.Store.ServiceStats <- &core.ServiceStats{
			ID:      id,
			Type:    c.FieldByName("Type").String(),
			Success: true,
			EndTime: time.Now().Unix(),
		}
		if ok := running.CompareAndSwap(id, true, false); !ok {
			logger.Error("This should not happen. Please report this issue with logs")
		}
	}()

	(*ctx).Run()
}
