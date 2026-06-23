package classic

import (
	"fmt"
	"strings"

	"github.com/UnipayFI/bitget-cli/common"
	"github.com/UnipayFI/go-bitget/classic/spot"
	"github.com/shopspring/decimal"
)

// ---- enum parsers --------------------------------------------------------

// ParseSpotSide resolves an order side string to the classic spot enum.
func ParseSpotSide(s string) (spot.Side, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "buy":
		return spot.SideBuy, nil
	case "sell":
		return spot.SideSell, nil
	default:
		return "", fmt.Errorf("invalid side %q: want buy or sell", s)
	}
}

// ParseSpotOrderType resolves an order-type string to the classic spot enum.
func ParseSpotOrderType(s string) (spot.OrderType, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "limit":
		return spot.OrderTypeLimit, nil
	case "market":
		return spot.OrderTypeMarket, nil
	default:
		return "", fmt.Errorf("invalid order type %q: want limit or market", s)
	}
}

// ParseSpotForce resolves a time-in-force string to the classic spot enum.
// An empty value defaults to gtc (the constructor requires a value, and it is
// ignored for market orders).
func ParseSpotForce(s string) (spot.Force, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "gtc":
		return spot.ForceGTC, nil
	case "post_only", "postonly":
		return spot.ForcePostOnly, nil
	case "fok":
		return spot.ForceFOK, nil
	case "ioc":
		return spot.ForceIOC, nil
	default:
		return "", fmt.Errorf("invalid time in force %q: want gtc, post_only, fok or ioc", s)
	}
}

// ---- service calls -------------------------------------------------------

func (c *SpotClient) GetAccountAssets(coin string) ([]spot.AccountAsset, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.sp.NewGetAccountAssetsService()
	if coin != "" {
		s.SetCoin(coin)
	}
	return s.Do(cx)
}

func (c *SpotClient) GetAccountInfo() (*spot.AccountInfo, error) {
	cx, cancel := ctx()
	defer cancel()
	return c.sp.NewGetAccountInfoService().Do(cx)
}

// SpotPlaceOrderParams collects the fields for a single classic spot order.
// Price is treated as unset when zero (required only for limit orders).
type SpotPlaceOrderParams struct {
	Symbol    string
	Side      spot.Side
	OrderType spot.OrderType
	Force     spot.Force
	Size      decimal.Decimal
	Price     decimal.Decimal
	ClientOid string
}

func (c *SpotClient) PlaceOrder(p SpotPlaceOrderParams) (*spot.PlaceOrderResponse, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.sp.NewPlaceOrderService(p.Symbol, p.Side, p.OrderType, p.Force, p.Size)
	if !p.Price.IsZero() {
		s.SetPrice(p.Price)
	}
	if p.ClientOid != "" {
		s.SetClientOid(p.ClientOid)
	}
	return s.Do(cx)
}

func (c *SpotClient) CancelOrder(symbol, orderID, clientOid string) (*spot.CancelOrderResponse, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.sp.NewCancelOrderService(symbol)
	if orderID != "" {
		s.SetOrderID(orderID)
	}
	if clientOid != "" {
		s.SetClientOid(clientOid)
	}
	return s.Do(cx)
}

func (c *SpotClient) GetOrderInfo(orderID, clientOid string) ([]spot.OrderInfo, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.sp.NewGetOrderInfoService()
	if orderID != "" {
		s.SetOrderID(orderID)
	}
	if clientOid != "" {
		s.SetClientOid(clientOid)
	}
	return s.Do(cx)
}

func (c *SpotClient) GetOpenOrders(symbol string, limit int) ([]spot.UnfilledOrder, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.sp.NewGetUnfilledOrdersService()
	if symbol != "" {
		s.SetSymbol(symbol)
	}
	if limit > 0 {
		s.SetLimit(limit)
	}
	return s.Do(cx)
}

// ---- table models --------------------------------------------------------

// SpotAssetRows renders the classic spot account's per-coin balances
// (non-zero only).
type SpotAssetRows []spot.AccountAsset

func (a SpotAssetRows) Header() []string {
	return []string{"Coin", "Available", "Frozen", "Locked", "Limit Available", "Updated"}
}

func (a SpotAssetRows) Row() [][]any {
	rows := [][]any{}
	for _, c := range a {
		if c.Available.IsZero() && c.Frozen.IsZero() && c.Locked.IsZero() {
			continue
		}
		rows = append(rows, []any{c.Coin, c.Available, c.Frozen, c.Locked, c.LimitAvailable, common.FormatTime(c.UTime)})
	}
	return rows
}

// SpotAccountInfoView renders the classic spot account identity and permissions.
type SpotAccountInfoView spot.AccountInfo

func (a *SpotAccountInfoView) Header() []string {
	return []string{"User ID", "Parent ID", "Trader Type", "Authorities", "IP List", "Register Time"}
}

func (a *SpotAccountInfoView) Row() [][]any {
	return [][]any{{
		a.UserId, a.ParentId, a.TraderType,
		strings.Join(a.Authorities, ","), strings.Join(a.Ips, ","), common.FormatTime(a.RegisTime),
	}}
}

// SpotOrderRefView renders the identifiers returned by place/cancel.
type SpotOrderRefView struct {
	OrderID   string
	ClientOid string
}

func (o *SpotOrderRefView) Header() []string {
	return []string{"Order ID", "Client Oid"}
}

func (o *SpotOrderRefView) Row() [][]any {
	return [][]any{{o.OrderID, o.ClientOid}}
}

// SpotOrderRows renders queried spot orders (order details).
type SpotOrderRows []spot.OrderInfo

func (o SpotOrderRows) Header() []string {
	return []string{"Order ID", "Symbol", "Side", "Type", "Status", "Price", "Size", "Avg Price", "Filled Base", "Filled Quote", "Created"}
}

func (o SpotOrderRows) Row() [][]any {
	rows := [][]any{}
	for _, ord := range o {
		rows = append(rows, []any{
			ord.OrderID, ord.Symbol, ord.Side, ord.OrderType, ord.Status,
			ord.Price, ord.Size, ord.PriceAvg, ord.BaseVolume, ord.QuoteVolume,
			common.FormatTime(ord.CTime),
		})
	}
	return rows
}

// SpotOpenOrderRows renders open (unfilled / partially filled) spot orders.
type SpotOpenOrderRows []spot.UnfilledOrder

func (o SpotOpenOrderRows) Header() []string {
	return []string{"Order ID", "Symbol", "Side", "Type", "Status", "Force", "Avg Price", "Size", "Filled Base", "Created"}
}

func (o SpotOpenOrderRows) Row() [][]any {
	rows := [][]any{}
	for _, ord := range o {
		rows = append(rows, []any{
			ord.OrderID, ord.Symbol, ord.Side, ord.OrderType, ord.Status, ord.Force,
			ord.PriceAvg, ord.Size, ord.BaseVolume, common.FormatTime(ord.CTime),
		})
	}
	return rows
}
