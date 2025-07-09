package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	apiKey := os.Getenv("TRELLO_API_KEY")
	token := os.Getenv("TRELLO_TOKEN")

	if strings.TrimSpace(apiKey) == "" || strings.TrimSpace(token) == "" {
		log.Fatal("Set TRELLO_API_KEY and TRELLO environment variables")
	}

	trello := NewTrelloClient(apiKey, token)
	boards := trello.GetBoards()

	fmt.Printf("Found %d boards:\n", len(boards))
	for i, board := range boards {
		fmt.Printf("\nBoard %d:\n", i+1)
		fmt.Printf("  Name: %v\n", board["name"])

		fmt.Println("  Keys:")
		for key := range board {
			fmt.Printf("    - %s\n", key)
		}
	}
}
