package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type TrelloClient struct {
	apiKey  string
	token   string
	baseUrl string
}

func NewTrelloClient(apiKey, token string) TrelloClient {
	return TrelloClient{
		token:   token,
		apiKey:  apiKey,
		baseUrl: "https://api.trello.com/1",
	}
}

func (t TrelloClient) GetBoards() []map[string]any {
	url := fmt.Sprintf("%s/members/me/boards?key=%s&token=%s", t.baseUrl, t.apiKey, t.token)

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

	return boards
}

func (t TrelloClient) GetLists(boardId string) []map[string]any {
	url := fmt.Sprintf("%s/boards/%s/lists?key=%s&token=%s", t.baseUrl, boardId, t.apiKey, t.token)

	r, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer r.Body.Close()

	var lists []map[string]any

	e = json.NewDecoder(r.Body).Decode(&lists)
	if e != nil {
		log.Fatal(e)
	}

	return lists
}

func (t TrelloClient) GetCards(listId string) []map[string]any {
	url := fmt.Sprintf("%s/lists/%s/cards?key=%s&token=%s", t.baseUrl, listId, t.apiKey, t.token)

	r, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer r.Body.Close()

	var cards []map[string]any

	e = json.NewDecoder(r.Body).Decode(&cards)
	if e != nil {
		log.Fatal(e)
	}

	return cards
}
