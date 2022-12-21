package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/http"
	"github.com/pteropackages/soar/input"
	"github.com/pteropackages/soar/util"
	"github.com/spf13/cobra"
)

var getUsersCmd = &cobra.Command{
	Use:   "users:get [--id id] [--external id] [--username name] [--email email] [--uuid id]",
	Short: "gets panel users",
	Long:  getUsersHelp,
	Run: func(cmd *cobra.Command, _ []string) {
		log.ApplyFlags(cmd.Flags())

		global, _ := cmd.Flags().GetBool("global")
		cfg, err := config.Get(global)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		single, query, err := parseUserQuery(cmd)
		if err != nil {
			log.Error("command error:").WithError(err)
			return
		}

		ctx := http.New(cfg, &cfg.Application, log)
		req := ctx.Request("GET", "/api/application/users"+query, nil)
		res, err := ctx.ExecuteWithFlags(req, cmd.Flags())
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

func parseUserQuery(cmd *cobra.Command) (bool, string, error) {
	var query strings.Builder
	var params []string
	single := false
	flags := cmd.Flags()

	if id, _ := flags.GetInt("id"); id != 0 {
		single = true
		query.WriteString(fmt.Sprintf("/%d", id))
	}

	if ext, _ := flags.GetString("external"); ext != "" {
		if query.Len() != 0 {
			return false, "", errors.New("id and external flags specified; pick one")
		}

		single = true
		query.WriteString("/external/" + ext)
	}

	if val, _ := flags.GetString("username"); val != "" {
		params = append(params, "filter[username]="+val)
	}

	if val, _ := flags.GetString("email"); val != "" {
		params = append(params, "filter[email]="+val)
	}

	if val, _ := flags.GetString("uuid"); val != "" {
		params = append(params, "filter[uuid]="+val)
	}

	if len(params) > 0 {
		query.WriteString("?" + params[0])

		if len(params) > 1 {
			for _, p := range params {
				query.WriteString("&" + p)
			}
		}
	}

	return single, query.String(), nil
}

var createUserCmd = &cobra.Command{
	Use:   "users:create --data[-file | -json] source",
	Short: "creates a user",
	Long:  createUserHelp,
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
			payload, err = parseInputSource(data)
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

			payload, err = parseJSONSource(v)
			if err != nil {
				log.WithError(err).Error("failed to parse json input")
				return
			}
		case js != "":
			payload, err = parseJSONSource([]byte(js))
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
		req := ctx.Request("POST", "/api/application/users", &body)
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

func parseInputSource(in string) ([]byte, error) {
	m, err := input.Parse(in)
	if err != nil {
		return nil, err
	}

	d := input.Definition{
		"username":    input.StringNode,
		"email":       input.StringNode,
		"external_id": input.NullStringNode,
		"first_name":  input.StringNode,
		"last_name":   input.StringNode,
		"root_admin":  input.BoolNode,
		"password":    input.NullStringNode,
	}

	return input.Marshal(d, m)
}

func parseJSONSource(d []byte) ([]byte, error) {
	var schema struct {
		Username   string `json:"username"`
		Email      string `json:"email"`
		ExternalID string `json:"external_id,omitempty"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		RootAdmin  bool   `json:"root_admin,omitempty"`
		Password   string `json:"password,omitempty"`
	}

	if err := json.Unmarshal(d, &schema); err != nil {
		return nil, err
	}

	v, _ := json.Marshal(schema)
	return v, nil
}

var deleteUserCmd = &cobra.Command{
	Use:   "users:delete id",
	Short: "deletes a user",
	Long:  "Deletes a user account from the panel by its ID.",
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
		req := ctx.Request("DELETE", "/api/application/users/"+args[0], nil)
		if _, err := ctx.Execute(req); err != nil {
			log.WithError(err)
		}
	},
}
