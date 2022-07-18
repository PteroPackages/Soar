package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/http"
	"github.com/spf13/cobra"
)

type node struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	LocationID         int    `json:"location_id"`
	Public             bool   `json:"public"`
	FQDN               string `json:"fqdn"`
	Scheme             string `json:"scheme"`
	BehindProxy        bool   `json:"behind_proxy"`
	Memory             int    `json:"memory"`
	MemoryOverallocate int    `json:"memory_overallocate"`
	Disk               int    `json:"disk"`
	DiskOverallocate   int    `json:"disk_overallocate"`
	DaemonBase         string `json:"daemon_base"`
	DaemonSFTP         int    `json:"daemon_sftp"`
	DaemonListen       int    `json:"daemon_listen"`
	MaintenanceMode    bool   `json:"maintenance_mode"`
	UploadSize         int    `json:"upload_size"`
}

type nodeAttrModel struct {
	O string `json:"object"`
	A *node  `json:"attributes"`
}

type nodeDataModel struct {
	O string          `json:"object"`
	D []nodeAttrModel `json:"data"`
}

var getNodesCmd = &cobra.Command{
	Use: "nodes:get",
	Run: func(cmd *cobra.Command, args []string) {
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

		ctx := http.New(cfg, &cfg.Application, log)
		if single {
			req := ctx.Request("GET", "/api/application/nodes"+query, nil)
			res, err := ctx.Execute(req)
			if err != nil {
				log.WithError(err)
				return
			}
			if res == nil {
				return
			}

			var model nodeAttrModel
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
			return
		}

		req := ctx.Request("GET", "/api/application/nodes", nil)
		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}
		if res == nil {
			return
		}

		var model nodeDataModel
		if err = json.Unmarshal(res, &model); err != nil {
			log.Error("failed to parse json:").WithError(err)
			return
		}

		var buf []byte
		if cfg.Http.ParseBody {
			inner := make([]*node, 0, len(model.D))
			for _, m := range model.D {
				inner = append(inner, m.A)
			}
			buf, err = json.MarshalIndent(inner, "", "  ")
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
