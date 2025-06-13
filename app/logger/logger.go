package logger

import (
	"log"
	"os"
	"wallet-app-server/app/config"

	"github.com/op/go-logging"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *logging.Logger

func Init() {
	loggingCfg := config.Cfg.Logging
	logLvl, err := logging.LogLevel(loggingCfg.LogLevel)
	if err != nil {
		log.Fatal("Log level is not valid:", err.Error())
	}

	stdoutBackend := logging.AddModuleLevel(logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stdout, "", 0),
		logging.MustStringFormatter(loggingCfg.LogFormat),
	))
	stdoutBackend.SetLevel(logLvl, "")

	logWriter := &lumberjack.Logger{
		Filename: loggingCfg.LogFilePath,
		MaxSize:  loggingCfg.LogFileMaxSizeInMB,
		MaxAge:   loggingCfg.LogFileRetentionInDays,
		Compress: true,
	}
	fileBackend := logging.AddModuleLevel(logging.NewBackendFormatter(
		logging.NewLogBackend(logWriter, "", 0),
		logging.MustStringFormatter(loggingCfg.LogFormat),
	))
	fileBackend.SetLevel(logLvl, "")

	logger = logging.MustGetLogger("fcca-fm-db-engine")
	logging.SetBackend(stdoutBackend, fileBackend)
}

func Debug(args ...any) {
	logger.Debug(args...)
}

func Debugf(format string, args ...any) {
	logger.Debugf(format, args...)
}

func Info(args ...any) {
	logger.Info(args...)
}

func Infof(format string, args ...any) {
	logger.Infof(format, args...)
}

func Warn(args ...any) {
	logger.Warning(args...)
}

func Warnf(format string, args ...any) {
	logger.Warningf(format, args...)
}

func Error(args ...any) {
	logger.Error(args...)
}

func Errorf(format string, args ...any) {
	logger.Errorf(format, args...)
}
