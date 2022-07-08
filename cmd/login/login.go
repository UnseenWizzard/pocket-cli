package login

import (
	"riedmann.dev/pocket-cli/pkg/login"
	"riedmann.dev/pocket-cli/pkg/util"
)

func Login(reset bool) {
	if reset {
		login.RemoveStoredCredentials()
	}
	login.AuthorizeApp(util.PocketAppId)
}
