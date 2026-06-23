package cmd

import (
	"github.com/UnipayFI/bitget-cli/cmd/account"
	"github.com/UnipayFI/bitget-cli/cmd/futures"
	"github.com/UnipayFI/bitget-cli/cmd/spot"
	"github.com/UnipayFI/bitget-cli/cmd/wallet"
	"github.com/spf13/cobra"
)

// utaCmd is the parent command for every Unified Trading Account (UTA) feature.
// It groups the v3 /api/v3/* endpoints under `bitget-cli UTA ...`; the classic
// (/api/v2/*) products live at the top level instead.
var utaCmd = &cobra.Command{
	Use:   "UTA",
	Short: "Unified Trading Account (UTA) commands",
	Long: `Unified Trading Account (UTA) commands, backed by Bitget's v3 private REST
API (/api/v3/*). The unified account serves spot and every futures line from one
client; the product is chosen per command via the --category flag.

Docs Link: https://www.bitget.com/api-doc/uta/intro`,
}

var (
	utaAccountCmd = &cobra.Command{
		Use:   "account",
		Short: "Account info, balances, settings and leverage",
	}

	utaSpotCmd = &cobra.Command{
		Use:   "spot",
		Short: "Spot trading (orders & fills)",
	}

	utaFuturesCmd = &cobra.Command{
		Use:   "futures",
		Short: "Futures trading (orders, positions & fills)",
		Long: `Futures trading commands. The product line is selected with the persistent
--category flag (default usdt-futures); it applies to every subcommand.`,
	}

	utaWalletCmd = &cobra.Command{
		Use:   "wallet",
		Short: "Funds: transfer, deposit and withdrawal",
	}
)

func init() {
	utaAccountCmd.AddCommand(account.InitCmds()...)

	utaSpotCmd.AddCommand(spot.InitCmds()...)

	utaFuturesCmd.PersistentFlags().StringP("category", "C", "usdt-futures",
		"futures category: usdt-futures (usdt), coin-futures (coin), usdc-futures (usdc)")
	utaFuturesCmd.AddCommand(futures.InitCmds()...)

	utaWalletCmd.AddCommand(wallet.InitCmds()...)

	utaCmd.AddCommand(utaAccountCmd, utaSpotCmd, utaFuturesCmd, utaWalletCmd)
	RootCmd.AddCommand(utaCmd)
}
