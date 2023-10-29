package worker

import (
	"github.com/jat001/ddns4cdn/worker"
)

func Worker(config []byte) {
	worker.Worker(config)
}
