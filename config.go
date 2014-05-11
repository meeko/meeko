// Copyright (c) 2013 The mk AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"errors"
	"io/ioutil"
	"os/user"
	"path/filepath"

	"gopkg.in/v1/yaml"
)

const ConfigFileName = ".meekorc"

type Config struct {
	Address         string `yaml:"endpoint_address"`
	AccessToken     string `yaml:"access_token"`
	ManagementToken []byte `yaml:"management_token"`
}

func (cfg *Config) Validate() error {
	switch {
	case len(cfg.Address) == 0:
		return errors.New("endpoint_address is not set")
	case len(cfg.AccessToken) == 0:
		return errors.New("access_token is not set")
	case len(cfg.ManagementToken) == 0:
		return errors.New("management_token is not set")
	}
	return nil
}

func DefaultConfigPath() (string, error) {
	me, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(me.HomeDir, ConfigFileName), nil
}

func MustDefaultConfigPath() string {
	path, err := DefaultConfigPath()
	if err != nil {
		panic(err)
	}
	return path
}

func LoadConfig(path string) (*Config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := yaml.Unmarshal(content, config); err != nil {
		return nil, err
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}
	return config, nil
}
