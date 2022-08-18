package client

import (
	"bytes"
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/http"
	"github.com/pteropackages/soar/util"
	"github.com/spf13/cobra"
)

type file struct {
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	Mode       string `json:"mode"`
	ModeBits   string `json:"mode_bits"`
	MimeType   string `json:"mimetype"`
	IsFile     bool   `json:"is_file"`
	IsSymlink  bool   `json:"is_symlink"`
	CreatedAt  string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
}

var listFilesCmd = &cobra.Command{
	Use:     "files:list identifier [-d | --dir] [-f | --file]",
	Aliases: []string{"files:ls", "files:dir"},
	Short:   "lists files on a server",
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

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("GET", "/api/client/servers/"+args[0]+"/files/list", nil)
		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}

		fileOnly, _ := cmd.Flags().GetBool("file")
		dirOnly, _ := cmd.Flags().GetBool("dir")

		if fileOnly || dirOnly {
			var model struct {
				O string `json:"object"`
				D []struct {
					O string `json:"object"`
					A *file  `json:"attributes"`
				} `json:"data"`
			}
			if err = json.Unmarshal(res, &model); err != nil {
				log.WithError(err)
				return
			}

			var filtered struct {
				O string `json:"object"`
				D []struct {
					O string `json:"object"`
					A *file  `json:"attributes"`
				} `json:"data"`
			}
			filtered.O = "list"

			if fileOnly {
				for _, file := range model.D {
					if file.A.MimeType == "inode/directory" {
						continue
					}

					filtered.D = append(filtered.D, file)
				}
			} else {
				for _, file := range model.D {
					if file.A.MimeType == "inode/directory" {
						filtered.D = append(filtered.D, file)
					}
				}
			}

			var buf []byte

			if cfg.Http.ParseBody {
				inner := make([]*file, 0, len(filtered.D))
				for _, file := range filtered.D {
					inner = append(inner, file.A)
				}

				if cfg.Http.ParseIndent {
					buf, err = json.MarshalIndent(inner, "", "  ")
				} else {
					buf, err = json.Marshal(inner)
				}
			} else {
				if cfg.Http.ParseIndent {
					buf, err = json.MarshalIndent(filtered, "", "  ")
				} else {
					buf, err = json.Marshal(filtered)
				}
			}

			if err != nil {
				log.WithError(err)
				return
			}

			log.LineB(buf)
		} else {
			buf, err := http.HandleDataResponse(res, cfg)
			if err != nil {
				log.WithError(err)
				return
			}

			log.LineB(buf)
		}
	},
}

var getFileContentsCmd = &cobra.Command{
	Use:     "files:contents identifier name",
	Aliases: []string{"files:cat"},
	Short:   "gets the contents of a file",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier", "name"}); err != nil {
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

		path := "/api/client/servers/" + args[0] + "/files/contents?file="
		path += url.QueryEscape(args[1])

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("GET", path, nil)
		req.Header.Set("Accept", "text/plain")

		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}

		log.LineB(res)
	},
}

var downloadFileCmd = &cobra.Command{
	Use:     "files:download identifier name [--dest path] [-U | --url-only]",
	Aliases: []string{"files:down"},
	Short:   "downloads a file or returns the url",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier", "name"}); err != nil {
			log.WithError(err)
			return
		}

		dest, _ := cmd.Flags().GetString("dest")
		if dest == "" {
			cwd, _ := os.Getwd()
			dest = filepath.Join(cwd, args[1])
		}

		skip, _ := cmd.Flags().GetBool("url-only")

		_, err := os.Stat(dest)
		if err == nil && !skip {
			log.Error("destination path already exists")
			return
		}

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		path := "/api/client/servers/" + args[0] + "/files/download?file="
		path += url.QueryEscape(args[1])

		ctx := http.New(cfg, &cfg.Client, log)
		urlReq := ctx.Request("GET", path, nil)
		res, err := ctx.Execute(urlReq)
		if err != nil {
			log.WithError(err)
			return
		}

		var model struct {
			Attributes struct {
				URL string `json:"url"`
			} `json:"attributes"`
		}
		if err = json.Unmarshal(res, &model); err != nil {
			log.WithError(err)
			return
		}

		if skip {
			log.Line(model.Attributes.URL)
			return
		}

		dlReq := http.Request("GET", model.Attributes.URL, nil)
		dlReq.Header.Set("Accept", "application/octet-stream")
		res, err = ctx.Execute(dlReq)
		if err != nil {
			log.WithError(err)
			return
		}

		file, err := os.Create(dest)
		if err != nil {
			log.Error("failed to create file:").WithError(err)
			return
		}

		log.Debug("attempting file write")
		defer file.Close()
		file.Write(res)
	},
}

var renameFileCmd = &cobra.Command{
	Use:   "files:rename identifier old-name new-name [--root path]",
	Short: "renames a file on the server",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier", "old-name", "new-name"}); err != nil {
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

		root, _ := cmd.Flags().GetString("root")
		info := map[string]string{"from": args[1], "to": args[2]}
		data, _ := json.Marshal(struct {
			Root  string              `json:"root"`
			Files []map[string]string `json:"files"`
		}{Root: root, Files: []map[string]string{info}})

		body := bytes.Buffer{}
		body.Write(data)

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("PUT", "/api/client/servers/"+args[0]+"/files/rename", &body)
		if _, err = ctx.Execute(req); err != nil {
			log.WithError(err)
		}
	},
}

var copyFileCmd = &cobra.Command{
	Use:     "files:copy identifier name",
	Aliases: []string{"files:cp"},
	Short:   "copies a file",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier", "name"}); err != nil {
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

		data, _ := json.Marshal(map[string]string{"location": args[1]})
		body := bytes.Buffer{}
		body.Write(data)

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("POST", "/api/client/servers/"+args[0]+"/files/copy", &body)
		if _, err = ctx.Execute(req); err != nil {
			log.WithError(err)
		}
	},
}

var writeFileCmd = &cobra.Command{
	Use:        "files:write identifier name content",
	Short:      "writes content to a file",
	SuggestFor: []string{"files:create"},
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgsOverflow(args, []string{"identifier", "name"}, 1); err != nil {
			log.WithError(err)
			return
		}

		if len(args) == 2 {
			log.Error("missing argument 'content'").Error("did you mean to run the 'files:create' command?")
			return
		}

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		body := bytes.Buffer{}
		body.Write([]byte(args[2]))

		path := "/api/client/servers/" + args[0] + "/files/write?file="
		path += url.QueryEscape(args[1])

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("POST", path, &body)
		req.Header.Set("Content-Type", "text/plain")
		if _, err = ctx.Execute(req); err != nil {
			log.WithError(err)
		}
	},
}

var createFileCmd = &cobra.Command{
	Use:        "files:create identifier name",
	Aliases:    []string{"files:touch"},
	Short:      "creates an empty file",
	SuggestFor: []string{"files:create"},
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgsOverflow(args, []string{"identifier", "name"}, 1); err != nil {
			log.WithError(err)
			return
		}

		if len(args) > 2 {
			log.Error("got %d more argument(s) than required (expected 2)", len(args)-2).
				Error("did you mean to run the 'files:create' command?")
			return
		}

		local, _ := cmd.Flags().GetBool("local")
		cfg, err := config.Get(local)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		path := "/api/client/servers/" + args[0] + "/files/write?file="
		path += url.QueryEscape(args[1])

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("POST", path, nil)
		if _, err = ctx.Execute(req); err != nil {
			log.WithError(err)
		}
	},
}
