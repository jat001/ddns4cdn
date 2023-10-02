package worker

import (
	"reflect"
	"time"

	"github.com/jat001/ddns4cdn/core"
	"github.com/jat001/ddns4cdn/services"
	"github.com/mitchellh/mapstructure"
)

type service struct {
	Logger   core.LogEntry
	Config   *core.Config
	Services map[string]*services.Services
}

func (ctx *service) parseConfig() {
	for k, v := range ctx.Config.Services {
		t := v["type"].(string)

		service := *new(services.Services)

		switch t {
		case "cloudflare":
			service = &services.Cloudflare{}

		default:
			ctx.Logger.Warn(k, ": unknown service type ", t)
			continue
		}

		mapstructure.Decode(v, &service)
		ctx.Services[k] = &service
	}

	ctx.Logger.Debug("Services:", ctx.Services)
}

func (ctx *service) run() {
	m := make(map[string]bool)
	c := make(chan string)
	for {
		addr6 := IP(ctx.Config.IP.IPV6, "tcp6")
		addr4 := IP(ctx.Config.IP.IPV4, "tcp4")

		for k, v := range ctx.Services {
			if val, ok := m[k]; ok && val {
				ctx.Logger.Debug("Skip running service: ", k)
				continue
			}
			// ptr -> interface -> ptr -> struct
			s := reflect.ValueOf(v).Elem().Elem().Elem()
			s.FieldByName("ADDR4").SetString(addr4)
			s.FieldByName("ADDR6").SetString(addr6)
			m[k] = true
			go services.Run(v, k, c)
		}

		select {
		case k := <-c:
			m[k] = false

		default:
			ctx.Logger.Debug("All services running")
		}

		ctx.Logger.Debugf("Wait %d seconds", ctx.Config.Service.Interval)
		time.Sleep(time.Second * time.Duration(ctx.Config.Service.Interval))
	}
}

func Service(config *core.Config) {
	ctx := service{
		Logger: core.Logger.WithFields(core.LogFields{
			"module":   "worker",
			"submoule": "service",
		}),
		Config:   config,
		Services: make(map[string]*services.Services),
	}

	ctx.parseConfig()
	ctx.run()
}
