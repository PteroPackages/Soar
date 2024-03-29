package client

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

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

type fractalFile struct {
	O string `json:"object"`
	A *file  `json:"attributes"`
}

type fractalFileList struct {
	O string        `json:"object"`
	D []fractalFile `json:"data"`
}

var listFilesCmd = &cobra.Command{
	Use:     "files:list identifier [-d | --dir] [-f | --file] [--root dir]",
	Aliases: []string{"files:ls", "files:dir"},
	Short:   "lists files on a server",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier"}); err != nil {
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

		root, _ := cmd.Flags().GetString("root")
		path := "/api/client/servers/" + args[0] + "/files/list?list&directory="
		path += url.QueryEscape(root)

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("GET", path, nil)
		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}

		fileOnly, _ := cmd.Flags().GetBool("file")
		dirOnly, _ := cmd.Flags().GetBool("dir")

		if fileOnly || dirOnly {
			var model fractalFileList

			if err = json.Unmarshal(res, &model); err != nil {
				log.WithError(err)
				return
			}

			var filtered fractalFileList
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

var getFileInfoCmd = &cobra.Command{
	Use:     "files:info identifier path",
	Aliases: []string{"files:stat"},
	Short:   "gets the file info for a specific file",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier", "path"}); err != nil {
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

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("GET", "/api/client/servers/"+args[0]+"/files/list", nil)
		res, err := ctx.Execute(req)
		if err != nil {
			log.WithError(err)
			return
		}

		var model fractalFileList
		if err = json.Unmarshal(res, &model); err != nil {
			log.WithError(err)
			return
		}

		target := fractalFile{}
		for _, file := range model.D {
			if file.A.Name == args[1] {
				target = file
				break
			}
		}

		if target == (fractalFile{}) {
			log.Error("file not found")
			return
		}

		var buf []byte

		if cfg.Http.ParseBody {
			if cfg.Http.ParseIndent {
				buf, err = json.MarshalIndent(target.A, "", "  ")
			} else {
				buf, err = json.Marshal(target.A)
			}
		} else {
			if cfg.Http.ParseIndent {
				buf, err = json.MarshalIndent(target, "", "  ")
			} else {
				buf, err = json.Marshal(target)
			}
		}

		if err != nil {
			log.WithError(err)
			return
		}

		log.LineB(buf)
	},
}

var getFileContentsCmd = &cobra.Command{
	Use:     "files:contents identifier path",
	Aliases: []string{"files:cat"},
	Short:   "gets the contents of a file",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier", "path"}); err != nil {
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
	Use:     "files:download identifier path [--dest path] [-U | --url-only]",
	Aliases: []string{"files:down"},
	Short:   "downloads a file or returns the url",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier", "path"}); err != nil {
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

		global, _ := cmd.Flags().GetBool("global")
		cfg, err := config.Get(global)
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
	Use:     "files:rename identifier old-name new-name [--root path]",
	Aliases: []string{"files:move", "files:mv"},
	Short:   "renames a file on the server",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier", "old-name", "new-name"}); err != nil {
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

		root, _ := cmd.Flags().GetString("root")
		info := map[string]string{"from": args[1], "to": args[2]}
		data, _ := json.Marshal(map[string]interface{}{"root": root, "files": info})

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
	Use:     "files:copy identifier path",
	Aliases: []string{"files:cp"},
	Short:   "copies a file",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier", "path"}); err != nil {
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
	Use:        "files:write identifier path content",
	Short:      "writes content to a file",
	SuggestFor: []string{"files:create"},
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgsOverflow(args, []string{"identifier", "path"}, 1); err != nil {
			log.WithError(err)
			return
		}

		if len(args) == 2 {
			log.Error("missing argument 'content'").Error("did you mean to run the 'files:create' command?")
			return
		}

		global, _ := cmd.Flags().GetBool("global")
		cfg, err := config.Get(global)
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
	Use:        "files:create identifier path",
	Aliases:    []string{"files:touch"},
	Short:      "creates an empty file",
	SuggestFor: []string{"files:create"},
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgsOverflow(args, []string{"identifier", "path"}, 1); err != nil {
			log.WithError(err)
			return
		}

		if len(args) > 2 {
			log.Error("got %d more argument(s) than required (expected 2)", len(args)-2).
				Error("did you mean to run the 'files:create' command?")
			return
		}

		global, _ := cmd.Flags().GetBool("global")
		cfg, err := config.Get(global)
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

var compressFilesCmd = &cobra.Command{
	Use:     "files:compress identifier files... [--root dir]",
	Aliases: []string{"files:cmp", "files:zip"},
	Short:   "compresses one or more files and folders",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgsOverflow(args, []string{"identifier"}, 20); err != nil {
			log.WithError(err)
			return
		}

		if len(args) == 1 {
			log.Error("at least one file must be specified to compress")
			return
		}

		global, _ := cmd.Flags().GetBool("global")
		cfg, err := config.Get(global)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		root, _ := cmd.Flags().GetString("root")
		data, _ := json.Marshal(map[string]interface{}{"root": root, "files": args[1:]})
		body := bytes.Buffer{}
		body.Write(data)

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("POST", "/api/client/servers/"+args[0]+"/files/compress", &body)
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

var decompressFileCmd = &cobra.Command{
	Use:     "files:decompress identifier name [--root dir]",
	Aliases: []string{"files:dcmp", "files:unzip"},
	Short:   "decompresses an archived file",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier", "name"}); err != nil {
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

		root, _ := cmd.Flags().GetString("root")
		data, _ := json.Marshal(map[string]string{"root": root, "file": args[1]})
		body := bytes.Buffer{}
		body.Write(data)

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("POST", "/api/client/servers/"+args[0]+"/files/decompress", &body)
		if _, err = ctx.Execute(req); err != nil {
			log.WithError(err)
		}
	},
}

var deleteFilesCmd = &cobra.Command{
	Use:     "files:delete identifer files... [--root dir]",
	Aliases: []string{"files:rm"},
	Short:   "deletes one or more files",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgsOverflow(args, []string{"identifier"}, 10); err != nil {
			log.WithError(err)
			return
		}

		if len(args) == 1 {
			log.Error("at least one file must be specified to delete")
			return
		}

		global, _ := cmd.Flags().GetBool("global")
		cfg, err := config.Get(global)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		root, _ := cmd.Flags().GetString("root")
		data, _ := json.Marshal(map[string]interface{}{"root": root, "files": args[1:]})
		body := bytes.Buffer{}
		body.Write(data)

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("POST", "/api/client/servers/"+args[0]+"/files/delete", &body)
		if _, err = ctx.Execute(req); err != nil {
			log.WithError(err)
		}
	},
}

var createFolderCmd = &cobra.Command{
	Use:     "files:folder identifier name [--root dir]",
	Aliases: []string{"files:mkdir"},
	Short:   "creates a folder",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier", "name"}); err != nil {
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

		root, _ := cmd.Flags().GetString("root")
		data, _ := json.Marshal(map[string]string{"root": root, "name": args[1]})
		body := bytes.Buffer{}
		body.Write(data)

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("POST", "/api/client/servers/"+args[0]+"/files/create-folder", &body)
		if _, err = ctx.Execute(req); err != nil {
			log.WithError(err)
		}
	},
}

var chmodFileCmd = &cobra.Command{
	Use:   "files:chmod identifier name mode [--root dir]",
	Short: "changes the permissions of a file",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier", "name", "mode"}); err != nil {
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

		mode, err := strconv.Atoi(args[2])
		if err != nil {
			log.Error("failed to parse mode bits:").WithError(err)
			return
		}

		root, _ := cmd.Flags().GetString("root")
		info := map[string]interface{}{"file": args[1], "mode": mode}
		data, _ := json.Marshal(map[string]interface{}{"root": root, "files": []interface{}{info}})
		body := bytes.Buffer{}
		body.Write(data)

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("POST", "/api/client/servers/"+args[0]+"/files/chmod", &body)
		if _, err = ctx.Execute(req); err != nil {
			log.WithError(err)
		}
	},
}

var pullFileCmd = &cobra.Command{
	Use:   "files:pull identifier url [--dest dir] [--name name]\n\t[--use-header] [-f | --foreground]",
	Short: "pulls a file from a remote source",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgs(args, []string{"identifier", "url"}); err != nil {
			log.WithError(err)
			return
		}

		source, err := url.Parse(args[1])
		if err != nil {
			log.Error("failed to parse url:").WithError(err)
			return
		}

		global, _ := cmd.Flags().GetBool("global")
		cfg, err := config.Get(global)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		dest, _ := cmd.Flags().GetString("dest")
		name, _ := cmd.Flags().GetString("name")
		useHeader, _ := cmd.Flags().GetBool("use-header")
		foreground, _ := cmd.Flags().GetBool("foreground")

		info := map[string]interface{}{
			"url":        source.String(),
			"use_header": useHeader,
			"foreground": foreground,
		}
		if dest != "" {
			info["directory"] = dest
		}
		if name != "" {
			info["filename"] = name
		}

		data, _ := json.Marshal(info)
		body := bytes.Buffer{}
		body.Write(data)

		ctx := http.New(cfg, &cfg.Client, log)
		req := ctx.Request("POST", "/api/client/servers/"+args[0]+"/files/pull", &body)
		if _, err = ctx.Execute(req); err != nil {
			log.WithError(err)
		}
	},
}

var uploadFilesCmd = &cobra.Command{
	Use:     "files:upload identifier files... [--dest path] [-U | --url-only]",
	Aliases: []string{"files:up"},
	Short:   "uploads one or more files to the server",
	Run: func(cmd *cobra.Command, args []string) {
		log.ApplyFlags(cmd.Flags())
		if err := util.RequireArgsOverflow(args, []string{"identifier"}, 10); err != nil {
			log.WithError(err)
			return
		}

		skip, _ := cmd.Flags().GetBool("url-only")

		if len(args) == 1 && !skip {
			log.Error("at least one file must be specified to upload")
			return
		}

		global, _ := cmd.Flags().GetBool("global")
		cfg, err := config.Get(global)
		if err != nil {
			config.HandleError(err, log)
			return
		}
		cfg.ApplyFlags(cmd.Flags())

		ctx := http.New(cfg, &cfg.Client, log)
		urlReq := ctx.Request("GET", "/api/client/servers/"+args[0]+"/files/upload", nil)
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

		var files []string

		for _, path := range args[1:] {
			info, err := os.Stat(path)
			if err != nil {
				if os.IsNotExist(err) {
					log.Warn("'%s' file does not exist", path)
				} else {
					log.Warn("'%s' raised an unexpected error, skipping", path)
				}

				continue
			}

			if info.IsDir() {
				log.Warn("'%s' cannot upload directories, skipping", path)
				continue
			}

			files = append(files, path)
		}

		log.Debug("%v", files)

		if len(files) == 0 {
			log.Error("no files found to upload")
			return
		}

		body := bytes.Buffer{}
		writer := multipart.NewWriter(&body)

		for _, path := range files {
			file, err := os.Open(path)
			if err != nil {
				log.Warn("could not open '%s', skipping", path)
				continue
			}

			part, _ := writer.CreateFormFile("files", file.Name())
			io.Copy(part, file)
			file.Close()
		}

		writer.Close()

		upReq := http.Request("POST", model.Attributes.URL, &body)
		upReq.Header.Set("Content-Type", writer.FormDataContentType())
		if _, err = ctx.Execute(upReq); err != nil {
			log.WithError(err)
			return
		}
	},
}
