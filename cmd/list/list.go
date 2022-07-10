package list

import (
	"fmt"
	"log"

	"github.com/UnseenWizzard/pocket-cli/pkg/login"
	"github.com/UnseenWizzard/pocket-cli/pkg/retrieve"
	"github.com/UnseenWizzard/pocket-cli/pkg/util"
	"github.com/manifoldco/promptui"
)

const count = 10

var offset = 0
var entries []listEntry

type listEntry struct {
	Title    string
	Excerpt  string
	ReadTime string
	Url      string
}

func ListArticles() {
	fetched := fetchArticles(util.PocketAppId, login.GetAccessToken(util.PocketAppId), count, offset)
	entries = append(entries, fetched...)

	templates := &promptui.SelectTemplates{
		Label:    "{{ . | bold }}",
		Active:   "\U0001F4D9 {{ .Title | bold }} {{if eq .ReadTime \"\"}} {{else}} ({{ .ReadTime | red }}) {{end}}",
		Inactive: "  {{ .Title | cyan }} {{if eq .ReadTime \"\"}} {{else}} ({{ .ReadTime | faint }}) {{end}}",
		Selected: "{{if eq .Title \"Load more ...\"}} {{ \"\U0001F504 Loading more articles...\" | red | bold}} {{else}} \U0001F4D6 {{ \"Opening...\" | bold}} {{ .Title | red | bold }} {{end}}",
		Details:  " {{.Excerpt | faint }}",
	}

	prompt := promptui.Select{
		Label:     "\U0001F4DA Which article do you want to read?",
		Items:     append(entries, listEntry{Title: "Load more ..."}),
		Templates: templates,
		Size:      11,
	}

	resIndex, _, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	if resIndex != len(entries) {
		util.OpenInBrowser(entries[resIndex].Url)
	} else {
		fetchMore()
	}
}

func fetchMore() {
	offset += 10
	ListArticles()
}

func fetchArticles(consumerKey string, accessToken string, count int, offset int) []listEntry {
	articles := retrieve.RetrieveUnread(consumerKey, accessToken, count, offset)

	var fetched = make([]listEntry, count)
	i := 0
	for _, a := range articles.List {

		fetched[i] = listEntry{
			Title:    beautifyTitle(a),
			Excerpt:  beautifyExcerpt(a),
			ReadTime: beautifyReadTime(a),
			Url:      a.Url,
		}
		i++
	}
	return fetched
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
	if len(a.Excerpt) == 0 {
		return "[No excerpt available]"
	}
	if len(a.Excerpt) > 120 {
		return a.Excerpt[:117] + "..."
	}
	return a.Excerpt
}
