package services

import (
	"github.com/go-resty/resty/v2"
	"github.com/jat001/ddns4cdn/core"
)

type Cloudflare core.Cloudflare

type cloudflare struct {
	Logger     *core.LogEntry
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

	_, err := ctx.Client.R().
		SetQueryParams(map[string]string{
			"name": ctx.Config.Zone,
		}).
		SetResult(&resp).
		Get(url)

	if err != nil {
		ctx.Logger.Error("Get zone id failed: ", err)
		return "", false
	}

	if len(resp.Result) == 0 {
		ctx.Logger.Error("Zone not found")
		return "", false
	}

	return resp.Result[0].ID, true
}

func (ctx *cloudflare) getRecord(zoneID, recordType string) (record, bool) {
	url := ctx.Entrypoint + "{zoneID}/dns_records"

	resp := struct {
		Result []record `json:"result"`
	}{}

	_, err := ctx.Client.R().
		SetPathParams(map[string]string{
			"zoneID": zoneID,
		}).
		SetQueryParams(map[string]string{
			"name": ctx.Config.Record + "." + ctx.Config.Zone,
			"type": recordType,
		}).
		SetResult(&resp).
		Get(url)

	if err != nil {
		ctx.Logger.Error("Get DNS records failed:", err)
		return record{}, false
	}

	if len(resp.Result) == 0 {
		ctx.Logger.Error("DNS record not found")
		return record{}, false
	}

	return resp.Result[0], true
}

func (ctx *cloudflare) updateRecord(record record, addr string) bool {
	url := ctx.Entrypoint + "{zoneID}/dns_records/{recordID}"

	_, err := ctx.Client.R().
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

	if err != nil {
		ctx.Logger.Error("Update DNS record failed: ", err)
		return false
	}

	return true
}

func (config *Cloudflare) Run() bool {
	client := resty.New().
		SetRetryCount(3).
		SetAuthToken(config.Token)

	ctx := cloudflare{
		Logger: core.Log.Logger.WithFields(core.LogFields{
			"module":  "services",
			"service": "cloudflare",
		}),
		Entrypoint: "https://api.cloudflare.com/client/v4/zones/",
		Config:     config,
		Client:     client,
	}

	ctx.Logger.Info("Run Cloudflare service")

	zoneID, ok := ctx.getZoneID()
	if !ok {
		return false
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
		ctx.Logger.Error("Unknown protocol: ", ctx.Config.Protocol)
		return false
	}

	record, ok := ctx.getRecord(zoneID, recordType)
	if !ok {
		return false
	}

	if record.Content == addr {
		ctx.Logger.Info("No need to update")
		return true
	}

	if ok := ctx.updateRecord(record, addr); !ok {
		return false
	}

	ctx.Logger.Info("Update DNS record success")
	return true
}
