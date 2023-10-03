package worker

import (
	"net"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/jat001/ddns4cdn/core"
)

type ip struct {
	Logger  *core.LogEntry
	Client  *resty.Client
	Network string
}

func (ctx *ip) getAddr(url string) (string, bool) {
	resp, err := ctx.Client.R().Get(url)
	if err != nil {
		ctx.Logger.Error("Get IP failed: ", err)
		return "", false
	}

	address := resp.String()
	ctx.Logger.Info(ctx.Network, " IP address: ", address)
	return address, true
}

func IP(url []string, network string) (string, bool) {
	client := resty.New().
		SetTransport(&http.Transport{
			Dial: func(_, addr string) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		}).
		SetRetryCount(3)

	ctx := ip{
		Logger: core.Log.Logger.WithFields(core.LogFields{
			"module":   "worker",
			"submoule": "ip",
		}),
		Client:  client,
		Network: network,
	}

	for _, v := range url {
		if addr, ok := ctx.getAddr(v); ok {
			return addr, true
		}
	}
	return "", false
}
