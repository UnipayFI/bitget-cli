package cmd

import (
	"github.com/UnipayFI/bitget-cli/cmd/account"
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Account info, balances, settings and leverage",
}

func init() {
	accountCmd.AddCommand(account.InitCmds()...)
	RootCmd.AddCommand(accountCmd)
}
