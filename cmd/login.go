package cmd

import (
	"github.com/UnseenWizzard/pocket-cli/pkg/auth"
	"github.com/UnseenWizzard/pocket-cli/pkg/auth/credentials"
	"github.com/UnseenWizzard/pocket-cli/pkg/util"
)

func Login(reset bool) {
	if reset {
		credentials.RemoveStoredCredentials()
	}
	auth.AuthorizeApp(util.PocketAppId)
}
