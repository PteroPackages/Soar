package app

import (
	"github.com/pteropackages/soar/logger"
	"github.com/pteropackages/soar/util"
	"github.com/spf13/cobra"
)

var log = logger.New()

func GroupCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app subcommand [options]",
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
	util.ApplyDefaultFlags(getNodeAllocationsCmd)
	util.ApplyDefaultFlags(createAllocationsCmd)
	util.ApplyDefaultFlags(deleteAllocationCmd)
	util.ApplyDefaultFlags(getLocationsCmd)
	util.ApplyDefaultFlags(createLocationCmd)
	util.ApplyDefaultFlags(deleteLocationCmd)
	util.ApplyDefaultFlags(getNestsCmd)
	util.ApplyDefaultFlags(getNestEggsCmd)

	util.ApplyFilterFlags(getUsersCmd)
	util.ApplyFilterFlags(getServersCmd)
	util.ApplyFilterFlags(getNodesCmd)
	util.ApplyFilterFlags(getLocationsCmd)
	util.ApplyFilterFlags(getNestsCmd)
	util.ApplyFilterFlags(getNestEggsCmd)

	getUsersCmd.Flags().Int("id", 0, "the id of the user")
	getUsersCmd.Flags().String("external", "", "the external id of the user")
	getUsersCmd.Flags().String("username", "", "filter by user username")
	getUsersCmd.Flags().String("email", "", "filter by user email")
	getUsersCmd.Flags().String("uuid", "", "filter by user uuid")
	createUserCmd.Flags().String("data", "", "a set of key-value pairs for the request")
	createUserCmd.Flags().String("data-file", "", "a file path to the json data")
	createUserCmd.Flags().String("data-json", "", "the json data for the request")
	getServersCmd.Flags().Int("id", 0, "the id of the server")
	getServersCmd.Flags().String("external", "", "the external id of the server")
	getServersCmd.Flags().String("name", "", "filter by server name")
	getServersCmd.Flags().String("desc", "", "filter by server description")
	getServersCmd.Flags().String("uuid", "", "filter by server uuid")
	getServersCmd.Flags().String("image", "", "filter by server docker image")
	deleteServerCmd.Flags().Bool("force", false, "force delete the server")
	getNodesCmd.Flags().Int("id", 0, "the id of the node")
	getNodesCmd.Flags().String("name", "", "filter by the node name")
	getNodesCmd.Flags().String("uuid", "", "filter by the node uuid")
	getNodesCmd.Flags().String("fqdn", "", "filter by the node fqdn")
	getNodesCmd.Flags().String("token", "", "filter by node token id")
	createAllocationsCmd.Flags().String("data", "", "a set of key-value pairs for the request")
	createAllocationsCmd.Flags().String("data-file", "", "a file path to the json data")
	createAllocationsCmd.Flags().String("data-json", "", "the json data for the request")
	getLocationsCmd.Flags().Int("id", 0, "the id of the location")
	createLocationCmd.Flags().String("data", "", "a set of key-value pairs for the request")
	createLocationCmd.Flags().String("data-file", "", "a file path to the json data")
	createLocationCmd.Flags().String("data-json", "", "the json data for the request")
	getNestsCmd.Flags().Int("id", 0, "the id of the nest")
	getNestEggsCmd.Flags().Int("id", 0, "the id of the egg")

	cmd.AddCommand(getUsersCmd)
	cmd.AddCommand(createUserCmd)
	cmd.AddCommand(deleteUserCmd)
	cmd.AddCommand(getServersCmd)
	cmd.AddCommand(suspendServerCmd)
	cmd.AddCommand(unsuspendServerCmd)
	cmd.AddCommand(reinstallServerCmd)
	cmd.AddCommand(deleteServerCmd)
	cmd.AddCommand(getNodesCmd)
	cmd.AddCommand(getNodeConfigCmd)
	cmd.AddCommand(getNodeAllocationsCmd)
	cmd.AddCommand(createAllocationsCmd)
	cmd.AddCommand(deleteAllocationCmd)
	cmd.AddCommand(getLocationsCmd)
	cmd.AddCommand(createLocationCmd)
	cmd.AddCommand(deleteLocationCmd)
	cmd.AddCommand(getNestsCmd)
	cmd.AddCommand(getNestEggsCmd)

	return cmd
}
