package config

import (
	"fun.flight/pkg/resiver"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Resivers []resiver.ResiverConfig `yaml:"resivers"`
}

func ConfigDefault() Config {
	return Config{
		Resivers: make([]resiver.ResiverConfig, 0, 0),
	}
}

func Load(data []byte) (Config, error) {
	conf := ConfigDefault()

	err := yaml.Unmarshal(data, &conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}
