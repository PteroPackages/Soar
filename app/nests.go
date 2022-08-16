package app

import (
	"fmt"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/http"
	"github.com/pteropackages/soar/util"
	"github.com/spf13/cobra"
)

var getNestsCmd = &cobra.Command{
	Use:   "nests:get",
	Short: "gets panel nests",
	Long:  getNestsHelp,
	Run: func(cmd *cobra.Command, _ []string) {
		log.ApplyFlags(cmd.Flags())

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		path := "/api/application/nests"
		id, _ := cmd.Flags().GetInt("id")
		if id != 0 {
			path += fmt.Sprintf("/%d", id)
		}

		ctx := http.New(cfg, &cfg.Application, log)
		req := ctx.Request("GET", path, nil)
		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}

		var buf []byte

		if id == 0 {
			buf, err = http.HandleDataResponse(res, cfg)
		} else {
			buf, err = http.HandleItemResponse(res, cfg)
		}
		if err != nil {
			log.WithError(err)
			return
		}

		log.LineB(buf)
	},
}

var getNestEggsCmd = &cobra.Command{
	Use:   "nests:eggs:get <id>",
	Short: "gets panel eggs for a nest",
	Long:  getNestEggsHelp,
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"nest id"}); err != nil {
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

		path := fmt.Sprintf("/api/application/nests/%s/eggs", args[0])
		id, _ := cmd.Flags().GetInt("id")
		if id != 0 {
			path += fmt.Sprintf("/%d", id)
		}

		ctx := http.New(cfg, &cfg.Application, log)
		req := ctx.Request("GET", path, nil)
		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}

		var buf []byte

		if id == 0 {
			buf, err = http.HandleDataResponse(res, cfg)
		} else {
			buf, err = http.HandleItemResponse(res, cfg)
		}
		if err != nil {
			log.WithError(err)
			return
		}

		log.LineB(buf)
	},
}
