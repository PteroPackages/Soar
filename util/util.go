package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func ApplyDefaultFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("debug", false, "print debug logs")
	cmd.Flags().Bool("no-color", false, "disable ansi color codes")
	cmd.Flags().BoolP("global", "g", false, "use the global config")
	cmd.Flags().BoolP("quiet", "q", false, "only print necessary logs")

	cmd.Flags().BoolP("retry-ratelimit", "r", false, "retry request on ratelimit")
	cmd.Flags().BoolP("no-retry-ratelimit", "R", false, "don't retry request on ratelimit")
	cmd.Flags().BoolP("parse-body", "b", false, "parse the response body")
	cmd.Flags().BoolP("no-parse-body", "B", false, "don't parse the response body")
	cmd.Flags().BoolP("parse-errors", "e", false, "parse the response errors")
	cmd.Flags().BoolP("no-parse-errors", "E", false, "don't parse the response errors")
	cmd.Flags().BoolP("parse-indent", "i", false, "indent the response body")
	cmd.Flags().BoolP("no-parse-indent", "I", false, "don't indent the response body")
}

func ApplyDataFlags(cmd *cobra.Command) {
	cmd.Flags().String("data", "", "a set of key-value pairs for the request")
	cmd.Flags().String("data-file", "", "a file path to the json data")
	cmd.Flags().String("data-json", "", "the json data for the request")
}

func ApplyFilterFlags(cmd *cobra.Command) {
	cmd.Flags().Int("page", 0, "the page to request from")
	cmd.Flags().Int("per-page", 0, "the number of results to return")
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

func RequireArgs(input, required []string) error {
	if len(input) == 0 {
		return fmt.Errorf("no arguments specified (expected %d)", len(required))
	}

	if len(input) < len(required) {
		missing := required[len(input)]
		include := ""
		if len(required)-len(input) > 1 {
			include = fmt.Sprintf(" and %d more", len(required)-1)
		}

		return fmt.Errorf("missing argument '%s'%s", missing, include)
	}

	if len(input) > len(required) {
		return fmt.Errorf("got %d more argument(s) than required (expected %d)", len(input)-len(required), len(required))
	}

	return nil
}

func RequireArgsOverflow(input, required []string, overflow int) error {
	if len(input) == 0 {
		return fmt.Errorf("no arguments specified (expected %d)", len(required)+overflow)
	}

	if len(input) < len(required) {
		missing := required[len(input)]
		include := ""
		if len(required)-len(input) > 1 {
			include = fmt.Sprintf(" and %d more", len(required)-1)
		}

		return fmt.Errorf("missing argument '%s'%s", missing, include)
	}

	if len(input) > len(required)+overflow {
		return fmt.Errorf("got %d more argument(s) than required (expected %d)", len(input)-len(required)-overflow, len(required)+overflow)
	}

	return nil
}

func ValidateSchema(in []byte, schema interface{}) ([]byte, error) {
	if err := json.Unmarshal(in, &schema); err != nil {
		return nil, err
	}

	return json.Marshal(schema)
}
