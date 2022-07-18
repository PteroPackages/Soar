package app

import (
	"encoding/json"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/http"
	"github.com/spf13/cobra"
)

type server struct {
	ID         int    `json:"id"`
	ExternalID string `json:"external_id"`
}

type serverAttrModel struct {
	O string  `json:"object"`
	A *server `json:"attributes"`
}

var getServersCmd = &cobra.Command{
	Use: "servers:get",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		ctx := http.New(cfg, &cfg.Application, log)
		req := ctx.Request("GET", "/api/application/servers", nil)
		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}
		if res == nil {
			return
		}

		var model serverAttrModel
		if err = json.Unmarshal(res, &model); err != nil {
			log.Error("failed to parse json:").WithError(err)
			return
		}

		var buf []byte
		if cfg.Http.ParseBody {
			buf, err = json.MarshalIndent(model.A, "", "  ")
		} else {
			buf, err = json.MarshalIndent(model, "", "  ")
		}
		if err != nil {
			log.Error("failed to parse response:").WithError(err)
			return
		}

		log.LineB(buf)
	},
}
