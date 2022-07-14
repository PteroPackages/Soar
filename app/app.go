package app

import (
	"github.com/pteropackages/soar/util"
	"github.com/spf13/cobra"
)

func GroupCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "app",
		Example: "soar app <group> <subcommands> [arguments]",
	}

	util.ApplyDefaultFlags(getUsersCmd)

	getUsersCmd.Flags().Int("id", 0, "the id of the user")
	getUsersCmd.Flags().String("external", "", "the external id of the user")
	getUsersCmd.Flags().String("username", "", "the username of the user")
	getUsersCmd.Flags().String("email", "", "the email of the user")
	getUsersCmd.Flags().String("uuid", "", "the uuid of the user")

	cmd.AddCommand(getUsersCmd)

	return cmd
}
