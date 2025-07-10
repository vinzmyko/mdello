package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
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

func (t TrelloClient) CreateBoard(boardName string) error {
	encodedName := url.QueryEscape(boardName)
	url := fmt.Sprintf("%s/boards/?name=%s&key=%s&token=%s", t.baseUrl, encodedName, t.apiKey, t.token)

	r, e := http.Post(url, "application/json", nil)
	if e != nil {
		log.Fatal(e)
	}
	defer r.Body.Close()

	fmt.Printf("%d\n", r.StatusCode)

	return nil
}

func (t TrelloClient) CreateList(listName, boardId string) error {
	encodedName := url.QueryEscape(listName)
	url := fmt.Sprintf("%s/lists/?name=%s&idBoard=%s&key=%s&token=%s", t.baseUrl, encodedName, boardId, t.apiKey, t.token)

	r, e := http.Post(url, "application/json", nil)
	if e != nil {
		log.Fatal(e)
	}
	defer r.Body.Close()

	fmt.Printf("%d\n", r.StatusCode)

	return nil
}

func (t TrelloClient) CreateCard(cardName, listId string) error {
	encodedName := url.QueryEscape(cardName)
	url := fmt.Sprintf("%s/cards/?name=%s&idList=%s&key=%s&token=%s", t.baseUrl, encodedName, listId, t.apiKey, t.token)

	r, e := http.Post(url, "application/json", nil)
	if e != nil {
		log.Fatal(e)
	}
	defer r.Body.Close()

	fmt.Printf("%d\n", r.StatusCode)

	return nil
}
