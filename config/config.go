package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string
}

func InitConfig() {
	fmt.Println("init config")
	viper.AutomaticEnv()
}
