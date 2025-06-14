package db

import (
	"fmt"
	"os"
	"wallet-app-server/app/config"
	"wallet-app-server/app/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// A global GORM DB
var DB *gorm.DB

// Init DB, must be called before using the DB
func Init() {
	dbConf := config.Cfg.DB
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=%s",
		dbConf.Host, dbConf.Port, dbConf.Username, dbConf.Password, dbConf.DBName, dbConf.Schema, dbConf.SSLMode)
	logger.Debug("DB connection string:", connString)
	db, err := gorm.Open(postgres.Open(connString))
	if err != nil {
		logger.Error("DB init error: ", err.Error())
		os.Exit(-1)
	}
	DB = db
	logger.Info("DB init sucess")
}
