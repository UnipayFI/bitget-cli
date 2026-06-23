package cmd

import (
	"github.com/UnipayFI/bitget-cli/cmd/classicspot"
	"github.com/spf13/cobra"
)

// spotCmd is the top-level (classic-account) spot command group.
var spotCmd = &cobra.Command{
	Use:   "spot",
	Short: "Classic spot trading (account, balances & orders)",
	Long: `Classic-account spot commands, backed by Bitget's v2 private REST API
(/api/v2/spot/*). For the unified account use "bitget-cli UTA spot" instead.`,
}

func init() {
	spotCmd.AddCommand(classicspot.InitCmds()...)
	RootCmd.AddCommand(spotCmd)
}
