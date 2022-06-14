package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

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
	var root string

	if local {
		root, _ = os.Getwd()
	} else {
		env, ok := os.LookupEnv("SOAR_PATH")
		if !ok {
			return nil, errors.New("environment variable 'SOAR_PATH' not set")
		}

		root = env
	}

	path := filepath.Join(root, "soar.yml")
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("soar path does not exist")
		}

		return nil, err
	}

	if info.Mode()&fs.FileMode(os.O_RDONLY) == 0o0 {
		return nil, errors.New("soar config file is not readable")
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
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		if info.Name() != "soar.yml" {
			return errors.New("refusing to overwrite non-soar config file")
		}

		if !force {
			return errors.New("a soar config already exists at this file path")
		}
	}

	perms := os.O_CREATE | os.O_RDWR | os.O_TRUNC
	if info.Mode()&fs.FileMode(perms) == 0o0 {
		return errors.New("missing read/write permissions for this file path")
	}

	file, err := os.OpenFile(path, perms, fs.FileMode(perms))
	if err != nil {
		return err
	}
	defer file.Close()

	buf, _ := yaml.Marshal(Config{})
	file.Write(buf)

	return nil
}
