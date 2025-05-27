package config

import (
	"sync"

	"github.com/pkg/errors"
)

var cfgLock = &sync.Mutex{}

type Config struct {
	App   *AppConfig
	Infra *InfraConfig
}

func App() *AppConfig {
	return instance().App
}

func Infra() *InfraConfig {
	return instance().Infra
}

var cfg *Config

func SetConfig(c *Config) {
	cfg = c
}

func instance() *Config {
	cfgLock.Lock()
	defer cfgLock.Unlock()

	if cfg == nil {
		cfg = &Config{
			App:   &AppConfig{},
			Infra: &InfraConfig{},
		}
	}

	return cfg
}

func Load() error {
	app, err := loadApplicationConfig()
	if err != nil {
		return errors.Wrap(err, "loading application config")
	}

	infra, err := loadInfraConfig()
	if err != nil {
		return errors.Wrap(err, "loading infra config")
	}

	cfg = &Config{
		App:   app,
		Infra: infra,
	}

	return nil
}
