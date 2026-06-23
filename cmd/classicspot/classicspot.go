// Package classicspot implements the top-level `spot` command group for the
// classic account (/api/v2/spot/*): account info, balances and basic order
// management.
package classicspot

import (
	"errors"

	"github.com/UnipayFI/bitget-cli/exchange/classic"
	"github.com/UnipayFI/bitget-cli/printer"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

const docBase = "https://www.bitget.com/api-doc/spot/"

var (
	accountCmd = &cobra.Command{
		Use:   "account",
		Short: "Spot account info and balances",
	}

	accountInfoCmd = &cobra.Command{
		Use:   "info",
		Short: "Show spot account identity and permissions",
		Long: `Show the classic spot account's identity and API permission metadata.

Docs Link: ` + docBase + "account/Get-Account-Info",
		RunE: showAccountInfo,
	}

	assetsCmd = &cobra.Command{
		Use:   "assets",
		Short: "Show spot per-coin balances (non-zero)",
		Long: `Show the classic spot account's per-coin balances. Only coins with a
non-zero available/frozen/locked amount are shown.

Docs Link: ` + docBase + "account/Get-Account-Assets",
		RunE: showAssets,
	}

	orderCmd = &cobra.Command{
		Use:   "order",
		Short: "Create, cancel and query spot orders",
	}

	createCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create a spot order",
		Long: `Place a new spot order.
* Required: --symbol, --side, --type, --size
* --price is required for limit orders
* For market buy orders --size is in quote currency (e.g. USDT)

Docs Link: ` + docBase + "trade/Place-Order",
		RunE: createOrder,
	}

	cancelCmd = &cobra.Command{
		Use:   "cancel",
		Short: "Cancel a spot order",
		Long: `Cancel a single spot order by --orderId or --clientOid.

Docs Link: ` + docBase + "trade/Cancel-Order",
		RunE: cancelOrder,
	}

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Query a single spot order",
		Long: `Query a single spot order by --orderId or --clientOid.

Docs Link: ` + docBase + "trade/Get-Order-Info",
		RunE: getOrder,
	}

	openCmd = &cobra.Command{
		Use:   "open",
		Short: "List open spot orders",
		Long: `List currently open (unfilled / partially filled) spot orders.

Docs Link: ` + docBase + "trade/Get-Unfilled-Orders",
		RunE: openOrders,
	}
)

// InitCmds registers flags and returns the classic spot subcommands.
func InitCmds() []*cobra.Command {
	assetsCmd.Flags().StringP("coin", "c", "", "coin filter, e.g. USDT")

	createCmd.Flags().StringP("symbol", "s", "", "symbol, e.g. BTCUSDT (required)")
	createCmd.Flags().StringP("side", "S", "", "buy or sell (required)")
	createCmd.Flags().StringP("type", "t", "", "limit or market (required)")
	createCmd.Flags().StringP("size", "q", "", "order size (decimal) (required)")
	createCmd.Flags().StringP("price", "p", "", "order price (required for limit)")
	createCmd.Flags().StringP("force", "f", "", "time in force: gtc, post_only, fok, ioc (default gtc)")
	createCmd.Flags().String("clientOid", "", "client order id")
	createCmd.MarkFlagRequired("symbol")
	createCmd.MarkFlagRequired("side")
	createCmd.MarkFlagRequired("type")
	createCmd.MarkFlagRequired("size")

	cancelCmd.Flags().StringP("symbol", "s", "", "symbol, e.g. BTCUSDT (required)")
	cancelCmd.Flags().StringP("orderId", "i", "", "order id")
	cancelCmd.Flags().StringP("clientOid", "c", "", "client order id")
	cancelCmd.MarkFlagRequired("symbol")

	getCmd.Flags().StringP("orderId", "i", "", "order id")
	getCmd.Flags().StringP("clientOid", "c", "", "client order id")

	openCmd.Flags().StringP("symbol", "s", "", "symbol filter")
	openCmd.Flags().IntP("limit", "l", 0, "max records")

	accountCmd.AddCommand(accountInfoCmd, assetsCmd)
	orderCmd.AddCommand(createCmd, cancelCmd, getCmd, openCmd)
	return []*cobra.Command{accountCmd, orderCmd}
}

func showAccountInfo(cmd *cobra.Command, _ []string) error {
	info, err := classic.NewSpotClient().GetAccountInfo()
	if err != nil {
		return err
	}
	view := classic.SpotAccountInfoView(*info)
	printer.Print(&view)
	return nil
}

func showAssets(cmd *cobra.Command, _ []string) error {
	coin, _ := cmd.Flags().GetString("coin")
	assets, err := classic.NewSpotClient().GetAccountAssets(coin)
	if err != nil {
		return err
	}
	printer.Print(classic.SpotAssetRows(assets))
	return nil
}

func createOrder(cmd *cobra.Command, _ []string) error {
	symbol, _ := cmd.Flags().GetString("symbol")
	sideRaw, _ := cmd.Flags().GetString("side")
	typeRaw, _ := cmd.Flags().GetString("type")
	sizeRaw, _ := cmd.Flags().GetString("size")
	forceRaw, _ := cmd.Flags().GetString("force")

	side, err := classic.ParseSpotSide(sideRaw)
	if err != nil {
		return err
	}
	orderType, err := classic.ParseSpotOrderType(typeRaw)
	if err != nil {
		return err
	}
	force, err := classic.ParseSpotForce(forceRaw)
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

	p := classic.SpotPlaceOrderParams{
		Symbol:    symbol,
		Side:      side,
		OrderType: orderType,
		Force:     force,
		Size:      size,
	}
	if priceRaw != "" {
		price, perr := decimal.NewFromString(priceRaw)
		if perr != nil {
			return errors.New("invalid --price: " + perr.Error())
		}
		p.Price = price
	}
	p.ClientOid, _ = cmd.Flags().GetString("clientOid")

	ref, err := classic.NewSpotClient().PlaceOrder(p)
	if err != nil {
		return err
	}
	printer.Print(&classic.SpotOrderRefView{OrderID: ref.OrderID, ClientOid: ref.ClientOid})
	return nil
}

func cancelOrder(cmd *cobra.Command, _ []string) error {
	symbol, _ := cmd.Flags().GetString("symbol")
	orderID, _ := cmd.Flags().GetString("orderId")
	clientOid, _ := cmd.Flags().GetString("clientOid")
	if orderID == "" && clientOid == "" {
		return errors.New("one of --orderId or --clientOid is required")
	}
	ref, err := classic.NewSpotClient().CancelOrder(symbol, orderID, clientOid)
	if err != nil {
		return err
	}
	printer.Print(&classic.SpotOrderRefView{OrderID: ref.OrderID, ClientOid: ref.ClientOid})
	return nil
}

func getOrder(cmd *cobra.Command, _ []string) error {
	orderID, _ := cmd.Flags().GetString("orderId")
	clientOid, _ := cmd.Flags().GetString("clientOid")
	if orderID == "" && clientOid == "" {
		return errors.New("one of --orderId or --clientOid is required")
	}
	orders, err := classic.NewSpotClient().GetOrderInfo(orderID, clientOid)
	if err != nil {
		return err
	}
	printer.Print(classic.SpotOrderRows(orders))
	return nil
}

func openOrders(cmd *cobra.Command, _ []string) error {
	symbol, _ := cmd.Flags().GetString("symbol")
	limit, _ := cmd.Flags().GetInt("limit")
	orders, err := classic.NewSpotClient().GetOpenOrders(symbol, limit)
	if err != nil {
		return err
	}
	printer.Print(classic.SpotOpenOrderRows(orders))
	return nil
}
