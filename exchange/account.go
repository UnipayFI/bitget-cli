package exchange

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/UnipayFI/bitget-cli/common"
	bguta "github.com/UnipayFI/go-bitget/uta"
	"github.com/shopspring/decimal"
)

// ---- service calls -------------------------------------------------------

func (c *Client) GetAccountAssets() (*bguta.AccountAssets, error) {
	cx, cancel := ctx()
	defer cancel()
	return c.uta.NewGetAccountAssetsService().Do(cx)
}

func (c *Client) GetAccountInfo() (*bguta.AccountInfo, error) {
	cx, cancel := ctx()
	defer cancel()
	return c.uta.NewGetAccountInfoService().Do(cx)
}

func (c *Client) GetAccountSettings() (*bguta.AccountSettings, error) {
	cx, cancel := ctx()
	defer cancel()
	return c.uta.NewGetAccountSettingsService().Do(cx)
}

func (c *Client) GetFeeRate(category bguta.Category, symbol string) (*bguta.AccountFeeRate, error) {
	cx, cancel := ctx()
	defer cancel()
	return c.uta.NewGetAccountFeeRateService(category, symbol).Do(cx)
}

func (c *Client) GetFundingAssets(coin string) ([]bguta.FundingAsset, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.uta.NewGetAccountFundingAssetsService()
	if coin != "" {
		s.SetCoin(coin)
	}
	return s.Do(cx)
}

func (c *Client) GetMaxTransferable(coin string) (*bguta.MaxTransferable, error) {
	cx, cancel := ctx()
	defer cancel()
	return c.uta.NewGetMaxTransferableService(coin).Do(cx)
}

func (c *Client) GetMaxWithdrawal(coin string) (*bguta.MaxWithdrawal, error) {
	cx, cancel := ctx()
	defer cancel()
	return c.uta.NewGetMaxWithdrawalService(coin).Do(cx)
}

// FinancialRecordsParams collects the optional filters for the ledger query.
type FinancialRecordsParams struct {
	Category   bguta.Category
	Coin       string
	RecordType string
	StartTime  time.Time
	EndTime    time.Time
	Limit      string
}

func (c *Client) GetFinancialRecords(p FinancialRecordsParams) (*bguta.FinancialRecords, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.uta.NewGetFinancialRecordsService(p.Category)
	if p.Coin != "" {
		s.SetCoin(p.Coin)
	}
	if p.RecordType != "" {
		s.SetType(p.RecordType)
	}
	if !p.StartTime.IsZero() {
		s.SetStartTime(p.StartTime)
	}
	if !p.EndTime.IsZero() {
		s.SetEndTime(p.EndTime)
	}
	if p.Limit != "" {
		n, err := strconv.Atoi(p.Limit)
		if err != nil {
			return nil, fmt.Errorf("invalid limit %q: %w", p.Limit, err)
		}
		s.SetLimit(n)
	}
	return s.Do(cx)
}

// SetLeverageParams collects the optional fields for a leverage change.
type SetLeverageParams struct {
	Category      bguta.Category
	Leverage      string
	Symbol        string
	Coin          string
	MarginMode    string
	PosSide       string
	LongLeverage  string
	ShortLeverage string
}

func (c *Client) SetLeverage(p SetLeverageParams) (string, error) {
	cx, cancel := ctx()
	defer cancel()
	lev, err := strconv.Atoi(p.Leverage)
	if err != nil {
		return "", fmt.Errorf("invalid leverage %q: %w", p.Leverage, err)
	}
	s := c.uta.NewSetLeverageService(p.Category, lev)
	if p.Symbol != "" {
		s.SetSymbol(p.Symbol)
	}
	if p.Coin != "" {
		s.SetCoin(p.Coin)
	}
	if p.MarginMode != "" {
		s.SetMarginMode(bguta.MarginMode(p.MarginMode))
	}
	if p.PosSide != "" {
		s.SetPosSide(bguta.PosSide(p.PosSide))
	}
	if p.LongLeverage != "" {
		s.SetLongLeverage(p.LongLeverage)
	}
	if p.ShortLeverage != "" {
		s.SetShortLeverage(p.ShortLeverage)
	}
	res, err := s.Do(cx)
	if err != nil {
		return "", err
	}
	return *res, nil
}

func (c *Client) SetHoldMode(holdMode string) (string, error) {
	cx, cancel := ctx()
	defer cancel()
	res, err := c.uta.NewSetHoldModeService(bguta.HoldMode(holdMode)).Do(cx)
	if err != nil {
		return "", err
	}
	return *res, nil
}

func (c *Client) SetMargin(category bguta.Category, symbol, posSide, operation string, amount decimal.Decimal) (string, error) {
	cx, cancel := ctx()
	defer cancel()
	res, err := c.uta.NewSetMarginService(category, symbol, bguta.PosSide(posSide), operation, amount).Do(cx)
	if err != nil {
		return "", err
	}
	return *res, nil
}

// ---- table models --------------------------------------------------------

// CoinAssets renders the unified account's per-coin balances (non-zero only).
type CoinAssets []bguta.CoinAsset

func (a CoinAssets) Header() []string {
	return []string{"Coin", "Equity", "Balance", "Available", "Locked", "Debt", "USD Value"}
}

func (a CoinAssets) Row() [][]any {
	rows := [][]any{}
	for _, c := range a {
		if c.Equity.IsZero() && c.Balance.IsZero() && c.Available.IsZero() && c.Debt.IsZero() {
			continue
		}
		rows = append(rows, []any{c.Coin, c.Equity, c.Balance, c.Available, c.Locked, c.Debt, c.USDValue})
	}
	return rows
}

// AccountSummary renders the unified account's aggregate equity and margin.
type AccountSummary bguta.AccountAssets

func (a *AccountSummary) Header() []string {
	return []string{"Account Equity", "USDT Equity", "BTC Equity", "Unrealised PNL", "Eff Equity", "IMR", "MMR", "Mgn Ratio"}
}

func (a *AccountSummary) Row() [][]any {
	return [][]any{{a.AccountEquity, a.USDTEquity, a.BtcEquity, a.UnrealizedPnL, a.EffEquity, a.Imr, a.Mmr, a.MgnRatio}}
}

// AccountHealthView renders the unified account's risk/health metrics: account
// and effective equity, unrealised PnL, initial/maintenance margin requirements
// (IMR/MMR) and the margin ratio. A margin ratio approaching 1 (100%) signals
// liquidation risk.
type AccountHealthView bguta.AccountAssets

func (a *AccountHealthView) Header() []string {
	return []string{"Account Equity", "Eff Equity", "Unrealised PNL", "IMR", "MMR", "Mgn Ratio"}
}

func (a *AccountHealthView) Row() [][]any {
	return [][]any{{a.AccountEquity, a.EffEquity, a.UnrealizedPnL, a.Imr, a.Mmr, a.MgnRatio}}
}

// AccountInfoView renders account identity and permission metadata.
type AccountInfoView bguta.AccountInfo

func (a *AccountInfoView) Header() []string {
	return []string{"User ID", "Parent ID", "Perm Type", "Permissions", "IP List", "Register Time"}
}

func (a *AccountInfoView) Row() [][]any {
	perms := strings.Join(a.Permissions, ",")
	return [][]any{{a.UserID, a.ParentID, a.PermType, perms, a.Ips, common.FormatTime(a.RegisTime)}}
}

// AccountSettingsView renders the account-mode header of the settings payload.
type AccountSettingsView bguta.AccountSettings

func (a *AccountSettingsView) Header() []string {
	return []string{"UID", "Account Mode", "Asset Mode", "Account Level", "Hold Mode", "STP Mode"}
}

func (a *AccountSettingsView) Row() [][]any {
	return [][]any{{a.UID, a.AccountMode, a.AssetMode, a.AccountLevel, a.HoldMode, a.StpMode}}
}

// LeverageConfigs renders the per-symbol leverage/margin configuration list.
type LeverageConfigs []bguta.SymbolLeverageConfig

func (l LeverageConfigs) Header() []string {
	return []string{"Category", "Symbol", "Margin Mode", "Leverage"}
}

func (l LeverageConfigs) Row() [][]any {
	rows := [][]any{}
	for _, c := range l {
		rows = append(rows, []any{c.Category, c.Symbol, c.MarginMode, c.Leverage})
	}
	return rows
}

// FeeRateView renders the maker/taker fee rate for a symbol.
type FeeRateView struct {
	Symbol string
	bguta.AccountFeeRate
}

func (f *FeeRateView) Header() []string {
	return []string{"Symbol", "Maker Fee Rate", "Taker Fee Rate"}
}

func (f *FeeRateView) Row() [][]any {
	return [][]any{{f.Symbol, f.MakerFeeRate, f.TakerFeeRate}}
}

// FundingAssets renders the funding (P2P) account balances.
type FundingAssets []bguta.FundingAsset

func (f FundingAssets) Header() []string {
	return []string{"Coin", "Balance", "Available", "Frozen"}
}

func (f FundingAssets) Row() [][]any {
	rows := [][]any{}
	for _, a := range f {
		rows = append(rows, []any{a.Coin, a.Balance, a.Available, a.Frozen})
	}
	return rows
}

// FinancialRecordList renders the account ledger records.
type FinancialRecordList []bguta.FinancialRecord

func (r FinancialRecordList) Header() []string {
	return []string{"Time", "Category", "Coin", "Type", "Amount", "Fee", "Balance", "Symbol"}
}

func (r FinancialRecordList) Row() [][]any {
	rows := [][]any{}
	for _, rec := range r {
		rows = append(rows, []any{common.FormatTime(rec.Ts), rec.Category, rec.Coin, rec.Type, rec.Amount, rec.Fee, rec.Balance, rec.Symbol})
	}
	return rows
}

// MaxTransferableView renders the max transferable amounts for a coin.
type MaxTransferableView bguta.MaxTransferable

func (m *MaxTransferableView) Header() []string {
	return []string{"Coin", "Max Transfer", "Borrow Max Transfer"}
}

func (m *MaxTransferableView) Row() [][]any {
	return [][]any{{m.Coin, m.MaxTransfer, m.BorrowMaxTransfer}}
}

// MaxWithdrawalView renders the max withdrawable amounts for a coin.
type MaxWithdrawalView bguta.MaxWithdrawal

func (m *MaxWithdrawalView) Header() []string {
	return []string{"Coin", "UTA Max", "Spot Max", "OTC Max", "Total Max"}
}

func (m *MaxWithdrawalView) Row() [][]any {
	return [][]any{{m.Coin, m.UtaMaxWithdrawal, m.SpotMaxWithdrawal, m.OTCMaxWithdrawal, m.TotalMaxWithdrawal}}
}
