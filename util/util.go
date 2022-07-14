package util

import "github.com/spf13/cobra"

func ApplyDefaultFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("debug", false, "print debug logs")
	cmd.Flags().Bool("no-color", false, "disable ansi color codes")
	cmd.Flags().BoolP("local", "l", false, "use the local config")
	cmd.Flags().BoolP("quiet", "q", false, "only output necessary logs")
}
