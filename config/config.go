package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Token string `mapstructure:"token"`
}

func GetConfig(name string) *Config {
	config := new(Config)

	viper.AddConfigPath(".")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.AddConfigPath("../../config")

	viper.SetConfigFile(name)
	viper.SetConfigType("json")

	viper.AutomaticEnv()
	//err := viper.BindEnv("TELEGRAM_TOKEN")
	//if err != nil {
	//	return nil
	//}

	//fmt.Println(viper.Get("TELEGRAM_TOKEN"))

	err := viper.ReadInConfig()
	if err != nil {
		return nil
	}

	err = viper.Unmarshal(config)
	if err != nil {
		return nil
	}

	return config
}