package main

import (
	"flag"
	"log"

	"github.com/Azumi67/DirectTunnelSec/servertomlserver_fob_toml" 
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "TOML config filepath needed")
	flag.Parse()

	if configFile == "" {
		log.Fatal("Config filepath is required")
	}

	servertoml.kharejServer(configFile)
}
