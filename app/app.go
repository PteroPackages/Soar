package app

import (
	"github.com/spf13/cobra"
)

func GroupCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "app",
		Example: "soar app <group> <subcommands> [arguments]",
	}

	getUsersCmd.Flags().Bool("debug", false, "print debug logs")
	getUsersCmd.Flags().BoolP("quiet", "q", false, "only output necessary logs")
	getUsersCmd.Flags().BoolP("local", "l", false, "use the local config")
	getUsersCmd.Flags().Int("id", 0, "the id of the user")
	getUsersCmd.Flags().String("external", "", "the external id of the user")

	cmd.AddCommand(getUsersCmd)

	return cmd
}
