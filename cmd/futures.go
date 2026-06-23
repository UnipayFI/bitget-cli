package cmd

import (
	"github.com/UnipayFI/bitget-cli/cmd/classicfutures"
	"github.com/spf13/cobra"
)

// futuresCmd is the top-level (classic-account) futures command group. The
// product line is selected with the persistent --product flag; it applies to
// every subcommand.
var futuresCmd = &cobra.Command{
	Use:   "futures",
	Short: "Classic futures trading (account, positions & orders)",
	Long: `Classic-account futures commands, backed by Bitget's v2 private REST API
(/api/v2/mix/*). The product line is selected with the persistent --product
flag (default usdt-futures). For the unified account use "bitget-cli UTA
futures" instead.`,
}

func init() {
	futuresCmd.PersistentFlags().StringP("product", "P", "usdt-futures",
		"futures product: usdt-futures (usdt), coin-futures (coin), usdc-futures (usdc)")
	futuresCmd.AddCommand(classicfutures.InitCmds()...)
	RootCmd.AddCommand(futuresCmd)
}
