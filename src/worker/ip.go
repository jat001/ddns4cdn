package worker

import (
	"net"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/jat001/ddns4cdn/core"
)

type ip struct {
	Logger core.LogEntry
}

func IP(url, network string) string {
	client := resty.New().
		SetTransport(&http.Transport{
			Dial: func(_, addr string) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		}).
		SetRetryCount(3)

	ctx := ip{
		Logger: core.Logger.WithFields(core.LogFields{
			"module":   "worker",
			"submoule": "ip",
		}),
	}
	resp, err := client.R().Get(url)
	if err != nil {
		ctx.Logger.Error("Get IP failed:", err)
		return ""
	}

	address := resp.String()
	ctx.Logger.Info("Get IP", address)
	return address
}
