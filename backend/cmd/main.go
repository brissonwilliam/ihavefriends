package cmd

import (
	"fmt"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// RootCommand returns the main command, properly configured
func RootCommand() *cobra.Command {
	return &cobra.Command{
		PreRun: initConfig,
		Short: "Start the backend web server for ihavefriends",
		Run: bootWebServer,
	}
}

func initConfig(cmd *cobra.Command, args []string) {
	viper.AutomaticEnv() // AutomaticEnv overrides everything
	viper.AddConfigPath(".")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func bootWebServer(cmd *cobra.Command, args []string) {
	err := server.NewWebServer().Start()
	if err != nil {
		fmt.Println("Could not start web server: " + err.Error())
		os.Exit(-1)
	}
}