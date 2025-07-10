package main

import (
	"fmt"
	"github.com/vinzmyko/mdello/trello"
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

	trelloClient := trello.NewTrelloClient(apiKey, token)

	boards := trelloClient.GetBoards()
	fmt.Printf("Found %d boards\n", len(boards))

	// Replace with desired board id
	boardId := "686e846f2aee13b00660b241"
	lists := trelloClient.GetLists(boardId)
	fmt.Printf("Found %d lists\n", len(lists))

	// Replace with desired list id
	listId := "686e846f2aee13b00660b281"
	cards := trelloClient.GetCards(listId)
	fmt.Printf("Found %d cards\n", len(cards))
}
