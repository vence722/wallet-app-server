package app

import (
	"fmt"
	"os"
	"wallet-app-server/app/config"
	"wallet-app-server/app/db"
	"wallet-app-server/app/logger"
	"wallet-app-server/app/redis"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func InitAndStart(configPath string) {
	// Init configuration
	config.LoadConfig(configPath)

	// Init logger
	logger.Init()

	// Init DB
	db.Init()

	// Init redis
	redis.Init()

	// Special setting for library github.com/shopspring/decimal
	// If set to true, the decimal value will be marshaled to number instead of string
	decimal.MarshalJSONWithoutQuotes = true

	// Create gin app
	r := gin.Default()

	// Config API routes
	configRoutes(r)

	// Start API server
	// Exit with code -1 if server cannot start
	if err := r.RunTLS(
		fmt.Sprintf("%s:%d", config.Cfg.Server.Host, config.Cfg.Server.Port),
		config.Cfg.Server.SSLCert,
		config.Cfg.Server.SSLKey,
	); err != nil {
		logger.Errorf("Failed to start server, err: %s", err.Error())
		os.Exit(-1)
	}
}
