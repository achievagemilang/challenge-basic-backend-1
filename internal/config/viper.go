package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var Viper *viper.Viper

// NewViper is a function to load config from config.json
// You can change the implementation, for example load from env file, consul, etcd, etc
func NewViper() *viper.Viper {
	config := viper.New()
	Viper = config

	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("./../../../../")
	config.AddConfigPath("./../../../")
	config.AddConfigPath("./../../")
	config.AddConfigPath("./../")
	config.AddConfigPath("./")

	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	return config
}
