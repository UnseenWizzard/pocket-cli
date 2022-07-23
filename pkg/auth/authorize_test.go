package auth

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/UnseenWizzard/pocket-cli/pkg/auth/credentials"
)

const testRequestToken = "request-token"
const testAccessToken = "auth-token"

func Test_authorizeApp(t *testing.T) {

	noStoredCredentialsExistFn := func() (credentials.Credentials, error) {
		return credentials.Credentials{}, fmt.Errorf("No Creds Exist")
	}

	defaultCreateReqTokenFn := func(s1, s2 string) (string, error) {
		return testRequestToken, nil
	}

	var storedCreds credentials.Credentials
	defaultStoreFn := func(c credentials.Credentials) {
		storedCreds = c
	}

	var urlToOpen string
	defaultOpenBrowserFn := func(url string) error {
		urlToOpen = url
		return nil
	}

	type args struct {
		appId                string
		readStoredCredsFn    func() (credentials.Credentials, error)
		storeCredsFn         func(credentials.Credentials)
		createRequestTokenFn func(string, string) (string, error)
		openBrowserFn        func(string) error
	}
	tests := []struct {
		name              string
		args              args
		wantStored        credentials.Credentials
		wantAuthUrlOpened string
		wantErr           bool
	}{
		{
			name: "PerformsFullAuthorization",
			args: args{
				appId:                "ID",
				readStoredCredsFn:    noStoredCredentialsExistFn,
				storeCredsFn:         defaultStoreFn,
				createRequestTokenFn: defaultCreateReqTokenFn,
				openBrowserFn:        defaultOpenBrowserFn,
			},
			wantStored:        credentials.Credentials{RequestToken: testRequestToken},
			wantAuthUrlOpened: fmt.Sprintf(authUrlTemplate, testRequestToken, redirectUri),
			wantErr:           false,
		},
		{
			name: "ReturnsStoredCredentialIfOneExists",
			args: args{
				appId: "ID",
				readStoredCredsFn: func() (credentials.Credentials, error) {
					return credentials.Credentials{RequestToken: testRequestToken}, nil
				},
				storeCredsFn:         defaultStoreFn,
				createRequestTokenFn: defaultCreateReqTokenFn,
				openBrowserFn:        defaultOpenBrowserFn,
			},
			wantStored:        credentials.Credentials{},
			wantAuthUrlOpened: "",
			wantErr:           false,
		},
		{
			name: "HandlesErrorCreatingRequestToken",
			args: args{
				appId:             "ID",
				readStoredCredsFn: noStoredCredentialsExistFn,
				storeCredsFn:      defaultStoreFn,
				createRequestTokenFn: func(s1, s2 string) (string, error) {
					return "", fmt.Errorf("failure creating token")
				},
				openBrowserFn: defaultOpenBrowserFn,
			},
			wantStored:        credentials.Credentials{},
			wantAuthUrlOpened: "",
			wantErr:           true,
		},
		{
			name: "HandlesErrorOpeningBrowser",
			args: args{
				appId:                "ID",
				readStoredCredsFn:    noStoredCredentialsExistFn,
				storeCredsFn:         defaultStoreFn,
				createRequestTokenFn: defaultCreateReqTokenFn,
				openBrowserFn: func(url string) error {
					return fmt.Errorf("failure opening url")
				},
			},
			wantStored:        credentials.Credentials{RequestToken: testRequestToken},
			wantAuthUrlOpened: "",
			wantErr:           true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storedCreds = credentials.Credentials{}
			urlToOpen = ""

			if err := authorizeApp(tt.args.appId, tt.args.readStoredCredsFn, tt.args.storeCredsFn, tt.args.createRequestTokenFn, tt.args.openBrowserFn); (err != nil) != tt.wantErr {
				t.Errorf("authorizeApp() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(storedCreds, tt.wantStored) {
				t.Errorf("authorizeApp - stored credentials %v, want %v", storedCreds, tt.wantStored)
			}

			if urlToOpen != tt.wantAuthUrlOpened {
				t.Errorf("authorizeApp - would open URL %v, want %v", urlToOpen, tt.wantAuthUrlOpened)
			}
		})
	}
}

func Test_getAccessToken(t *testing.T) {

	noStoredAccessTokenExistFn := func() (credentials.Credentials, error) {
		return credentials.Credentials{RequestToken: testRequestToken}, fmt.Errorf("No Creds Exist")
	}

	defaultCreateAccessTokenFn := func(s1, s2 string) (string, error) {
		return testAccessToken, nil
	}

	var storedCreds credentials.Credentials
	defaultStoreFn := func(c credentials.Credentials) {
		storedCreds = c
	}

	type args struct {
		appId               string
		readStoredCredsFn   func() (credentials.Credentials, error)
		storeCredsFn        func(credentials.Credentials)
		createAccessTokenFn func(string, string) (string, error)
	}
	tests := []struct {
		name       string
		args       args
		want       string
		wantStored credentials.Credentials
		wantErr    bool
	}{
		{
			name: "CreatesAndStoresNewToken",
			args: args{
				appId:               "id",
				readStoredCredsFn:   noStoredAccessTokenExistFn,
				storeCredsFn:        defaultStoreFn,
				createAccessTokenFn: defaultCreateAccessTokenFn,
			},
			want:       testAccessToken,
			wantStored: credentials.Credentials{RequestToken: testRequestToken, AccessToken: testAccessToken},
			wantErr:    false,
		},
		{
			name: "ReturnsStoredTokenIfOneExists",
			args: args{
				appId: "id",
				readStoredCredsFn: func() (credentials.Credentials, error) {
					return credentials.Credentials{RequestToken: testRequestToken, AccessToken: testAccessToken}, nil
				},
				storeCredsFn:        defaultStoreFn,
				createAccessTokenFn: defaultCreateAccessTokenFn,
			},
			want:       testAccessToken,
			wantStored: credentials.Credentials{},
			wantErr:    false,
		},
		{
			name: "ReturnsErrorIfRequestTokenIsMissing",
			args: args{
				appId: "id",
				readStoredCredsFn: func() (credentials.Credentials, error) {
					return credentials.Credentials{}, fmt.Errorf("No Creds Exist")
				},
				storeCredsFn:        defaultStoreFn,
				createAccessTokenFn: defaultCreateAccessTokenFn,
			},
			want:       "",
			wantStored: credentials.Credentials{},
			wantErr:    true,
		},
		{
			name: "ReturnsErrorIfCreatingTokenFails",
			args: args{
				appId:             "id",
				readStoredCredsFn: noStoredAccessTokenExistFn,
				storeCredsFn:      defaultStoreFn,
				createAccessTokenFn: func(s1, s2 string) (string, error) {
					return "", fmt.Errorf("failure creating token")
				},
			},
			want:       "",
			wantStored: credentials.Credentials{},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storedCreds = credentials.Credentials{}

			got, err := getAccessToken(tt.args.appId, tt.args.readStoredCredsFn, tt.args.storeCredsFn, tt.args.createAccessTokenFn)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getAccessToken() = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(storedCreds, tt.wantStored) {
				t.Errorf("getAccessToken - stored credentials %v, want %v", storedCreds, tt.wantStored)
			}
		})
	}
}
