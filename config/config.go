package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var Perms = os.O_CREATE | os.O_RDWR | os.O_TRUNC

type Auth struct {
	URL string
	Key string
}

type HttpConfig struct {
	MaxBodySize    bool
	PromptSave     bool
	SaveRequests   bool
	SendFullBody   bool
	RetryRateLimit bool
}

type LogConfig struct {
	UseColor       bool
	UseDebug       bool
	IgnoreWarnings bool
}

type Config struct {
	Application Auth
	Client      Auth
	Http        HttpConfig
	Logs        LogConfig
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

		path = filepath.Join(root, "config.yml")
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

	if info.Mode()&fs.FileMode(os.O_RDONLY) == 0o0 {
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

func Create(path string, force bool) error {
	root, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	if path == "" {
		path = filepath.Join(root, "config.yml")
	}

	if !filepath.IsAbs(path) {
		return errors.New("file path is not absolute")
	}

	info, err := os.Stat(path)
	if err == nil {
		if info.Name() != "config.yml" && info.Name() != "soar.yml" {
			return errors.New("refusing to overwrite non-soar config file")
		}

		if !force {
			return errors.New("a soar config already exists at this file path")
		}

		if info.Mode()&fs.FileMode(Perms) == 0o0 {
			return errors.New("missing read/write permissions for this file path")
		}
	}

	file, err := os.OpenFile(path, Perms, fs.FileMode(Perms))
	if err != nil {
		return err
	}
	defer file.Close()

	buf, _ := yaml.Marshal(Config{})
	file.Write(buf)

	return nil
}
