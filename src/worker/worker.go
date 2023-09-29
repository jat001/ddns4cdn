package worker

import (
	"github.com/jat001/ddns4cdn/core"
	"github.com/pelletier/go-toml/v2"
)

type worker struct {
	Logger core.LogEntry
}

func Start(raw []byte) {
	core.Log.SetFormatter()

	ctx := worker{
		Logger: core.Logger.WithFields(core.LogFields{
			"module":   "worker",
			"submoule": "run",
		}),
	}

	ctx.Logger.Info("Run worker")

	config := core.Config{}
	err := toml.Unmarshal(raw, &config)
	if err != nil {
		ctx.Logger.Fatal(err)
		return
	}
	core.Log.SetLevel(config.Log.Level)

	addr6 := IP(config.IP.IPV6, "tcp6")
	addr4 := IP(config.IP.IPV4, "tcp4")

	Service(&config, addr4, addr6)
}
