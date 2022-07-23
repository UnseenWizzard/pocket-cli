package cmd

import (
	"log"
	"os"

	"github.com/UnseenWizzard/pocket-cli/pkg/api/modify"
	"github.com/UnseenWizzard/pocket-cli/pkg/auth"
	"github.com/UnseenWizzard/pocket-cli/pkg/util"
	"github.com/manifoldco/promptui"
)

func ModifyItemPrompt(id string) {

	actions := []struct {
		Title string
		act   modify.ModifyFunction
	}{
		{
			"Mark as read",
			modify.ArchiveArticle,
		},
		{
			"Mark as favorite",
			modify.FavoriteArticle,
		},
		{
			"Exit",
			nil,
		},
	}

	prompt := promptui.Select{
		Label: "What's next?",
		Items: actions,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . | bold }}",
			Active:   "{{ .Title | bold }}",
			Inactive: "  {{ .Title | cyan }}",
			Selected: "{{ .Title | red | bold }}",
		},
	}

	resIndex, _, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	selection := actions[resIndex]

	if selection.Title == "Exit" {
		os.Exit(0)
	}
	err = selection.act(util.PocketAppId, auth.GetAccessToken(util.PocketAppId), id)
	if err != nil {
		log.Fatal(err)
	}
}
