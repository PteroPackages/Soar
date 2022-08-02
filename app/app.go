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
		Use:   "app <subcommand> [options]",
		Short: "application api management",
		Long:  "Commands for interacting with the application API.",
	}

	util.ApplyDefaultFlags(getUsersCmd)
	util.ApplyDefaultFlags(createUserCmd)
	util.ApplyDefaultFlags(deleteUserCmd)
	util.ApplyDefaultFlags(getServersCmd)
	util.ApplyDefaultFlags(suspendServerCmd)
	util.ApplyDefaultFlags(unsuspendServerCmd)
	util.ApplyDefaultFlags(reinstallServerCmd)
	util.ApplyDefaultFlags(deleteServerCmd)
	util.ApplyDefaultFlags(getNodesCmd)
	util.ApplyDefaultFlags(getNodeConfigCmd)
	util.ApplyDefaultFlags(getLocationsCmd)
	util.ApplyDefaultFlags(getNestsCmd)
	util.ApplyDefaultFlags(getNestEggsCmd)

	getUsersCmd.Flags().Int("id", 0, "the id of the user")
	getUsersCmd.Flags().String("external", "", "the external id of the user")
	getUsersCmd.Flags().String("username", "", "the username of the user")
	getUsersCmd.Flags().String("email", "", "the email of the user")
	getUsersCmd.Flags().String("uuid", "", "the uuid of the user")
	createUserCmd.Flags().String("src", "", "the json file to read from")
	getServersCmd.Flags().Int("id", 0, "the id of the server")
	getServersCmd.Flags().String("external", "", "the external id of the server")
	deleteServerCmd.Flags().Bool("force", false, "force delete the server")
	getNodesCmd.Flags().Int("id", 0, "the id of the node")
	getLocationsCmd.Flags().Int("id", 0, "the id of the location")
	getNestsCmd.Flags().Int("id", 0, "the id of the nest")
	getNestEggsCmd.Flags().Int("id", 0, "the id of the egg")

	cmd.AddCommand(getUsersCmd)
	cmd.AddCommand(createUserCmd)
	cmd.AddCommand(deleteUserCmd)
	cmd.AddCommand(getServersCmd)
	cmd.AddCommand(suspendServerCmd)
	cmd.AddCommand(unsuspendServerCmd)
	cmd.AddCommand(reinstallServerCmd)
	cmd.AddCommand(getNodesCmd)
	cmd.AddCommand(getNodeConfigCmd)
	cmd.AddCommand(getLocationsCmd)
	cmd.AddCommand(getNestsCmd)
	cmd.AddCommand(getNestEggsCmd)

	return cmd
}
