package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jat001/ddns4cdn/worker"
)

func main() {
	config := flag.String("c", "config.toml", "config file path")
	data, err := os.ReadFile(*config)
	if err != nil {
		fmt.Println(err)
		return
	}
	worker.Worker(data)
}
