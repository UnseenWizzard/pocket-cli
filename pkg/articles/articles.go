package articles

import (
	"fmt"
	"strings"

	"github.com/UnseenWizzard/pocket-cli/pkg/api/login"
	"github.com/UnseenWizzard/pocket-cli/pkg/api/retrieve"
	"github.com/UnseenWizzard/pocket-cli/pkg/util"
)

type Article struct {
	Id       string
	Title    string
	Excerpt  string
	ReadTime string
	Url      string
}

func Fetch(count int, offset int) ([]Article, error) {
	return fetch(util.PocketAppId, login.GetAccessToken(util.PocketAppId), count, offset, retrieve.RetrieveUnread)
}

func fetch(consumerKey string, accessToken string, count int, offset int, retrieveFn func(string, string, int, int) (retrieve.ResponsePayload, error)) ([]Article, error) {
	articles, err := retrieveFn(consumerKey, accessToken, count, offset)
	if err != nil {
		return []Article{}, fmt.Errorf("failed to retrieve articles: %w", err)
	}
	if articles.Status != 1 {
		return []Article{}, fmt.Errorf("failed to retrieve articles: API returned failure status %v", articles.Status)
	}

	var fetched = make([]Article, len(articles.List))
	i := 0
	for _, a := range articles.List {

		fetched[i] = Article{
			Id:       a.Id,
			Title:    beautifyTitle(a),
			Excerpt:  beautifyExcerpt(a),
			ReadTime: beautifyReadTime(a),
			Url:      a.Url,
		}
		i++
	}
	return fetched, nil
}

func beautifyTitle(a retrieve.Article) string {
	title := a.Title
	if title == "" {
		title = a.GivenTitle
	}
	return title
}

func beautifyReadTime(a retrieve.Article) string {
	time := "?"
	if a.ReadTime > 0 {
		time = fmt.Sprintf("%v min", a.ReadTime)
	}
	return time
}

func beautifyExcerpt(a retrieve.Article) string {
	if len(strings.TrimSpace(a.Excerpt)) == 0 {
		return "[No excerpt available]"
	}
	if len(a.Excerpt) > 120 {
		return a.Excerpt[:117] + "..."
	}
	return a.Excerpt
}
