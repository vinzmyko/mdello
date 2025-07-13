package main

import (
	"fmt"
	"github.com/vinzmyko/mdello/trello"
	"log"
	"strings"
)

func main() {
	if strings.TrimSpace(trelloAPIKey) == "" {
		log.Fatal("API key not set in secrets.go")
	}
	config, err := loadConfig()
	if err != nil {
		fmt.Println("No config found. Please run 'mdello init'.")
		Execute()
		return
	}
	Execute()

	trelloClient, err := trello.NewTrelloClient(trelloAPIKey, config.Token)
	if err != nil {
		log.Fatal(err)
	}

	boards, err := trelloClient.GetBoards()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d boards\n", len(boards))

	// Replace with desired board id
	boardId := "686e846f2aee13b00660b241"
	lists, err := trelloClient.GetLists(boardId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d lists\n", len(lists))

	// Replace with desired list id
	listId := "686e846f2aee13b00660b281"
	cards, err := trelloClient.GetCards(listId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d cards\n", len(cards))
}
