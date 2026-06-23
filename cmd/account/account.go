package account

import (
	"github.com/UnipayFI/bitget-cli/cmd/cmdutil"
	"github.com/UnipayFI/bitget-cli/exchange"
	"github.com/UnipayFI/bitget-cli/printer"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

const docBase = "https://www.bitget.com/api-doc/uta/account/"

var (
	assetsCmd = &cobra.Command{
		Use:   "assets",
		Short: "Show unified account per-coin balances (non-zero)",
		Long: `Show the unified trading account's per-coin balances. Only coins with a
non-zero equity/balance/available/debt are shown.

Docs Link: ` + docBase + "Get-Account",
		RunE: showAssets,
	}

	equityCmd = &cobra.Command{
		Use:   "equity",
		Short: "Show unified account aggregate equity and margin",
		Long: `Show the unified trading account's aggregate equity, unrealised PnL and
margin metrics (IMR / MMR / margin ratio).

Docs Link: ` + docBase + "Get-Account",
		RunE: showEquity,
	}

	healthCmd = &cobra.Command{
		Use:   "health",
		Short: "Show unified account health (margin ratio & risk)",
		Long: `Show the unified trading account's risk/health metrics: account and
effective equity, unrealised PnL, initial/maintenance margin requirements
(IMR/MMR) and the margin ratio. A margin ratio approaching 1 (100%) signals
liquidation risk.

Docs Link: ` + docBase + "Get-Account",
		RunE: showHealth,
	}

	infoCmd = &cobra.Command{
		Use:   "info",
		Short: "Show account identity and permissions",
		Long: `Show the calling account's identity and API permission metadata.

Docs Link: ` + docBase + "Get-Account-Info",
		RunE: showInfo,
	}

	settingsCmd = &cobra.Command{
		Use:   "settings",
		Short: "Show account mode settings",
		Long: `Show the unified account's mode settings (account mode, asset mode,
account level, hold mode, STP mode).

Docs Link: ` + docBase + "Get-Account-Setting",
		RunE: showSettings,
	}

	leverageConfigCmd = &cobra.Command{
		Use:   "leverage-config",
		Short: "Show per-symbol leverage / margin configuration",
		Long: `Show the per-symbol leverage and margin-mode configuration list from the
account settings.

Docs Link: ` + docBase + "Get-Account-Setting",
		RunE: showLeverageConfig,
	}

	feeRateCmd = &cobra.Command{
		Use:   "fee-rate",
		Short: "Show maker/taker fee rate for a symbol",
		Long: `Show the account's maker/taker trading fee rate for a symbol.

Docs Link: ` + docBase + "Get-Account-Fee-Rate",
		RunE: showFeeRate,
	}

	fundingAssetsCmd = &cobra.Command{
		Use:   "funding-assets",
		Short: "Show funding (P2P) account balances",
		Long: `Show the funding (P2P) account balances, optionally filtered to one coin.

Docs Link: ` + docBase + "Get-Account-Funding-Assets",
		RunE: showFundingAssets,
	}

	billsCmd = &cobra.Command{
		Use:   "bills",
		Short: "Show account financial (ledger) records",
		Long: `Show the unified account's financial (ledger) records for a category,
bounded to a 90-day lookback window.

Docs Link: ` + docBase + "Get-Financial-Records",
		RunE: showBills,
	}

	maxTransferableCmd = &cobra.Command{
		Use:   "max-transferable",
		Short: "Show max transferable amount for a coin",
		Long: `Show the maximum amount of a coin transferable out of the unified account.

Docs Link: ` + docBase + "Get-Max-Transferable",
		RunE: showMaxTransferable,
	}

	maxWithdrawalCmd = &cobra.Command{
		Use:   "max-withdrawal",
		Short: "Show max withdrawable amount for a coin",
		Long: `Show the maximum withdrawable amount of a coin, broken down by account type.

Docs Link: ` + docBase + "Get-Max-Withdrawal",
		RunE: showMaxWithdrawal,
	}

	setLeverageCmd = &cobra.Command{
		Use:   "set-leverage",
		Short: "Set leverage for a coin or futures symbol",
		Long: `Adjust the leverage for a margin coin or futures symbol.
* Futures: set --symbol; Margin: set --coin.

Docs Link: ` + docBase + "Change-Leverage",
		RunE: setLeverage,
	}

	setHoldModeCmd = &cobra.Command{
		Use:   "set-hold-mode",
		Short: "Set futures position hold mode (one-way / hedge)",
		Long: `Switch the futures position holding mode.
* --holdMode one_way_mode | hedge_mode

Docs Link: ` + docBase + "Change-Position-Mode",
		RunE: setHoldMode,
	}

	setMarginCmd = &cobra.Command{
		Use:   "set-margin",
		Short: "Add/reduce isolated position margin",
		Long: `Add or reduce isolated-position margin for a futures symbol.
* --operation add | reduce

Docs Link: ` + docBase + "Set-Margin",
		RunE: setMargin,
	}
)

// InitCmds registers flags and returns the account subcommands.
func InitCmds() []*cobra.Command {
	feeRateCmd.Flags().StringP("category", "C", "", "category: spot, usdt-futures, coin-futures, usdc-futures (required)")
	feeRateCmd.Flags().StringP("symbol", "s", "", "symbol, e.g. BTCUSDT (required)")
	feeRateCmd.MarkFlagRequired("category")
	feeRateCmd.MarkFlagRequired("symbol")

	fundingAssetsCmd.Flags().StringP("coin", "c", "", "coin filter, e.g. USDT")

	billsCmd.Flags().StringP("category", "C", "", "category: spot, usdt-futures, coin-futures, usdc-futures (required)")
	billsCmd.Flags().StringP("coin", "c", "", "coin filter")
	billsCmd.Flags().StringP("type", "t", "", "record type filter")
	billsCmd.Flags().StringP("startTime", "a", "", "start time (unix ms or \"YYYY-MM-DD HH:MM:SS\")")
	billsCmd.Flags().StringP("endTime", "e", "", "end time (unix ms or \"YYYY-MM-DD HH:MM:SS\")")
	billsCmd.Flags().StringP("limit", "l", "", "max records")
	billsCmd.MarkFlagRequired("category")

	maxTransferableCmd.Flags().StringP("coin", "c", "", "coin, e.g. USDT (required)")
	maxTransferableCmd.MarkFlagRequired("coin")

	maxWithdrawalCmd.Flags().StringP("coin", "c", "", "coin, e.g. USDT (required)")
	maxWithdrawalCmd.MarkFlagRequired("coin")

	setLeverageCmd.Flags().StringP("category", "C", "", "category: usdt-futures, coin-futures, usdc-futures, margin (required)")
	setLeverageCmd.Flags().StringP("leverage", "L", "", "leverage multiple, e.g. 10 (required)")
	setLeverageCmd.Flags().StringP("symbol", "s", "", "futures symbol, e.g. BTCUSDT")
	setLeverageCmd.Flags().StringP("coin", "c", "", "margin coin")
	setLeverageCmd.Flags().StringP("marginMode", "m", "", "margin mode: crossed, isolated")
	setLeverageCmd.Flags().StringP("posSide", "p", "", "position side: long, short")
	setLeverageCmd.Flags().String("longLeverage", "", "long leverage (isolated, two-way)")
	setLeverageCmd.Flags().String("shortLeverage", "", "short leverage (isolated, two-way)")
	setLeverageCmd.MarkFlagRequired("category")
	setLeverageCmd.MarkFlagRequired("leverage")

	setHoldModeCmd.Flags().StringP("holdMode", "H", "", "one_way_mode or hedge_mode (required)")
	setHoldModeCmd.MarkFlagRequired("holdMode")

	setMarginCmd.Flags().StringP("category", "C", "", "category: usdt-futures, coin-futures, usdc-futures (required)")
	setMarginCmd.Flags().StringP("symbol", "s", "", "futures symbol (required)")
	setMarginCmd.Flags().StringP("posSide", "p", "long", "position side: long, short")
	setMarginCmd.Flags().StringP("operation", "o", "", "add or reduce (required)")
	setMarginCmd.Flags().StringP("amount", "m", "", "margin amount (decimal) (required)")
	setMarginCmd.MarkFlagRequired("category")
	setMarginCmd.MarkFlagRequired("symbol")
	setMarginCmd.MarkFlagRequired("operation")
	setMarginCmd.MarkFlagRequired("amount")

	return []*cobra.Command{
		assetsCmd, equityCmd, healthCmd, infoCmd, settingsCmd, leverageConfigCmd,
		feeRateCmd, fundingAssetsCmd, billsCmd, maxTransferableCmd, maxWithdrawalCmd,
		setLeverageCmd, setHoldModeCmd, setMarginCmd,
	}
}

func showAssets(cmd *cobra.Command, _ []string) error {
	assets, err := exchange.NewClient().GetAccountAssets()
	if err != nil {
		return err
	}
	printer.Print(exchange.CoinAssets(assets.Assets))
	return nil
}

func showEquity(cmd *cobra.Command, _ []string) error {
	assets, err := exchange.NewClient().GetAccountAssets()
	if err != nil {
		return err
	}
	summary := exchange.AccountSummary(*assets)
	printer.Print(&summary)
	return nil
}

func showHealth(cmd *cobra.Command, _ []string) error {
	assets, err := exchange.NewClient().GetAccountAssets()
	if err != nil {
		return err
	}
	view := exchange.AccountHealthView(*assets)
	printer.Print(&view)
	return nil
}

func showInfo(cmd *cobra.Command, _ []string) error {
	info, err := exchange.NewClient().GetAccountInfo()
	if err != nil {
		return err
	}
	view := exchange.AccountInfoView(*info)
	printer.Print(&view)
	return nil
}

func showSettings(cmd *cobra.Command, _ []string) error {
	settings, err := exchange.NewClient().GetAccountSettings()
	if err != nil {
		return err
	}
	view := exchange.AccountSettingsView(*settings)
	printer.Print(&view)
	return nil
}

func showLeverageConfig(cmd *cobra.Command, _ []string) error {
	settings, err := exchange.NewClient().GetAccountSettings()
	if err != nil {
		return err
	}
	printer.Print(exchange.LeverageConfigs(settings.SymbolConfigList))
	return nil
}

func showFeeRate(cmd *cobra.Command, _ []string) error {
	categoryRaw, _ := cmd.Flags().GetString("category")
	symbol, _ := cmd.Flags().GetString("symbol")
	category, err := exchange.ParseCategory(categoryRaw)
	if err != nil {
		return err
	}
	rate, err := exchange.NewClient().GetFeeRate(category, symbol)
	if err != nil {
		return err
	}
	printer.Print(&exchange.FeeRateView{Symbol: symbol, AccountFeeRate: *rate})
	return nil
}

func showFundingAssets(cmd *cobra.Command, _ []string) error {
	coin, _ := cmd.Flags().GetString("coin")
	assets, err := exchange.NewClient().GetFundingAssets(coin)
	if err != nil {
		return err
	}
	printer.Print(exchange.FundingAssets(assets))
	return nil
}

func showBills(cmd *cobra.Command, _ []string) error {
	categoryRaw, _ := cmd.Flags().GetString("category")
	category, err := exchange.ParseCategory(categoryRaw)
	if err != nil {
		return err
	}
	p := exchange.FinancialRecordsParams{Category: category}
	p.Coin, _ = cmd.Flags().GetString("coin")
	p.RecordType, _ = cmd.Flags().GetString("type")
	p.Limit, _ = cmd.Flags().GetString("limit")
	if p.StartTime, err = cmdutil.ParseTime(cmd, "startTime"); err != nil {
		return err
	}
	if p.EndTime, err = cmdutil.ParseTime(cmd, "endTime"); err != nil {
		return err
	}
	records, err := exchange.NewClient().GetFinancialRecords(p)
	if err != nil {
		return err
	}
	printer.Print(exchange.FinancialRecordList(records.List))
	return nil
}

func showMaxTransferable(cmd *cobra.Command, _ []string) error {
	coin, _ := cmd.Flags().GetString("coin")
	res, err := exchange.NewClient().GetMaxTransferable(coin)
	if err != nil {
		return err
	}
	view := exchange.MaxTransferableView(*res)
	printer.Print(&view)
	return nil
}

func showMaxWithdrawal(cmd *cobra.Command, _ []string) error {
	coin, _ := cmd.Flags().GetString("coin")
	res, err := exchange.NewClient().GetMaxWithdrawal(coin)
	if err != nil {
		return err
	}
	view := exchange.MaxWithdrawalView(*res)
	printer.Print(&view)
	return nil
}

func setLeverage(cmd *cobra.Command, _ []string) error {
	categoryRaw, _ := cmd.Flags().GetString("category")
	category, err := exchange.ParseCategory(categoryRaw)
	if err != nil {
		return err
	}
	p := exchange.SetLeverageParams{Category: category}
	p.Leverage, _ = cmd.Flags().GetString("leverage")
	p.Symbol, _ = cmd.Flags().GetString("symbol")
	p.Coin, _ = cmd.Flags().GetString("coin")
	p.MarginMode, _ = cmd.Flags().GetString("marginMode")
	p.PosSide, _ = cmd.Flags().GetString("posSide")
	p.LongLeverage, _ = cmd.Flags().GetString("longLeverage")
	p.ShortLeverage, _ = cmd.Flags().GetString("shortLeverage")
	res, err := exchange.NewClient().SetLeverage(p)
	if err != nil {
		return err
	}
	printer.Print(map[string]string{"result": res})
	return nil
}

func setHoldMode(cmd *cobra.Command, _ []string) error {
	holdMode, _ := cmd.Flags().GetString("holdMode")
	res, err := exchange.NewClient().SetHoldMode(holdMode)
	if err != nil {
		return err
	}
	printer.Print(map[string]string{"result": res})
	return nil
}

func setMargin(cmd *cobra.Command, _ []string) error {
	categoryRaw, _ := cmd.Flags().GetString("category")
	category, err := exchange.ParseCategory(categoryRaw)
	if err != nil {
		return err
	}
	symbol, _ := cmd.Flags().GetString("symbol")
	posSide, _ := cmd.Flags().GetString("posSide")
	operation, _ := cmd.Flags().GetString("operation")
	amountRaw, _ := cmd.Flags().GetString("amount")
	amount, err := decimal.NewFromString(amountRaw)
	if err != nil {
		return err
	}
	res, err := exchange.NewClient().SetMargin(category, symbol, posSide, operation, amount)
	if err != nil {
		return err
	}
	printer.Print(map[string]string{"result": res})
	return nil
}
