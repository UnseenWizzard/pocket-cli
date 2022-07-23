package cmd

import (
	"log"

	"github.com/UnseenWizzard/pocket-cli/pkg/articles"
	"github.com/UnseenWizzard/pocket-cli/pkg/util"
	"github.com/manifoldco/promptui"
)

const count = 10

var offset = 0
var entries []articles.Article

func ListArticles() {
	fetched, err := articles.Fetch(count, offset)
	if err != nil {
		log.Fatal("Failed to fetch articles: %w", err)
	}
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
		Items:     append(entries, articles.Article{Title: "Load more ..."}),
		Templates: templates,
		Size:      11,
	}

	resIndex, _, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	if resIndex != len(entries) {
		e := entries[resIndex]
		err := util.OpenInBrowser(e.Url)
		if err != nil {
			log.Fatal("Failed to open browser: %w", err)
		}
		ModifyItemPrompt(e.Id)
	} else {
		fetchMore()
	}
}

func fetchMore() {
	offset += 10
	ListArticles()
}
