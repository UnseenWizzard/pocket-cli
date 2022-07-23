package login

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/UnseenWizzard/pocket-cli/pkg/util"
)

const requestTokenApi = "https://getpocket.com/v3/oauth/request" //nolint:gosec
const authorizeApi = "https://getpocket.com/v3/oauth/authorize"

func CreateRequestToken(appId string, redirectUri string) (string, error) {
	return createRequestToken(requestTokenApi, appId, redirectUri)
}

func createRequestToken(apiUrl string, appId string, redirectUri string) (string, error) {
	payload := url.Values{
		"consumer_key": {appId},
		"redirect_uri": {redirectUri},
	}

	res, err := http.Post(apiUrl, "application/x-www-form-urlencoded", strings.NewReader(payload.Encode())) //nolint:gosec
	if err != nil || !util.IsHttpSuccess(res.StatusCode) {
		if err == nil {
			err = fmt.Errorf("%s", res.Status)
		}
		return "", fmt.Errorf("failed to make http request: %w", err)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return "", fmt.Errorf("failed to read http response: %w", err)
	}

	split := strings.Split(string(body), "=")
	if len(split) != 2 || split[0] != "code" {
		return "", fmt.Errorf("unexpected API response: %v", body)
	}
	token := split[1]

	return token, nil
}

type accessTokenRequest struct {
	ConsumerKey string `json:"consumer_key"`
	RequestCode string `json:"code"`
}

func CreateAccessToken(appId string, reqCode string) (string, error) {
	return createAccessToken(authorizeApi, appId, reqCode)
}

func createAccessToken(apiUrl string, appId string, reqCode string) (string, error) {
	payload := accessTokenRequest{
		ConsumerKey: appId,
		RequestCode: reqCode,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal http request body: %w", err)
	}

	res, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(b)) //nolint:gosec

	if err == nil && res.StatusCode == 403 {
		return "", fmt.Errorf("failed to request AccessCode - if this persist please re-authorize using login --reset")
	}

	if err != nil || !util.IsHttpSuccess(res.StatusCode) {
		if err == nil {
			err = fmt.Errorf("%s", res.Status)
		}
		return "", fmt.Errorf("failed to make http request: %w", err)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return "", fmt.Errorf("failed to read http response: %w", err)
	}

	token, user, err := parseAccessTokenResponse(string(body))
	if err != nil {
		return "", err
	}

	fmt.Printf("Acquired new access token for user %s\n", user)

	return token, nil
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
