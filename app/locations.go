package app

import (
	"bytes"
	"fmt"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/http"
	"github.com/pteropackages/soar/input"
	"github.com/pteropackages/soar/util"
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
		res, err := ctx.ExecuteWithFlags(req, cmd.Flags())
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

var createLocationCmd = &cobra.Command{
	Use:   "locations:create --data[-file | -json] source",
	Short: "creates a location",
	Run: func(cmd *cobra.Command, _ []string) {
		log.ApplyFlags(cmd.Flags())

		global, _ := cmd.Flags().GetBool("global")
		cfg, err := config.Get(global)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		var payload []byte
		data, _ := cmd.Flags().GetString("data")
		file, _ := cmd.Flags().GetString("file")
		js, _ := cmd.Flags().GetString("json")

		switch {
		case data != "":
			m, err := input.Parse(data)
			if err != nil {
				log.WithError(err).Error("failed to parse data input")
				return
			}

			payload, err = input.Marshal(input.Definition{
				"short": input.StringNode,
				"long":  input.StringNode,
			}, m)
			if err != nil {
				log.WithError(err).Error("failed to parse data input")
				return
			}
		case file != "":
			v, err := util.SafeReadFile(file)
			if err != nil {
				log.WithError(err)
				return
			}

			payload, err = util.ValidateSchema(v, struct {
				Short string `json:"short"`
				Long  string `json:"long"`
			}{})
			if err != nil {
				log.WithError(err).Error("failed to parse json input")
				return
			}
		case js != "":
			payload, err = util.ValidateSchema([]byte(js), struct {
				Short string `json:"short"`
				Long  string `json:"long"`
			}{})
			if err != nil {
				log.WithError(err).Error("failed to parse json input")
				return
			}
		default:
			log.Error("no data source provided").Error("'--data', '--data-file' or '--data-json' must be specified")
			return
		}

		body := bytes.Buffer{}
		body.Write(payload)

		ctx := http.New(cfg, &cfg.Application, log)
		req := ctx.Request("POST", "/api/application/locations", &body)
		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}

		buf, err := http.HandleItemResponse(res, cfg)
		if err != nil {
			log.WithError(err)
			return
		}

		log.LineB(buf)
	},
}

var deleteLocationCmd = &cobra.Command{
	Use:   "locations:delete id",
	Short: "deletes a location",
	Long:  "Deletes a location",
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
		req := ctx.Request("DELETE", "/api/application/locations/"+args[0], nil)
		if _, err := ctx.Execute(req); err != nil {
			log.WithError(err)
		}
	},
}
