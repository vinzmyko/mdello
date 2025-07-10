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

// We should add it so that the input shouldn't be boardName but the actual struct, when I add the types later
// I would then need to automatically create the apiUrl based on if a param is nil or not
func (t TrelloClient) UpdateBoard(boardId, newBoardName string) error {
	encodedName := url.QueryEscape(newBoardName)
	url := fmt.Sprintf("%s/boards/%s?name=%s&key=%s&token=%s", t.baseUrl, boardId, encodedName, t.apiKey, t.token)

	req, e := http.NewRequest("PUT", url, nil)
	if e != nil {
		log.Fatal(e)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(e)
	}

	fmt.Printf("%d\n", response.StatusCode)

	return nil
}

// should be something like what list/board/card you want to update and then the struct with all the information
func (t TrelloClient) UpdateList(listId, newListName string) error {
	encodedName := url.QueryEscape(newListName)
	url := fmt.Sprintf("%s/lists/%s?name=%s&key=%s&token=%s", t.baseUrl, listId, encodedName, t.apiKey, t.token)

	req, e := http.NewRequest("PUT", url, nil)
	if e != nil {
		log.Fatal(e)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	fmt.Printf("%d\n", response.StatusCode)

	return nil
}

func (t TrelloClient) UpdateCard(cardId, newCardName string) error {
	encodedName := url.QueryEscape(newCardName)
	url := fmt.Sprintf("%s/cards/%s?name=%s&key=%s&token=%s", t.baseUrl, cardId, encodedName, t.apiKey, t.token)

	req, e := http.NewRequest("PUT", url, nil)
	if e != nil {
		log.Fatal(e)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	fmt.Printf("%d\n", response.StatusCode)

	return nil
}

func (t TrelloClient) DeleteBoard(boardId string) error {
	url := fmt.Sprintf("%s/boards/%s?key=%s&token=%s", t.baseUrl, boardId, t.apiKey, t.token)

	req, e := http.NewRequest("DELETE", url, nil)
	if e != nil {
		log.Fatal(e)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	fmt.Printf("%d\n", response.StatusCode)

	return nil
}

// this has a value key for if you want to archive it or not
func (t TrelloClient) ArchiveList(listId string, setArchive bool) error {
	url := fmt.Sprintf("%s/lists/%s/closed?value=%t&key=%s&token=%s", t.baseUrl, listId, setArchive, t.apiKey, t.token)

	req, e := http.NewRequest("PUT", url, nil)
	if e != nil {
		log.Fatal(e)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	fmt.Printf("%d\n", response.StatusCode)

	return nil
}

func (t TrelloClient) DeleteCard(cardId string) error {
	url := fmt.Sprintf("%s/cards/%s?key=%s&token=%s", t.baseUrl, cardId, t.apiKey, t.token)

	req, e := http.NewRequest("DELETE", url, nil)
	if e != nil {
		log.Fatal(e)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	fmt.Printf("%d\n", response.StatusCode)

	return nil
}
