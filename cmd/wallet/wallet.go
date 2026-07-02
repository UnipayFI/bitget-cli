package wallet

import (
	"errors"
	"time"

	"github.com/UnipayFI/bitget-cli/cmd/cmdutil"
	"github.com/UnipayFI/bitget-cli/exchange"
	"github.com/UnipayFI/bitget-cli/printer"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

const docBase = "https://www.bitget.com/api-doc/uta/account/"

var (
	transferCmd = &cobra.Command{
		Use:   "transfer",
		Short: "Transfer a coin between account types",
		Long: `Transfer a coin between two account types within the same account.
* Account types: spot, p2p, coin_futures, usdt_futures, usdc_futures, crossed_margin, isolated_margin, uta

Docs Link: ` + docBase + "transfer",
		RunE: doTransfer,
	}

	transferableCoinsCmd = &cobra.Command{
		Use:   "transferable-coins",
		Short: "List coins transferable between two account types",
		Long: `List the coins that can be transferred between the given account types.

Docs Link: ` + docBase + "transfer/Get-Transfer-Coins",
		RunE: transferableCoins,
	}

	depositCmd = &cobra.Command{
		Use:   "deposit",
		Short: "Deposit address and records",
	}

	depositAddressCmd = &cobra.Command{
		Use:   "address",
		Short: "Get on-chain deposit address",
		Long: `Get the on-chain deposit address for a coin, optionally on a chain.

Docs Link: ` + docBase + "deposit/Get-Deposit-Address",
		RunE: depositAddress,
	}

	depositRecordsCmd = &cobra.Command{
		Use:   "records",
		Short: "List deposit records",
		Long: `List deposit records within a time window (default: last 30 days).

Docs Link: ` + docBase + "deposit/Get-Deposit-Records",
		RunE: depositRecords,
	}

	withdrawCmd = &cobra.Command{
		Use:   "withdraw",
		Short: "Withdraw, list records and address book",
	}

	withdrawCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"do"},
		Short:   "Submit a withdrawal",
		Long: `Submit a withdrawal (on-chain or internal).
* Required: --coin, --transferType (on_chain | internal_transfer), --address, --size
* --chain is required for on-chain withdrawals

Docs Link: ` + docBase + "withdrawal",
		RunE: doWithdraw,
	}

	withdrawRecordsCmd = &cobra.Command{
		Use:   "records",
		Short: "List withdrawal records",
		Long: `List withdrawal records within a time window (default: last 30 days).

Docs Link: ` + docBase + "withdrawal/Get-Withdrawal-Records",
		RunE: withdrawRecords,
	}

	withdrawAddressCmd = &cobra.Command{
		Use:   "address",
		Short: "List saved withdrawal addresses",
		Long: `List the saved withdrawal address book entries.

Docs Link: ` + docBase + "withdrawal/Get-Withdraw-Address",
		RunE: withdrawAddress,
	}

	withdrawCancelCmd = &cobra.Command{
		Use:   "cancel",
		Short: "Cancel a pending withdrawal",
		Long: `Cancel a withdrawal still within its cooling-off period.
* Identify by --orderId or --clientOid.

Docs Link: ` + docBase + "withdrawal/Cancel-Withdrawal",
		RunE: cancelWithdraw,
	}
)

// InitCmds registers flags and returns the wallet subcommands.
func InitCmds() []*cobra.Command {
	transferCmd.Flags().StringP("fromType", "f", "", "source account type (required)")
	transferCmd.Flags().StringP("toType", "t", "", "target account type (required)")
	transferCmd.Flags().StringP("coin", "c", "", "coin, e.g. USDT (required)")
	transferCmd.Flags().StringP("amount", "m", "", "amount (decimal) (required)")
	transferCmd.Flags().StringP("symbol", "s", "", "isolated spot-margin symbol")
	transferCmd.Flags().String("allowBorrow", "", "auto-borrow when insufficient: yes or no")
	transferCmd.Flags().String("clientOid", "", "client transaction id")
	transferCmd.MarkFlagRequired("fromType")
	transferCmd.MarkFlagRequired("toType")
	transferCmd.MarkFlagRequired("coin")
	transferCmd.MarkFlagRequired("amount")

	transferableCoinsCmd.Flags().StringP("fromType", "f", "", "source account type (required)")
	transferableCoinsCmd.Flags().StringP("toType", "t", "", "target account type (required)")
	transferableCoinsCmd.MarkFlagRequired("fromType")
	transferableCoinsCmd.MarkFlagRequired("toType")

	depositAddressCmd.Flags().StringP("coin", "c", "", "coin, e.g. USDT (required)")
	depositAddressCmd.Flags().String("chain", "", "chain, e.g. trc20")
	depositAddressCmd.MarkFlagRequired("coin")

	depositRecordsCmd.Flags().StringP("coin", "c", "", "coin filter")
	depositRecordsCmd.Flags().StringP("startTime", "a", "", "start time (unix ms or \"YYYY-MM-DD HH:MM:SS\")")
	depositRecordsCmd.Flags().StringP("endTime", "e", "", "end time (unix ms or \"YYYY-MM-DD HH:MM:SS\")")
	depositRecordsCmd.Flags().StringP("limit", "l", "", "max records")

	withdrawCreateCmd.Flags().StringP("coin", "c", "", "coin, e.g. USDT (required)")
	withdrawCreateCmd.Flags().StringP("transferType", "T", "", "on_chain or internal_transfer (required)")
	withdrawCreateCmd.Flags().StringP("address", "d", "", "destination address / UID / email / mobile (required)")
	withdrawCreateCmd.Flags().StringP("size", "m", "", "amount (decimal) (required)")
	withdrawCreateCmd.Flags().String("chain", "", "chain, e.g. trc20 (required for on-chain)")
	withdrawCreateCmd.Flags().String("tag", "", "address tag/memo")
	withdrawCreateCmd.Flags().String("innerToType", "", "internal address type: uid, email, mobile")
	withdrawCreateCmd.Flags().String("remark", "", "remark")
	withdrawCreateCmd.Flags().String("clientOid", "", "client order id")
	withdrawCreateCmd.MarkFlagRequired("coin")
	withdrawCreateCmd.MarkFlagRequired("transferType")
	withdrawCreateCmd.MarkFlagRequired("address")
	withdrawCreateCmd.MarkFlagRequired("size")

	withdrawRecordsCmd.Flags().StringP("coin", "c", "", "coin filter")
	withdrawRecordsCmd.Flags().StringP("startTime", "a", "", "start time (unix ms or \"YYYY-MM-DD HH:MM:SS\")")
	withdrawRecordsCmd.Flags().StringP("endTime", "e", "", "end time (unix ms or \"YYYY-MM-DD HH:MM:SS\")")
	withdrawRecordsCmd.Flags().StringP("limit", "l", "", "max records")

	withdrawAddressCmd.Flags().StringP("coin", "c", "", "coin filter")
	withdrawAddressCmd.Flags().StringP("type", "t", "", "address type: EVM, regular, universal, internal")

	withdrawCancelCmd.Flags().StringP("orderId", "i", "", "withdrawal order id")
	withdrawCancelCmd.Flags().StringP("clientOid", "c", "", "client order id")

	depositCmd.AddCommand(depositAddressCmd, depositRecordsCmd)
	withdrawCmd.AddCommand(withdrawCreateCmd, withdrawRecordsCmd, withdrawAddressCmd, withdrawCancelCmd)
	return []*cobra.Command{transferCmd, transferableCoinsCmd, depositCmd, withdrawCmd}
}

// timeWindow resolves start/end time flags, defaulting to the last 30 days.
func timeWindow(cmd *cobra.Command) (start, end time.Time, err error) {
	if start, err = cmdutil.ParseTime(cmd, "startTime"); err != nil {
		return
	}
	if end, err = cmdutil.ParseTime(cmd, "endTime"); err != nil {
		return
	}
	if end.IsZero() {
		end = time.Now()
	}
	if start.IsZero() {
		start = end.AddDate(0, 0, -30)
	}
	return
}

func doTransfer(cmd *cobra.Command, _ []string) error {
	amountRaw, _ := cmd.Flags().GetString("amount")
	amount, err := decimal.NewFromString(amountRaw)
	if err != nil {
		return errors.New("invalid --amount: " + err.Error())
	}
	p := exchange.TransferParams{Amount: amount}
	p.FromType, _ = cmd.Flags().GetString("fromType")
	p.ToType, _ = cmd.Flags().GetString("toType")
	p.Coin, _ = cmd.Flags().GetString("coin")
	p.Symbol, _ = cmd.Flags().GetString("symbol")
	p.AllowBorrow, _ = cmd.Flags().GetString("allowBorrow")
	p.ClientOid, _ = cmd.Flags().GetString("clientOid")
	res, err := exchange.NewClient().Transfer(p)
	if err != nil {
		return err
	}
	view := exchange.TransferResultView(*res)
	printer.Print(&view)
	return nil
}

func transferableCoins(cmd *cobra.Command, _ []string) error {
	fromType, _ := cmd.Flags().GetString("fromType")
	toType, _ := cmd.Flags().GetString("toType")
	coins, err := exchange.NewClient().TransferableCoins(fromType, toType)
	if err != nil {
		return err
	}
	printer.Print(exchange.CoinList(coins))
	return nil
}

func depositAddress(cmd *cobra.Command, _ []string) error {
	coin, _ := cmd.Flags().GetString("coin")
	chain, _ := cmd.Flags().GetString("chain")
	addr, err := exchange.NewClient().GetDepositAddress(coin, chain)
	if err != nil {
		return err
	}
	view := exchange.DepositAddressView(*addr)
	printer.Print(&view)
	return nil
}

func depositRecords(cmd *cobra.Command, _ []string) error {
	start, end, err := timeWindow(cmd)
	if err != nil {
		return err
	}
	coin, _ := cmd.Flags().GetString("coin")
	limit, _ := cmd.Flags().GetString("limit")
	records, err := exchange.NewClient().GetDepositRecords(coin, start, end, limit)
	if err != nil {
		return err
	}
	printer.Print(exchange.DepositRecords(records))
	return nil
}

func doWithdraw(cmd *cobra.Command, _ []string) error {
	sizeRaw, _ := cmd.Flags().GetString("size")
	size, err := decimal.NewFromString(sizeRaw)
	if err != nil {
		return errors.New("invalid --size: " + err.Error())
	}
	p := exchange.WithdrawParams{Size: size}
	p.Coin, _ = cmd.Flags().GetString("coin")
	p.TransferType, _ = cmd.Flags().GetString("transferType")
	p.Address, _ = cmd.Flags().GetString("address")
	p.Chain, _ = cmd.Flags().GetString("chain")
	p.Tag, _ = cmd.Flags().GetString("tag")
	p.InnerToType, _ = cmd.Flags().GetString("innerToType")
	p.Remark, _ = cmd.Flags().GetString("remark")
	p.ClientOid, _ = cmd.Flags().GetString("clientOid")
	res, err := exchange.NewClient().Withdraw(p)
	if err != nil {
		return err
	}
	view := exchange.WithdrawResultView(*res)
	printer.Print(&view)
	return nil
}

func withdrawRecords(cmd *cobra.Command, _ []string) error {
	start, end, err := timeWindow(cmd)
	if err != nil {
		return err
	}
	coin, _ := cmd.Flags().GetString("coin")
	limit, _ := cmd.Flags().GetString("limit")
	records, err := exchange.NewClient().GetWithdrawalRecords(coin, start, end, limit)
	if err != nil {
		return err
	}
	printer.Print(exchange.WithdrawalRecords(records))
	return nil
}

func withdrawAddress(cmd *cobra.Command, _ []string) error {
	coin, _ := cmd.Flags().GetString("coin")
	addressType, _ := cmd.Flags().GetString("type")
	book, err := exchange.NewClient().GetWithdrawAddress(coin, addressType)
	if err != nil {
		return err
	}
	printer.Print(exchange.WithdrawAddresses(book.AddressList))
	return nil
}

func cancelWithdraw(cmd *cobra.Command, _ []string) error {
	orderID, _ := cmd.Flags().GetString("orderId")
	clientOid, _ := cmd.Flags().GetString("clientOid")
	if orderID == "" && clientOid == "" {
		return errors.New("one of --orderId or --clientOid is required")
	}
	res, err := exchange.NewClient().CancelWithdrawal(orderID, clientOid)
	if err != nil {
		return err
	}
	printer.Print(map[string]string{"result": res})
	return nil
}
