package util

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func ApplyDefaultFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("debug", false, "print debug logs")
	cmd.Flags().Bool("no-color", false, "disable ansi color codes")
	cmd.Flags().BoolP("local", "l", false, "use the local config")
	cmd.Flags().BoolP("quiet", "q", false, "only print necessary logs")

	cmd.Flags().BoolP("retry-ratelimit", "r", false, "retry request on ratelimit")
	cmd.Flags().BoolP("no-retry-ratelimit", "R", false, "don't retry request on ratelimit")
	cmd.Flags().Int("max-body", 0, "the maximum body size to accept")
	cmd.Flags().BoolP("parse", "p", false, "parse the response body")
	cmd.Flags().BoolP("no-parse", "P", false, "don't parse the response body")
}

func SafeReadFile(path string) ([]byte, error) {
	if !filepath.IsAbs(path) {
		root, _ := os.Getwd()
		path = filepath.Join(root, path)
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("file path does not exist")
		}

		return nil, err
	}

	if info.IsDir() {
		return nil, errors.New("invalid file path")
	}

	if info.Mode()&0o644 == 0 {
		return nil, errors.New("file path is not readable")
	}

	return os.ReadFile(path)
}
