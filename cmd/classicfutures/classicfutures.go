// Package classicfutures implements the top-level `futures` command group for
// the classic account (/api/v2/mix/*): account info, account health,
// positions and basic order management. The product line is selected with the
// persistent --product flag.
package classicfutures

import (
	"errors"
	"strings"

	"github.com/UnipayFI/bitget-cli/exchange/classic"
	"github.com/UnipayFI/bitget-cli/printer"
	"github.com/UnipayFI/go-bitget/classic/mix"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

const docBase = "https://www.bitget.com/api-doc/contract/"

var (
	accountCmd = &cobra.Command{
		Use:   "account",
		Short: "Show futures account balances and equity",
		Long: `Show the classic futures account balances and equity per margin coin.

Docs Link: ` + docBase + "account/Get-Account-List",
		RunE: showAccount,
	}

	healthCmd = &cobra.Command{
		Use:   "health",
		Short: "Show futures account health (risk rate & maintenance margin)",
		Long: `Show the classic futures account-health / risk picture per margin coin:
equity, crossed risk rate, maintenance margin and unrealised PnL. A higher
crossed risk rate means a higher risk of liquidation.

Docs Link: ` + docBase + "account/Get-Account-List",
		RunE: showHealth,
	}

	positionCmd = &cobra.Command{
		Use:   "position",
		Short: "Query futures positions",
	}

	positionListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List open positions",
		Long: `List the account's open futures positions in the product line, optionally
filtered by --marginCoin and/or --symbol.

Docs Link: ` + docBase + "position/get-all-position",
		RunE: listPositions,
	}

	orderCmd = &cobra.Command{
		Use:   "order",
		Short: "Create, cancel and query futures orders",
	}

	createCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create a futures order",
		Long: `Place a new futures order.
* Required: --symbol, --side, --type, --size
* --price is required for limit orders
* --marginCoin defaults to USDT; set it for coin/usdc lines
* --tradeSide (open/close) and --reduceOnly apply in hedge mode

Docs Link: ` + docBase + "trade/Place-Order",
		RunE: createOrder,
	}

	cancelCmd = &cobra.Command{
		Use:   "cancel",
		Short: "Cancel a futures order",
		Long: `Cancel a single futures order by --orderId or --clientOid.

Docs Link: ` + docBase + "trade/Cancel-Order",
		RunE: cancelOrder,
	}

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Query a single futures order",
		Long: `Query a single futures order by --orderId or --clientOid.

Docs Link: ` + docBase + "trade/Get-Order-Details",
		RunE: getOrder,
	}

	openCmd = &cobra.Command{
		Use:   "open",
		Short: "List open futures orders",
		Long: `List currently open (unfilled / partially filled) futures orders.

Docs Link: ` + docBase + "trade/Get-Orders-Pending",
		RunE: openOrders,
	}
)

// InitCmds registers flags and returns the classic futures subcommands.
func InitCmds() []*cobra.Command {
	positionListCmd.Flags().StringP("marginCoin", "m", "", "margin coin filter, e.g. USDT")
	positionListCmd.Flags().StringP("symbol", "s", "", "symbol filter, e.g. BTCUSDT")

	createCmd.Flags().StringP("symbol", "s", "", "symbol, e.g. BTCUSDT (required)")
	createCmd.Flags().StringP("side", "S", "", "buy or sell (required)")
	createCmd.Flags().StringP("type", "t", "", "limit or market (required)")
	createCmd.Flags().StringP("size", "q", "", "order size (decimal) (required)")
	createCmd.Flags().StringP("price", "p", "", "order price (required for limit)")
	createCmd.Flags().StringP("marginCoin", "m", "USDT", "margin coin")
	createCmd.Flags().StringP("marginMode", "M", "crossed", "margin mode: crossed, isolated")
	createCmd.Flags().StringP("force", "f", "", "time in force: gtc, post_only, fok, ioc")
	createCmd.Flags().StringP("tradeSide", "T", "", "trade side (hedge mode): open, close")
	createCmd.Flags().StringP("reduceOnly", "r", "", "reduce-only: yes or no")
	createCmd.Flags().String("clientOid", "", "client order id")
	createCmd.MarkFlagRequired("symbol")
	createCmd.MarkFlagRequired("side")
	createCmd.MarkFlagRequired("type")
	createCmd.MarkFlagRequired("size")

	cancelCmd.Flags().StringP("symbol", "s", "", "symbol, e.g. BTCUSDT (required)")
	cancelCmd.Flags().StringP("marginCoin", "m", "", "margin coin")
	cancelCmd.Flags().StringP("orderId", "i", "", "order id")
	cancelCmd.Flags().StringP("clientOid", "c", "", "client order id")
	cancelCmd.MarkFlagRequired("symbol")

	getCmd.Flags().StringP("symbol", "s", "", "symbol, e.g. BTCUSDT (required)")
	getCmd.Flags().StringP("orderId", "i", "", "order id")
	getCmd.Flags().StringP("clientOid", "c", "", "client order id")
	getCmd.MarkFlagRequired("symbol")

	openCmd.Flags().StringP("symbol", "s", "", "symbol filter")
	openCmd.Flags().IntP("limit", "l", 0, "max records")

	positionCmd.AddCommand(positionListCmd)
	orderCmd.AddCommand(createCmd, cancelCmd, getCmd, openCmd)
	return []*cobra.Command{accountCmd, healthCmd, positionCmd, orderCmd}
}

// resolveProductType reads the inherited persistent --product flag.
func resolveProductType(cmd *cobra.Command) (mix.ProductType, error) {
	raw, _ := cmd.Flags().GetString("product")
	if raw == "" {
		raw = "usdt-futures"
	}
	return classic.ParseProductType(raw)
}

func showAccount(cmd *cobra.Command, _ []string) error {
	pt, err := resolveProductType(cmd)
	if err != nil {
		return err
	}
	accounts, err := classic.NewFuturesClient().GetAccountList(pt)
	if err != nil {
		return err
	}
	printer.Print(classic.FuturesAccountRows(accounts))
	return nil
}

func showHealth(cmd *cobra.Command, _ []string) error {
	pt, err := resolveProductType(cmd)
	if err != nil {
		return err
	}
	accounts, err := classic.NewFuturesClient().GetAccountList(pt)
	if err != nil {
		return err
	}
	printer.Print(classic.FuturesHealthRows(accounts))
	return nil
}

func listPositions(cmd *cobra.Command, _ []string) error {
	pt, err := resolveProductType(cmd)
	if err != nil {
		return err
	}
	marginCoin, _ := cmd.Flags().GetString("marginCoin")
	symbol, _ := cmd.Flags().GetString("symbol")
	positions, err := classic.NewFuturesClient().GetAllPositions(pt, marginCoin)
	if err != nil {
		return err
	}
	if symbol != "" {
		filtered := positions[:0]
		for _, p := range positions {
			if strings.EqualFold(p.Symbol, symbol) {
				filtered = append(filtered, p)
			}
		}
		positions = filtered
	}
	printer.Print(classic.FuturesPositionRows(positions))
	return nil
}

func createOrder(cmd *cobra.Command, _ []string) error {
	pt, err := resolveProductType(cmd)
	if err != nil {
		return err
	}
	symbol, _ := cmd.Flags().GetString("symbol")
	sideRaw, _ := cmd.Flags().GetString("side")
	typeRaw, _ := cmd.Flags().GetString("type")
	sizeRaw, _ := cmd.Flags().GetString("size")
	marginModeRaw, _ := cmd.Flags().GetString("marginMode")

	side, err := classic.ParseFuturesSide(sideRaw)
	if err != nil {
		return err
	}
	orderType, err := classic.ParseFuturesOrderType(typeRaw)
	if err != nil {
		return err
	}
	marginMode, err := classic.ParseMarginMode(marginModeRaw)
	if err != nil {
		return err
	}
	size, err := decimal.NewFromString(sizeRaw)
	if err != nil {
		return errors.New("invalid --size: " + err.Error())
	}
	priceRaw, _ := cmd.Flags().GetString("price")
	if orderType == "limit" && priceRaw == "" {
		return errors.New("--price is required for limit orders")
	}

	p := classic.FuturesPlaceOrderParams{
		Symbol:      symbol,
		ProductType: pt,
		MarginMode:  marginMode,
		Size:        size,
		Side:        side,
		OrderType:   orderType,
	}
	p.MarginCoin, _ = cmd.Flags().GetString("marginCoin")
	if priceRaw != "" {
		price, perr := decimal.NewFromString(priceRaw)
		if perr != nil {
			return errors.New("invalid --price: " + perr.Error())
		}
		p.Price = price
	}
	p.Force, _ = cmd.Flags().GetString("force")
	p.TradeSide, _ = cmd.Flags().GetString("tradeSide")
	p.ReduceOnly, _ = cmd.Flags().GetString("reduceOnly")
	p.ClientOid, _ = cmd.Flags().GetString("clientOid")

	ref, err := classic.NewFuturesClient().PlaceOrder(p)
	if err != nil {
		return err
	}
	view := classic.FuturesOrderRefView(*ref)
	printer.Print(&view)
	return nil
}

func cancelOrder(cmd *cobra.Command, _ []string) error {
	pt, err := resolveProductType(cmd)
	if err != nil {
		return err
	}
	symbol, _ := cmd.Flags().GetString("symbol")
	marginCoin, _ := cmd.Flags().GetString("marginCoin")
	orderID, _ := cmd.Flags().GetString("orderId")
	clientOid, _ := cmd.Flags().GetString("clientOid")
	if orderID == "" && clientOid == "" {
		return errors.New("one of --orderId or --clientOid is required")
	}
	ref, err := classic.NewFuturesClient().CancelOrder(symbol, pt, marginCoin, orderID, clientOid)
	if err != nil {
		return err
	}
	view := classic.FuturesOrderRefView(*ref)
	printer.Print(&view)
	return nil
}

func getOrder(cmd *cobra.Command, _ []string) error {
	pt, err := resolveProductType(cmd)
	if err != nil {
		return err
	}
	symbol, _ := cmd.Flags().GetString("symbol")
	orderID, _ := cmd.Flags().GetString("orderId")
	clientOid, _ := cmd.Flags().GetString("clientOid")
	if orderID == "" && clientOid == "" {
		return errors.New("one of --orderId or --clientOid is required")
	}
	order, err := classic.NewFuturesClient().GetOrderDetail(symbol, pt, orderID, clientOid)
	if err != nil {
		return err
	}
	printer.Print(classic.FuturesOrderRows{*order})
	return nil
}

func openOrders(cmd *cobra.Command, _ []string) error {
	pt, err := resolveProductType(cmd)
	if err != nil {
		return err
	}
	symbol, _ := cmd.Flags().GetString("symbol")
	limit, _ := cmd.Flags().GetInt("limit")
	list, err := classic.NewFuturesClient().GetOpenOrders(pt, symbol, limit)
	if err != nil {
		return err
	}
	printer.Print(classic.FuturesOrderRows(list.EntrustedList))
	return nil
}
