package main

import (
	"flag"
	"fmt"
	"github.com/Azumi67/DirectTunnelSec/clienttoml"
	"github.com/sirupsen/logrus"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "TOML config filepath needed")
	flag.Parse()

	if configFile == "" {
		logrus.Fatal("Config filepath is required")
	}

	clienttoml.iranClient(configFile)
}
