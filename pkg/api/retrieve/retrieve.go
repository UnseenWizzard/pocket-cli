package retrieve

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
	Sort        string `json:"sort"`
	Count       int    `json:"count"`
	Offset      int    `json:"offset"`
}

type ResponsePayload struct {
	Status int                `json:"status"`
	List   map[string]Article `json:"list"`
}

type Article struct {
	Id         string `json:"item_id"`
	Title      string `json:"resolved_title"`
	GivenTitle string `json:"given_title"`
	Url        string `json:"given_url"`
	Excerpt    string `json:"excerpt"`
	ReadTime   int    `json:"time_to_read"`
}

func RetrieveUnread(consumerKey string, accessToken string, count int, offset int) (ResponsePayload, error) {
	return retrieveUnread("https://getpocket.com/v3/get", consumerKey, accessToken, count, offset)
}

func retrieveUnread(url string, consumerKey string, accessToken string, count int, offset int) (ResponsePayload, error) {
	payload := requestPayload{
		ConsumerKey: consumerKey,
		AccessToken: accessToken,
		Sort:        "newest",
		Count:       count,
		Offset:      offset,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return ResponsePayload{}, fmt.Errorf("failed to marshal http request body: %w", err)
	}

	log.Printf("Fetching %v articles from API - offset:%v, sort:%v", payload.Count, payload.Offset, payload.Sort)

	res, err := http.Post(url, "application/json", bytes.NewBuffer(b)) //nolint:gosec
	if err != nil || (res.StatusCode < 200 || res.StatusCode >= 300) {
		if err == nil {
			err = fmt.Errorf("%s", res.Status)
		}
		return ResponsePayload{}, fmt.Errorf("failed to make http request: %w", err)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return ResponsePayload{}, fmt.Errorf("failed to read http response: %w", err)
	}

	var list ResponsePayload

	err = json.Unmarshal(body, &list)
	if err != nil {
		return ResponsePayload{}, fmt.Errorf("failed to unmarshal http response: %w", err)
	}

	return list, nil
}
