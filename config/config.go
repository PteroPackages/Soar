package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/pteropackages/soar/logger"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

type Auth struct {
	URL string `validate:"required,url" yaml:"url"`
	Key string `validate:"required" yaml:"key"`
}

type HttpConfig struct {
	ParseBody      bool `yaml:"parse_body"`
	ParseIndent    bool `yaml:"parse_indent"`
	RetryRateLimit bool `yaml:"retry_rate_limit"`
}

type LogConfig struct {
	UseColor       bool `yaml:"use_color"`
	UseDebug       bool `yaml:"use_debug"`
	IgnoreWarnings bool `yaml:"ignore_warnings"`
}

type Config struct {
	Application Auth       `validate:"required" yaml:"application"`
	Client      Auth       `validate:"required" yaml:"client"`
	Http        HttpConfig `validate:"required" yaml:"http"`
	Logs        LogConfig  `yaml:"logs"`
}

func (c *Config) Format() string {
	fmt, _ := yaml.Marshal(c)

	return string(fmt)
}

func (c *Config) ApplyFlags(flags *pflag.FlagSet) {
	if ok, _ := flags.GetBool("retry-ratelimit"); ok {
		c.Http.RetryRateLimit = true
	}
	if ok, _ := flags.GetBool("no-retry-ratelimit"); ok {
		c.Http.RetryRateLimit = false
	}

	if ok, _ := flags.GetBool("parse-body"); ok {
		c.Http.ParseBody = true
	}
	if ok, _ := flags.GetBool("no-parse-body"); ok {
		c.Http.ParseBody = false
	}

	if ok, _ := flags.GetBool("parse-indent"); ok {
		c.Http.ParseIndent = true
	}
	if ok, _ := flags.GetBool("no-parse-indent"); ok {
		c.Http.ParseIndent = false
	}
}

func GetStatic(global bool) (*Config, error) {
	var path string

	if !global {
		root, _ := os.Getwd()
		path = filepath.Join(root, ".soar.yml")

		if _, err := os.Stat(path); err != nil {
			path = ""
		}
	}

	if path == "" {
		root, err := os.UserConfigDir()
		if err != nil {
			return nil, err
		}

		if _, err = os.Stat(root); err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("user config directory not found (path: %s)", root)
			}

			return nil, err
		}

		path = filepath.Join(root, ".soar", "config.yml")
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("file path does not exist")
		}

		return nil, err
	}

	if info.IsDir() {
		return nil, errors.New("invalid file path, cannot be a directory")
	}

	if info.Mode()&0o644 == 0 {
		return nil, errors.New("file path is not readable")
	}

	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg *Config
	if err = yaml.Unmarshal(buf, &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func Get(global bool) (*Config, error) {
	cfg, err := GetStatic(global)
	if err != nil {
		return nil, err
	}

	validate := validator.New()
	err = validate.Struct(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func Create(path string, force bool) (string, error) {
	root, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	if _, err = os.Stat(root); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("user config directory not found (path: %s)", root)
		}

		return "", err
	}

	if path == "" {
		base := filepath.Join(root, ".soar")
		if _, err = os.Stat(base); err != nil {
			if err = os.MkdirAll(base, 0o755); err != nil {
				return "", errors.New("failed to create .soar directory at config directory")
			}
		}

		path = filepath.Join(base, "config.yml")
	} else {
		cwd, _ := os.Getwd()
		path = filepath.Join(cwd, path)
	}

	if !filepath.IsAbs(path) {
		return "", errors.New("file path is not absolute")
	}

	writeFile := func() error {
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
		if err != nil {
			return err
		}
		defer file.Close()

		buf, _ := yaml.Marshal(Config{})
		file.Write(buf)

		return nil
	}

	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			path = filepath.Join(path, ".soar.yml")
			if err = writeFile(); err != nil {
				return "", err
			}

			return path, nil
		}

		if !strings.HasSuffix(path, "config.yml") && !strings.HasSuffix(path, ".soar.yml") {
			return "", errors.New("refusing to overwrite non-soar config file")
		}

		if !force {
			return "", errors.New("a soar config already exists at this file path")
		}

		if info.Mode()&fs.FileMode(os.O_RDWR) == 0 {
			return "", errors.New("missing read/write permissions for this file path")
		}
	}

	if err = writeFile(); err != nil {
		return "", err
	}

	return path, nil
}

func HandleError(err error, log *logger.Logger) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		log.Error("failed to validate config, %d error(s):", len(errs))

		for _, e := range errs {
			log.Error(fmt.Sprintf("field %s didn't satisfy the '%s' tag", e.Namespace(), e.Tag()))
		}
	} else {
		log.Error("failed to get config:").WithError(err).WithCmd("soar config --help")
	}
}
