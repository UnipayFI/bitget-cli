package classic

import (
	"fmt"
	"strings"

	"github.com/UnipayFI/bitget-cli/common"
	"github.com/UnipayFI/go-bitget/classic/mix"
	"github.com/shopspring/decimal"
)

// ---- enum parsers --------------------------------------------------------

// ParseFuturesSide resolves an order side string to the classic mix enum.
func ParseFuturesSide(s string) (mix.Side, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "buy":
		return mix.SideBuy, nil
	case "sell":
		return mix.SideSell, nil
	default:
		return "", fmt.Errorf("invalid side %q: want buy or sell", s)
	}
}

// ParseFuturesOrderType resolves an order-type string to the classic mix enum.
func ParseFuturesOrderType(s string) (mix.OrderType, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "limit":
		return mix.OrderTypeLimit, nil
	case "market":
		return mix.OrderTypeMarket, nil
	default:
		return "", fmt.Errorf("invalid order type %q: want limit or market", s)
	}
}

// ParseMarginMode resolves a margin-mode string to the classic mix enum.
// An empty value defaults to crossed.
func ParseMarginMode(s string) (mix.MarginMode, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "crossed", "cross":
		return mix.MarginModeCrossed, nil
	case "isolated":
		return mix.MarginModeIsolated, nil
	default:
		return "", fmt.Errorf("invalid margin mode %q: want crossed or isolated", s)
	}
}

// ---- service calls -------------------------------------------------------

func (c *FuturesClient) GetAccountList(pt mix.ProductType) ([]mix.AccountListItem, error) {
	cx, cancel := ctx()
	defer cancel()
	return c.mx.NewGetAccountListService(pt).Do(cx)
}

func (c *FuturesClient) GetAllPositions(pt mix.ProductType, marginCoin string) ([]mix.Position, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.mx.NewGetAllPositionService(pt)
	if marginCoin != "" {
		s.SetMarginCoin(marginCoin)
	}
	return s.Do(cx)
}

func (c *FuturesClient) GetSinglePosition(pt mix.ProductType, symbol, marginCoin string) ([]mix.Position, error) {
	cx, cancel := ctx()
	defer cancel()
	return c.mx.NewGetSinglePositionService(pt, symbol, marginCoin).Do(cx)
}

// FuturesPlaceOrderParams collects the fields for a single classic futures
// order. Price is treated as unset when zero (required only for limit orders);
// the optional fields are applied only when non-empty.
type FuturesPlaceOrderParams struct {
	Symbol      string
	ProductType mix.ProductType
	MarginMode  mix.MarginMode
	MarginCoin  string
	Size        decimal.Decimal
	Side        mix.Side
	OrderType   mix.OrderType
	Price       decimal.Decimal
	Force       string
	TradeSide   string
	ReduceOnly  string
	ClientOid   string
}

func (c *FuturesClient) PlaceOrder(p FuturesPlaceOrderParams) (*mix.OrderRef, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.mx.NewPlaceOrderService(p.Symbol, p.ProductType, p.MarginMode, p.MarginCoin, p.Size, p.Side, p.OrderType)
	if !p.Price.IsZero() {
		s.SetPrice(p.Price)
	}
	if p.Force != "" {
		s.SetForce(mix.Force(strings.ToLower(p.Force)))
	}
	switch strings.ToLower(strings.TrimSpace(p.TradeSide)) {
	case "open":
		s.SetTradeSide(mix.TradeSideOpen)
	case "close":
		s.SetTradeSide(mix.TradeSideClose)
	}
	switch strings.ToLower(strings.TrimSpace(p.ReduceOnly)) {
	case "yes", "true":
		s.SetReduceOnly(mix.ReduceOnlyYes)
	case "no", "false":
		s.SetReduceOnly(mix.ReduceOnlyNo)
	}
	if p.ClientOid != "" {
		s.SetClientOid(p.ClientOid)
	}
	return s.Do(cx)
}

func (c *FuturesClient) CancelOrder(symbol string, pt mix.ProductType, marginCoin, orderID, clientOid string) (*mix.OrderRef, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.mx.NewCancelOrderService(symbol, pt)
	if marginCoin != "" {
		s.SetMarginCoin(marginCoin)
	}
	if orderID != "" {
		s.SetOrderId(orderID)
	}
	if clientOid != "" {
		s.SetClientOid(clientOid)
	}
	return s.Do(cx)
}

func (c *FuturesClient) GetOrderDetail(symbol string, pt mix.ProductType, orderID, clientOid string) (*mix.MixOrder, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.mx.NewGetOrderDetailService(symbol, pt)
	if orderID != "" {
		s.SetOrderId(orderID)
	}
	if clientOid != "" {
		s.SetClientOid(clientOid)
	}
	return s.Do(cx)
}

func (c *FuturesClient) GetOpenOrders(pt mix.ProductType, symbol string, limit int) (*mix.MixOrderList, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.mx.NewGetOrdersPendingService(pt)
	if symbol != "" {
		s.SetSymbol(symbol)
	}
	if limit > 0 {
		s.SetLimit(limit)
	}
	return s.Do(cx)
}

// ---- table models --------------------------------------------------------

// FuturesAccountRows renders the classic futures per-margin-coin account
// balances and equity.
type FuturesAccountRows []mix.AccountListItem

func (a FuturesAccountRows) Header() []string {
	return []string{"Margin Coin", "Account Equity", "USDT Equity", "Available", "Locked", "Crossed Max Avail", "Unrealized PL", "Max Transfer Out"}
}

func (a FuturesAccountRows) Row() [][]any {
	rows := [][]any{}
	for _, ac := range a {
		rows = append(rows, []any{
			ac.MarginCoin, ac.AccountEquity, ac.UsdtEquity, ac.Available, ac.Locked,
			ac.CrossedMaxAvailable, ac.UnrealizedPL, ac.MaxTransferOut,
		})
	}
	return rows
}

// FuturesHealthRows renders the classic futures account-health / risk picture:
// equity, maintenance margin, crossed risk rate and unrealised PnL per margin
// coin. A higher crossed risk rate means a higher risk of liquidation.
type FuturesHealthRows []mix.AccountListItem

func (h FuturesHealthRows) Header() []string {
	return []string{"Margin Coin", "Account Equity", "Crossed Risk Rate", "Union Maint. Margin", "Union Total Margin", "Union Available", "Unrealized PL", "Coupon"}
}

func (h FuturesHealthRows) Row() [][]any {
	rows := [][]any{}
	for _, ac := range h {
		rows = append(rows, []any{
			ac.MarginCoin, ac.AccountEquity, ac.CrossedRiskRate, ac.UnionMm,
			ac.UnionTotalMargin, ac.UnionAvailable, ac.UnrealizedPL, ac.Coupon,
		})
	}
	return rows
}

// FuturesPositionRows renders the account's open futures positions.
type FuturesPositionRows []mix.Position

func (p FuturesPositionRows) Header() []string {
	return []string{"Symbol", "Hold Side", "Margin Mode", "Leverage", "Total", "Available", "Open Avg", "Mark Price", "Liq Price", "Unrealized PL", "Margin Ratio", "Margin Coin"}
}

func (p FuturesPositionRows) Row() [][]any {
	rows := [][]any{}
	for _, pos := range p {
		rows = append(rows, []any{
			pos.Symbol, pos.HoldSide, pos.MarginMode, pos.Leverage, pos.Total, pos.Available,
			pos.OpenPriceAvg, pos.MarkPrice, pos.LiquidationPrice, pos.UnrealizedPL, pos.MarginRatio, pos.MarginCoin,
		})
	}
	return rows
}

// FuturesOrderRefView renders the identifiers returned by place/cancel.
type FuturesOrderRefView mix.OrderRef

func (o *FuturesOrderRefView) Header() []string {
	return []string{"Order ID", "Client Oid"}
}

func (o *FuturesOrderRefView) Row() [][]any {
	return [][]any{{o.OrderId, o.ClientOid}}
}

// FuturesOrderRows renders queried futures orders (detail / open list).
type FuturesOrderRows []mix.MixOrder

func (o FuturesOrderRows) Header() []string {
	return []string{"Order ID", "Symbol", "Side", "Type", "Status", "Price", "Size", "Avg Price", "Filled Base", "Leverage", "Margin Mode", "Reduce", "Created"}
}

func (o FuturesOrderRows) Row() [][]any {
	rows := [][]any{}
	for _, ord := range o {
		status := ord.Status
		if status == "" {
			status = ord.State
		}
		rows = append(rows, []any{
			ord.OrderId, ord.Symbol, ord.Side, ord.OrderType, status,
			ord.Price, ord.Size, ord.PriceAvg, ord.BaseVolume, ord.Leverage,
			ord.MarginMode, ord.ReduceOnly, common.FormatTime(ord.CTime),
		})
	}
	return rows
}
