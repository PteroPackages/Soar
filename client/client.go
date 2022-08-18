package client

import (
	"github.com/pteropackages/soar/logger"
	"github.com/pteropackages/soar/util"
	"github.com/spf13/cobra"
)

var log = logger.New()

func GroupCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "client subcommand [options]",
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
	util.ApplyDefaultFlags(getServerWSCmd)
	util.ApplyDefaultFlags(getServerResourcesCmd)
	util.ApplyDefaultFlags(getServerActivityCmd)
	util.ApplyDefaultFlags(sendServerCommandCmd)
	util.ApplyDefaultFlags(setServerPowerStateCmd)
	util.ApplyDefaultFlags(getDatabasesCmd)
	util.ApplyDefaultFlags(listFilesCmd)
	util.ApplyDefaultFlags(getFileInfoCmd)
	util.ApplyDefaultFlags(getFileContentsCmd)
	util.ApplyDefaultFlags(downloadFileCmd)
	util.ApplyDefaultFlags(renameFileCmd)
	util.ApplyDefaultFlags(copyFileCmd)
	util.ApplyDefaultFlags(writeFileCmd)
	util.ApplyDefaultFlags(createFileCmd)

	listFilesCmd.Flags().BoolP("dir", "d", false, "only list directories")
	listFilesCmd.Flags().BoolP("file", "f", false, "only list files")
	downloadFileCmd.Flags().String("dest", "", "the path to save the file at")
	downloadFileCmd.Flags().BoolP("url-only", "U", false, "only return the url")
	renameFileCmd.Flags().String("root", "/", "the root directory of the file")

	cmd.AddCommand(getAccountCmd)
	cmd.AddCommand(getPermissionsCmd)
	cmd.AddCommand(getServersCmd)
	cmd.AddCommand(getTwoFactorCodeCmd)
	cmd.AddCommand(enableTwoFactorCmd)
	cmd.AddCommand(disableTwoFactorCmd)
	cmd.AddCommand(getAccountActivityCmd)
	cmd.AddCommand(getAPIKeysCmd)
	cmd.AddCommand(deleteAPIKeyCmd)
	cmd.AddCommand(getServerWSCmd)
	cmd.AddCommand(getServerResourcesCmd)
	cmd.AddCommand(getServerActivityCmd)
	cmd.AddCommand(sendServerCommandCmd)
	cmd.AddCommand(setServerPowerStateCmd)
	cmd.AddCommand(getDatabasesCmd)
	cmd.AddCommand(listFilesCmd)
	cmd.AddCommand(getFileInfoCmd)
	cmd.AddCommand(getFileContentsCmd)
	cmd.AddCommand(downloadFileCmd)
	cmd.AddCommand(renameFileCmd)
	cmd.AddCommand(copyFileCmd)
	cmd.AddCommand(writeFileCmd)
	cmd.AddCommand(createFileCmd)

	return cmd
}
