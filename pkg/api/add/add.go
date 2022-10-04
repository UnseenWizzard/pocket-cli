package add

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type requestPayload struct {
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
	Url         string `json:"url"`
}

type responsePayload struct {
	Status int `json:"status"`
}

const addApiUrl = "https://getpocket.com/v3/add"

func AddArticle(consumerKey string, accessToken string, articleUrl string) error {
	return addArticle(addApiUrl, consumerKey, accessToken, articleUrl)
}

func addArticle(apiUrl string, consumerKey string, accessToken string, articleUrl string) error {
	b, err := buildRequestPayload(consumerKey, accessToken, articleUrl)
	if err != nil {
		return fmt.Errorf("failed to build http request body: %w", err)
	}

	log.Printf("Adding new article (%v)...", articleUrl)

	res, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(b)) //nolint:gosec
	if err != nil || (res.StatusCode < 200 || res.StatusCode >= 300) {
		if err == nil {
			err = fmt.Errorf("%s", res.Status)
		}
		return fmt.Errorf("failed to make http request: %w", err)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to read http response: %w", err)
	}

	return checkResponseForErrors(body)
}

func buildRequestPayload(consumerKey string, accessToken string, articleUrl string) ([]byte, error) {
	payload := requestPayload{
		ConsumerKey: consumerKey,
		AccessToken: accessToken,
		Url:         articleUrl,
	}
	return json.Marshal(payload)
}

func checkResponseForErrors(responseBody []byte) error {
	var result responsePayload

	err := json.Unmarshal(responseBody, &result)
	if err != nil {
		return fmt.Errorf("failed to unmarshal http response: %w", err)
	}

	if result.Status != 1 {
		return fmt.Errorf("API failed to excute action")
	}

	return nil
}
