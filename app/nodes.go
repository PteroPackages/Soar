package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/http"
	"github.com/pteropackages/soar/util"
	"github.com/spf13/cobra"
)

var getNodesCmd = &cobra.Command{
	Use:   "nodes:get",
	Short: "gets panel nodes",
	Long:  getNodesHelp,
	Run: func(cmd *cobra.Command, _ []string) {
		log.ApplyFlags(cmd.Flags())

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
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

		path := "/api/application/nodes"
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

func parseNodeQuery(cmd *cobra.Command) (bool, string, error) {
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

var getNodeConfigCmd = &cobra.Command{
	Use:   "nodes:config",
	Short: "gets a node config",
	Long:  "Gets the configuration for a specified node.",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"node id"}); err != nil {
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
