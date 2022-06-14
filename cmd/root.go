package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:     "soar",
	Example: "soar [options] <command> [arguments]",
	Version: Version,
}

func Execute() {
	rootCmd.Execute()
}
