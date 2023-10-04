package worker

import (
	"time"

	"github.com/jat001/ddns4cdn/core"
	"github.com/jat001/ddns4cdn/services"
	"github.com/mitchellh/mapstructure"
)

type service struct {
	Logger   *core.LogEntry
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

		case "alibaba":
			service = &services.Alibaba{}

		case "tencent":
			service = &services.Tencent{}

		default:
			ctx.Logger.Error(k, ": unknown service type ", t)
			continue
		}

		c := mapstructure.DecoderConfig{
			Squash: true,
			Result: &service,
		}
		d, e := mapstructure.NewDecoder(&c)
		if e == nil {
			e = d.Decode(v)
		}
		if e != nil {
			ctx.Logger.Error(e)
			continue
		}

		ctx.Services[k] = &service
	}

	ctx.Logger.Debug("Services:", ctx.Services)
}

func (ctx *service) run() {
	for {
		addr6, ok6 := IP(ctx.Config.IP.IPV6, "tcp6")
		addr4, ok4 := IP(ctx.Config.IP.IPV4, "tcp4")
		if !ok6 || !ok4 {
			ctx.Logger.Error("Get IP failed")
			return
		}

		for k, v := range ctx.Services {
			if a, _ := core.Store.RunningService.LoadOrStore(k, false); a.(bool) {
				ctx.Logger.Debug("Skip running service: ", k)
				continue
			}

			s := core.GetRealStruct(v)
			s.FieldByName("ADDR4").SetString(addr4)
			s.FieldByName("ADDR6").SetString(addr6)

			go services.Run(v, k)
		}

		s := time.Duration(ctx.Config.Service.Interval)
		ctx.Logger.Debugf("Wait %d seconds", s)
		time.Sleep(time.Second * s)
	}
}

func Service(config *core.Config) {
	ctx := service{
		Logger: core.Log.Logger.WithFields(core.LogFields{
			"module":   "worker",
			"submoule": "service",
		}),
		Config:   config,
		Services: make(map[string]*services.Services),
	}

	ctx.parseConfig()
	ctx.run()
}
