package cmd

import (
	"github.com/UnseenWizzard/pocket-cli/pkg/articles"
	"github.com/UnseenWizzard/pocket-cli/pkg/util"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type entry struct {
	title   string
	Excerpt string
	Url     string
}

func (a entry) Title() string       { return a.title }
func (a entry) Description() string { return a.Excerpt }
func (a entry) FilterValue() string { return a.title }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {

		case key.Matches(msg, m.list.KeyMap.ForceQuit):
			return m, tea.Quit

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			url := arts[m.list.Index()].Url
			util.OpenInBrowser(url)
		}
		if msg.String() == "l" {
			offset += 10
			a, _ := articles.Fetch(count, offset)
			arts = append(arts, a...)
			var items []list.Item
			for _, a := range arts {
				items = append(items, entry{
					title:   a.Title,
					Excerpt: a.Excerpt,
					Url:     a.Url,
				})
			}
			m.list.SetItems(items)
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	m.list.Cursor()
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

var arts []articles.Article

func GetProgram() *tea.Program {
	var err error
	arts, err = articles.Fetch(count, offset)
	if err != nil {
		log.Fatal("Failed to fetch articles: %w", err)
	}

	var items []list.Item
	for _, a := range arts {
		items = append(items, entry{
			title:   a.Title,
			Excerpt: a.Excerpt,
			Url:     a.Url,
		})
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "My Fave Things"

	return tea.NewProgram(m, tea.WithAltScreen())
}
