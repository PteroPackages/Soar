package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/pteropackages/soar/app"
	"github.com/pteropackages/soar/client"
	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/logger"
	"github.com/pteropackages/soar/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var log = logger.New()

var rootCmd = &cobra.Command{
	Use:     "soar subcommand [options] arguments",
	Short:   "Commands for interacting with Pterodactyl via the API",
	Version: Version,
}

var initConfigCmd = &cobra.Command{
	Use: "init [--dir path] [-f | --force]",
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

var copyConfigCmd = &cobra.Command{
	Use:   "copy scope",
	Short: "copies a global or local config to the corresponding destination",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())

		if err := util.RequireArgs(args, []string{"scope"}); err != nil {
			log.WithError(err)
			return
		}

		var global bool

		switch args[0] {
		case "global":
			global = true
		case "local":
			global = false
		default:
			log.Error("invalid config scope; must be global or local")
			return
		}

		var path string

		if global {
			root, err := os.UserConfigDir()
			if err != nil {
				log.Error("failed to get global config:").WithError(err)
			}

			if _, err = os.Stat(root); err != nil {
				if os.IsNotExist(err) {
					log.Error(fmt.Sprintf("user config directory not found (path: %s)", err))
					return
				}

				log.WithError(err)
				return
			}
			path = filepath.Join(root, ".soar", "config.yml")
		} else {
			root, _ := os.Getwd()
			path = filepath.Join(root, ".soar.yml")
		}

		cfg, err := config.GetStatic(global)
		if err != nil {
			log.WithError(err)
			return
		}

		file, err := os.Create(path)
		if err != nil {
			log.WithError(err)
			return
		}
		defer file.Close()

		buf, _ := yaml.Marshal(cfg)
		file.Write(buf)
	},
}

var configCmd = &cobra.Command{
	Use:   "config [init] [-g | --global] [-v | --validate]",
	Short: "manages the soar config",
	Long:  "Manages the soar config for HTTP and logging",
	Run: func(cmd *cobra.Command, _ []string) {
		log.ApplyFlags(cmd.Flags())

		var cfg *config.Config
		var err error

		validate, _ := cmd.Flags().GetBool("validate")
		global, _ := cmd.Flags().GetBool("global")

		if validate {
			cfg, err = config.Get(global)
		} else {
			cfg, err = config.GetStatic(global)
		}
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
	configCmd.AddCommand(copyConfigCmd)
	configCmd.Flags().BoolP("global", "g", false, "use the global config")
	configCmd.Flags().Bool("no-color", false, "disable ansi color codes")
	configCmd.Flags().BoolP("validate", "v", false, "validate the config")

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
