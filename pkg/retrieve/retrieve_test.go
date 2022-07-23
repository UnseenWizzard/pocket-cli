package retrieve

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRetrieveUnread(t *testing.T) {
	type args struct {
		consumerKey string
		accessToken string
		count       int
		offset      int
	}
	tests := []struct {
		name           string
		args           args
		serverResponse func(w http.ResponseWriter, r *http.Request)
		want           ResponsePayload
		wantErr        bool
	}{
		{
			"Sunny Case",
			args{
				"key",
				"token",
				1,
				0,
			},
			func(w http.ResponseWriter, r *http.Request) {
				res := ResponsePayload{
					Status:   1,
					Complete: 1,
					List: map[string]Article{
						"42": {
							Id:         "42",
							Title:      "Test",
							GivenTitle: "Test",
							Url:        "test.com/test",
							Excerpt:    "A summary",
							ReadTime:   500,
						},
					},
				}
				w.WriteHeader(200)
				b, _ := json.Marshal(res)
				_, _ = w.Write(b)
			},
			ResponsePayload{
				Status:   1,
				Complete: 1,
				List: map[string]Article{
					"42": {
						Id:         "42",
						Title:      "Test",
						GivenTitle: "Test",
						Url:        "test.com/test",
						Excerpt:    "A summary",
						ReadTime:   500,
					},
				},
			},
			false,
		},
		{
			"HTTP Error",
			args{
				"key",
				"token",
				1,
				0,
			},
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(401)
			},
			ResponsePayload{},
			true,
		},
		{
			"HTTP Error",
			args{
				"key",
				"token",
				1,
				0,
			},
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(502)
			},
			ResponsePayload{},
			true,
		},
		{
			"Invalid response",
			args{
				"key",
				"token",
				1,
				0,
			},
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(201)
			},
			ResponsePayload{},
			true,
		},
	}
	for _, tt := range tests {
		ts := httptest.NewServer(http.HandlerFunc(tt.serverResponse))

		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := retrieveUnread(ts.URL, tt.args.consumerKey, tt.args.accessToken, tt.args.count, tt.args.offset)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("RetrieveUnread() Error = %v, wanted err? %v", gotErr, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RetrieveUnread() = %v, want %v", got, tt.want)
			}
		})

		ts.Close()
	}
}
