package spot

import (
	"errors"

	"github.com/UnipayFI/bitget-cli/cmd/cmdutil"
	"github.com/UnipayFI/bitget-cli/exchange"
	"github.com/UnipayFI/bitget-cli/printer"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

const docBase = "https://www.bitget.com/api-doc/uta/trade/"

// category is fixed for the spot command group.
const category = exchange.CategorySpot

var (
	orderCmd = &cobra.Command{
		Use:   "order",
		Short: "Create, modify, cancel and query spot orders",
	}

	createCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create a spot order",
		Long: `Place a new spot order.
* Required: --symbol, --side, --type, --qty
* --price is required for limit orders

Docs Link: ` + docBase + "Place-Order",
		RunE: createOrder,
	}

	cancelCmd = &cobra.Command{
		Use:   "cancel",
		Short: "Cancel a spot order",
		Long: `Cancel a single spot order by --orderId or --clientOid.

Docs Link: ` + docBase + "Cancel-Order",
		RunE: cancelOrder,
	}

	modifyCmd = &cobra.Command{
		Use:   "modify",
		Short: "Modify a spot order",
		Long: `Amend a spot order's quantity and/or price.
* Identify by --orderId or --clientOid; supply at least one of --qty / --price.

Docs Link: ` + docBase + "Modify-Order",
		RunE: modifyOrder,
	}

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Query a single spot order",
		Long: `Query a single spot order by --orderId or --clientOid.

Docs Link: ` + docBase + "Get-Order-Details",
		RunE: getOrder,
	}

	openCmd = &cobra.Command{
		Use:   "open",
		Short: "List open spot orders",
		Long: `List currently open (unfilled / partially filled) spot orders.

Docs Link: ` + docBase + "Get-Order-Pending",
		RunE: openOrders,
	}

	historyCmd = &cobra.Command{
		Use:   "history",
		Short: "List historical spot orders",
		Long: `List historical spot orders, bounded to a 90-day access window.

Docs Link: ` + docBase + "Get-Order-History",
		RunE: orderHistory,
	}

	cancelAllCmd = &cobra.Command{
		Use:   "cancel-all",
		Short: "Cancel all open spot orders",
		Long: `Cancel all open spot orders, optionally limited to one --symbol.

Docs Link: ` + docBase + "Cancel-All-Order",
		RunE: cancelAll,
	}

	fillsCmd = &cobra.Command{
		Use:   "fills",
		Short: "List spot trade fills",
		Long: `List spot trade fills, bounded to a 90-day access window.

Docs Link: ` + docBase + "Get-Order-Fills",
		RunE: listFills,
	}
)

// InitCmds registers flags and returns the spot subcommands.
func InitCmds() []*cobra.Command {
	createCmd.Flags().StringP("symbol", "s", "", "symbol, e.g. BTCUSDT (required)")
	createCmd.Flags().StringP("side", "S", "", "buy or sell (required)")
	createCmd.Flags().StringP("type", "t", "", "limit or market (required)")
	createCmd.Flags().StringP("qty", "q", "", "order quantity (decimal) (required)")
	createCmd.Flags().StringP("price", "p", "", "order price (required for limit)")
	createCmd.Flags().StringP("tif", "T", "", "time in force: gtc, post_only, fok, ioc")
	createCmd.Flags().String("clientOid", "", "client order id")
	createCmd.MarkFlagRequired("symbol")
	createCmd.MarkFlagRequired("side")
	createCmd.MarkFlagRequired("type")
	createCmd.MarkFlagRequired("qty")

	cancelCmd.Flags().StringP("orderId", "i", "", "order id")
	cancelCmd.Flags().StringP("clientOid", "c", "", "client order id")

	modifyCmd.Flags().StringP("orderId", "i", "", "order id")
	modifyCmd.Flags().StringP("clientOid", "c", "", "client order id")
	modifyCmd.Flags().StringP("symbol", "s", "", "symbol, e.g. BTCUSDT")
	modifyCmd.Flags().StringP("qty", "q", "", "new quantity (decimal)")
	modifyCmd.Flags().StringP("price", "p", "", "new price (decimal)")

	getCmd.Flags().StringP("orderId", "i", "", "order id")
	getCmd.Flags().StringP("clientOid", "c", "", "client order id")

	openCmd.Flags().StringP("symbol", "s", "", "symbol filter")
	openCmd.Flags().StringP("limit", "l", "", "max records")

	historyCmd.Flags().StringP("symbol", "s", "", "symbol filter")
	historyCmd.Flags().StringP("startTime", "a", "", "start time (unix ms or \"YYYY-MM-DD HH:MM:SS\")")
	historyCmd.Flags().StringP("endTime", "e", "", "end time (unix ms or \"YYYY-MM-DD HH:MM:SS\")")
	historyCmd.Flags().StringP("limit", "l", "", "max records")

	cancelAllCmd.Flags().StringP("symbol", "s", "", "symbol filter (cancel all in spot when omitted)")

	fillsCmd.Flags().StringP("orderId", "i", "", "order id filter")
	fillsCmd.Flags().StringP("startTime", "a", "", "start time (unix ms or \"YYYY-MM-DD HH:MM:SS\")")
	fillsCmd.Flags().StringP("endTime", "e", "", "end time (unix ms or \"YYYY-MM-DD HH:MM:SS\")")
	fillsCmd.Flags().StringP("limit", "l", "", "max records")

	orderCmd.AddCommand(createCmd, cancelCmd, modifyCmd, getCmd, openCmd, historyCmd, cancelAllCmd)
	return []*cobra.Command{orderCmd, fillsCmd}
}

func createOrder(cmd *cobra.Command, _ []string) error {
	symbol, _ := cmd.Flags().GetString("symbol")
	sideRaw, _ := cmd.Flags().GetString("side")
	typeRaw, _ := cmd.Flags().GetString("type")
	qtyRaw, _ := cmd.Flags().GetString("qty")

	side, err := exchange.ParseSide(sideRaw)
	if err != nil {
		return err
	}
	orderType, err := exchange.ParseOrderType(typeRaw)
	if err != nil {
		return err
	}
	qty, err := decimal.NewFromString(qtyRaw)
	if err != nil {
		return errors.New("invalid --qty: " + err.Error())
	}

	p := exchange.PlaceOrderParams{
		Category:  category,
		Symbol:    symbol,
		Side:      side,
		OrderType: orderType,
		Qty:       qty,
	}
	if priceRaw, _ := cmd.Flags().GetString("price"); priceRaw != "" {
		price, perr := decimal.NewFromString(priceRaw)
		if perr != nil {
			return errors.New("invalid --price: " + perr.Error())
		}
		p.Price = price
	}
	p.TimeInForce, _ = cmd.Flags().GetString("tif")
	p.ClientOid, _ = cmd.Flags().GetString("clientOid")

	ref, err := exchange.NewClient().PlaceOrder(p)
	if err != nil {
		return err
	}
	view := exchange.OrderRefView(*ref)
	printer.Print(&view)
	return nil
}

func cancelOrder(cmd *cobra.Command, _ []string) error {
	orderID, _ := cmd.Flags().GetString("orderId")
	clientOid, _ := cmd.Flags().GetString("clientOid")
	if orderID == "" && clientOid == "" {
		return errors.New("one of --orderId or --clientOid is required")
	}
	ref, err := exchange.NewClient().CancelOrder(category, orderID, clientOid)
	if err != nil {
		return err
	}
	view := exchange.OrderRefView(*ref)
	printer.Print(&view)
	return nil
}

func modifyOrder(cmd *cobra.Command, _ []string) error {
	p := exchange.ModifyOrderParams{Category: category}
	p.OrderId, _ = cmd.Flags().GetString("orderId")
	p.ClientOid, _ = cmd.Flags().GetString("clientOid")
	p.Symbol, _ = cmd.Flags().GetString("symbol")
	if p.OrderId == "" && p.ClientOid == "" {
		return errors.New("one of --orderId or --clientOid is required")
	}
	if qtyRaw, _ := cmd.Flags().GetString("qty"); qtyRaw != "" {
		qty, err := decimal.NewFromString(qtyRaw)
		if err != nil {
			return errors.New("invalid --qty: " + err.Error())
		}
		p.Qty = qty
	}
	if priceRaw, _ := cmd.Flags().GetString("price"); priceRaw != "" {
		price, err := decimal.NewFromString(priceRaw)
		if err != nil {
			return errors.New("invalid --price: " + err.Error())
		}
		p.Price = price
	}
	ref, err := exchange.NewClient().ModifyOrder(p)
	if err != nil {
		return err
	}
	view := exchange.OrderRefView(*ref)
	printer.Print(&view)
	return nil
}

func getOrder(cmd *cobra.Command, _ []string) error {
	orderID, _ := cmd.Flags().GetString("orderId")
	clientOid, _ := cmd.Flags().GetString("clientOid")
	if orderID == "" && clientOid == "" {
		return errors.New("one of --orderId or --clientOid is required")
	}
	order, err := exchange.NewClient().GetOrderInfo(orderID, clientOid)
	if err != nil {
		return err
	}
	printer.Print(exchange.OrderRows{*order})
	return nil
}

func openOrders(cmd *cobra.Command, _ []string) error {
	symbol, _ := cmd.Flags().GetString("symbol")
	limit, _ := cmd.Flags().GetString("limit")
	list, err := exchange.NewClient().GetOpenOrders(category, symbol, limit)
	if err != nil {
		return err
	}
	printer.Print(exchange.OrderRows(list.List))
	return nil
}

func orderHistory(cmd *cobra.Command, _ []string) error {
	p := exchange.HistoryParams{Category: category}
	p.Symbol, _ = cmd.Flags().GetString("symbol")
	p.Limit, _ = cmd.Flags().GetString("limit")
	var err error
	if p.StartTime, err = cmdutil.ParseTime(cmd, "startTime"); err != nil {
		return err
	}
	if p.EndTime, err = cmdutil.ParseTime(cmd, "endTime"); err != nil {
		return err
	}
	list, err := exchange.NewClient().GetOrderHistory(p)
	if err != nil {
		return err
	}
	printer.Print(exchange.OrderRows(list.List))
	return nil
}

func cancelAll(cmd *cobra.Command, _ []string) error {
	symbol, _ := cmd.Flags().GetString("symbol")
	res, err := exchange.NewClient().CancelSymbolOrders(category, symbol)
	if err != nil {
		return err
	}
	printer.Print(exchange.CancelResults(res.List))
	return nil
}

func listFills(cmd *cobra.Command, _ []string) error {
	p := exchange.HistoryParams{Category: category}
	p.OrderId, _ = cmd.Flags().GetString("orderId")
	p.Limit, _ = cmd.Flags().GetString("limit")
	var err error
	if p.StartTime, err = cmdutil.ParseTime(cmd, "startTime"); err != nil {
		return err
	}
	if p.EndTime, err = cmdutil.ParseTime(cmd, "endTime"); err != nil {
		return err
	}
	list, err := exchange.NewClient().GetFills(p)
	if err != nil {
		return err
	}
	printer.Print(exchange.FillRows(list.List))
	return nil
}
