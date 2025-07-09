package main

import (
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
}
