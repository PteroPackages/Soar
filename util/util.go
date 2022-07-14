package util

import "github.com/spf13/cobra"

func ApplyDefaultFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("debug", false, "print debug logs")
	cmd.Flags().Bool("no-color", false, "disable ansi color codes")
	cmd.Flags().BoolP("local", "l", false, "use the local config")
	cmd.Flags().BoolP("quiet", "q", false, "only output necessary logs")

	cmd.Flags().BoolP("retry-ratelimit", "r", false, "retry request on ratelimit")
	cmd.Flags().BoolP("no-retry-ratelimit", "R", false, "opposite of retry-ratelimit")
	cmd.Flags().Int("max-body", 0, "the maximum body size to accept")
	cmd.Flags().BoolP("parse", "p", false, "parse the response body")
	cmd.Flags().BoolP("no-parse", "P", false, "opposite of parse")
}
