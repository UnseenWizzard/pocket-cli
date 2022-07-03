package list

import (
	"fmt"
	"log"

	"github.com/dixonwille/wmenu"
	"riedmann.dev/pocket-cli/pkg/login"
	"riedmann.dev/pocket-cli/pkg/retrieve"
	"riedmann.dev/pocket-cli/pkg/util"
)

const count = 10
var offset = 0
var articles []retrieve.Article

func ListArticles() {
	menu := wmenu.NewMenu("Which article do you want to read?")
	menu.ClearOnMenuRun()

	fetched := fetchArticles(util.PocketAppId, login.GetAccessToken(util.PocketAppId), count, offset)
	articles = append(articles, fetched...)

	for _, a := range articles {
		title := a.Title
		if title == "" {
			title = a.GivenTitle
		}
		info := fmt.Sprintf("%s (%v min)", title, a.ReadTime)
		menu.Option(info, a.Url, false, openArticle)
	}

	menu.Option("load more ...", &articles, true, fetchMore)

	err := menu.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func fetchMore(opt wmenu.Opt) error {
	offset += 10
	ListArticles()
	return nil
}

func fetchArticles(consumerKey string, accessToken string, count int, offset int) []retrieve.Article {
	articles := retrieve.RetrieveUnread(consumerKey, accessToken, count, offset)

	var fetched = make([]retrieve.Article, count)
	i := 0
	for _, a := range articles.List {
		fetched[i] = a
		i++
	}
	return fetched
} 

func openArticle(opt wmenu.Opt) error {
	url, ok := opt.Value.(string)
	if !ok {
		return fmt.Errorf("expected an url value as string")
	}
	util.OpenInBrowser(url)
	return nil
}

