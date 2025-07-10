package trello

import (
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

func (t TrelloClient) HealthCheck() error {
	url := fmt.Sprintf("%s/members/me/?key=%s&token=%s", t.baseUrl, t.apiKey, t.token)

	r, e := http.Get(url)
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
