package services

import (
	"encoding/json"
	"strings"

	"github.com/jat001/ddns4cdn/core"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

type Tencent core.Tencent

type tencent struct {
	Logger *core.LogEntry
	Config *Tencent
	Client *cdn.Client
}

func (ctx *tencent) getDomainInfo() (*cdn.BriefDomain, bool) {
	request := cdn.NewDescribeDomainsRequest()
	request.Filters = []*cdn.DomainFilter{{
		Name:  common.StringPtr("domain"),
		Value: common.StringPtrs([]string{ctx.Config.Domain}),
	}}

	resp, err := ctx.Client.DescribeDomains(request)

	if err != nil {
		ctx.Logger.Warn("DescribeDomains failed: ", err)
		return &cdn.BriefDomain{}, false
	}

	if *resp.Response.TotalNumber == 0 {
		ctx.Logger.Warn("Domain not found")
		return &cdn.BriefDomain{}, false
	}

	return resp.Response.Domains[0], true
}

func (ctx *tencent) updateOrigin(addr string) bool {
	v, _ := json.Marshal(map[string][]string{"update": {addr}})

	request := cdn.NewModifyDomainConfigRequest()
	request.Domain = common.StringPtr(ctx.Config.Domain)
	request.Route = common.StringPtr("Origin.Origins")
	request.Value = common.StringPtr(string(v))

	_, err := ctx.Client.ModifyDomainConfig(request)

	if err != nil {
		ctx.Logger.Warn("ModifyDomainConfig failed: ", err)
		return false
	}

	return true
}

func (config *Tencent) Run() bool {
	credential := common.NewCredential(
		config.AccessKey,
		config.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cdn.tencentcloudapi.com"
	client, _ := cdn.NewClient(credential, "", cpf)

	ctx := tencent{
		Logger: core.Log.Logger.WithFields(core.LogFields{
			"module":  "services",
			"service": "tencent",
		}),
		Config: config,
		Client: client,
	}

	ctx.Logger.Info("Run Tencent service")

	originType := ""
	addr := ""

	switch ctx.Config.Protocol {
	case "ipv4":
		originType = "ip"
		addr = ctx.Config.ADDR4

	case "ipv6":
		originType = "ipv6"
		addr = ctx.Config.ADDR6

	default:
		ctx.Logger.Error("Unknown protocol: ", ctx.Config.Protocol)
		return false
	}

	info, ok := ctx.getDomainInfo()
	if !ok {
		return false
	}

	if *info.Origin.OriginType != originType {
		ctx.Logger.Warn("Origin type not match")
		return false
	}

	if len(info.Origin.Origins) != 1 {
		ctx.Logger.Warn("Only support one origin")
		return false
	}

	sep := ":"
	if originType == "ipv6" {
		sep = "]:"
	}

	s := strings.Split(*info.Origin.Origins[0], sep)
	if len(s) == 2 {
		if originType == "ipv6" {
			addr = "[" + addr + "]"
		}
		addr = addr + ":" + s[1]
	}

	if *info.Origin.Origins[0] == addr {
		ctx.Logger.Info("No need to update")
		return true
	}

	if ok := ctx.updateOrigin(addr); !ok {
		return false
	}

	ctx.Logger.Info("Update origin host success")
	return true
}
