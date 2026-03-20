package config

import (
	"fmt"
	"os"
	"strings"
)

type ConnectMethod string

const (
	ConnectQR          ConnectMethod = "qr"
	ConnectPairingCode ConnectMethod = "pairing"
)

type Config struct {
	// Connection
	ConnectMethod ConnectMethod
	PhoneNumber   string // Required for pairing code (e.g., "6281234567890")

	// Storage
	DBPath string

	// Bot behavior
	Prefix      string
	OwnerJID    string
	BotName     string
	AutoReconnect bool
}

func Load() (*Config, error) {
	cfg := &Config{
		ConnectMethod: ConnectMethod(getEnv("CONNECT_METHOD", "qr")),
		PhoneNumber:   getEnv("PHONE_NUMBER", ""),
		DBPath:        getEnv("DB_PATH", "./data/sessions.db"),
		Prefix:        getEnv("BOT_PREFIX", "!"),
		OwnerJID:      getEnv("OWNER_JID", ""),
		BotName:       getEnv("BOT_NAME", "GoBot"),
		AutoReconnect: getEnv("AUTO_RECONNECT", "true") == "true",
	}

	// Normalize connect method
	cfg.ConnectMethod = ConnectMethod(strings.ToLower(string(cfg.ConnectMethod)))

	if cfg.ConnectMethod == ConnectPairingCode && cfg.PhoneNumber == "" {
		return nil, fmt.Errorf("PHONE_NUMBER is required when using pairing code method")
	}

	if cfg.ConnectMethod != ConnectQR && cfg.ConnectMethod != ConnectPairingCode {
		return nil, fmt.Errorf("CONNECT_METHOD must be 'qr' or 'pairing', got: %s", cfg.ConnectMethod)
	}

	// Ensure DB directory exists
	dbDir := cfg.DBPath[:strings.LastIndex(cfg.DBPath, "/")]
	if dbDir != "" {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create DB directory: %w", err)
		}
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
