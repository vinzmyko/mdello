package trello

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func paramsToURLValues(params any) (url.Values, error) {
	values := url.Values{}
	structData := reflect.ValueOf(params)

	if structData.Kind() == reflect.Ptr {
		structData = structData.Elem()
	}

	if structData.Kind() != reflect.Struct {
		return nil, fmt.Errorf("params must be a struct. Got %T", params)
	}

	structInfo := structData.Type()
	for i := 0; i < structData.NumField(); i++ {
		// structData contains the actual values in a list, structInfo contains the corresponding types
		field := structData.Field(i)
		fieldType := structInfo.Field(i)

		tag := fieldType.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}

		key := strings.Split(tag, ",")[0]
		if key == "id" {
			continue
		}

		// If field is nil we do not want to add it to strValue
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				continue
			}
			field = field.Elem()
		}

		var strValue string
		switch field.Kind() {
		case reflect.String:
			strValue = field.String()
		case reflect.Bool:
			strValue = strconv.FormatBool(field.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			strValue = strconv.FormatInt(field.Int(), 10)
		default:
			continue
		}

		values.Set(key, strValue)
	}

	return values, nil
}

func (t *TrelloClient) doRequest(method, path string, queryParams url.Values, result any) error {
	fullURL, err := url.JoinPath(t.baseUrl, path)
	if err != nil {
		return fmt.Errorf("failed to create URL path: %w", err)
	}

	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	query := req.URL.Query()
	query.Set("key", t.apiKey)
	query.Set("token", t.token)

	for key, vals := range queryParams {
		for _, val := range vals {
			query.Add(key, val)
		}
	}
	req.URL.RawQuery = query.Encode()

	response, err := t.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("network error: %w", err)
	}
	defer response.Body.Close()

	if err := handleHTTPResponse(response); err != nil {
		return fmt.Errorf("failed to handle http response: %w", err)
	}

	if result != nil {
		if err = json.NewDecoder(response.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

func (t *TrelloClient) GetBoards() ([]Board, error) {
	var boards []Board

	err := t.doRequest("GET", "/members/me/boards", nil, &boards)
	if err != nil {
		return nil, fmt.Errorf("failed to get user boards: %w", err)
	}
	return boards, nil
}

func (t *TrelloClient) GetLists(boardId string) ([]List, error) {
	if boardId == "" {
		return nil, errors.New("boardId is required to get lists")
	}

	var lists []List

	path := fmt.Sprintf("/boards/%s/lists", boardId)
	err := t.doRequest("GET", path, nil, &lists)
	if err != nil {
		return nil, fmt.Errorf("failed to get lists for board %s: %w", boardId, err)
	}
	return lists, nil
}

func (t *TrelloClient) GetCards(listId string) ([]Card, error) {
	if listId == "" {
		return nil, errors.New("listId is required to get cards")
	}
	var cards []Card

	path := fmt.Sprintf("/lists/%s/cards", listId)
	err := t.doRequest("GET", path, nil, &cards)
	if err != nil {
		return nil, fmt.Errorf("failed to get cards for list %s: %w", listId, err)
	}
	return cards, nil
}

func (t *TrelloClient) CreateBoard(params *CreateBoardParams) (*Board, error) {
	if params == nil || params.Name == "" {
		return nil, errors.New("CreateBoardParams with a Name is required")
	}

	queryParams, err := paramsToURLValues(params)
	if err != nil {
		return nil, fmt.Errorf("could not process create board params: %w", err)
	}

	var board Board

	err = t.doRequest("POST", "/boards", queryParams, &board)
	if err != nil {
		return nil, fmt.Errorf("failed to create trello board: %w", err)
	}
	return &board, nil
}

func (t *TrelloClient) CreateList(params *CreateListParams) (*List, error) {
	if params == nil || params.Name == "" || params.IdBoard == "" {
		return nil, errors.New("CreateListParams required Name and IdBoard")
	}

	queryParams, err := paramsToURLValues(params)
	if err != nil {
		return nil, fmt.Errorf("could not process create list params: %w", err)
	}

	var list List

	err = t.doRequest("POST", "/lists", queryParams, &list)
	if err != nil {
		return nil, fmt.Errorf("failed to create trello list: %w", err)
	}
	return &list, nil
}

func (t *TrelloClient) CreateCard(params *CreateCardParams) (*Card, error) {
	if params == nil || params.IdList == "" {
		return nil, errors.New("CreateCardParams with a IdList is required")
	}

	queryParams, err := paramsToURLValues(params)
	if err != nil {
		return nil, fmt.Errorf("could not process create card params: %w", err)
	}

	var card Card

	err = t.doRequest("POST", "/cards", queryParams, &card)
	if err != nil {
		return nil, fmt.Errorf("failed to create trello card: %w", err)
	}
	return &card, nil
}

func (t *TrelloClient) UpdateBoard(params *UpdateBoardParams) (*Board, error) {
	if params == nil || params.ID == "" {
		return nil, errors.New("UpdateBoardParams with a valid ID is required")
	}

	queryParams, err := paramsToURLValues(params)
	if err != nil {
		return nil, fmt.Errorf("could not process update board params: %w", err)
	}

	path := fmt.Sprintf("/boards/%s", params.ID)

	var board Board
	err = t.doRequest("PUT", path, queryParams, &board)
	if err != nil {
		return nil, fmt.Errorf("failed to update trello board: %w", err)
	}

	return &board, nil
}

func (t *TrelloClient) UpdateList(params *UpdateListParams) (*List, error) {
	if params == nil || params.ID == "" {
		return nil, errors.New("UpdateListParams with a valid ID is required")
	}

	queryParams, err := paramsToURLValues(params)
	if err != nil {
		return nil, fmt.Errorf("could not process update list params: %w", err)
	}

	path := fmt.Sprintf("/lists/%s", params.ID)

	var list List
	err = t.doRequest("PUT", path, queryParams, &list)
	if err != nil {
		return nil, fmt.Errorf("failed to update trello list: %w", err)
	}

	return &list, nil
}

func (t *TrelloClient) UpdateCard(params *UpdateCardParams) (*Card, error) {
	if params == nil || params.ID == "" {
		return nil, errors.New("UpdateCardParams with a valid ID is required")
	}

	queryParams, err := paramsToURLValues(params)
	if err != nil {
		return nil, fmt.Errorf("could not process update card params: %w", err)
	}

	path := fmt.Sprintf("/cards/%s", params.ID)

	var card Card
	err = t.doRequest("PUT", path, queryParams, &card)
	if err != nil {
		return nil, fmt.Errorf("failed to update trello card: %w", err)
	}

	return &card, nil
}

func (t *TrelloClient) DeleteBoard(boardId string) error {
	if boardId == "" {
		return errors.New("boardId is required to delete a trello board")
	}
	path := fmt.Sprintf("/boards/%s", boardId)
	err := t.doRequest("DELETE", path, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete trello board %s: %w", boardId, err)
	}

	return nil
}

// this has a value key for if you want to archive it or not
func (t *TrelloClient) ArchiveList(params *ArchiveListParams) (*List, error) {
	if params == nil || params.ID == "" || params.Value == nil {
		return nil, errors.New("ArchiveListParams requires an ID and a Value")
	}

	queryParams, err := paramsToURLValues(params)
	if err != nil {
		return nil, fmt.Errorf("could not process archive list params: %w", err)
	}

	path := fmt.Sprintf("/lists/%s/closed", params.ID)

	var list List
	err = t.doRequest("PUT", path, queryParams, &list)
	if err != nil {
		return nil, fmt.Errorf("failed to update archive status for list %s: %w", params.ID, err)
	}

	return &list, nil
}

func (t *TrelloClient) DeleteCard(cardId string) error {
	if cardId == "" {
		return errors.New("cardId is required to delete a trello board")
	}

	path := fmt.Sprintf("/cards/%s", cardId)
	err := t.doRequest("DELETE", path, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete trello card: %w", err)
	}

	return nil
}
