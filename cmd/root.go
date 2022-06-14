package cmd

import (
	"fmt"

	"github.com/pteropackages/soar/config"
	"github.com/spf13/cobra"
)

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
			fmt.Println(err.Error())
			return
		}

		fmt.Println(path)
	},
}

var configCmd = &cobra.Command{
	Use:     "config",
	Example: "soar config init --force",
	Run: func(cmd *cobra.Command, _ []string) {
		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(cfg.Format())
	},
}

func init() {
	initConfigCmd.Flags().String("dir", "", "the directory to create the config in")
	initConfigCmd.Flags().BoolP("force", "f", false, "force overwrite the config")
	configCmd.AddCommand(initConfigCmd)
	configCmd.Flags().BoolP("local", "l", false, "get the local config")

	rootCmd.AddCommand(configCmd)
}

func Execute() {
	rootCmd.Execute()
}
