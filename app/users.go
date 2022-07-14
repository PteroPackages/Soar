package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/http"
	"github.com/pteropackages/soar/logger"
	"github.com/spf13/cobra"
)

var log = logger.New()

type user struct {
	ID         int    `json:"id"`
	ExternalID string `json:"external_id"`
	UUID       string `json:"uuid"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Language   string `json:"language"`
	RootAdmin  bool   `json:"root_admin"`
	TwoFactor  bool   `json:"2fa"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

var getUsersCmd = &cobra.Command{
	Use: "users:get",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
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
		if single {
			req := ctx.Request("GET", "/api/application/users"+query)
			buf, err := ctx.Execute(req)
			if err != nil {
				log.WithError(err)
				return
			}
			if buf == nil {
				return
			}

			var model struct {
				O string `json:"object"`
				A user   `json:"attributes"`
			}
			if err = json.Unmarshal(buf, &model); err != nil {
				log.Error("failed to parse json:").WithError(err)
				return
			}

			var fmt []byte
			if cfg.Http.ParseBody {
				fmt, err = json.MarshalIndent(model.A, "", "  ")
			} else {
				fmt, err = json.MarshalIndent(model, "", "  ")
			}
			if err != nil {
				log.Error("failed to parse response:").WithError(err)
				return
			}

			log.LineB(fmt)
			return
		}

		req := ctx.Request("GET", "/api/application/users")
		buf, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}
		if buf == nil {
			return
		}

		var model struct {
			O string `json:"object"`
			D []struct {
				O string `json:"object"`
				A user   `json:"attributes"`
			} `json:"data"`
		}
		if err = json.Unmarshal(buf, &model); err != nil {
			log.Error("failed to parse json:").WithError(err)
			return
		}

		var fmt []byte
		if cfg.Http.ParseBody {
			var inner []user
			for _, u := range model.D {
				inner = append(inner, u.A)
			}
			fmt, err = json.MarshalIndent(inner, "", "  ")
		} else {
			fmt, err = json.MarshalIndent(model, "", "  ")
		}
		if err != nil {
			log.Error("failed to parse response:").WithError(err)
			return
		}

		log.LineB(fmt)
	},
}

func parseUserQuery(cmd *cobra.Command) (bool, string, error) {
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
