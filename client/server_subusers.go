package client

import (
	"bytes"
	"encoding/json"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/http"
	"github.com/pteropackages/soar/util"
	"github.com/spf13/cobra"
)

var getSubUsersCmd = &cobra.Command{
	Use:   "subusers:get identifier [--uuid id]",
	Short: "gets the server subusers",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier"}); err != nil {
			log.WithError(err)
			return
		}

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		uuid, _ := cmd.Flags().GetString("uuid")
		path := "/api/client/servers/" + args[0] + "/users"
		if uuid != "" {
			path += "/" + uuid
		}

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("GET", path, nil)
		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}

		var buf []byte

		if uuid != "" {
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

var addSubUserCmd = &cobra.Command{
	Use:   "subusers:add identifier email permissions...",
	Short: "adds a subuser to the server",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgsOverflow(args, []string{"identifier", "email"}, 20); err != nil {
			log.WithError(err)
			return
		}

		if len(args) == 2 {
			log.Error("at least one permission must be specified to create")
			return
		}

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		data, _ := json.Marshal(map[string]interface{}{"email": args[1], "permissions": args[2:]})
		body := bytes.Buffer{}
		body.Write(data)

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("POST", "/api/client/servers/"+args[0]+"/users", &body)
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

var removeSubUserCmd = &cobra.Command{
	Use:   "subusers:remove identifier uuid",
	Short: "removes a subuser from the server",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier", "uuid"}); err != nil {
			log.WithError(err)
			return
		}

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("DELETE", "/api/client/servers/"+args[0]+"/users/"+args[1], nil)
		if _, err = ctx.Execute(req); err != nil {
			log.WithError(err)
		}
	},
}
