package modify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestModifyApiCall(t *testing.T) {
	type args struct {
		consumerKey string
		accessToken string
		action      string
		articleId   string
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
				"archive",
				"42",
			},
			func(w http.ResponseWriter, r *http.Request) {
				res := responsePayload{
					ActionResults: []bool{true},
					Status:        1,
				}
				w.WriteHeader(200)
				b, _ := json.Marshal(res)
				w.Write(b)
			},
			false,
		},
		{
			"OverallErrorResponse_ReturnsError",
			args{
				"key",
				"token",
				"archive",
				"42",
			},
			func(w http.ResponseWriter, r *http.Request) {
				res := responsePayload{
					ActionResults: []bool{true},
					Status:        0,
				}
				w.WriteHeader(200)
				b, _ := json.Marshal(res)
				w.Write(b)
			},
			true,
		},
		{
			"SingleActionErrorResponse_ReturnsError",
			args{
				"key",
				"token",
				"archive",
				"42",
			},
			func(w http.ResponseWriter, r *http.Request) {
				res := responsePayload{
					ActionResults: []bool{false},
					Status:        1,
				}
				w.WriteHeader(200)
				b, _ := json.Marshal(res)
				w.Write(b)
			},
			true,
		},
		{
			"HTTP_401_ReturnsError",
			args{
				"key",
				"token",
				"archive",
				"42",
			},
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(401)
			},
			true,
		},
		{
			"HTTP_502_ReturnsError",
			args{
				"key",
				"token",
				"archive",
				"42",
			},
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(502)
			},
			true,
		},
		{
			"InvalidResponse_ReturnsError",
			args{
				"key",
				"token",
				"archive",
				"42",
			},
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(201)
			},
			true,
		},
	}
	for _, tt := range tests {
		ts := httptest.NewServer(http.HandlerFunc(tt.serverResponse))

		t.Run(tt.name, func(t *testing.T) {
			gotErr := makeModifyApiCall(ts.URL, tt.args.consumerKey, tt.args.accessToken, tt.args.action, tt.args.articleId)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("RetrieveUnread() Error = %v, wanted err? %v", gotErr, tt.wantErr)
			}
		})

		ts.Close()
	}
}
