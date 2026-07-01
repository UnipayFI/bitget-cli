package exchange

import (
	"github.com/UnipayFI/bitget-cli/common"
	bguta "github.com/UnipayFI/go-bitget/uta"
)

// ---- service calls -------------------------------------------------------

func (c *Client) GetPositions(category bguta.Category, symbol, posSide string) ([]bguta.Position, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.uta.NewGetPositionService(category)
	if symbol != "" {
		s.SetSymbol(symbol)
	}
	if posSide != "" {
		s.SetPosSide(bguta.PosSide(posSide))
	}
	return s.Do(cx)
}

func (c *Client) GetPositionHistory(p HistoryParams) (*bguta.PositionHistory, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.uta.NewGetPositionHistoryService(p.Category)
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

func (c *Client) GetPositionADLRank() ([]bguta.PositionADLRank, error) {
	cx, cancel := ctx()
	defer cancel()
	return c.uta.NewGetPositionADLRankService().Do(cx)
}

func (c *Client) ClosePositions(category bguta.Category, symbol, posSide string) (*bguta.ClosePositionsResult, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.uta.NewClosePositionsService(category)
	if symbol != "" {
		s.SetSymbol(symbol)
	}
	if posSide != "" {
		s.SetPosSide(bguta.PosSide(posSide))
	}
	return s.Do(cx)
}

// ---- table models --------------------------------------------------------

// Positions renders the account's open futures positions.
type Positions []bguta.Position

func (p Positions) Header() []string {
	return []string{"Symbol", "Pos Side", "Margin Mode", "Leverage", "Total", "Available", "Avg Price", "Mark Price", "Liq Price", "Unrealised PNL", "Profit Rate", "Margin Coin"}
}

func (p Positions) Row() [][]any {
	rows := [][]any{}
	for _, pos := range p {
		rows = append(rows, []any{
			pos.Symbol, pos.PosSide, pos.MarginMode, pos.Leverage, pos.Total, pos.Available,
			pos.AvgPrice, pos.MarkPrice, pos.LiquidationPrice, pos.UnrealizedPnL, pos.ProfitRate, pos.MarginCoin,
		})
	}
	return rows
}

// PositionHistoryRows renders closed/historical futures positions.
type PositionHistoryRows []bguta.HistoryPosition

func (p PositionHistoryRows) Header() []string {
	return []string{"Symbol", "Pos Side", "Open Avg", "Close Avg", "Open Qty", "Close Qty", "Realised PNL", "Net Profit", "Funding", "Created", "Updated"}
}

func (p PositionHistoryRows) Row() [][]any {
	rows := [][]any{}
	for _, pos := range p {
		rows = append(rows, []any{
			pos.Symbol, pos.PosSide, pos.OpenPriceAvg, pos.ClosePriceAvg, pos.OpenTotalPos, pos.CloseTotalPos,
			pos.CumRealisedPnL, pos.NetProfit, pos.TotalFunding, common.FormatTime(pos.CreatedTime), common.FormatTime(pos.UpdatedTime),
		})
	}
	return rows
}

// ADLRankRows renders the auto-deleveraging queue ranking per position.
type ADLRankRows []bguta.PositionADLRank

func (a ADLRankRows) Header() []string {
	return []string{"Symbol", "Margin Coin", "Hold Side", "ADL Rank"}
}

func (a ADLRankRows) Row() [][]any {
	rows := [][]any{}
	for _, r := range a {
		rows = append(rows, []any{r.Symbol, r.MarginCoin, r.HoldSide, r.ADLRank})
	}
	return rows
}

// OrderResults renders the per-order outcomes returned by batch/close endpoints.
type OrderResults []bguta.OrderResult

func (o OrderResults) Header() []string {
	return []string{"Order ID", "Client Oid", "Code", "Msg"}
}

func (o OrderResults) Row() [][]any {
	rows := [][]any{}
	for _, r := range o {
		rows = append(rows, []any{r.OrderID, r.ClientOrderID, r.Code, r.Msg})
	}
	return rows
}
