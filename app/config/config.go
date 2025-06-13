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
		Host    string `toml:"host"`
		Port    int    `toml:"port"`
		SSLCert string `toml:"ssl-cert"`
		SSLKey  string `toml:"ssl-key"`
	}
	Logging struct {
		LogLevel               string `toml:"log-level"`
		LogFilePath            string `toml:"log-file-path"`
		LogFormat              string `toml:"log-format"`
		LogFileMaxSizeInMB     int    `toml:"log-file-max-size-in-mb"`
		LogFileRetentionInDays int    `toml:"log-file-retention-in-days"`
	}
	DB struct {
		Host     string `toml:"host"`
		Port     int    `toml:"port"`
		Schema   string `toml:"schema"`
		Username string `toml:"username"`
		Password string `toml:"password"`
		SSLMode  string `toml:"sslmode"`
	}
	Redis struct {
		Addr     string `toml:"addr"`
		Password string `toml:"password"`
		DB       int    `toml:"db"`
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
