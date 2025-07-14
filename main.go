package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vinzmyko/mdello/cli"
	"github.com/vinzmyko/mdello/config"
)

func main() {
	if strings.TrimSpace(trelloAPIKey) == "" {
		log.Fatal("API key not set in secrets.go")
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		if len(os.Args) > 1 && os.Args[1] == "init" {
			cli.Execute(trelloAPIKey, nil)
		} else {
			fmt.Println("No config found. Please run 'mdello init'.")
			cli.Execute(trelloAPIKey, nil)
		}
		return
	}

	cli.Execute(trelloAPIKey, cfg)
}
