package config

import (
	"os"
	"strconv"
)

// Config holds the runtime configuration assembled from environment variables
// and global CLI flags. Credentials come from the Bitget API-management page;
// all three (key, secret, passphrase) are required to sign private requests.
var Config struct {
	APIKey     string
	APISecret  string
	Passphrase string
	Proxy      string
	Locale     string
	BaseURL    string
	Demo       bool
	OutputJSON bool
}

func init() {
	Config.APIKey = os.Getenv("BITGET_API_KEY")
	Config.APISecret = os.Getenv("BITGET_API_SECRET")
	Config.Passphrase = os.Getenv("BITGET_PASSPHRASE")
	Config.Proxy = os.Getenv("BITGET_PROXY")
	Config.Locale = os.Getenv("BITGET_LOCALE")
	Config.BaseURL = os.Getenv("BITGET_BASE_URL")
	if v := os.Getenv("BITGET_DEMO"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			Config.Demo = b
		}
	}
}
