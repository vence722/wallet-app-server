package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// The configuration struct
// All the fields should align to config.toml
type Config struct {
	Server struct {
		Host    string
		Port    int
		SSLCert string
		SSLKey  string
	}
}

// The global configuration
// DO NOT change the value when using it
// They're supposed to be loaded when the app is initialized
var Cfg Config

// Load server configuration
// If failed, server cannot be started, and will print error in stdout
func LoadConfig(configPath string) {
	if _, err := toml.DecodeFile(configPath, &Cfg); err != nil {
		// FATAL ERROR: configuration can't be loaded, server exit
		fmt.Printf("[FATAL] server start failed since configuration cannot be loaded, err: %s\n", err.Error())
		os.Exit(-1)
	}
}
