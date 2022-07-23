package auth

import (
	"fmt"
	"log"

	"github.com/UnseenWizzard/pocket-cli/pkg/api/login"
	"github.com/UnseenWizzard/pocket-cli/pkg/auth/credentials"
	"github.com/UnseenWizzard/pocket-cli/pkg/util"
)

const redirectUri = "https://riedmann.dev"
const authUrlTemplate = "https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s"

func AuthorizeApp(appId string) {
	err := authorizeApp(appId, credentials.ReadStoredCredentials, credentials.StoreCredentials, login.CreateRequestToken, util.OpenInBrowser)
	if err != nil {
		log.Fatalf("Encountered error authorizing app: %v", err)
	}
}

func authorizeApp(
	appId string,
	readStoredCredsFn func() (credentials.Credentials, error),
	storeCredsFn func(credentials.Credentials),
	createRequestTokenFn func(string, string) (string, error),
	openBrowserFn func(string) error,
) error {
	creds, err := readStoredCredsFn()
	if err == nil && creds.RequestToken != "" {
		log.Println("Already authorized")
		return nil
	}

	reqToken, err := createRequestTokenFn(appId, redirectUri)
	if err != nil {
		return err
	}

	storeCredsFn(credentials.Credentials{RequestToken: reqToken})

	authUrl := fmt.Sprintf(authUrlTemplate, reqToken, redirectUri)
	log.Println("Please authorize app in browser - then run commands like `pocket-cli list`")

	err = openBrowserFn(authUrl)
	if err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}

	return nil
}

func GetAccessToken(appId string) string {
	token, err := getAccessToken(appId, credentials.ReadStoredCredentials, credentials.StoreCredentials, login.CreateAccessToken)
	if err != nil {
		log.Fatalf("Encountered error getting access token: %v", err)
	}
	return token
}

func getAccessToken(
	appId string,
	readStoredCredsFn func() (credentials.Credentials, error),
	storeCredsFn func(credentials.Credentials),
	createAccessTokenFn func(string, string) (string, error),
) (string, error) {
	creds, err := readStoredCredsFn()
	if err == nil && len(creds.AccessToken) > 0 {
		return creds.AccessToken, nil
	}
	if len(creds.AccessToken) == 0 && len(creds.RequestToken) == 0 {
		return "", fmt.Errorf("application is not authorized - please run 'pocket-cli login'")
	}
	log.Println("Did not find stored access token - requesting new one. (This should only happen once after login)")

	token, err := createAccessTokenFn(appId, creds.RequestToken)
	if err != nil {
		return "", err
	}

	storeCredsFn(credentials.Credentials{
		RequestToken: creds.RequestToken,
		AccessToken:  token,
	})

	return token, nil
}
