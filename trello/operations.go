package trello

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func (t TrelloClient) GetBoards() ([]Board, error) {
	url := fmt.Sprintf("%s/members/me/boards?key=%s&token=%s", t.baseUrl, t.apiKey, t.token)

	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	defer r.Body.Close()

	var boards []Board

	e = json.NewDecoder(r.Body).Decode(&boards)
	if e != nil {
		return nil, e
	}

	return boards, nil
}

func (t TrelloClient) GetLists(boardId string) ([]List, error) {
	url := fmt.Sprintf("%s/boards/%s/lists?key=%s&token=%s", t.baseUrl, boardId, t.apiKey, t.token)

	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	defer r.Body.Close()

	var lists []List

	e = json.NewDecoder(r.Body).Decode(&lists)
	if e != nil {
		return nil, e
	}

	return lists, nil
}

func (t TrelloClient) GetCards(listId string) ([]Card, error) {
	url := fmt.Sprintf("%s/lists/%s/cards?key=%s&token=%s", t.baseUrl, listId, t.apiKey, t.token)

	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	defer r.Body.Close()

	var cards []Card

	e = json.NewDecoder(r.Body).Decode(&cards)
	if e != nil {
		return nil, e
	}

	return cards, nil
}

func (t TrelloClient) CreateBoard(boardName string) (*Board, error) {
	encodedName := url.QueryEscape(boardName)
	url := fmt.Sprintf("%s/boards/?name=%s&key=%s&token=%s", t.baseUrl, encodedName, t.apiKey, t.token)

	r, e := http.Post(url, "application/json", nil)
	if e != nil {
		return nil, e
	}
	defer r.Body.Close()

	var board Board
	e = json.NewDecoder(r.Body).Decode(&board)
	if e != nil {
		return nil, e
	}

	return &board, nil
}

func (t TrelloClient) CreateList(listName, boardId string) (*List, error) {
	encodedName := url.QueryEscape(listName)
	url := fmt.Sprintf("%s/lists/?name=%s&idBoard=%s&key=%s&token=%s", t.baseUrl, encodedName, boardId, t.apiKey, t.token)

	r, e := http.Post(url, "application/json", nil)
	if e != nil {
		return nil, e
	}
	defer r.Body.Close()

	var list List
	e = json.NewDecoder(r.Body).Decode(&list)
	if e != nil {
		return nil, e
	}

	return &list, nil
}

func (t TrelloClient) CreateCard(cardName, listId string) (*Card, error) {
	encodedName := url.QueryEscape(cardName)
	url := fmt.Sprintf("%s/cards/?name=%s&idList=%s&key=%s&token=%s", t.baseUrl, encodedName, listId, t.apiKey, t.token)

	r, e := http.Post(url, "application/json", nil)
	if e != nil {
		return nil, e
	}
	defer r.Body.Close()

	var card Card
	e = json.NewDecoder(r.Body).Decode(&card)
	if e != nil {
		return nil, e
	}

	return &card, nil
}

// We should add it so that the input shouldn't be boardName but the actual struct, when I add the types later
// I would then need to automatically create the apiUrl based on if a param is nil or not
func (t TrelloClient) UpdateBoard(boardId, newBoardName string) (*Board, error) {
	encodedName := url.QueryEscape(newBoardName)
	url := fmt.Sprintf("%s/boards/%s?name=%s&key=%s&token=%s", t.baseUrl, boardId, encodedName, t.apiKey, t.token)

	req, e := http.NewRequest("PUT", url, nil)
	if e != nil {
		return nil, e
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var board Board
	err = json.NewDecoder(response.Body).Decode(&board)
	if err != nil {
		return nil, err
	}

	return &board, nil
}

// should be something like what list/board/card you want to update and then the struct with all the information
func (t TrelloClient) UpdateList(listId, newListName string) (*List, error) {
	encodedName := url.QueryEscape(newListName)
	url := fmt.Sprintf("%s/lists/%s?name=%s&key=%s&token=%s", t.baseUrl, listId, encodedName, t.apiKey, t.token)

	req, e := http.NewRequest("PUT", url, nil)
	if e != nil {
		return nil, e
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var list List
	err = json.NewDecoder(response.Body).Decode(&list)
	if err != nil {
		return nil, err
	}

	return &list, nil
}

func (t TrelloClient) UpdateCard(cardId, newCardName string) (*Card, error) {
	encodedName := url.QueryEscape(newCardName)
	url := fmt.Sprintf("%s/cards/%s?name=%s&key=%s&token=%s", t.baseUrl, cardId, encodedName, t.apiKey, t.token)

	req, e := http.NewRequest("PUT", url, nil)
	if e != nil {
		return nil, e
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var card Card
	err = json.NewDecoder(response.Body).Decode(&card)
	if err != nil {
		return nil, err
	}

	return &card, nil
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
