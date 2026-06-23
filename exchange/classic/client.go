// Package classic wraps the go-bitget classic-account SDK (the /api/v2/*
// endpoints) with the thin, CLI-friendly surface the cmd layer consumes. Unlike
// the unified account, classic splits products into separate clients: Spot and
// Futures (Mix) each get their own authenticated client and TableWriter models.
// Credential and transport handling is shared with the UTA wrapper via
// exchange.BuildOptions.
package classic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/UnipayFI/bitget-cli/exchange"
	bitget "github.com/UnipayFI/go-bitget"
	"github.com/UnipayFI/go-bitget/classic/mix"
	"github.com/UnipayFI/go-bitget/classic/spot"
)

// requestTimeout bounds every REST call so a stuck connection fails the command
// instead of hanging the terminal.
const requestTimeout = 30 * time.Second

// SpotClient is the authenticated wrapper around the classic Spot SDK client.
type SpotClient struct {
	sp *spot.SpotClient
}

// NewSpotClient builds an authenticated classic Spot client from the global
// config and best-effort syncs the server clock.
func NewSpotClient() *SpotClient {
	c := bitget.NewSpotClient(exchange.BuildOptions()...)
	cx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_ = c.SyncServerTime(cx)
	return &SpotClient{sp: c}
}

// FuturesClient is the authenticated wrapper around the classic Futures (Mix)
// SDK client.
type FuturesClient struct {
	mx *mix.MixClient
}

// NewFuturesClient builds an authenticated classic Futures (Mix) client from the
// global config and best-effort syncs the server clock.
func NewFuturesClient() *FuturesClient {
	c := bitget.NewMixClient(exchange.BuildOptions()...)
	cx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_ = c.SyncServerTime(cx)
	return &FuturesClient{mx: c}
}

// ctx returns a request-scoped context with the standard timeout. Callers must
// defer the returned cancel func.
func ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), requestTimeout)
}

// ---- shared enum parsers -------------------------------------------------

// ParseProductType resolves a futures product-line string (with friendly
// aliases) to the classic SDK's ProductType enum.
func ParseProductType(s string) (mix.ProductType, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "usdt", "usdt-futures", "usdtfutures":
		return mix.ProductTypeUSDTFutures, nil
	case "coin", "coin-futures", "coinfutures":
		return mix.ProductTypeCoinFutures, nil
	case "usdc", "usdc-futures", "usdcfutures":
		return mix.ProductTypeUSDCFutures, nil
	default:
		return "", fmt.Errorf("invalid product type %q: want one of usdt-futures (usdt), coin-futures (coin), usdc-futures (usdc)", s)
	}
}
