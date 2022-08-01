package app

import (
	"encoding/json"
	"fmt"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/http"
	"github.com/spf13/cobra"
)

type location struct {
	ID        int    `json:"id"`
	Short     string `json:"short"`
	Long      string `json:"long"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type locAttrModel struct {
	O string    `json:"object"`
	A *location `json:"attributes"`
}

type locDataModel struct {
	O string         `json:"object"`
	D []locAttrModel `json:"data"`
}

var getLocationsCmd = &cobra.Command{
	Use: "locations:get",
	Run: func(cmd *cobra.Command, _ []string) {
		log.ApplyFlags(cmd.Flags())

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
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
		if res == nil {
			return
		}

		if id == 0 {
			var model locDataModel
			if err = json.Unmarshal(res, &model); err != nil {
				log.Error("failed to parse json:").WithError(err)
				return
			}

			var buf []byte
			if cfg.Http.ParseBody {
				inner := make([]*location, 0, len(model.D))
				for _, m := range model.D {
					inner = append(inner, m.A)
				}
				buf, err = json.MarshalIndent(inner, "", "  ")
			} else {
				buf, err = json.MarshalIndent(model.D, "", "  ")
			}
			if err != nil {
				log.Error("failed to parse response:").WithError(err)
				return
			}

			log.LineB(buf)
			return
		}

		var model locAttrModel
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
