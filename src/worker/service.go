package worker

import (
	"reflect"
	"sync"
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

func (ctx *service) run(m *sync.Map) {
	for {
		addr6, ok6 := IP(ctx.Config.IP.IPV6, "tcp6")
		addr4, ok4 := IP(ctx.Config.IP.IPV4, "tcp4")
		if !ok6 || !ok4 {
			ctx.Logger.Error("Get IP failed")
			return
		}

		for k, v := range ctx.Services {
			if a, _ := m.LoadOrStore(k, false); a.(bool) {
				ctx.Logger.Debug("Skip running service: ", k)
				continue
			}
			// ptr -> interface -> ptr -> struct
			s := reflect.ValueOf(v).Elem().Elem().Elem()
			s.FieldByName("ADDR4").SetString(addr4)
			s.FieldByName("ADDR6").SetString(addr6)
			go services.Run(v, k, m)
		}

		s := time.Duration(ctx.Config.Service.Interval)
		ctx.Logger.Debugf("Wait %d seconds", s)
		time.Sleep(time.Second * s)
	}
}

func Service(config *core.Config, m *sync.Map) {
	ctx := service{
		Logger: core.Logger.WithFields(core.LogFields{
			"module":   "worker",
			"submoule": "service",
		}),
		Config:   config,
		Services: make(map[string]*services.Services),
	}

	ctx.parseConfig()
	ctx.run(m)
}
