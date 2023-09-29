package worker

import (
	"time"

	"github.com/jat001/ddns4cdn/core"
	"github.com/jat001/ddns4cdn/services"
	"github.com/mitchellh/mapstructure"
)

type service struct {
	logger   core.LogEntry
	config   *core.Config
	services map[string]*services.Services
}

func (ctx *service) parseConfig() {
	for k, v := range ctx.config.Services {
		ctx.logger.Debug("Service name:", k)
		t := v["type"].(string)

		service := *new(services.Services)

		switch t {
		case "cloudflare":
			service = &services.Cloudflare{}

		default:
			ctx.logger.Error("Unknown service type:", t)
			continue
		}

		mapstructure.Decode(v, &service)
		ctx.services[k] = &service
	}

	ctx.logger.Debug("Services:", ctx.services)
}

func (ctx *service) run() {
	m := make(map[string]bool)
	c := make(chan string)
	for {
		for k, v := range ctx.services {
			if val, ok := m[k]; ok && val {
				ctx.logger.Debug("Skip running service:", k)
				continue
			}
			m[k] = true
			go services.Start(v, k, c)
		}

		select {
		case k := <-c:
			m[k] = false

		default:
			ctx.logger.Debug("All services running")
		}

		ctx.logger.Debugf("Wait %d seconds", ctx.config.Service.Interval)
		time.Sleep(time.Second * time.Duration(ctx.config.Service.Interval))
	}
}

func Service(config *core.Config, addr4, addr6 string) {
	ctx := service{
		logger: core.Logger.WithFields(core.LogFields{
			"module":   "worker",
			"submoule": "service",
		}),
		config:   config,
		services: make(map[string]*services.Services),
	}

	ctx.parseConfig()
	ctx.run()
}
