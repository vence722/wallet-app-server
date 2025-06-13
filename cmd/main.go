package main

import (
	"flag"
	"wallet-app-server/app"
)

func main() {
	// Parse commandline flags
	var configPath string
	flag.StringVar(&configPath, "-c", "config.toml", "Configutation file path")
	flag.Parse()

	// Init & start application
	app.InitAndStart(configPath)
}
