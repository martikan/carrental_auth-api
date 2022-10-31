package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	PostgreDriver   string `mapstructure:"POSTGRES_DRIVER"`
	PostgreHost     string `mapstructure:"POSTGRES_HOST"`
	PostgrePort     string `mapstructure:"POSTGRES_PORT"`
	PostgreUser     string `mapstructure:"POSTGRES_USER"`
	PostgrePassword string `mapstructure:"POSTGRES_PASSWORD"`
	PostgreDb       string `mapstructure:"POSTGRES_DB"`
	SSLMode         string `mapstructure:"SSL_MODE"`

	Port string `mapstructure:"PORT"`

	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	TokenSymetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
