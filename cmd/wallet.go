package cmd

import (
	"github.com/UnipayFI/bitget-cli/cmd/wallet"
	"github.com/spf13/cobra"
)

var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Funds: transfer, deposit and withdrawal",
}

func init() {
	walletCmd.AddCommand(wallet.InitCmds()...)
	RootCmd.AddCommand(walletCmd)
}
