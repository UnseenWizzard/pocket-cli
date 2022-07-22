package login

import "testing"

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
