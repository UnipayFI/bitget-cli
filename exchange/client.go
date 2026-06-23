// Package exchange wraps the go-bitget UTA SDK with the thin, CLI-friendly
// surface the cmd layer consumes: a single authenticated client constructor,
// request helpers, and TableWriter response models. Bitget's Unified Trading
// Account serves spot and every futures line from one client; the product is
// selected per call via the category parameter.
package exchange

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/UnipayFI/bitget-cli/config"
	bitget "github.com/UnipayFI/go-bitget"
	"github.com/UnipayFI/go-bitget/client"
	bguta "github.com/UnipayFI/go-bitget/uta"
)

// requestTimeout bounds every REST call so a stuck connection fails the command
// instead of hanging the terminal.
const requestTimeout = 30 * time.Second

// Category re-exports the SDK product-line enum so the cmd layer need not import
// the SDK package directly.
type Category = bguta.Category

const (
	CategorySpot        = bguta.CategorySpot
	CategoryMargin      = bguta.CategoryMargin
	CategoryUSDTFutures = bguta.CategoryUSDTFutures
	CategoryCoinFutures = bguta.CategoryCoinFutures
	CategoryUSDCFutures = bguta.CategoryUSDCFutures
)

// Client is the authenticated wrapper around the SDK's UTA client.
type Client struct {
	uta *bguta.UTAClient
}

// BuildOptions assembles the shared SDK client options (auth, proxy, locale,
// base URL, demo, silent logging) from the global config. Both the UTA and the
// classic-account client wrappers build on top of this so credential and
// transport handling stay identical across account systems.
func BuildOptions() []client.Options {
	opts := []client.Options{
		client.WithAuth(config.Config.APIKey, config.Config.APISecret, config.Config.Passphrase),
		client.WithLogger(silentLogger{}),
	}
	if config.Config.Proxy != "" {
		opts = append(opts, client.WithProxy(config.Config.Proxy))
	}
	if config.Config.Locale != "" {
		opts = append(opts, client.WithLocale(config.Config.Locale))
	}
	if config.Config.BaseURL != "" {
		opts = append(opts, client.WithBaseURL(config.Config.BaseURL))
	}
	if config.Config.Demo {
		opts = append(opts, client.WithDemoTrading(true))
	}
	return opts
}

// NewClient builds an authenticated UTA client from the global config and
// best-effort syncs the server clock so signed requests carry an accepted
// timestamp.
func NewClient() *Client {
	c := bitget.NewUTAClient(BuildOptions()...)
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_ = c.SyncServerTime(ctx)
	return &Client{uta: c}
}

// UTA exposes the underlying SDK client for direct service calls.
func (c *Client) UTA() *bguta.UTAClient { return c.uta }

// ctx returns a request-scoped context with the standard timeout. Callers must
// defer the returned cancel func.
func ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), requestTimeout)
}

// ParseCategory resolves a user-supplied category string (with friendly
// aliases) to the SDK's Category enum.
func ParseCategory(s string) (bguta.Category, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "spot":
		return bguta.CategorySpot, nil
	case "margin":
		return bguta.CategoryMargin, nil
	case "usdt", "usdt-futures", "usdtfutures":
		return bguta.CategoryUSDTFutures, nil
	case "coin", "coin-futures", "coinfutures":
		return bguta.CategoryCoinFutures, nil
	case "usdc", "usdc-futures", "usdcfutures":
		return bguta.CategoryUSDCFutures, nil
	default:
		return "", fmt.Errorf("invalid category %q: want one of spot, margin, usdt-futures (usdt), coin-futures (coin), usdc-futures (usdc)", s)
	}
}

// ParseSide resolves an order side string to the SDK enum (case-insensitive).
func ParseSide(s string) (bguta.Side, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "buy":
		return bguta.SideBuy, nil
	case "sell":
		return bguta.SideSell, nil
	default:
		return "", fmt.Errorf("invalid side %q: want buy or sell", s)
	}
}

// ParseOrderType resolves an order-type string to the SDK enum (case-insensitive).
func ParseOrderType(s string) (bguta.OrderType, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "limit":
		return bguta.OrderTypeLimit, nil
	case "market":
		return bguta.OrderTypeMarket, nil
	default:
		return "", fmt.Errorf("invalid order type %q: want limit or market", s)
	}
}

// silentLogger discards all SDK logging. The SDK would otherwise echo request
// errors (which the cmd layer already surfaces cleanly) and info chatter to
// stdout, polluting output — especially under --json.
type silentLogger struct{}

func (silentLogger) Infof(string, ...any)  {}
func (silentLogger) Warnf(string, ...any)  {}
func (silentLogger) Debugf(string, ...any) {}
func (silentLogger) Errorf(string, ...any) {}
