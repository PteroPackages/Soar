package app

import (
	"github.com/spf13/cobra"
)

func GroupCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "app",
		Example: "soar app <group> <subcommands> [arguments]",
	}

	getUsersCmd.Flags().BoolP("local", "l", false, "use the local config")

	cmd.AddCommand(getUsersCmd)

	return cmd
}
