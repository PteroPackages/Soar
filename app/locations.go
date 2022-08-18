package app

import (
	"fmt"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/http"
	"github.com/spf13/cobra"
)

var getLocationsCmd = &cobra.Command{
	Use:   "locations:get",
	Short: "gets panel node locations",
	Long:  getLocationsHelp,
	Run: func(cmd *cobra.Command, _ []string) {
		log.ApplyFlags(cmd.Flags())

		global, _ := cmd.Flags().GetBool("global")
		cfg, err := config.Get(global)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		path := "/api/application/locations"
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
