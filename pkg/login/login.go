package login

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"riedmann.dev/pocket-cli/pkg/util"
)

const redirectUri = "https://riedmann.dev"

func AuthorizeApp(appId string) {
	creds, err := ReadStoredCredentials()
	if err == nil && creds.RequestToken != "" {
		log.Println("Already authorized")
		return
	}

	reqToken := getRequestToken(appId)
	authUrl := fmt.Sprintf("https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s", reqToken, redirectUri)
	util.OpenInBrowser(authUrl)
}

func getRequestToken(appId string) string {
	payload := url.Values{
		"consumer_key": {appId},
		"redirect_uri": {redirectUri},
	}

	res, err := http.Post("https://getpocket.com/v3/oauth/request", "application/x-www-form-urlencoded", strings.NewReader(payload.Encode()))
	if err != nil || !util.IsHttpSuccess(res.StatusCode) {
		log.Printf("Failed to make http request (%v)\n", res.Status)
		panic(err)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		println("Failed to read http response")
		panic(err)
	}

	split := strings.Split(string(body), "=")
	if len(split) != 2 && split[0] != "code" {
		err := fmt.Errorf("unexpected API response: %v", body)
		panic(err)
	}
	token := split[1]

	StoreCredentials(credentials{RequestToken: token})

	return token
}

func GetAccessToken(appId string) string {
	creds, err := ReadStoredCredentials()
	if err == nil && len(creds.AccessToken) > 0 {
		return creds.AccessToken
	}
	if len(creds.AccessToken) == 0 && len(creds.RequestToken) == 0 {
		log.Fatalln("Application is not authorized - please run 'pocket-cli login'!")
	}
	log.Println("Did not find stored access token - requesting new one.\n\tThis should only happen once after login.")
	return getAccessToken(appId, creds.RequestToken)
}

type accessTokenRequest struct {
	ConsumerKey string `json:"consumer_key"`
	RequestCode string `json:"code"`
}

func getAccessToken(appId string, reqCode string) string {
	payload := accessTokenRequest{
		ConsumerKey: appId,
		RequestCode: reqCode,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		println("Failed to marshal http request body")
		panic(err)
	}

	res, err := http.Post("https://getpocket.com/v3/oauth/authorize", "application/json", bytes.NewBuffer(b))

	log.Println(string(b))
	if err == nil && res.StatusCode == 403 {
		log.Println("Failed to request AccessCode - if this persist please re-authorize using login --reset")
		panic(err)
	}

	if err != nil || !util.IsHttpSuccess(res.StatusCode) {
		log.Printf("Failed to make http request (%v)\n", res.Status)
		panic(err)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		println("Failed to read http response")
		panic(err)
	}

	token, user, err := parseAccessTokenResponse(string(body))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Acquired new access token for user %s\n", user)

	StoreCredentials(credentials{
		AccessToken:  token,
		RequestToken: reqCode,
	})

	return token
}

func parseAccessTokenResponse(resp string) (accessToken string, forUser string, err error) {
	split := strings.Split(resp, "&")
	if len(split) != 2 {
		return "", "", fmt.Errorf("unexpected API response: %v", resp)
	}
	tokenSplit := strings.Split(split[0], "=")
	if len(tokenSplit) != 2 && tokenSplit[0] != "access_token" {
		return "", "", fmt.Errorf("unexpected API response: %v", resp)
	}
	accessToken = tokenSplit[1]

	userSplit := strings.Split(split[1], "=")
	if len(userSplit) != 2 && userSplit[0] != "username" {
		return "", "", fmt.Errorf("unexpected API response: %v", resp)
	}
	forUser = userSplit[1]

	return accessToken, forUser, err
}
