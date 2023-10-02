package services

import (
	"github.com/go-resty/resty/v2"
	"github.com/jat001/ddns4cdn/core"
)

type Cloudflare struct {
	ADDR4    string
	ADDR6    string
	Type     string
	Protocol string
	Token    string
	Zone     string
	Record   string
}

type cloudflare struct {
	Logger     core.LogEntry
	Entrypoint string
	Config     *Cloudflare
	Client     *resty.Client
}

type record struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Proxied bool   `json:"proxied"`
	TTL     int    `json:"ttl"`
	ZoneID  string `json:"zone_id"`
}

func (ctx *cloudflare) getZoneID() (string, bool) {
	url := ctx.Entrypoint

	resp := struct {
		Result []struct {
			ID string `json:"id"`
		} `json:"result"`
	}{}

	r, err := ctx.Client.R().
		SetQueryParams(map[string]string{
			"name": ctx.Config.Zone,
		}).
		SetResult(&resp).
		Get(url)

	ctx.Logger.Debug(r.String())

	if err != nil {
		ctx.Logger.Error("Get zone id failed:", err)
		return "", false
	}

	if len(resp.Result) == 0 {
		ctx.Logger.Error("Zone not found:", r.Body())
		return "", false
	}

	return resp.Result[0].ID, true
}

func (ctx *cloudflare) getRecord(zoneID, recordType string) (record, bool) {
	url := ctx.Entrypoint + "{zoneID}/dns_records"

	resp := struct {
		Result []record `json:"result"`
	}{}

	r, err := ctx.Client.R().
		SetPathParams(map[string]string{
			"zoneID": zoneID,
		}).
		SetQueryParams(map[string]string{
			"name": ctx.Config.Record + "." + ctx.Config.Zone,
			"type": recordType,
		}).
		SetResult(&resp).
		Get(url)

	ctx.Logger.Debug(r.String())

	if err != nil {
		ctx.Logger.Error("Get DNS records failed:", err)
		return record{}, false
	}

	if len(resp.Result) == 0 {
		ctx.Logger.Error("Get DNS records failed:", ctx.Config.Zone)
		return record{}, false
	}

	return resp.Result[0], true
}

func (ctx *cloudflare) updateRecord(record record, addr string) {
	url := ctx.Entrypoint + "{zoneID}/dns_records/{recordID}"

	r, err := ctx.Client.R().
		SetPathParams(map[string]string{
			"zoneID":   record.ZoneID,
			"recordID": record.ID,
		}).
		SetBody(map[string]interface{}{
			"type":    record.Type,
			"name":    record.Name,
			"content": addr,
			"proxied": record.Proxied,
			"ttl":     record.TTL,
		}).
		Put(url)

	ctx.Logger.Debug(r.String())

	if err != nil {
		ctx.Logger.Error("Update DNS record failed:", err)
		return
	}
}

func (config *Cloudflare) Run() {
	client := resty.New().
		SetRetryCount(3).
		SetAuthToken(config.Token)

	ctx := cloudflare{
		Logger: core.Logger.WithFields(core.LogFields{
			"module":  "service",
			"service": "cloudflare",
		}),
		Entrypoint: "https://api.cloudflare.com/client/v4/zones/",
		Config:     config,
		Client:     client,
	}

	ctx.Logger.Info("Run Cloudflare service")

	zoneID, ok := ctx.getZoneID()
	if !ok {
		return
	}

	recordType := ""
	addr := ""

	switch ctx.Config.Protocol {
	case "ipv4":
		recordType = "A"
		addr = ctx.Config.ADDR4

	case "ipv6":
		recordType = "AAAA"
		addr = ctx.Config.ADDR6

	default:
		ctx.Logger.Error("Unknown protocol:", ctx.Config.Protocol)
		return
	}

	record, ok := ctx.getRecord(zoneID, recordType)
	if !ok {
		return
	}
	if record.Content == addr {
		ctx.Logger.Info("No need to update")
		return
	}

	ctx.updateRecord(record, addr)

	ctx.Logger.Info("Update DNS record success")
}
