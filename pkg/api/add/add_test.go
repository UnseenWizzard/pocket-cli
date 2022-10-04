package add

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddApiCall(t *testing.T) {
	serverCalled := false
	type args struct {
		consumerKey string
		accessToken string
		articleUrl  string
	}
	tests := []struct {
		name           string
		args           args
		serverResponse func(w http.ResponseWriter, r *http.Request)
		wantErr        bool
	}{
		{
			"SunnyCase",
			args{
				"key",
				"token",
				"http://www.some.article/url",
			},
			func(w http.ResponseWriter, r *http.Request) {
				serverCalled = true
				res := responsePayload{
					Status: 1,
				}
				w.WriteHeader(200)
				b, _ := json.Marshal(res)
				_, _ = w.Write(b)
			},
			false,
		},
		{
			"OverallErrorResponse_ReturnsError",
			args{
				"key",
				"token",
				"http://www.some.article/url",
			},
			func(w http.ResponseWriter, r *http.Request) {
				serverCalled = true
				res := responsePayload{
					Status: 0,
				}
				w.WriteHeader(200)
				b, _ := json.Marshal(res)
				_, _ = w.Write(b)
			},
			true,
		},
		{
			"HTTP_401_ReturnsError",
			args{
				"key",
				"token",
				"http://www.some.article/url",
			},
			func(w http.ResponseWriter, r *http.Request) {
				serverCalled = true
				w.WriteHeader(401)
			},
			true,
		},
		{
			"HTTP_502_ReturnsError",
			args{
				"key",
				"token",
				"http://www.some.article/url",
			},
			func(w http.ResponseWriter, r *http.Request) {
				serverCalled = true
				w.WriteHeader(502)
			},
			true,
		},
		{
			"InvalidResponse_ReturnsError",
			args{
				"key",
				"token",
				"http://www.some.article/url",
			},
			func(w http.ResponseWriter, r *http.Request) {
				serverCalled = true
				w.WriteHeader(305)
			},
			true,
		},
	}
	for _, tt := range tests {
		ts := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
		serverCalled = false

		t.Run(tt.name, func(t *testing.T) {
			gotErr := addArticle(ts.URL, tt.args.consumerKey, tt.args.accessToken, tt.args.articleUrl)
			if !serverCalled {
				t.Error("No HTTP request sent to test sever")
			}
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("AddArticle() Error = %v, wanted err? %v", gotErr, tt.wantErr)
			}
		})

		ts.Close()
	}
}
