package cli

import (
	"fmt"
	"os"

	"example.com/m/v2/fetch"
	tea "github.com/charmbracelet/bubbletea"
)

type Status struct {
	Type string
}
type Event struct {
	Id     int
	Slug   string
	status Status
}
type League struct {
	Id     int
	Name   string
	events []Event
}

type model struct {
	choices  []fetch.Choice   // items on the list
	leagues  []fetch.League   // items on the list
	cursor   int              // which item our cursor is pointing at
	selected map[int]struct{} // which items are selected
}

func initialModel() model {
	jsonData := fetch.GetJson()
	leagueData := fetch.GetActiveLeagues(jsonData)

	return model{
		choices:  leagueData,
		leagues:  jsonData,
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "What game would you like to watch?\n\n"
	checked := " " //not selected
	selectedLeague := 0
	for i, choice := range m.choices {
		cursor := " " //no cursor
		if m.cursor == i {
			cursor = ">"
		}

		if _, ok := m.selected[i]; ok {
			checked = m.choices[i].Name
			selectedLeague = m.choices[i].Id
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice.Name)
	}

	if checked != " " {
		// send current league Id to new func
		// and get current games in progress
		s = fmt.Sprintf("%s HAS BEEN SELECTED", checked)

		fetch.GetGames(selectedLeague, m.leagues)
		return s
	}

	s += "\nPress q to quit.\n"
	return s
}

func Run() {
	p := tea.NewProgram(initialModel())

	if err := p.Start(); err != nil {
		fmt.Print("Alas there's been an error %v", err)
		os.Exit(1)
	}
}
