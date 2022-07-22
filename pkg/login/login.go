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

	"github.com/UnseenWizzard/pocket-cli/pkg/util"
)

const redirectUri = "https://riedmann.dev"
const requestTokenApi = "https://getpocket.com/v3/oauth/request"
const authorizeApi = "https://getpocket.com/v3/oauth/authorize"

func AuthorizeApp(appId string) {
	creds, err := readStoredCredentials()
	if err == nil && creds.RequestToken != "" {
		log.Println("Already authorized")
		return
	}

	reqToken := getRequestToken(requestTokenApi, appId, storeCredentials)
	authUrl := fmt.Sprintf("https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s", reqToken, redirectUri)
	log.Println("Please authorize app in browser - then run commands like `pocket-cli list`")
	util.OpenInBrowser(authUrl)
}

func getRequestToken(apiUrl string, appId string, storeCredentialsFn func(credentials)) string {
	payload := url.Values{
		"consumer_key": {appId},
		"redirect_uri": {redirectUri},
	}

	res, err := http.Post(apiUrl, "application/x-www-form-urlencoded", strings.NewReader(payload.Encode()))
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

	storeCredentialsFn(credentials{RequestToken: token})

	return token
}

func GetAccessToken(appId string) string {
	creds, err := readStoredCredentials()
	if err == nil && len(creds.AccessToken) > 0 {
		return creds.AccessToken
	}
	if len(creds.AccessToken) == 0 && len(creds.RequestToken) == 0 {
		log.Fatalln("Application is not authorized - please run 'pocket-cli login'!")
	}
	log.Println("Did not find stored access token - requesting new one. (This should only happen once after login)")
	return getAccessToken(authorizeApi, appId, creds.RequestToken, storeCredentials)
}

type accessTokenRequest struct {
	ConsumerKey string `json:"consumer_key"`
	RequestCode string `json:"code"`
}

func getAccessToken(apiUrl string, appId string, reqCode string, storeCredentialsFn func(credentials)) string {
	payload := accessTokenRequest{
		ConsumerKey: appId,
		RequestCode: reqCode,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		println("Failed to marshal http request body")
		panic(err)
	}

	res, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(b))

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

	storeCredentialsFn(credentials{
		AccessToken:  token,
		RequestToken: reqCode,
	})

	return token
}

func parseAccessTokenResponse(resp string) (accessToken string, forUser string, err error) {
	split := strings.Split(resp, "&")
	if len(split) < 2 {
		return "", "", fmt.Errorf("unexpected API response: %v", resp)
	}

	for _, entry := range split {
		kvSplit := strings.Split(entry, "=")
		if len(kvSplit) != 2 {
			return "", "", fmt.Errorf("unable parse API response entry: %v (%v)", entry, resp)
		}
		key := kvSplit[0]
		val := kvSplit[1]
		if key == "access_token" {
			accessToken = val
		}
		if key == "username" {
			forUser, _ = url.QueryUnescape(val)
		}
	}
	if len(accessToken) == 0 {
		return "", "", fmt.Errorf("api response did not contain access_token")
	}
	if len(forUser) == 0 {
		return "", "", fmt.Errorf("api response did not contain username")
	}

	return accessToken, forUser, err
}
