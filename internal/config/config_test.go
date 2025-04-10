package config

import (
	"os"
	"testing"
)

// TestLoadConfig tests the LoadConfig function.
// Boundary conditions:
// - Ensure .env file is loaded if present.
// - Validate required environment variables (GOURBOT_OPENAI_KEY, GOURBOT_TGBOT_TOKEN).
// - Test default values for optional environment variables.
func TestLoadConfig(t *testing.T) {
	// Set up environment variables for testing
	os.Setenv("GOURBOT_OPENAI_KEY", "test_openai_key")
	os.Setenv("GOURBOT_TGBOT_TOKEN", "test_tg_bot_token")
	os.Setenv("GOURBOT_MASTER_UID", "12345")
	os.Setenv("GOURBOT_LOG_MAX_SIZE", "20")
	os.Setenv("GOURBOT_LOG_STDOUT", "true")

	// Clean up after test
	defer func() {
		os.Unsetenv("GOURBOT_OPENAI_KEY")
		os.Unsetenv("GOURBOT_TGBOT_TOKEN")
		os.Unsetenv("GOURBOT_MASTER_UID")
		os.Unsetenv("GOURBOT_LOG_MAX_SIZE")
		os.Unsetenv("GOURBOT_LOG_STDOUT")
	}()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Validate loaded configuration
	if config.OpenAIKey != "test_openai_key" {
		t.Errorf("Expected OpenAIKey to be 'test_openai_key', got '%s'", config.OpenAIKey)
	}
	if config.TGBotToken != "test_tg_bot_token" {
		t.Errorf("Expected TGBotToken to be 'test_tg_bot_token', got '%s'", config.TGBotToken)
	}
	if config.MasterUID != 12345 {
		t.Errorf("Expected MasterUID to be 12345, got %d", config.MasterUID)
	}
	if config.LogMaxSize != 20 {
		t.Errorf("Expected LogMaxSize to be 20, got %d", config.LogMaxSize)
	}
	if !config.LogStdout {
		t.Errorf("Expected LogStdout to be true, got false")
	}
}

// TestGetEnvOrDefault tests the getEnvOrDefault helper function.
// Boundary conditions:
// - Environment variable exists.
// - Environment variable does not exist.
func TestGetEnvOrDefault(t *testing.T) {
	key := "TEST_ENV_VAR"
	defaultValue := "default_value"

	// Test when environment variable exists
	os.Setenv(key, "test_value")
	if value := getEnvOrDefault(key, defaultValue); value != "test_value" {
		t.Errorf("Expected '%s', got '%s'", "test_value", value)
	}
	os.Unsetenv(key)

	// Test when environment variable does not exist
	if value := getEnvOrDefault(key, defaultValue); value != defaultValue {
		t.Errorf("Expected '%s', got '%s'", defaultValue, value)
	}
}

// TestGetEnvAsInt tests the getEnvAsInt helper function.
// Boundary conditions:
// - Environment variable is a valid integer.
// - Environment variable is not set.
// - Environment variable is invalid.
func TestGetEnvAsInt(t *testing.T) {
	key := "TEST_ENV_INT"
	defaultValue := 42

	// Test valid integer
	os.Setenv(key, "100")
	if value := getEnvAsInt(key, defaultValue); value != 100 {
		t.Errorf("Expected 100, got %d", value)
	}
	os.Unsetenv(key)

	// Test unset environment variable
	if value := getEnvAsInt(key, defaultValue); value != defaultValue {
		t.Errorf("Expected %d, got %d", defaultValue, value)
	}

	// Test invalid integer
	os.Setenv(key, "invalid")
	if value := getEnvAsInt(key, defaultValue); value != defaultValue {
		t.Errorf("Expected %d, got %d", defaultValue, value)
	}
	os.Unsetenv(key)
}

// TestGetEnvAsBool tests the getEnvAsBool helper function.
// Boundary conditions:
// - Environment variable is a valid boolean.
// - Environment variable is not set.
// - Environment variable is invalid.
func TestGetEnvAsBool(t *testing.T) {
	key := "TEST_ENV_BOOL"
	defaultValue := true

	// Test valid boolean
	os.Setenv(key, "false")
	if value := getEnvAsBool(key, defaultValue); value != false {
		t.Errorf("Expected false, got %t", value)
	}
	os.Unsetenv(key)

	// Test unset environment variable
	if value := getEnvAsBool(key, defaultValue); value != defaultValue {
		t.Errorf("Expected %t, got %t", defaultValue, value)
	}

	// Test invalid boolean
	os.Setenv(key, "invalid")
	if value := getEnvAsBool(key, defaultValue); value != defaultValue {
		t.Errorf("Expected %t, got %t", defaultValue, value)
	}
	os.Unsetenv(key)
}

// TestGetEnvAsBoolFromFirstChar tests the getEnvAsBoolFromFirstChar helper function.
// Boundary conditions:
// - Environment variable starts with a valid character.
// - Environment variable is not set.
// - Environment variable starts with an invalid character.
func TestGetEnvAsBoolFromFirstChar(t *testing.T) {
	key := "TEST_ENV_BOOL_CHAR"
	defaultValue := false

	// Test valid true character
	os.Setenv(key, "Y")
	if value := getEnvAsBoolFromFirstChar(key, defaultValue); value != true {
		t.Errorf("Expected true, got %t", value)
	}
	os.Unsetenv(key)

	// Test valid false character
	os.Setenv(key, "N")
	if value := getEnvAsBoolFromFirstChar(key, defaultValue); value != false {
		t.Errorf("Expected false, got %t", value)
	}
	os.Unsetenv(key)

	// Test unset environment variable
	if value := getEnvAsBoolFromFirstChar(key, defaultValue); value != defaultValue {
		t.Errorf("Expected %t, got %t", defaultValue, value)
	}

	// Test invalid character
	os.Setenv(key, "invalid")
	if value := getEnvAsBoolFromFirstChar(key, defaultValue); value != defaultValue {
		t.Errorf("Expected %t, got %t", defaultValue, value)
	}
	os.Unsetenv(key)
}
