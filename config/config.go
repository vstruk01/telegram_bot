package config

import (
	"github.com/spf13/viper"
)

type PostgresConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Name        string `mapstructure:"name"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	AppName     string `mapstructure:"app_name"`
	SourceFiles string `mapstructure:"source_files"`
}

type Config struct {
	Token string         `mapstructure:"TOKEN"`
	DB    PostgresConfig `mapstructure:"db"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app.toml")
	viper.SetConfigType("toml")
	err = viper.BindEnv("TOKEN")
	if err != nil {
		return Config{}, err
	}

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
