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
	Config.Proxy = proxyFromEnv()
	Config.Locale = os.Getenv("BITGET_LOCALE")
	Config.BaseURL = os.Getenv("BITGET_BASE_URL")
	if v := os.Getenv("BITGET_DEMO"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			Config.Demo = b
		}
	}
}

// proxyFromEnv resolves the REST proxy from the standard proxy environment
// variables. All Bitget REST traffic is HTTPS, so the scheme-specific
// HTTPS_PROXY wins, then the ALL_PROXY catch-all, then HTTP_PROXY as a
// last resort; each is honored in both upper- and lower-case spelling.
func proxyFromEnv() string {
	for _, name := range []string{
		"HTTPS_PROXY", "https_proxy",
		"ALL_PROXY", "all_proxy",
		"HTTP_PROXY", "http_proxy",
	} {
		if v := os.Getenv(name); v != "" {
			return v
		}
	}
	return ""
}
