package app

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/http"
	"github.com/pteropackages/soar/util"
	"github.com/spf13/cobra"
)

var getNodesCmd = &cobra.Command{
	Use:   "nodes:get [--id id] [--name name] [--uuid id] [--fqdn name] [--token token]",
	Short: "gets panel nodes",
	Long:  getNodesHelp,
	Run: func(cmd *cobra.Command, _ []string) {
		log.ApplyFlags(cmd.Flags())

		global, _ := cmd.Flags().GetBool("global")
		cfg, err := config.Get(global)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		single, query, err := parseNodeQuery(cmd)
		if err != nil {
			log.Error("command error:").WithError(err)
			return
		}

		ctx := http.New(cfg, &cfg.Application, log)
		req := ctx.Request("GET", "/api/application/nodes"+query, nil)
		res, err := ctx.ExecuteWithFlags(req, cmd.Flags())
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

func parseNodeQuery(cmd *cobra.Command) (bool, string, error) {
	var query strings.Builder
	var params []string
	single := false
	flags := cmd.Flags()

	if id, _ := flags.GetInt("id"); id != 0 {
		single = true
		query.WriteString(fmt.Sprintf("/%d", id))
	}

	if value, _ := flags.GetString("name"); value != "" {
		params = append(params, "filter[name]="+value)
	}

	if value, _ := flags.GetString("uuid"); value != "" {
		params = append(params, "filter[uuid]="+value)
	}

	if value, _ := flags.GetString("fqdn"); value != "" {
		params = append(params, "filter[fqdn]="+value)
	}

	if value, _ := flags.GetString("token"); value != "" {
		params = append(params, "filter[daemon_token_id]="+value)
	}

	if len(params) > 0 {
		query.WriteString("?" + params[0])

		if len(params) > 1 {
			for _, p := range params {
				query.WriteString("&" + p)
			}
		}
	}

	return single, query.String(), nil
}

var getNodeConfigCmd = &cobra.Command{
	Use:   "nodes:config id",
	Short: "gets a node config",
	Long:  "Gets the configuration for a specified node.",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"id"}); err != nil {
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

		ctx := http.New(cfg, &cfg.Application, log)
		req := ctx.Request("GET", "/api/application/nodes/"+args[0]+"/configuration", nil)
		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}

		var buf []byte

		if cfg.Http.ParseIndent {
			var model interface{}
			if err = json.Unmarshal(res, &model); err != nil {
				log.WithError(err)
				return
			}

			buf, err = json.MarshalIndent(model, "", "  ")
			if err != nil {
				log.WithError(err)
				return
			}
		}

		log.LineB(buf)
	},
}

var getNodeAllocationsCmd = &cobra.Command{
	Use:   "nodes:alloc:get id",
	Short: "gets node allocations",
	Long:  "Gets the allocations for a specified node.",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"id"}); err != nil {
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

		ctx := http.New(cfg, &cfg.Application, log)
		req := ctx.Request("GET", "/api/application/nodes/"+args[0]+"/allocations", nil)
		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}

		log.LineB(res)
	},
}
