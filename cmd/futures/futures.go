package futures

import (
	"errors"

	"github.com/UnipayFI/bitget-cli/cmd/cmdutil"
	"github.com/UnipayFI/bitget-cli/exchange"
	"github.com/UnipayFI/bitget-cli/printer"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

const docBase = "https://www.bitget.com/api-doc/uta/trade/"

var (
	orderCmd = &cobra.Command{
		Use:   "order",
		Short: "Create, modify, cancel and query futures orders",
	}

	createCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create a futures order",
		Long: `Place a new futures order.
* Required: --symbol, --side, --type, --qty
* --price is required for limit orders
* --posSide (long/short) is required in hedge mode

Docs Link: ` + docBase + "Place-Order",
		RunE: createOrder,
	}

	cancelCmd = &cobra.Command{
		Use:   "cancel",
		Short: "Cancel a futures order",
		Long: `Cancel a single futures order by --orderId or --clientOid.

Docs Link: ` + docBase + "Cancel-Order",
		RunE: cancelOrder,
	}

	modifyCmd = &cobra.Command{
		Use:   "modify",
		Short: "Modify a futures order",
		Long: `Amend a futures order's quantity and/or price.
* Identify by --orderId or --clientOid; supply at least one of --qty / --price.

Docs Link: ` + docBase + "Modify-Order",
		RunE: modifyOrder,
	}

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Query a single futures order",
		Long: `Query a single futures order by --orderId or --clientOid.

Docs Link: ` + docBase + "Get-Order-Details",
		RunE: getOrder,
	}

	openCmd = &cobra.Command{
		Use:   "open",
		Short: "List open futures orders",
		Long: `List currently open (unfilled / partially filled) futures orders.

Docs Link: ` + docBase + "Get-Open-Orders",
		RunE: openOrders,
	}

	historyCmd = &cobra.Command{
		Use:   "history",
		Short: "List historical futures orders",
		Long: `List historical futures orders, bounded to a 90-day access window.

Docs Link: ` + docBase + "Get-Order-History",
		RunE: orderHistory,
	}

	cancelAllCmd = &cobra.Command{
		Use:   "cancel-all",
		Short: "Cancel all open futures orders",
		Long: `Cancel all open futures orders in the category, optionally limited to one --symbol.

Docs Link: ` + docBase + "Cancel-All-Orders",
		RunE: cancelAll,
	}

	fillsCmd = &cobra.Command{
		Use:   "fills",
		Short: "List futures trade fills",
		Long: `List futures trade fills, bounded to a 90-day access window.

Docs Link: ` + docBase + "Get-Fill-History",
		RunE: listFills,
	}

	positionCmd = &cobra.Command{
		Use:   "position",
		Short: "Query and close futures positions",
	}

	positionListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List open positions",
		Long: `List the account's open futures positions in the category.

Docs Link: ` + docBase + "Get-Position-Info",
		RunE: listPositions,
	}

	positionHistoryCmd = &cobra.Command{
		Use:   "history",
		Short: "List closed positions",
		Long: `List closed/historical futures positions, bounded to a 90-day window.

Docs Link: ` + docBase + "Get-Positions-History",
		RunE: positionHistory,
	}

	adlRankCmd = &cobra.Command{
		Use:   "adl-rank",
		Short: "Show ADL ranking per position",
		Long: `Show the auto-deleveraging (ADL) rank for each open position.

Docs Link: ` + docBase + "Get-Position-ADL-Rank",
		RunE: adlRank,
	}

	closeCmd = &cobra.Command{
		Use:   "close",
		Short: "Market-close positions",
		Long: `Market-close positions. Without --symbol closes all in the category;
without --posSide closes both sides.

Docs Link: ` + docBase + "Close-All-Positions",
		RunE: closePositions,
	}
)

// InitCmds registers flags and returns the futures subcommands.
func InitCmds() []*cobra.Command {
	createCmd.Flags().StringP("symbol", "s", "", "symbol, e.g. BTCUSDT (required)")
	createCmd.Flags().StringP("side", "S", "", "buy or sell (required)")
	createCmd.Flags().StringP("type", "t", "", "limit or market (required)")
	createCmd.Flags().StringP("qty", "q", "", "order quantity (decimal) (required)")
	createCmd.Flags().StringP("price", "p", "", "order price (required for limit)")
	createCmd.Flags().StringP("tif", "T", "", "time in force: gtc, post_only, fok, ioc")
	createCmd.Flags().StringP("posSide", "P", "", "position side: long, short (hedge mode)")
	createCmd.Flags().StringP("reduceOnly", "r", "", "reduce-only: yes or no")
	createCmd.Flags().StringP("marginMode", "m", "", "margin mode: crossed, isolated")
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

	cancelAllCmd.Flags().StringP("symbol", "s", "", "symbol filter")

	fillsCmd.Flags().StringP("orderId", "i", "", "order id filter")
	fillsCmd.Flags().StringP("startTime", "a", "", "start time (unix ms or \"YYYY-MM-DD HH:MM:SS\")")
	fillsCmd.Flags().StringP("endTime", "e", "", "end time (unix ms or \"YYYY-MM-DD HH:MM:SS\")")
	fillsCmd.Flags().StringP("limit", "l", "", "max records")

	positionListCmd.Flags().StringP("symbol", "s", "", "symbol filter")
	positionListCmd.Flags().StringP("posSide", "P", "", "position side: long, short")

	positionHistoryCmd.Flags().StringP("symbol", "s", "", "symbol filter")
	positionHistoryCmd.Flags().StringP("startTime", "a", "", "start time (unix ms or \"YYYY-MM-DD HH:MM:SS\")")
	positionHistoryCmd.Flags().StringP("endTime", "e", "", "end time (unix ms or \"YYYY-MM-DD HH:MM:SS\")")
	positionHistoryCmd.Flags().StringP("limit", "l", "", "max records")

	closeCmd.Flags().StringP("symbol", "s", "", "symbol filter")
	closeCmd.Flags().StringP("posSide", "P", "", "position side: long, short")

	orderCmd.AddCommand(createCmd, cancelCmd, modifyCmd, getCmd, openCmd, historyCmd, cancelAllCmd)
	positionCmd.AddCommand(positionListCmd, positionHistoryCmd, adlRankCmd, closeCmd)
	return []*cobra.Command{orderCmd, positionCmd, fillsCmd}
}

// resolveCategory reads the inherited persistent --category flag.
func resolveCategory(cmd *cobra.Command) (exchange.Category, error) {
	raw, _ := cmd.Flags().GetString("category")
	if raw == "" {
		raw = "usdt-futures"
	}
	return exchange.ParseCategory(raw)
}

func createOrder(cmd *cobra.Command, _ []string) error {
	category, err := resolveCategory(cmd)
	if err != nil {
		return err
	}
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
	p.PosSide, _ = cmd.Flags().GetString("posSide")
	p.ReduceOnly, _ = cmd.Flags().GetString("reduceOnly")
	p.MarginMode, _ = cmd.Flags().GetString("marginMode")
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
	category, err := resolveCategory(cmd)
	if err != nil {
		return err
	}
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
	category, err := resolveCategory(cmd)
	if err != nil {
		return err
	}
	p := exchange.ModifyOrderParams{Category: category}
	p.OrderId, _ = cmd.Flags().GetString("orderId")
	p.ClientOid, _ = cmd.Flags().GetString("clientOid")
	p.Symbol, _ = cmd.Flags().GetString("symbol")
	if p.OrderId == "" && p.ClientOid == "" {
		return errors.New("one of --orderId or --clientOid is required")
	}
	if qtyRaw, _ := cmd.Flags().GetString("qty"); qtyRaw != "" {
		qty, qerr := decimal.NewFromString(qtyRaw)
		if qerr != nil {
			return errors.New("invalid --qty: " + qerr.Error())
		}
		p.Qty = qty
	}
	if priceRaw, _ := cmd.Flags().GetString("price"); priceRaw != "" {
		price, perr := decimal.NewFromString(priceRaw)
		if perr != nil {
			return errors.New("invalid --price: " + perr.Error())
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
	category, err := resolveCategory(cmd)
	if err != nil {
		return err
	}
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
	category, err := resolveCategory(cmd)
	if err != nil {
		return err
	}
	p := exchange.HistoryParams{Category: category}
	p.Symbol, _ = cmd.Flags().GetString("symbol")
	p.Limit, _ = cmd.Flags().GetString("limit")
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
	category, err := resolveCategory(cmd)
	if err != nil {
		return err
	}
	symbol, _ := cmd.Flags().GetString("symbol")
	res, err := exchange.NewClient().CancelSymbolOrders(category, symbol)
	if err != nil {
		return err
	}
	printer.Print(exchange.CancelResults(res.List))
	return nil
}

func listFills(cmd *cobra.Command, _ []string) error {
	category, err := resolveCategory(cmd)
	if err != nil {
		return err
	}
	p := exchange.HistoryParams{Category: category}
	p.OrderId, _ = cmd.Flags().GetString("orderId")
	p.Limit, _ = cmd.Flags().GetString("limit")
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

func listPositions(cmd *cobra.Command, _ []string) error {
	category, err := resolveCategory(cmd)
	if err != nil {
		return err
	}
	symbol, _ := cmd.Flags().GetString("symbol")
	posSide, _ := cmd.Flags().GetString("posSide")
	positions, err := exchange.NewClient().GetPositions(category, symbol, posSide)
	if err != nil {
		return err
	}
	printer.Print(exchange.Positions(positions))
	return nil
}

func positionHistory(cmd *cobra.Command, _ []string) error {
	category, err := resolveCategory(cmd)
	if err != nil {
		return err
	}
	p := exchange.HistoryParams{Category: category}
	p.Symbol, _ = cmd.Flags().GetString("symbol")
	p.Limit, _ = cmd.Flags().GetString("limit")
	if p.StartTime, err = cmdutil.ParseTime(cmd, "startTime"); err != nil {
		return err
	}
	if p.EndTime, err = cmdutil.ParseTime(cmd, "endTime"); err != nil {
		return err
	}
	history, err := exchange.NewClient().GetPositionHistory(p)
	if err != nil {
		return err
	}
	printer.Print(exchange.PositionHistoryRows(history.List))
	return nil
}

func adlRank(cmd *cobra.Command, _ []string) error {
	ranks, err := exchange.NewClient().GetPositionADLRank()
	if err != nil {
		return err
	}
	printer.Print(exchange.ADLRankRows(ranks))
	return nil
}

func closePositions(cmd *cobra.Command, _ []string) error {
	category, err := resolveCategory(cmd)
	if err != nil {
		return err
	}
	symbol, _ := cmd.Flags().GetString("symbol")
	posSide, _ := cmd.Flags().GetString("posSide")
	res, err := exchange.NewClient().ClosePositions(category, symbol, posSide)
	if err != nil {
		return err
	}
	printer.Print(exchange.OrderResults(res.List))
	return nil
}
