package main

import (
	"flag"
	"fmt"
	"log"
	"wallet-app-server/app/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// Parse commandline flags
	var configPath string
	flag.StringVar(&configPath, "-c", "config.toml", "Configutation file path")
	flag.Parse()

	// Init configuration
	config.LoadConfig(configPath)

	// Create gin app
	r := gin.Default()

	// Start API server
	log.Fatal(r.RunTLS(
		fmt.Sprintf("%s:%d", config.Cfg.Server.Host, config.Cfg.Server.Port),
		config.Cfg.Server.SSLCert,
		config.Cfg.Server.SSLKey,
	))
}
