package cmd

import (
	"errors"

	"github.com/UnipayFI/bitget-cli/config"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "bitget-cli",
	Short: "Bitget API CLI for the Unified (UTA) and classic accounts",
	Long: `bitget-cli is a command-line client for Bitget's private REST APIs, covering
both of Bitget's account systems:

  UTA      Unified Trading Account (v3 /api/v3/*): run "bitget-cli UTA ..."
  spot     Classic spot account (v2 /api/v2/spot/*)
  futures  Classic futures account (v2 /api/v2/mix/*)

The top-level spot/futures commands target the classic account; the unified
account lives under the UTA subcommand.

Credentials are read from the environment:
  BITGET_API_KEY, BITGET_API_SECRET, BITGET_PASSPHRASE   (required)
  BITGET_LOCALE, BITGET_BASE_URL, BITGET_DEMO   (optional)
  HTTPS_PROXY / ALL_PROXY / HTTP_PROXY   (optional, standard proxy vars)

Use --json on any command for the raw API response.

Docs Link: https://www.bitget.com/api-doc/common/intro`,
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
