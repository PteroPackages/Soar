package client

import (
	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/http"
	"github.com/spf13/cobra"
)

var getAccountCmd = &cobra.Command{
	Use:   "account:get",
	Short: "gets account information",
	Run: func(cmd *cobra.Command, _ []string) {
		log.ApplyFlags(cmd.Flags())

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("GET", "/api/client/account", nil)
		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}

		buf, err := http.HandleItemResponse(res, cfg)
		if err != nil {
			log.WithError(err)
			return
		}

		log.LineB(buf)
	},
}

var getPermissionsCmd = &cobra.Command{
	Use:     "account:perms",
	Aliases: []string{"account:p"},
	Short:   "gets system permissions",
	Run: func(cmd *cobra.Command, _ []string) {
		log.ApplyFlags(cmd.Flags())

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("GET", "/api/client/permissions", nil)
		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}
		if res == nil {
			return
		}

		buf, err := http.HandleItemResponse(res, cfg)
		if err != nil {
			log.WithError(err)
			return
		}

		log.LineB(buf)
	},
}

var getServersCmd = &cobra.Command{
	Use:   "servers:get",
	Short: "gets account servers",
	Run: func(cmd *cobra.Command, _ []string) {
		log.ApplyFlags(cmd.Flags())

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("GET", "/api/client", nil)
		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}

		buf, err := http.HandleDataResponse(res, cfg)
		if err != nil {
			log.WithError(err)
			return
		}

		log.LineB(buf)
	},
}
