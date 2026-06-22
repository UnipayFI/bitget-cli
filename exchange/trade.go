package exchange

import (
	"fmt"
	"strings"
	"time"

	"github.com/UnipayFI/bitget-cli/common"
	bguta "github.com/UnipayFI/go-bitget/uta"
	"github.com/shopspring/decimal"
)

// ---- service calls -------------------------------------------------------

// PlaceOrderParams collects the fields for a single order. Price is treated as
// unset when zero (required only for limit orders); the futures-only fields are
// applied only when non-empty.
type PlaceOrderParams struct {
	Category    bguta.Category
	Symbol      string
	Side        bguta.Side
	OrderType   bguta.OrderType
	Qty         decimal.Decimal
	Price       decimal.Decimal
	TimeInForce string
	PosSide     string
	ReduceOnly  string
	MarginMode  string
	ClientOid   string
}

func (c *Client) PlaceOrder(p PlaceOrderParams) (*bguta.OrderRef, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.uta.NewPlaceOrderService(p.Category, p.Symbol, p.Qty, p.Side, p.OrderType)
	if !p.Price.IsZero() {
		s.SetPrice(p.Price)
	}
	if p.TimeInForce != "" {
		s.SetTimeInForce(bguta.TimeInForce(p.TimeInForce))
	}
	if p.PosSide != "" {
		s.SetPosSide(bguta.PosSide(p.PosSide))
	}
	if p.ReduceOnly != "" {
		s.SetReduceOnly(p.ReduceOnly)
	}
	if p.MarginMode != "" {
		s.SetMarginMode(bguta.MarginMode(p.MarginMode))
	}
	if p.ClientOid != "" {
		s.SetClientOid(p.ClientOid)
	}
	return s.Do(cx)
}

// ModifyOrderParams collects the fields for an order amendment. Identify the
// order by OrderId or ClientOid; supply at least one of Qty/Price.
type ModifyOrderParams struct {
	Category  bguta.Category
	Symbol    string
	OrderId   string
	ClientOid string
	Qty       decimal.Decimal
	Price     decimal.Decimal
}

func (c *Client) ModifyOrder(p ModifyOrderParams) (*bguta.OrderRef, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.uta.NewModifyOrderService().SetCategory(p.Category)
	if p.Symbol != "" {
		s.SetSymbol(p.Symbol)
	}
	if p.OrderId != "" {
		s.SetOrderId(p.OrderId)
	}
	if p.ClientOid != "" {
		s.SetClientOid(p.ClientOid)
	}
	if !p.Qty.IsZero() {
		s.SetQty(p.Qty)
	}
	if !p.Price.IsZero() {
		s.SetPrice(p.Price)
	}
	return s.Do(cx)
}

func (c *Client) CancelOrder(category bguta.Category, orderId, clientOid string) (*bguta.OrderRef, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.uta.NewCancelOrderService().SetCategory(category)
	if orderId != "" {
		s.SetOrderId(orderId)
	}
	if clientOid != "" {
		s.SetClientOid(clientOid)
	}
	return s.Do(cx)
}

func (c *Client) CancelSymbolOrders(category bguta.Category, symbol string) (*bguta.CancelSymbolOrderResult, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.uta.NewCancelSymbolOrderService(category)
	if symbol != "" {
		s.SetSymbol(symbol)
	}
	return s.Do(cx)
}

func (c *Client) GetOrderInfo(orderId, clientOid string) (*bguta.Order, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.uta.NewGetOrderInfoService()
	if orderId != "" {
		s.SetOrderID(orderId)
	}
	if clientOid != "" {
		s.SetClientOid(clientOid)
	}
	return s.Do(cx)
}

func (c *Client) GetOpenOrders(category bguta.Category, symbol, limit string) (*bguta.OrderList, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.uta.NewGetOpenOrdersService().SetCategory(category)
	if symbol != "" {
		s.SetSymbol(symbol)
	}
	if limit != "" {
		s.SetLimit(limit)
	}
	return s.Do(cx)
}

// HistoryParams collects the common filters for the order/fill history queries.
type HistoryParams struct {
	Category  bguta.Category
	Symbol    string
	OrderId   string
	StartTime time.Time
	EndTime   time.Time
	Limit     string
}

func (c *Client) GetOrderHistory(p HistoryParams) (*bguta.OrderList, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.uta.NewGetOrderHistoryService(p.Category)
	if p.Symbol != "" {
		s.SetSymbol(p.Symbol)
	}
	if !p.StartTime.IsZero() {
		s.SetStartTime(p.StartTime)
	}
	if !p.EndTime.IsZero() {
		s.SetEndTime(p.EndTime)
	}
	if p.Limit != "" {
		s.SetLimit(p.Limit)
	}
	return s.Do(cx)
}

func (c *Client) GetFills(p HistoryParams) (*bguta.FillList, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.uta.NewGetFillHistoryService().SetCategory(p.Category)
	if p.OrderId != "" {
		s.SetOrderID(p.OrderId)
	}
	if !p.StartTime.IsZero() {
		s.SetStartTime(p.StartTime)
	}
	if !p.EndTime.IsZero() {
		s.SetEndTime(p.EndTime)
	}
	if p.Limit != "" {
		s.SetLimit(p.Limit)
	}
	return s.Do(cx)
}

// ---- table models --------------------------------------------------------

// formatFees joins an order/fill fee breakdown into a single cell.
func formatFees(details []bguta.FeeDetail) string {
	parts := make([]string, 0, len(details))
	for _, d := range details {
		if d.Fee.IsZero() {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s %s", d.Fee, d.FeeCoin))
	}
	return strings.Join(parts, ", ")
}

// OrderRows renders a collection of orders (info / open / history).
type OrderRows []bguta.Order

func (o OrderRows) Header() []string {
	return []string{"Order ID", "Symbol", "Side", "Type", "Status", "Price", "Qty", "Filled", "Avg Price", "TIF", "Pos Side", "Reduce", "Created"}
}

func (o OrderRows) Row() [][]any {
	rows := [][]any{}
	for _, ord := range o {
		rows = append(rows, []any{
			ord.OrderID, ord.Symbol, ord.Side, ord.OrderType, ord.OrderStatus,
			ord.Price, ord.Qty, ord.CumExecQty, ord.AvgPrice, ord.TimeInForce,
			ord.PosSide, ord.ReduceOnly, common.FormatTime(ord.CreatedTime),
		})
	}
	return rows
}

// OrderRefView renders the identifiers returned by place/modify/cancel.
type OrderRefView bguta.OrderRef

func (o *OrderRefView) Header() []string {
	return []string{"Order ID", "Client Oid"}
}

func (o *OrderRefView) Row() [][]any {
	return [][]any{{o.OrderId, o.ClientOid}}
}

// CancelResults renders the per-order outcome of a bulk cancellation.
type CancelResults []bguta.CancelResult

func (c CancelResults) Header() []string {
	return []string{"Order ID", "Client Oid", "Code", "Msg"}
}

func (c CancelResults) Row() [][]any {
	rows := [][]any{}
	for _, r := range c {
		rows = append(rows, []any{r.OrderId, r.ClientOid, r.Code, r.Msg})
	}
	return rows
}

// FillRows renders a collection of trade fills.
type FillRows []bguta.Fill

func (f FillRows) Header() []string {
	return []string{"Exec ID", "Order ID", "Symbol", "Side", "Exec Price", "Exec Qty", "Exec Value", "Scope", "Fee", "PNL", "Time"}
}

func (f FillRows) Row() [][]any {
	rows := [][]any{}
	for _, fill := range f {
		rows = append(rows, []any{
			fill.ExecID, fill.OrderID, fill.Symbol, fill.Side, fill.ExecPrice,
			fill.ExecQty, fill.ExecValue, fill.TradeScope, formatFees(fill.FeeDetail),
			fill.ExecPnl, common.FormatTime(fill.CreatedTime),
		})
	}
	return rows
}
