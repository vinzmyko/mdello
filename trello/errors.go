package trello

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	StatusCode int    `json:"-"` // Don't try to unmarshal from JSON
	Message    string `json:"message,omitempty"`
	ErrorMsg   string `json:"error,omitempty"`
}

// Implementing the Error() method to satisfy the error interface. Called when needing to convert Error to string
func (e Error) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("Trello API Error %d: %s", e.StatusCode, e.Message)
	}
	if e.ErrorMsg != "" {
		return fmt.Sprintf("Trello API Error %d: %s", e.StatusCode, e.ErrorMsg)
	}
	return fmt.Sprintf("Trello API Error %d: %s", e.StatusCode, getGenericErrorMessage(e.StatusCode))
}

func (e Error) IsAuthError() bool {
	return e.StatusCode == http.StatusUnauthorized
}

func (e Error) IsNotFoundError() bool {
	return e.StatusCode == http.StatusNotFound
}

func (e Error) IsRateLimitError() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

func (e Error) IsBadRequestError() bool {
	return e.StatusCode == http.StatusBadRequest
}

func handleHTTPResponse(response *http.Response) error {
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return nil
	}

	var apiError Error
	err := json.NewDecoder(response.Body).Decode(&apiError)
	if err != nil {
		apiError = Error{
			StatusCode: response.StatusCode,
			Message:    getGenericErrorMessage(response.StatusCode),
		}
	}

	apiError.StatusCode = response.StatusCode

	return apiError
}

func getGenericErrorMessage(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "bad request - check your request parameters"
	case http.StatusUnauthorized:
		return "unauthorised - check your API key and token"
	case http.StatusForbidden:
		return "forbidden - you don't have permission for this action"
	case http.StatusNotFound:
		return "not found - the requested resource doesn't exist"
	case http.StatusConflict:
		return "conflict - the request conflicts with current state"
	case http.StatusTooManyRequests:
		return "too many requests - you've exceeded the rate limit"
	case http.StatusInternalServerError:
		return "internal server error - please try again later"
	case http.StatusServiceUnavailable:
		return "service unavailable - Trello may be experiencing issues"
	default:
		return fmt.Sprintf("HTTP %d error", statusCode)
	}
}
