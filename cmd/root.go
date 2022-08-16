package cmd

import (
	"runtime/debug"

	"github.com/pteropackages/soar/app"
	"github.com/pteropackages/soar/client"
	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/logger"
	"github.com/spf13/cobra"
)

var log = logger.New()

var rootCmd = &cobra.Command{
	Use:     "soar subcommand [options] arguments",
	Short:   "Commands for interacting with Pterodactyl via the API",
	Version: Version,
}

var initConfigCmd = &cobra.Command{
	Use: "init [--dir path] [--force]",
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
	Use:   "config [init]",
	Short: "manages the soar config",
	Long:  "Manages the soar config for HTTP and logging",
	Run: func(cmd *cobra.Command, _ []string) {
		log.ApplyFlags(cmd.Flags())

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
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
	rootCmd.AddCommand(client.GroupCommands())
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() {
	defer func() {
		if state := recover(); state != nil {
			stack := string(debug.Stack())
			log.SetLevel(2).Error("a fatal error occurred:").Line(stack)
		}
	}()

	rootCmd.Execute()
}
