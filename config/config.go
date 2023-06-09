package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Schema struct {
	Environment string `mapstructure:"environment"`
	Port        int    `mapstructure:"port"`
	AuthSecret  string `mapstructure:"auth_secret"`
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
	TestingEnv    = "testing"
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
	// Handle errors reading the config file
	if err != nil && viper.GetString("environment") != TestingEnv {
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
