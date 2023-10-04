package services

import (
	"encoding/json"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dcdn "github.com/alibabacloud-go/dcdn-20180115/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/jat001/ddns4cdn/core"
)

type Alibaba core.Alibaba

type alibaba struct {
	Logger *core.LogEntry
	Config *Alibaba
	Client *dcdn.Client
}

func (ctx *alibaba) getDomainInfo() (*dcdn.DescribeDcdnDomainDetailResponseBodyDomainDetail, bool) {
	resp, err := ctx.Client.DescribeDcdnDomainDetail(&dcdn.DescribeDcdnDomainDetailRequest{
		DomainName: tea.String(ctx.Config.Domain),
	})
	s, _ := json.Marshal(resp)
	ctx.Logger.Debug(string(s))
	if err != nil {
		ctx.Logger.Warn("DescribeDcdnDomainDetail failed: ", err)
		return &dcdn.DescribeDcdnDomainDetailResponseBodyDomainDetail{}, false
	}
	return resp.Body.DomainDetail, true
}

func (ctx *alibaba) updateSource(sources string) bool {
	resp, err := ctx.Client.UpdateDcdnDomain(&dcdn.UpdateDcdnDomainRequest{
		DomainName: tea.String(ctx.Config.Domain),
		Sources:    tea.String(sources),
	})
	s, _ := json.Marshal(resp)
	ctx.Logger.Debug(string(s))
	if err != nil {
		ctx.Logger.Warn("UpdateDcdnDomain failed: ", err)
		return false
	}
	return true
}

func (config *Alibaba) Run() bool {
	credential := &openapi.Config{
		AccessKeyId:     &config.AccessKey,
		AccessKeySecret: &config.SecretKey,
	}
	credential.Endpoint = tea.String("dcdn.aliyuncs.com")
	client, _ := dcdn.NewClient(credential)

	ctx := alibaba{
		Logger: core.Log.Logger.WithFields(core.LogFields{
			"module":  "services",
			"service": "alibaba",
		}),
		Config: config,
		Client: client,
	}

	addr := ""

	switch ctx.Config.Protocol {
	case "ipv4":
		addr = ctx.Config.ADDR4

	case "ipv6":
		addr = ctx.Config.ADDR6

	default:
		ctx.Logger.Warn("Unknown protocol:", ctx.Config.Protocol)
		return false
	}

	info, ok := ctx.getDomainInfo()
	if !ok {
		return false
	}

	if config.Protocol == "ipv6" && len(info.Sources.Source) != 2 {
		ctx.Logger.Warn("Only support two sources on ipv6")
		return false
	}
	if config.Protocol == "ipv4" && len(info.Sources.Source) != 1 {
		ctx.Logger.Warn("Only support one source on ipv4")
		return false
	}

	main := -1
	for i, source := range info.Sources.Source {
		if *source.Type != "ipaddr" {
			ctx.Logger.Warn("Only support ipaddr Source")
			return false
		}

		// 20 = main source, 30 = backup source
		if *source.Priority == "20" {
			if main != -1 {
				ctx.Logger.Warn("Only support one main source")
				return false
			}
			main = i
		}
	}

	if main == -1 {
		ctx.Logger.Warn("No main source")
		return false
	}

	if *info.Sources.Source[main].Content == addr {
		ctx.Logger.Info("No need to update")
		return true
	}
	info.Sources.Source[main].Content = tea.String(addr)

	sources, _ := json.Marshal(info.Sources.Source)
	if ok := ctx.updateSource(string(sources)); !ok {
		return false
	}

	return true
}
