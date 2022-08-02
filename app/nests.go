package app

import (
	"encoding/json"
	"fmt"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/http"
	"github.com/spf13/cobra"
)

type nest struct {
	ID          int    `json:"id"`
	UUID        string `json:"uuid"`
	Author      string `json:"author"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type egg struct {
	ID           int               `json:"id"`
	UUID         string            `json:"uuid"`
	Name         string            `json:"name"`
	Author       string            `json:"author"`
	Description  string            `json:"description"`
	Nest         int               `json:"nest"`
	DockerImage  string            `json:"docker_image"`
	DockerImages map[string]string `json:"docker_images"`
	Config       struct {
		Files    interface{}            `json:"files"`
		Startup  map[string]string      `json:"startup"`
		Stop     string                 `json:"stop"`
		Logs     map[string]interface{} `json:"logs"`
		Denylist []string               `json:"file_denylist"`
		Extends  string                 `json:"extends"`
	} `json:"config"`
	Startup string `json:"startup"`
	Script  struct {
		Privileged bool   `json:"bool"`
		Install    string `json:"install"`
		Entry      string `json:"entry"`
		Container  string `json:"container"`
		Extends    string `json:"extends"`
	} `json:"script"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type nestAttrModel struct {
	O string `json:"object"`
	A *nest  `json:"attributes"`
}

type nestDataModel struct {
	O string          `json:"object"`
	D []nestAttrModel `json:"data"`
}

type eggAttrModel struct {
	O string `json:"object"`
	A *egg   `json:"attributes"`
}

type eggDataModel struct {
	O string         `json:"object"`
	D []eggAttrModel `json:"data"`
}

var getNestsCmd = &cobra.Command{
	Use: "nests:get",
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
		if res == nil {
			return
		}

		if id == 0 {
			var model nestDataModel
			if err = json.Unmarshal(res, &model); err != nil {
				log.Error("failed to parse json:").WithError(err)
				return
			}

			var buf []byte
			if cfg.Http.ParseBody {
				inner := make([]*nest, 0, len(model.D))
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

		var model nestAttrModel
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

var getNestEggsCmd = &cobra.Command{
	Use: "nests:eggs:get",
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
			log.Error("no nest id specified")
			return
		} else if len(args) > 1 {
			log.Error("more than one nest id argument specified").WithCmd("soar app nests:eggs:get --help")
			return
		}

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
		if res == nil {
			return
		}

		if id == 0 {
			var model eggDataModel
			if err = json.Unmarshal(res, &model); err != nil {
				log.Error("failed to parse json:").WithError(err)
				return
			}

			var buf []byte
			if cfg.Http.ParseBody {
				inner := make([]*egg, 0, len(model.D))
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

		var model eggAttrModel
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
