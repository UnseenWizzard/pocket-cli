package login

import (
	"github.com/UnseenWizzard/pocket-cli/pkg/login"
	"github.com/UnseenWizzard/pocket-cli/pkg/util"
)

func Login(reset bool) {
	if reset {
		login.RemoveStoredCredentials()
	}
	login.AuthorizeApp(util.PocketAppId)
}
