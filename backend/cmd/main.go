package cmd

import (
	"fmt"
	"github.com/brissonwilliam/ihavefriends/backend/config"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/server"
	"github.com/spf13/cobra"
)

// RootCommand returns the main command, properly configured
func RootCommand() *cobra.Command {
	return &cobra.Command{
		PreRun: config.InitConfig,
		Short: "Start the backend web server for ihavefriends",
		Run: bootWebServer,
	}
}

func bootWebServer(cmd *cobra.Command, args []string) {
	err := server.NewWebServer().Start()
	if err != nil {
		fmt.Println("Could not start web server: " + err.Error())
	}
}