package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds the application configuration.
type Config struct {
	OpenAIKey     string
	TGBotToken    string
	MasterUID     int64
	LogFilename   string
	LogMaxSize    int
	LogMaxBackups int
	LogMaxAge     int
	LogCompress   bool
	LogStdout     bool
	DbPath        string
}

// LoadConfig loads configuration from environment variables or .env file.
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	// Determine default prefix for log and database filenames
	execPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	defaultPrefix := strings.TrimSuffix(execPath, filepath.Ext(execPath))

	masterUID, _ := strconv.ParseInt(os.Getenv("GOURBOT_MASTER_UID"), 10, 64)

	config := &Config{
		OpenAIKey:     os.Getenv("GOURBOT_OPENAI_KEY"),
		TGBotToken:    os.Getenv("GOURBOT_TGBOT_TOKEN"),
		MasterUID:     masterUID,
		LogFilename:   getEnvOrDefault("GOURBOT_LOG_FILENAME", defaultPrefix+".log"),
		LogMaxSize:    getEnvAsInt("GOURBOT_LOG_MAX_SIZE", 10),
		LogMaxBackups: getEnvAsInt("GOURBOT_LOG_MAX_BACKUPS", 3),
		LogMaxAge:     getEnvAsInt("GOURBOT_LOG_MAX_AGE", 28),
		LogCompress:   getEnvAsBool("GOURBOT_LOG_COMPRESS", true),
		LogStdout:     getEnvAsBoolFromFirstChar("GOURBOT_LOG_STDOUT", false),
		DbPath:        getEnvOrDefault("GOURBOT_DB_PATH", defaultPrefix+".sqlite"),
	}

	// Validate required fields
	if config.OpenAIKey == "" || config.TGBotToken == "" {
		return nil, fmt.Errorf("missing required environment variables: GOURBOT_OPENAI_KEY, GOURBOT_TGBOT_TOKEN")
	}

	return config, nil
}

// Helper functions to get environment variables with defaults
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsBoolFromFirstChar(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists && len(value) > 0 {
		switch strings.ToUpper(string(value[0])) {
		case "1", "T", "Y":
			return true
		case "0", "F", "N":
			return false
		}
	}
	return defaultValue
}
