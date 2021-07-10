package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Path    string `yaml:"path"`
	AppFile string `yaml:"appFile"`
	// To help WinAPI search the app
	WindowTitle  string `yaml:"windowTitle"`
	HWKey        bool   `yaml:"hardwareKey"`
	ScreenWidth  int    `yaml:"screenWidth"`
	ScreenHeight int    `yaml:"screenHeight"`
	IsWindowMode *bool  `yaml:"isWindowMode"`
	// Frontend plugin
	HasChat   bool   `yaml:"hasChat"`
	PageTitle string `yaml:"pageTitle"`
}

func ReadConfig(path string) (Config, error) {
	cfgyml, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	err = yaml.Unmarshal(cfgyml, &cfg)

	if cfg.ScreenWidth == 0 {
		cfg.ScreenWidth = 800
	}
	if cfg.ScreenHeight == 0 {
		cfg.ScreenHeight = 600
	}
	return cfg, err
}
