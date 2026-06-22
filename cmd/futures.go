package cmd

import (
	"github.com/UnipayFI/bitget-cli/cmd/futures"
	"github.com/spf13/cobra"
)

var futuresCmd = &cobra.Command{
	Use:   "futures",
	Short: "Futures trading (orders, positions & fills)",
	Long: `Futures trading commands. The product line is selected with the persistent
--category flag (default usdt-futures); it applies to every subcommand.`,
}

func init() {
	futuresCmd.PersistentFlags().StringP("category", "C", "usdt-futures",
		"futures category: usdt-futures (usdt), coin-futures (coin), usdc-futures (usdc)")
	futuresCmd.AddCommand(futures.InitCmds()...)
	RootCmd.AddCommand(futuresCmd)
}
