package logger

import (
	"io"
	"os"

	"gourbot/internal/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger initializes the logger with log rotation using the provided Config.
func InitLogger(cfg *config.Config) *logrus.Logger {
	logger := logrus.New()

	// Set log output to a rotating file
	logOutputs := []io.Writer{
		&lumberjack.Logger{
			Filename:   cfg.LogFilename,
			MaxSize:    cfg.LogMaxSize,    // Max megabytes before log is rotated
			MaxBackups: cfg.LogMaxBackups, // Max number of old log files to keep
			MaxAge:     cfg.LogMaxAge,     // Max number of days to retain old log files
			Compress:   cfg.LogCompress,   // Compress the old log files
		},
	}

	// If LogStdout is enabled, add stdout to the log outputs
	if cfg.LogStdout {
		logOutputs = append(logOutputs, os.Stdout)
	}

	logger.SetOutput(io.MultiWriter(logOutputs...))

	// Set log format to JSON
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Set log level to Info
	logger.SetLevel(logrus.InfoLevel)

	return logger
}
