package retrieve

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type requestPayload struct {
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
	Sort string `json:"sort"`
	Count int `json:"count"`
	Offset int `json:"offset"`
}

type ResponsePayload struct {
	Status int `json:"status"`
	Complete int `json:"complete"`
	List map[string]Article `json:"list"`
}

type Article struct {
	Title string `json:"resolved_title"`
	GivenTitle string `json:"given_title"`
	Url string `json:"given_url"`
	Excerpt string `json:"excerpt"`
	ReadTime int `json:"time_to_read"`
}

func RetrieveUnread(consumerKey string, accessToken string, count int, offset int) ResponsePayload {
	payload := requestPayload {
		ConsumerKey: consumerKey, 
		AccessToken: accessToken,
		Sort: "newest",
		Count: count,
		Offset: offset,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		println("Failed to marshal http request body")
		panic(err)
	}

	log.Printf("Fetching %v articles from API - offset:%v, sort:%v", payload.Count, payload.Offset, payload.Sort)

	res, err := http.Post("https://getpocket.com/v3/get", "application/json", bytes.NewBuffer(b))
	if err != nil || (res.StatusCode < 200 || res.StatusCode >= 300) {
		println("Failed to make http request")
		panic(err)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		println("Failed to read http response")
		panic(err)
	}

	var list ResponsePayload 

	err = json.Unmarshal(body, &list)
	if err != nil {
		println("Failed to unmarshal http response")
		panic(err)
	}

	return list
}