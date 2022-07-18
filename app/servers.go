package app

import (
	"encoding/json"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/http"
	"github.com/spf13/cobra"
)

type server struct {
	ID            int           `json:"id"`
	ExternalID    string        `json:"external_id"`
	UUID          string        `json:"uuid"`
	Identifier    string        `json:"identifer"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Status        string        `json:"status"`
	Suspended     bool          `json:"suspended"`
	Limits        Limits        `json:"limits"`
	FeatureLimits FeatureLimits `json:"feature_limits"`
	User          int           `json:"user"`
	Node          int           `json:"node"`
	Allocation    int           `json:"allocation"`
	Nest          int           `json:"nest"`
	Egg           int           `json:"egg"`
	Container     struct {
		StartupCommand string                 `json:"startup_command"`
		Image          string                 `json:"image"`
		Installed      int                    `json:"installed"`
		Environment    map[string]interface{} `json:"environment"`
	} `json:"container"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type serverAttrModel struct {
	O string  `json:"object"`
	A *server `json:"attributes"`
}

type serverDataModel struct {
	O string            `json:"object"`
	D []serverAttrModel `json:"data"`
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

		var model serverDataModel
		if err = json.Unmarshal(res, &model); err != nil {
			log.Error("failed to parse json:").WithError(err)
			return
		}

		var buf []byte
		if cfg.Http.ParseBody {
			buf, err = json.MarshalIndent(model.D, "", "  ")
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
