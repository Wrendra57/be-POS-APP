package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	DatabaseURL string
}

func InitConfig() {

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("./")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
}
