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
	Jobs          Jobs          `mapstructure:"JOBS"`
}

type Environment struct {
	Name string `mapstructure:"NAME"`
}

type Logging struct {
	Level string `mapstructure:"LEVEL"`
}

type MessageQueueProcessor struct {
	MaxRetries int32         `mapstructure:"MAX_RETRIES"`
	Limit      int32         `mapstructure:"LIMIT"`
	Wait       time.Duration `mapstructure:"WAIT"`
}

type Nats struct {
	Publishers Publishers `mapstructure:"PUBLISHERS"`
}

type Publishers struct {
	Articles Articles `mapstructure:"ARTICLES"`
}

type Articles struct {
	Update ArticleUpdate `mapstructure:"UPDATE"`
}

type ArticleUpdate struct {
	Subject string `mapstructure:"SUBJECT"`
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

type Job struct {
	Enabled       bool   `mapstructure:"ENABLED"`
	Type          string `mapstructure:"TYPE"`
	Interval      string `mapstructure:"INTERVAL"`
	UseSeconds    bool   `mapstructure:"USESECONDS"`
	ExternalAddrs string `mapstructure:"EXTERNALADDRESS"`
	Retry         Retry  `mapstructure:"RETRY"`
}

type Retry struct {
	MaxAttempts int    `mapstructure:"MAXATTEMPTS"`
	Duration    string `mapstructure:"DURATION"`
}

type Jobs map[string]Job

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
