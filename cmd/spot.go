package cmd

import (
	"github.com/UnipayFI/bitget-cli/cmd/spot"
	"github.com/spf13/cobra"
)

var spotCmd = &cobra.Command{
	Use:   "spot",
	Short: "Spot trading (orders & fills)",
}

func init() {
	spotCmd.AddCommand(spot.InitCmds()...)
	RootCmd.AddCommand(spotCmd)
}
