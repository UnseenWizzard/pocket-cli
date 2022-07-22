package modify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type ModifyFunction func(string, string, string) error

type simpleModifyAction struct {
	Action string `json:"action"`
	ItemId string `json:"item_id"`
}

type requestPayload struct {
	ConsumerKey string               `json:"consumer_key"`
	AccessToken string               `json:"access_token"`
	Actions     []simpleModifyAction `json:"actions"`
}

type responsePayload struct {
	ActionResults []bool `json:"action_results"`
	Status        int    `json:"status"`
}

const sendApiUrl = "https://getpocket.com/v3/send"

func ArchiveArticle(consumerKey string, accessToken string, articleId string) error {
	return makeModifyApiCall(sendApiUrl, consumerKey, accessToken, "archive", articleId)
}

func FavoriteArticle(consumerKey string, accessToken string, articleId string) error {
	return makeModifyApiCall(sendApiUrl, consumerKey, accessToken, "favorite", articleId)
}

func makeModifyApiCall(url string, consumerKey string, accessToken string, action string, articleId string) error {

	b, err := buildRequestPayload(consumerKey, accessToken, action, articleId)
	if err != nil {
		return fmt.Errorf("failed to build http request body: %w", err)
	}

	log.Printf("Marking article (%v) as %s-ed...", articleId, action)

	res, err := http.Post(url, "application/json", bytes.NewBuffer(b))
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

func buildRequestPayload(consumerKey string, accessToken string, action string, articleId string) ([]byte, error) {
	payload := requestPayload{
		ConsumerKey: consumerKey,
		AccessToken: accessToken,
		Actions: []simpleModifyAction{
			{
				action,
				articleId,
			},
		},
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
