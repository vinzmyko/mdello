package trello

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type TrelloClient struct {
	apiKey     string
	token      string
	baseUrl    string
	httpClient *http.Client
}

func NewTrelloClient(apiKey, token string) (*TrelloClient, error) {
	trelloClient := TrelloClient{
		token:   token,
		apiKey:  apiKey,
		baseUrl: "https://api.trello.com/1",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	err := trelloClient.HealthCheck()
	if err != nil {
		return nil, fmt.Errorf("Invalid Trello Token: %w", err)
	}

	return &trelloClient, nil
}

func (t *TrelloClient) SetTimeout(timeout time.Duration) {
	t.httpClient.Timeout = timeout
}

func (t TrelloClient) HealthCheck() error {
	url := fmt.Sprintf("%s/members/me/?key=%s&token=%s", t.baseUrl, t.apiKey, t.token)

	r, e := t.httpClient.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		err := fmt.Sprintf("API unavailable: status %d", r.StatusCode)
		log.Fatal(err)
	}

	return nil
}
