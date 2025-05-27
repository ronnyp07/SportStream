package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	environmentAppConfig = "/etc/config/app/"
	localAppConfig       = "config/app"
)

type AppConfig struct {
	Env           Environment   `mapstructure:"ENVIRONMENT"`
	Observability Observability `mapstructure:"OBSERVABILITY"`
	Http          Http          `mapstructure:"HTTP"`
}

type Environment struct {
	Name string `mapstructure:"NAME"`
}

type Logging struct {
	Level string `mapstructure:"LEVEL"`
}

type Observability struct {
	Tracing Tracing `mapstructure:"TRACING"`
	Logging Logging `mapstructure:"LOGGING"`
}

type Tracing struct {
	SampleRate float64 `mapstructure:"SAMPLE"`
	TraceHost  string  `mapstructure:"OTEL_TRACE_ENDPOINT"`
	TracePath  string  `mapstructure:"OTEL_TRACE_URL_PATH"`
}

type Http struct {
	HostAddress  string        `mapstructure:"HOST_ADDRESS"`
	ReadTimeout  time.Duration `mapstructure:"READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"WRITE_TIMEOUT"`
	BasePath     string        `mapstructure:"BASE_PATH"`
}

func loadApplicationConfig() (config *AppConfig, err error) {
	config = &AppConfig{}
	name := "config"

	if _, err := os.Stat(fmt.Sprintf("%v%v.yaml", environmentAppConfig, name)); errors.Is(err, os.ErrNotExist) {
		viper.AddConfigPath(localAppConfig)
	} else {
		viper.AddConfigPath(environmentAppConfig)
	}

	viper.SetConfigName(name)
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	var viperErr error
	viper.OnConfigChange(func(e fsnotify.Event) {
		err = viper.Unmarshal(&config)
		if err != nil {
			viperErr = errors.Wrap(err, "loading new config")
		}
	})

	viper.WatchConfig()

	if viperErr != nil {
		return config, viperErr
	}

	return config, nil

}
