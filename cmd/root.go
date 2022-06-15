package cmd

import (
	"runtime/debug"
	"strings"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/logger"
	"github.com/spf13/cobra"
)

var log = logger.New()

var rootCmd = &cobra.Command{
	Use:     "soar",
	Version: Version,
}

var initConfigCmd = &cobra.Command{
	Use:     "init",
	Example: "soar config init --dir=/",
	Run: func(cmd *cobra.Command, _ []string) {
		force, _ := cmd.Flags().GetBool("force")
		dir, _ := cmd.Flags().GetString("dir")

		path, err := config.Create(dir, force)
		if err != nil {
			log.Error("failed to initialize config:").WithError(err).Log()
			return
		}

		log.Line(path).Log()
	},
}

var configCmd = &cobra.Command{
	Use:     "config",
	Example: "soar config init --force",
	Run: func(cmd *cobra.Command, _ []string) {
		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			log.Error("failed to get config:").WithError(err).WithTip("soar config --help").Log()
			return
		}

		log.Line(cfg.Format()).Log()
	},
}

func init() {
	initConfigCmd.Flags().String("dir", "", "the directory to create the config in")
	initConfigCmd.Flags().BoolP("force", "f", false, "force overwrite the config")
	initConfigCmd.Flags().BoolVar(&log.NoColor, "no-color", false, "disable ansi color codes")
	initConfigCmd.Flags().BoolVar(&log.Persist, "save", false, "save the command and request logs")
	configCmd.AddCommand(initConfigCmd)
	configCmd.Flags().BoolP("local", "l", false, "get the local config")
	configCmd.Flags().BoolVar(&log.NoColor, "no-color", false, "disable ansi color codes")
	configCmd.Flags().BoolVar(&log.Persist, "save", false, "save the command and request logs")

	rootCmd.AddCommand(configCmd)
}

func Execute() {
	defer func() {
		if state := recover(); state != nil {
			stack := string(debug.Stack())
			entry := log.Error("an unknown error occurred:")
			for _, line := range strings.Split(stack, "\n") {
				entry.Line(line)
			}
		}
	}()

	rootCmd.Execute()
}
