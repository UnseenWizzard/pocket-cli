package login

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_parseAccessTokenResponse(t *testing.T) {
	tests := []struct {
		name            string
		responseInput   string
		wantAccessToken string
		wantForUser     string
		wantErr         bool
	}{
		{
			"ParsesSimpleResponse",
			"access_token=abc42&username=malcom.reynolds",
			"abc42",
			"malcom.reynolds",
			false,
		},
		{
			"ParsesValidResponseWithExtraData",
			"access_token=abc42&username=malcom.reynolds&useless_addition=this",
			"abc42",
			"malcom.reynolds",
			false,
		},
		{
			"ParsesResponseInUnexpectedOrder",
			"username=malcom.reynolds&access_token=abc42",
			"abc42",
			"malcom.reynolds",
			false,
		},
		{
			"FailsIfInvalidQueryResponse",
			"username=&&",
			"",
			"",
			true,
		},
		{
			"FailsIfNoUsernameKeyFound",
			"access_token=abc42",
			"",
			"",
			true,
		},
		{
			"FailsIfNoUsernameValueFound",
			"username&access_token=abc42",
			"",
			"",
			true,
		},
		{
			"FailsIfUsernameValueEmpty",
			"username=&access_token=abc42",
			"",
			"",
			true,
		},
		{
			"FailsIfNoTokenKeyFound",
			"username=malcom.reynolds",
			"",
			"",
			true,
		},
		{
			"FailsIfNoTokenValueFound",
			"username=malcom.reynolds&access_token",
			"",
			"",
			true,
		},
		{
			"FailsIfTokenValueEmpty",
			"username=malcom.reynolds&access_token=",
			"",
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAccessToken, gotForUser, err := parseAccessTokenResponse(tt.responseInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAccessTokenResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAccessToken != tt.wantAccessToken {
				t.Errorf("parseAccessTokenResponse() gotAccessToken = %v, want %v", gotAccessToken, tt.wantAccessToken)
			}
			if gotForUser != tt.wantForUser {
				t.Errorf("parseAccessTokenResponse() gotForUser = %v, want %v", gotForUser, tt.wantForUser)
			}
		})
	}
}

func Test_getRequestToken(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		want           string
		wantErr        bool
	}{
		{
			"GetsRequestToken",
			func(w http.ResponseWriter, r *http.Request) {
				res := "code=a-request-token"
				w.WriteHeader(200)
				_, _ = io.WriteString(w, res)
			},
			"a-request-token",
			false,
		},
		{
			"HTTP_401_ReturnsError",
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(401)
			},
			"",
			true,
		},
		{
			"HTTP_502_ReturnsError",
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(502)
			},
			"",
			true,
		},
		{
			"EmptyResponse_ReturnsError",
			func(w http.ResponseWriter, r *http.Request) {
				res := ""
				w.WriteHeader(200)
				_, _ = io.WriteString(w, res)
			},
			"",
			true,
		},
		{
			"InvalidResponse_ReturnsError",
			func(w http.ResponseWriter, r *http.Request) {
				res := "some-key=some-val"
				w.WriteHeader(200)
				_, _ = io.WriteString(w, res)
			},
			"",
			true,
		},
		{
			"InvalidResponseSayingCode_ReturnsError",
			func(w http.ResponseWriter, r *http.Request) {
				res := "not_expected_code=a-request-token"
				w.WriteHeader(200)
				_, _ = io.WriteString(w, res)
			},
			"",
			true,
		},
		{
			"TooLongResponse_ReturnsError",
			func(w http.ResponseWriter, r *http.Request) {
				res := "not_expected_code=a-request-token&some-other=entry"
				w.WriteHeader(200)
				_, _ = io.WriteString(w, res)
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(tt.serverResponse))

			var storedCreds credentials
			storeCredsFn := func(c credentials) {
				storedCreds = c
			}

			got, gotErr := getRequestToken(ts.URL, "app-id", storeCredsFn)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("getRequestToken() Error = %v, wanted err? %v", gotErr, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("getRequestToken() = %v, want %v", got, tt.want)
			}

			if storedCreds.RequestToken != tt.want {
				t.Errorf("failed to store request token. stored value = %v", storedCreds)
			}

			ts.Close()
		})
	}
}

func Test_getAccessToken(t *testing.T) {

	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		want           string
		wantErr        bool
	}{
		{
			"GetsAccessToken",
			func(w http.ResponseWriter, r *http.Request) {
				res := "access_token=an-access-token&username=malcom.reynolds"
				w.WriteHeader(200)
				_, _ = io.WriteString(w, res)
			},
			"an-access-token",
			false,
		},
		{
			"GetsAccessToken_ValidReplyWithAdditionalData",
			func(w http.ResponseWriter, r *http.Request) {
				res := "username=malcom.reynolds&some-extra=key&access_token=an-access-token"
				w.WriteHeader(200)
				_, _ = io.WriteString(w, res)
			},
			"an-access-token",
			false,
		},
		{
			"HTTP_401_ReturnsError",
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(401)
			},
			"",
			true,
		},
		{
			"HTTP_403_ReturnsError",
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(403)
			},
			"",
			true,
		},
		{
			"HTTP_502_ReturnsError",
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(502)
			},
			"",
			true,
		},
		{
			"EmptyResponse_ReturnsError",
			func(w http.ResponseWriter, r *http.Request) {
				res := ""
				w.WriteHeader(200)
				_, _ = io.WriteString(w, res)
			},
			"",
			true,
		},
		{
			"InvalidResponse_ReturnsError",
			func(w http.ResponseWriter, r *http.Request) {
				res := "some-key=some-val"
				w.WriteHeader(200)
				_, _ = io.WriteString(w, res)
			},
			"",
			true,
		},
		{
			"PartialResponse_ReturnsError",
			func(w http.ResponseWriter, r *http.Request) {
				res := "username=malcom.reynolds"
				w.WriteHeader(200)
				_, _ = io.WriteString(w, res)
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(tt.serverResponse))

			var storedCreds credentials
			storeCredsFn := func(c credentials) {
				storedCreds = c
			}

			got, gotErr := getAccessToken(ts.URL, "app-id", "a-request-code", storeCredsFn)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("getAccessToken() Error = %v, wanted err? %v", gotErr, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("getAccessToken() = %v, want %v", got, tt.want)
			}

			if storedCreds.AccessToken != tt.want {
				t.Errorf("failed to store access token. stored value = %v", storedCreds)
			}

			ts.Close()
		})
	}
}
