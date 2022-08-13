package app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/http"
	"github.com/spf13/cobra"
)

var getServersCmd = &cobra.Command{
	Use:   "servers:get",
	Short: "gets panel servers",
	Long:  getServersHelp,
	Run: func(cmd *cobra.Command, _ []string) {
		log.ApplyFlags(cmd.Flags())

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		single, query, err := parseServerQuery(cmd)
		if err != nil {
			log.Error("command error:").WithError(err)
			return
		}

		path := "/api/application/servers"
		if single {
			path += query
		}

		ctx := http.New(cfg, &cfg.Application, log)
		req := ctx.Request("GET", path, nil)
		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}

		var buf []byte

		if single {
			buf, err = http.HandleItemResponse(res, cfg)
		} else {
			buf, err = http.HandleDataResponse(res, cfg)
		}
		if err != nil {
			log.WithError(err)
			return
		}

		log.LineB(buf)
	},
}

func parseServerQuery(cmd *cobra.Command) (bool, string, error) {
	var query strings.Builder
	single := false
	flags := cmd.Flags()

	if id, _ := flags.GetInt("id"); id != 0 {
		single = true
		query.WriteString(fmt.Sprintf("/%d", id))
	}

	if ext, _ := flags.GetString("external"); ext != "" {
		if query.Len() != 0 {
			return false, "", errors.New("id an external flags specified; pick one")
		}

		query.WriteString("/external/" + ext)
	}

	return single, query.String(), nil
}

var suspendServerCmd = &cobra.Command{
	Use:   "servers:suspend <id>",
	Short: "suspends a server",
	Long:  "Suspends a server on the panel by its ID.",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		if len(args) == 0 {
			log.Error("no server id specified")
			return
		} else if len(args) > 1 {
			log.Error("more than one server id argument specified").WithCmd("soar app servers:suspend --help")
			return
		}

		ctx := http.New(cfg, &cfg.Application, log)
		req := ctx.Request("POST", "/api/application/servers/"+args[0]+"/suspend", nil)
		if _, err = ctx.Execute(req); err != nil {
			log.WithError(err)
		}
	},
}

var unsuspendServerCmd = &cobra.Command{
	Use:   "servers:unsuspend <id>",
	Short: "unsuspends a server",
	Long:  "Unsuspends a server on the panel by its ID.",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		if len(args) == 0 {
			log.Error("no server id specified")
			return
		} else if len(args) > 1 {
			log.Error("more than one server id argument specified").WithCmd("soar app servers:unsuspend --help")
			return
		}

		ctx := http.New(cfg, &cfg.Application, log)
		req := ctx.Request("POST", "/api/application/servers/"+args[0]+"/unsuspend", nil)
		if _, err = ctx.Execute(req); err != nil {
			log.WithError(err)
		}
	},
}

var reinstallServerCmd = &cobra.Command{
	Use:   "servers:reinstall <id>",
	Short: "reinstalls a server",
	Long:  "Triggers the reinstall process for a server by its ID.",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		if len(args) == 0 {
			log.Error("no server id specified")
			return
		} else if len(args) > 1 {
			log.Error("more than one server id argument specified").WithCmd("soar app servers:reinstall --help")
			return
		}

		ctx := http.New(cfg, &cfg.Application, log)
		req := ctx.Request("POST", "/api/application/servers/"+args[0]+"/reinstall", nil)
		if _, err = ctx.Execute(req); err != nil {
			log.WithError(err)
		}
	},
}

var deleteServerCmd = &cobra.Command{
	Use:   "servers:delete <id>",
	Short: "deletes a server",
	Long:  "Deletes a server on the panel by its ID (supports the --force flag).",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		if len(args) == 0 {
			log.Error("no server id specified")
			return
		} else if len(args) > 1 {
			log.Error("more than one server id argument specified").WithCmd("soar app servers:delete --help")
			return
		}

		path := "/api/application/servers/" + args[0]
		if force, _ := cmd.Flags().GetBool("force"); force {
			path += "/force"
		}

		ctx := http.New(cfg, &cfg.Application, log)
		req := ctx.Request("DELETE", path, nil)
		if _, err = ctx.Execute(req); err != nil {
			log.WithError(err)
		}
	},
}
