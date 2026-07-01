package exchange

import (
	"time"

	"github.com/UnipayFI/bitget-cli/common"
	bguta "github.com/UnipayFI/go-bitget/uta"
	"github.com/shopspring/decimal"
)

// ---- service calls -------------------------------------------------------

// TransferParams collects the fields for an intra-account transfer.
type TransferParams struct {
	FromType    string
	ToType      string
	Coin        string
	Amount      decimal.Decimal
	Symbol      string
	AllowBorrow string
	ClientOid   string
}

func (c *Client) Transfer(p TransferParams) (*bguta.TransferResult, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.uta.NewTransferService(p.FromType, p.ToType, p.Coin, p.Amount)
	if p.Symbol != "" {
		s.SetSymbol(p.Symbol)
	}
	if p.AllowBorrow != "" {
		s.SetAllowBorrow(p.AllowBorrow)
	}
	if p.ClientOid != "" {
		s.SetClientOrderID(p.ClientOid)
	}
	return s.Do(cx)
}

func (c *Client) TransferableCoins(fromType, toType string) ([]string, error) {
	cx, cancel := ctx()
	defer cancel()
	return c.uta.NewGetTransferableCoinsService(fromType, toType).Do(cx)
}

func (c *Client) GetDepositAddress(coin, chain string) (*bguta.DepositAddress, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.uta.NewGetDepositAddressService(coin)
	if chain != "" {
		s.SetChain(chain)
	}
	return s.Do(cx)
}

func (c *Client) GetDepositRecords(coin string, start, end time.Time, limit string) ([]bguta.DepositRecord, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.uta.NewGetDepositRecordsService(start, end)
	if coin != "" {
		s.SetCoin(coin)
	}
	if limit != "" {
		s.SetLimit(limit)
	}
	return s.Do(cx)
}

// WithdrawParams collects the fields for a withdrawal. Chain is required for
// on-chain withdrawals; the identity/area fields apply to specific corridors.
type WithdrawParams struct {
	Coin         string
	TransferType string
	Address      string
	Size         decimal.Decimal
	Chain        string
	Tag          string
	InnerToType  string
	Remark       string
	ClientOid    string
}

func (c *Client) Withdraw(p WithdrawParams) (*bguta.WithdrawResult, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.uta.NewWithdrawService(p.Coin, p.TransferType, p.Address, p.Size)
	if p.Chain != "" {
		s.SetChain(p.Chain)
	}
	if p.Tag != "" {
		s.SetTag(p.Tag)
	}
	if p.InnerToType != "" {
		s.SetInnerToType(p.InnerToType)
	}
	if p.Remark != "" {
		s.SetRemark(p.Remark)
	}
	if p.ClientOid != "" {
		s.SetClientOrderID(p.ClientOid)
	}
	return s.Do(cx)
}

func (c *Client) GetWithdrawalRecords(coin string, start, end time.Time, limit string) ([]bguta.WithdrawalRecord, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.uta.NewGetWithdrawalRecordsService(start, end)
	if coin != "" {
		s.SetCoin(coin)
	}
	if limit != "" {
		s.SetLimit(limit)
	}
	return s.Do(cx)
}

func (c *Client) GetWithdrawAddress(coin, addressType string) (*bguta.WithdrawAddressBook, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.uta.NewGetWithdrawAddressService()
	if coin != "" {
		s.SetCoin(coin)
	}
	if addressType != "" {
		s.SetType(addressType)
	}
	return s.Do(cx)
}

func (c *Client) CancelWithdrawal(orderId, clientOid string) (string, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.uta.NewCancelWithdrawalService()
	if orderId != "" {
		s.SetOrderID(orderId)
	}
	if clientOid != "" {
		s.SetClientOrderID(clientOid)
	}
	res, err := s.Do(cx)
	if err != nil {
		return "", err
	}
	return *res, nil
}

// ---- table models --------------------------------------------------------

// TransferResultView renders the identifiers returned by a transfer.
type TransferResultView bguta.TransferResult

func (t *TransferResultView) Header() []string {
	return []string{"Transfer ID", "Client Oid"}
}

func (t *TransferResultView) Row() [][]any {
	return [][]any{{t.TransferID, t.ClientOrderID}}
}

// CoinList renders a bare list of coin names (e.g. transferable coins).
type CoinList []string

func (c CoinList) Header() []string {
	return []string{"Coin"}
}

func (c CoinList) Row() [][]any {
	rows := [][]any{}
	for _, coin := range c {
		rows = append(rows, []any{coin})
	}
	return rows
}

// DepositAddressView renders a single deposit address.
type DepositAddressView bguta.DepositAddress

func (d *DepositAddressView) Header() []string {
	return []string{"Coin", "Chain", "Address", "Tag"}
}

func (d *DepositAddressView) Row() [][]any {
	return [][]any{{d.Coin, d.Chain, d.Address, d.Tag}}
}

// DepositRecords renders the account's deposit history.
type DepositRecords []bguta.DepositRecord

func (d DepositRecords) Header() []string {
	return []string{"Order ID", "Coin", "Size", "Status", "Chain", "From", "To", "Created"}
}

func (d DepositRecords) Row() [][]any {
	rows := [][]any{}
	for _, r := range d {
		rows = append(rows, []any{r.OrderID, r.Coin, r.Size, r.Status, r.Chain, r.FromAddress, r.ToAddress, common.FormatTime(r.CreatedTime)})
	}
	return rows
}

// WithdrawResultView renders the identifiers returned by a withdrawal.
type WithdrawResultView bguta.WithdrawResult

func (w *WithdrawResultView) Header() []string {
	return []string{"Order ID", "Client Oid"}
}

func (w *WithdrawResultView) Row() [][]any {
	return [][]any{{w.OrderID, w.ClientOrderID}}
}

// WithdrawalRecords renders the account's withdrawal history.
type WithdrawalRecords []bguta.WithdrawalRecord

func (w WithdrawalRecords) Header() []string {
	return []string{"Order ID", "Coin", "Size", "Fee", "Status", "Chain", "To", "Created"}
}

func (w WithdrawalRecords) Row() [][]any {
	rows := [][]any{}
	for _, r := range w {
		rows = append(rows, []any{r.OrderID, r.Coin, r.Size, r.Fee, r.Status, r.Chain, r.ToAddress, common.FormatTime(r.CreatedTime)})
	}
	return rows
}

// WithdrawAddresses renders the saved withdrawal address book.
type WithdrawAddresses []bguta.WithdrawAddress

func (w WithdrawAddresses) Header() []string {
	return []string{"Coin", "Chain", "Address", "Tag", "Label", "Type"}
}

func (w WithdrawAddresses) Row() [][]any {
	rows := [][]any{}
	for _, a := range w {
		rows = append(rows, []any{a.Coin, a.Chain, a.Address, a.Memo, a.Label, a.Type})
	}
	return rows
}
