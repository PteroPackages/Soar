package client

import (
	"bytes"
	"encoding/json"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/http"
	"github.com/pteropackages/soar/util"
	"github.com/spf13/cobra"
)

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

var getServerWSCmd = &cobra.Command{
	Use:     "servers:websocket identifier",
	Aliases: []string{"servers:ws"},
	Short:   "gets the server websocket data",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier"}); err != nil {
			log.WithError(err)
			return
		}

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("GET", "/api/client/servers/"+args[0]+"/websocket", nil)
		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}

		var model struct {
			Data struct {
				Token  string `json:"token"`
				Socket string `json:"socket"`
			} `json:"data"`
		}
		if err = json.Unmarshal(res, &model); err != nil {
			log.WithError(err)
			return
		}

		var buf []byte

		if cfg.Http.ParseBody {
			if cfg.Http.ParseIndent {
				buf, err = json.MarshalIndent(model.Data, "", "  ")
			} else {
				buf, err = json.Marshal(model.Data)
			}
		} else {
			if cfg.Http.ParseIndent {
				buf, err = json.MarshalIndent(model, "", "  ")
			} else {
				buf, err = json.Marshal(model)
			}
		}
		if err != nil {
			log.WithError(err)
			return
		}

		log.LineB(buf)
	},
}

var getServerResourcesCmd = &cobra.Command{
	Use:     "servers:resources identifier",
	Aliases: []string{"servers:usage"},
	Short:   "gets server resource usage",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier"}); err != nil {
			log.WithError(err)
			return
		}

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("GET", "/api/client/servers/"+args[0]+"/resources", nil)
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

var getServerActivityCmd = &cobra.Command{
	Use:   "servers:activity identifier",
	Short: "gets the server activity logs",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier"}); err != nil {
			log.WithError(err)
			return
		}

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("GET", "/api/client/servers/"+args[0]+"/activity", nil)
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

var sendServerCommandCmd = &cobra.Command{
	Use:     "servers:command identifier command",
	Aliases: []string{"servers:cmd"},
	Short:   "sends a command to the server console",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier", "command"}); err != nil {
			log.WithError(err)
			return
		}

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		data, _ := json.Marshal(map[string]string{"command": args[1]})
		body := bytes.Buffer{}
		body.Write(data)

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("POST", "/api/client/servers/"+args[0]+"/command", &body)
		if _, err = ctx.Execute(req); err != nil {
			log.WithError(err)
		}
	},
}

var setServerPowerStateCmd = &cobra.Command{
	Use:     "servers:power identifier state",
	Aliases: []string{"servers:state", "servers:status", "servers:toggle"},
	Short:   "sets the server power state",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier", "state"}); err != nil {
			log.WithError(err)
			return
		}

		switch args[1] {
		case "start":
		case "stop":
		case "restart":
		case "kill":
		default:
			log.Error("invalid power state '%s'", args[1])
			return
		}

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		data, _ := json.Marshal(map[string]string{"signal": args[1]})
		body := bytes.Buffer{}
		body.Write(data)

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("POST", "/api/client/servers/"+args[0]+"/power", &body)
		if _, err = ctx.Execute(req); err != nil {
			log.WithError(err)
		}
	},
}
