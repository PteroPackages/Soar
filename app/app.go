package app

import (
	"github.com/pteropackages/soar/logger"
	"github.com/pteropackages/soar/util"
	"github.com/spf13/cobra"
)

var log = logger.New()

type FeatureLimits struct {
	Allocations int `json:"allocations"`
	Backups     int `json:"backups"`
	Databases   int `json:"databases"`
}

type Limits struct {
	Memory      int    `json:"memory"`
	Disk        int    `json:"disk"`
	Swap        int    `json:"swap"`
	IO          int    `json:"io"`
	CPU         int    `json:"cpu"`
	Threads     string `json:"threads"`
	OOMDisabled bool   `json:"oom_disabled"`
}

func GroupCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "app",
		Example: "soar app <group> <subcommands> [arguments]",
	}

	util.ApplyDefaultFlags(getUsersCmd)
	util.ApplyDefaultFlags(createUserCmd)
	util.ApplyDefaultFlags(deleteUserCmd)
	util.ApplyDefaultFlags(getServersCmd)
	util.ApplyDefaultFlags(suspendServerCmd)
	util.ApplyDefaultFlags(unsuspendServerCmd)
	util.ApplyDefaultFlags(reinstallServerCmd)

	getUsersCmd.Flags().Int("id", 0, "the id of the user")
	getUsersCmd.Flags().String("external", "", "the external id of the user")
	getUsersCmd.Flags().String("username", "", "the username of the user")
	getUsersCmd.Flags().String("email", "", "the email of the user")
	getUsersCmd.Flags().String("uuid", "", "the uuid of the user")

	createUserCmd.Flags().String("src", "", "the json file to read from")

	getServersCmd.Flags().Int("id", 0, "the id of the server")
	getServersCmd.Flags().String("external", "", "the external id of the server")

	cmd.AddCommand(getUsersCmd)
	cmd.AddCommand(createUserCmd)
	cmd.AddCommand(deleteUserCmd)
	cmd.AddCommand(getServersCmd)
	cmd.AddCommand(suspendServerCmd)
	cmd.AddCommand(unsuspendServerCmd)
	cmd.AddCommand(reinstallServerCmd)

	return cmd
}
