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

		single, query, err := parseServerQuery(cmd)
		if err != nil {
			log.Error("command error:").WithError(err)
			return
		}

		ctx := http.New(cfg, &cfg.Application, log)
		if single {
			req := ctx.Request("GET", "/api/application/servers"+query, nil)
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
			return
		}

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
			inner := make([]*server, 0, len(model.D))
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
	Use: "servers:suspend",
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
	Use: "servers:unsuspend",
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
	Use: "servers:reinstall",
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
