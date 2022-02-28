package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type CompleteConfig struct {
	DB     DB
	Logger Logger
	Web    Web
}

func GetConfig() CompleteConfig {
	return CompleteConfig{
		DB:     GetDB(),
		Logger: GetLogger(),
		Web:    GetWeb(),
	}
}

func InitConfig(cmd *cobra.Command, args []string) {
	viper.AutomaticEnv() // AutomaticEnv overrides everything
	viper.AddConfigPath(".")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Errorf("Fatal error config file: %s \n", err)
		}
	}
}
