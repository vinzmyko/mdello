package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	apiKey := os.Getenv("TRELLO_API_KEY")
	token := os.Getenv("TRELLO_TOKEN")

	if strings.TrimSpace(apiKey) == "" || strings.TrimSpace(token) == "" {
		log.Fatal("Set TRELLO_API_KEY and TRELLO environment variables")
	}

	url := fmt.Sprintf("https://api.trello.com/1/members/me/boards?key=%s&token=%s", apiKey, token)

	r, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer r.Body.Close()

	var boards []map[string]any

	e = json.NewDecoder(r.Body).Decode(&boards)
	if e != nil {
		log.Fatal(e)
	}

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
