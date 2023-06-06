package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Schema struct {
	Environment string `mapstructure:"environment"`
	Port        int    `mapstructure:"port"`
	DatabaseURI string `mapstructure:"database_uri"`

	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
		Database int    `mapstructure:"database"`
	} `mapstructure:"redis"`

	Cache struct {
		Enable     bool `mapstructure:"enable"`
		ExpiryTime int  `mapstructure:"expiry_time"`
	} `mapstructure:"cache"`
}

var (
	ProductionEnv = "production"
	cfg           Schema
)

func init() {
	config := viper.New()
	config.SetConfigName("config")
	config.AddConfigPath(".")          // Look for config in current directory
	config.AddConfigPath("config/")    // Optionally look for config in the working directory.
	config.AddConfigPath("../config/") // Look for config needed for tests.
	config.AddConfigPath("../")        // Look for config needed for tests.

	config.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	config.AutomaticEnv()

	err := config.ReadInConfig() // Find and read the config file
	if err != nil {              // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = config.Unmarshal(&cfg)
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// fmt.Printf("Current Config: %+v", Config)
}

func GetConfig() *Schema {
	return &cfg
}
