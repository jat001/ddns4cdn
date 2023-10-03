package services

import (
	"sync"

	"github.com/jat001/ddns4cdn/core"
)

type Services interface {
	Run()
}

func Run(ctx *Services, id string, m *sync.Map) {
	logger := core.Logger.WithFields(core.LogFields{
		"module":   "services",
		"submoule": "services",
	})

	if ok := m.CompareAndSwap(id, false, true); !ok {
		logger.Warn("Skip running service: ", id)
		return
	}

	(*ctx).Run()

	m.Store(id, false)
}
