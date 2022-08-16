package client

import (
	"encoding/json"

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
	Use:   "files:list identifier [-d | --dir] [-f | --file]",
	Short: "lists files on a server",
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
