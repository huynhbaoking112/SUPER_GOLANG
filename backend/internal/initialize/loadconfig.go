package initialize

import (
	"fmt"
	"go-backend-v2/global"

	"github.com/spf13/viper"
)

func LoadConfig() {
	viper := viper.New()
	viper.AddConfigPath("configs")
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("ERROR READING CONFIG FILE: %v", err))
	}

	if err := viper.Unmarshal(&global.Config); err != nil {
		panic(fmt.Errorf("ERROR UNMARSHALING CONFIG FILE: %v", err))
	}
}
