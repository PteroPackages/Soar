package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Auth struct {
	URL string `yaml:"url"`
	Key string `yaml:"key"`
}

type HttpConfig struct {
	MaxBodySize    int  `yaml:"max_body_size"`
	ParseBody      bool `yaml:"parse_body"`
	SaveRequests   bool `yaml:"save_requests"`
	RetryRateLimit bool `yaml:"retry_rate_limit"`
}

type LogConfig struct {
	UseColor       bool `yaml:"use_color"`
	UseDebug       bool `yaml:"use_debug"`
	IgnoreWarnings bool `yaml:"ignore_warnings"`
}

type Config struct {
	Application Auth       `yaml:"application"`
	Client      Auth       `yaml:"client"`
	Http        HttpConfig `yaml:"http"`
	Logs        LogConfig  `yaml:"logs"`
}

func (c *Config) Format() string {
	fmt, _ := yaml.Marshal(c)

	return string(fmt)
}

func Get(local bool) (*Config, error) {
	var path string

	if local {
		root, _ := os.Getwd()
		path = filepath.Join(root, "soar.yml")
	} else {
		root, err := os.UserConfigDir()
		if err != nil {
			return nil, err
		}

		path = filepath.Join(root, "soar", "config.yml")
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

func Create(path string, force bool) (string, error) {
	root, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	if path == "" {
		path = filepath.Join(root, "soar", "config.yml")
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
			path = filepath.Join(path, "soar.yml")
			if err = writeFile(); err != nil {
				return "", err
			}

			return path, nil
		}

		if !strings.HasSuffix(path, "config.yml") && !strings.HasSuffix(path, "soar.yml") {
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
