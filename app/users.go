package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/http"
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
	Use:   "users:create --src path",
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

		src, _ := cmd.Flags().GetString("src")
		if src == "" {
			log.Error("a source file must be provided")
			return
		}

		input, err := util.SafeReadFile(src)
		if err != nil {
			log.WithError(err)
			return
		}

		var schema struct {
			Username   string `json:"username"`
			Email      string `json:"email"`
			FirstName  string `json:"first_name"`
			LastName   string `json:"last_name"`
			Password   string `json:"password,omitempty"`
			RootAdmin  bool   `json:"root_admin,omitempty"`
			ExternalID string `json:"external_id,omitempty"`
		}
		if err = json.Unmarshal(input, &schema); err != nil {
			log.Error("failed to parse json:").WithError(err)
			return
		}

		data, _ := json.Marshal(schema)
		body := bytes.Buffer{}
		body.Write(data)

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
