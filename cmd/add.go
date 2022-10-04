package cmd

import (
	"github.com/UnseenWizzard/pocket-cli/pkg/api/add"
	"github.com/UnseenWizzard/pocket-cli/pkg/auth"
	"github.com/UnseenWizzard/pocket-cli/pkg/util"
	"log"
)

func AddArticle(url string) {
	err := add.AddArticle(util.PocketAppId, auth.GetAccessToken(util.PocketAppId), url)
	if err != nil {
		log.Fatalf("Failed to add article (%s): %v", url, err)
	}
}
