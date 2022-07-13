package cmd

import (
	"runtime/debug"
	"strings"

	"github.com/pteropackages/soar/app"
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
		log.ApplyFlags(cmd.Flags())
		force, _ := cmd.Flags().GetBool("force")
		dir, _ := cmd.Flags().GetString("dir")

		path, err := config.Create(dir, force)
		if err != nil {
			log.Error("failed to initialize config:").WithError(err)
			return
		}

		log.Line(path)
	},
}

var configCmd = &cobra.Command{
	Use:     "config",
	Example: "soar config init --force",
	Run: func(cmd *cobra.Command, _ []string) {
		log.ApplyFlags(cmd.Flags())
		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			log.Error("failed to get config:").WithError(err).WithCmd("soar config --help")
			return
		}

		log.Line(cfg.Format())
	},
}

func init() {
	initConfigCmd.Flags().String("dir", "", "the directory to create the config in")
	initConfigCmd.Flags().BoolP("force", "f", false, "force overwrite the config")
	initConfigCmd.Flags().Bool("no-color", false, "disable ansi color codes")

	configCmd.AddCommand(initConfigCmd)
	configCmd.Flags().BoolP("local", "l", false, "get the local config")
	configCmd.Flags().Bool("no-color", false, "disable ansi color codes")

	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(app.GroupCommands())
}

func Execute() {
	defer func() {
		if state := recover(); state != nil {
			stack := string(debug.Stack())
			log.SetLevel(2).Error("an unknown error occurred:")
			for _, line := range strings.Split(stack, "\n") {
				log.Line(line)
			}
		}
	}()

	rootCmd.Execute()
}
