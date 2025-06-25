package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	Database []struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		DbName   string `mapstructure:"dbName"`
	} `mapstructure:"database"`
}

func main() {

	viper := viper.New()
	viper.AddConfigPath("configs")
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("ERROR READING CONFIG FILE: %v", err))
	}

	// read server config
	fmt.Println("Server Port: ", viper.GetInt("server.port"))

	// config struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("ERROR UNMARSHALLING CONFIG: %v", err))
	}

	fmt.Println("Database: ", config.Database)
}
