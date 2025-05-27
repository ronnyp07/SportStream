package config

import (
	"github.com/spf13/viper"
)

const (
	infraEnviromentConfig = "/opt/poller/"
)

type InfraConfig struct {
	Nats         NATSInfrastructure         `mapstructure:"NATS"`
	MessageQueue MessageQueueInfrastructure `mapstructure:"MESSAGE_QUEUE"`
}

type MessageQueueInfrastructure struct {
	Enabled bool `mapstructure:"NATS_ENABLED"`
}

type NATSInfrastructure struct {
	Port string `mapstructure:"NATS_PORT"`
	Host string `mapstructure:"NATS_HOST"`
}

func loadInfraConfig() (config *InfraConfig, err error) {
	name := "infra"
	config = &InfraConfig{}

	viper.AddConfigPath(".")
	viper.SetConfigName(name)
	viper.SetConfigType("env")
	//viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return
}
