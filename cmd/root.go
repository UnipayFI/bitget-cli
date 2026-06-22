package cmd

import (
	"errors"

	"github.com/UnipayFI/bitget-cli/config"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "bitget-cli",
	Short: "Bitget UTA (Unified Trading Account) API for CLI",
	Long: `bitget-cli is a command-line client for Bitget's Unified Trading Account
(UTA) private REST API: spot and futures trading, account, positions and funds.

Credentials are read from the environment:
  BITGET_API_KEY, BITGET_API_SECRET, BITGET_PASSPHRASE   (required)
  BITGET_PROXY, BITGET_LOCALE, BITGET_BASE_URL, BITGET_DEMO   (optional)

Use --json on any command for the raw API response.

Docs Link: https://www.bitget.com/api-doc/uta/intro`,
	PersistentPreRunE: checkCredentials,
	SilenceUsage:      true,
	SilenceErrors:     true,
}

func init() {
	initCommandConfig()
}

func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func initCommandConfig() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.PersistentFlags().BoolVar(&config.Config.OutputJSON, "json", false, "Output JSON instead of a table")
}

func checkCredentials(cmd *cobra.Command, args []string) error {
	if config.Config.APIKey == "" || config.Config.APISecret == "" || config.Config.Passphrase == "" {
		return errors.New("missing credentials: set BITGET_API_KEY, BITGET_API_SECRET and BITGET_PASSPHRASE")
	}
	return nil
}
