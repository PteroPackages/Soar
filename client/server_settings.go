package client

import (
	"bytes"
	"encoding/json"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/http"
	"github.com/pteropackages/soar/util"
	"github.com/spf13/cobra"
)

var renameServerCmd = &cobra.Command{
	Use:   "settings:rename identifier name",
	Short: "renames a server",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier", "name"}); err != nil {
			log.WithError(err)
			return
		}

		global, _ := cmd.Flags().GetBool("global")
		cfg, err := config.Get(global)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		data, _ := json.Marshal(map[string]string{"name": args[1]})
		body := bytes.Buffer{}
		body.Write(data)

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("POST", "/api/client/servers/"+args[0]+"/settings/rename", &body)
		if _, err := ctx.Execute(req); err != nil {
			log.WithError(err)
		}
	},
}

var reinstallServerCmd = &cobra.Command{
	Use:   "settings:reinstall identifier",
	Short: "reinstalls a server",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier"}); err != nil {
			log.WithError(err)
			return
		}

		global, _ := cmd.Flags().GetBool("global")
		cfg, err := config.Get(global)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("POST", "/api/client/servers/"+args[0]+"/settings/reinstall", nil)
		if _, err := ctx.Execute(req); err != nil {
			log.WithError(err)
		}
	},
}

var setDockerImageCmd = &cobra.Command{
	Use:   "settings:image identifier image",
	Short: "sets the docker image for a server",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier", "image"}); err != nil {
			log.WithError(err)
			return
		}

		global, _ := cmd.Flags().GetBool("global")
		cfg, err := config.Get(global)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		data, _ := json.Marshal(map[string]string{"docker_image": args[1]})
		body := bytes.Buffer{}
		body.Write(data)

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("PUT", "/api/client/servers/"+args[0]+"/settings/docker-image", &body)
		if _, err := ctx.Execute(req); err != nil {
			log.WithError(err)
		}
	},
}
