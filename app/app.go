package app

import (
	"fmt"
	"os"
	"wallet-app-server/app/config"
	"wallet-app-server/app/db"
	"wallet-app-server/app/logger"
	"wallet-app-server/app/redis"

	"github.com/gin-gonic/gin"
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
