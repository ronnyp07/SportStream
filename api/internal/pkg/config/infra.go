package config

import (
	"github.com/spf13/viper"
)

const (
	infraEnviromentConfig = "/opt/worker/"
)

type InfraConfig struct {
	Nats         NATSInfrastructure         `mapstructure:"NATS"`
	MessageQueue MessageQueueInfrastructure `mapstructure:"MESSAGE_QUEUE"`
	MongoDB      Mongostructure             `mapstructure:"MONGODB"`
}

type MessageQueueInfrastructure struct {
	Enabled bool `mapstructure:"NATS_ENABLED"`
}

type NATSInfrastructure struct {
	Port int    `mapstructure:"NATS_PORT"`
	Host string `mapstructure:"NATS_HOST"`
}

type Mongostructure struct {
	Url        string `mapstructure:"URL"`
	DataBase   string `mapstructure:"DATABASE"`
	UserName   string `mapstructure:"USERNAME"`
	PassWord   string `mapstructure:"PASSWORD"`
	AuthSource string `mapstructure:"AUTHSOURCE"`
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
