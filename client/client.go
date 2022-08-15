package client

import (
	"github.com/pteropackages/soar/logger"
	"github.com/pteropackages/soar/util"
	"github.com/spf13/cobra"
)

var log = logger.New()

func GroupCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "client <subcommand> [options]",
		Short: "client api management",
		Long:  "Commands for interacting with the CLient API.",
	}

	util.ApplyDefaultFlags(getAccountCmd)
	util.ApplyDefaultFlags(getPermissionsCmd)
	util.ApplyDefaultFlags(getServersCmd)
	util.ApplyDefaultFlags(getTwoFactorCodeCmd)
	util.ApplyDefaultFlags(enableTwoFactorCmd)
	util.ApplyDefaultFlags(disableTwoFactorCmd)
	util.ApplyDefaultFlags(getAccountActivityCmd)
	util.ApplyDefaultFlags(getAPIKeysCmd)
	util.ApplyDefaultFlags(deleteAPIKeyCmd)

	cmd.AddCommand(getAccountCmd)
	cmd.AddCommand(getPermissionsCmd)
	cmd.AddCommand(getServersCmd)
	cmd.AddCommand(getTwoFactorCodeCmd)
	cmd.AddCommand(enableTwoFactorCmd)
	cmd.AddCommand(disableTwoFactorCmd)
	cmd.AddCommand(getAccountActivityCmd)
	cmd.AddCommand(getAPIKeysCmd)
	cmd.AddCommand(deleteAPIKeyCmd)

	return cmd
}
