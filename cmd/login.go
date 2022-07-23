package cmd

import (
	"github.com/UnseenWizzard/pocket-cli/pkg/api/login"
	"github.com/UnseenWizzard/pocket-cli/pkg/util"
)

func Login(reset bool) {
	if reset {
		login.RemoveStoredCredentials()
	}
	login.AuthorizeApp(util.PocketAppId)
}
