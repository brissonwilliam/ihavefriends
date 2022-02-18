package main

import (
	"github.com/brissonwilliam/ihavefriends/backend/cmd"
	"os"
)

func main() {
	if err := cmd.RootCommand().Execute(); err != nil {
		os.Exit(-1)
	}
}
