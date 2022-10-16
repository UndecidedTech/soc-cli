// package cli

// import (
// 	"fmt"
// 	"os"

// 	"example.com/m/v2/fetch"
// 	tea "github.com/charmbracelet/bubbletea"
// )

// type Status struct {
// 	Type string
// }
// type Event struct {
// 	Id     int
// 	Slug   string
// 	status Status
// }
// type League struct {
// 	Id     int
// 	Name   string
// 	events []Event
// }

// type model struct {
// 	games   []fetch.Event
// 	choices []fetch.Choice // items on the list
// 	leagues []fetch.League // items on the list
// 	cursor  int            // which item our cursor is pointing at
// 	chosen  bool           // which items is selected
// 	choice  int            // which items is selected
// }

// func initialModel() model {
// 	jsonData := fetch.GetJson()
// 	leagueData := fetch.GetActiveLeagues(jsonData)

// 	return model{
// 		choices: leagueData,
// 		leagues: jsonData,
// 		chosen:  false,
// 		choice:  0,
// 	}
// }

// func (m model) Init() tea.Cmd {
// 	return nil
// }

// func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "ctrl+c", "q":
// 			return m, tea.Quit
// 		case "up", "k":
// 			if m.cursor > 0 {
// 				m.cursor--
// 			}
// 		case "down", "j":
// 			if m.cursor < len(m.choices)-1 {
// 				m.cursor++
// 			}
// 		case "enter", " ":
// 			m.chosen = true
// 			// _, ok := m.selected[m.cursor]
// 			// if ok {
// 			// 	delete(m.selected, m.cursor)
// 			// } else {
// 			// 	m.selected[m.cursor] = struct{}{}
// 			// }
// 		}
// 	}
// 	return m, nil
// }

// func chosenView(m model) string {
// 	s := "What game would you like to watch?\n\n"
// 	checked := " "    //not selected
// 	selectedGame := 0 // not using this yet
// 	_ = selectedGame
// 	for i, game := range m.games {
// 		cursor := " " //no cursor
// 		if m.cursor == i {
// 			cursor = ">"
// 		}

// 		if _, ok := m.selected[i]; ok {
// 			checked = m.games[i].Slug
// 			selectedGame = m.games[i].Id
// 		}

// 		s += fmt.Sprintf("%s %s\n", cursor, game.Slug)
// 	}

// 	if checked != " " {
// 		// send current league Id to new func
// 		// and get current games in progress
// 		s = fmt.Sprintf("%s HAS BEEN SELECTED", checked)

// 		return s
// 	}

// 	s += "\nPress q to quit.\n"
// 	return s
// }

// func choicesView(m model) string {
// 	s := "What league would you like to watch?\n\n"
// 	checked := " " //not selected
// 	selectedLeague := 0
// 	for i, choice := range m.choices {
// 		cursor := " " //no cursor
// 		if m.cursor == i {
// 			cursor = ">"
// 		}

// 		if _, ok := m.selected[i]; ok {
// 			checked = m.choices[i].Name
// 			selectedLeague = m.choices[i].Id
// 		}

// 		s += fmt.Sprintf("%s %s\n", cursor, choice.Name)
// 	}

// 	if checked != " " {
// 		// send current league Id to new func
// 		// and get current games in progress
// 		s = fmt.Sprintf("%s HAS BEEN SELECTED", checked)
// 		games := fetch.GetGames(selectedLeague, m.leagues)

// 		m.games = append(m.games, games...)
// 		return s
// 	}

// 	s += "\nPress q to quit.\n"
// 	return s
// }

// func (m model) View() string {

// 	if m.selected == false {
// 		return choicesView(m)
// 	} else {
// 		return chosenView(m)
// 	}
// }

// func Run() {
// 	p := tea.NewProgram(initialModel())

// 	if err := p.Start(); err != nil {
// 		fmt.Print("Alas there's been an error %v", err)
// 		os.Exit(1)
// 	}
// }

package cli

// An example demonstrating an application with multiple views.
//
// Note that this example was produced before the Bubbles progress component
// was available (github.com/charmbracelet/bubbles/progress) and thus, we're
// implementing a progress bar from scratch here.

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"example.com/m/v2/fetch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/termenv"
)

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
)

// General stuff for styling the view
var (
	term          = termenv.EnvColorProfile()
	keyword       = makeFgStyle("211")
	subtle        = makeFgStyle("241")
	progressEmpty = subtle(progressEmptyChar)
	dot           = colorFg(" • ", "236")

	// Gradient colors we'll use for the progress bar
	ramp = makeRamp("#B14FFF", "#00FFA3", progressBarWidth)
)

func Start() {
	jsonData := fetch.GetJson()
	// leagueData := fetch.GetActiveLeagues(jsonData)

	initialModel := model{
		Choice:   0,
		Choice2:  0,
		Chosen:   false,
		Leagues:  jsonData,
		Ticks:    10,
		Frames:   0,
		Progress: 0,
		Loaded:   false,
		Quitting: false,
	}
	p := tea.NewProgram(initialModel)

	if err := p.Start(); err != nil {
		fmt.Print("Alas there's been an error %v", err)
		os.Exit(1)
	}

}

type (
	tickMsg  struct{}
	frameMsg struct{}
)

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return frameMsg{}
	})
}

type model struct {
	Choice   int
	Choice2  int
	Chosen   bool
	Leagues  []fetch.League
	Ticks    int
	Frames   int
	Progress float64
	Loaded   bool
	Quitting bool
}

func (m model) Init() tea.Cmd {
	return tick()
}

// Main update function.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	if !m.Chosen {
		return updateChoices(msg, m)
	}
	return updateChosen(msg, m)
}

// The main view, which just calls the appropriate sub-view
func (m model) View() string {
	var s string
	if m.Quitting {
		return "\n  See you later!\n\n"
	}
	if !m.Chosen {
		s = choicesView(m)
	} else {
		s = chosenView(m)
	}
	return indent.String("\n"+s+"\n\n", 2)
}

// Sub-update functions

// Update loop for the first view where you're choosing a task.
func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.Choice++
			if m.Choice > len(m.Leagues) {
				m.Choice = len(m.Leagues) - 1
			}
		case "k", "up":
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 0
			}
		case "enter":
			m.Chosen = true
			return m, frame()
		}

	case tickMsg:
		if m.Ticks == 0 {
			m.Quitting = true
			return m, tea.Quit
		}
		m.Ticks--
		return m, tick()
	}

	return m, nil
}

func updateChosen(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.Choice2++
			if m.Choice2 > len(m.Leagues[m.Choice].Events) {
				m.Choice2 = len(m.Leagues[m.Choice].Events) - 1
			}
		case "k", "up":
			m.Choice2--
			if m.Choice2 < 0 {
				m.Choice2 = 0
			}
		case "enter":
			m.Chosen = true
			return m, frame()
		}

	case tickMsg:
		if m.Ticks == 0 {
			m.Quitting = true
			return m, tea.Quit
		}
		m.Ticks--
		return m, tick()
	}

	return m, nil
}

// Update loop for the second view after a choice has been made
// func updateChosen(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
// 	switch msg.(type) {
// 	case frameMsg:
// 		if !m.Loaded {
// 			m.Frames++
// 			m.Progress = ease.OutBounce(float64(m.Frames) / float64(100))
// 			if m.Progress >= 1 {
// 				m.Progress = 1
// 				m.Loaded = true
// 				m.Ticks = 3
// 				return m, tick()
// 			}
// 			return m, frame()
// 		}

// 	case tickMsg:
// 		if m.Loaded {
// 			if m.Ticks == 0 {
// 				m.Quitting = true
// 				return m, tea.Quit
// 			}
// 			m.Ticks--
// 			return m, tick()
// 		}
// 	}

// 	return m, nil
// }

// Sub-views

// The first view, where you're choosing a task
func choicesView(m model) string {
	// c := m.Choice

	tpl := "What League do you want to watch?\n\n"
	tpl += "%s\n\n"
	tpl += "Program quits in %s seconds\n\n"
	tpl += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")

	choices := ""
	for i, v := range m.Leagues {
		choice := checkbox(v.Name, m.Choice == i)
		newString := fmt.Sprintf("%s\n", choice)
		choices += newString
	}

	// choices := fmt.Sprintf(
	// 	replace,
	// 	options[:(len(options)-1)],
	// "%s\n%s\n%s\n%s",
	// checkboxList(m),
	// checkbox("", m.Choice == 0),
	// checkbox("Go to the market", m.Choice == 1),
	// checkbox("Read something", m.Choice == 2),
	// checkbox("See friends", m.Choice == 3),
	// )

	return fmt.Sprintf(tpl, choices, colorFg(strconv.Itoa(m.Ticks), "79"))
}

// The second view, after a task has been chosen
func chosenView(m model) string {
	tpl := "What game do you want to watch?\n\n"
	tpl += "%s\n\n"
	tpl += "Program quits in %s seconds\n\n"
	tpl += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")

	choices := ""
	for i, v := range m.Leagues[m.Choice].Events {
		choice := checkbox(v.Slug, m.Choice2 == i)
		newString := fmt.Sprintf("%s\n", choice)
		choices += newString
	}

	// var msg string

	// switch m.Choice {
	// case 0:
	// 	msg = fmt.Sprintf("Carrot planting?\n\nCool, we'll need %s and %s...", keyword("libgarden"), keyword("vegeutils"))
	// case 1:
	// 	msg = fmt.Sprintf("A trip to the market?\n\nOkay, then we should install %s and %s...", keyword("marketkit"), keyword("libshopping"))
	// case 2:
	// 	msg = fmt.Sprintf("Reading time?\n\nOkay, cool, then we’ll need a library. Yes, an %s.", keyword("actual library"))
	// default:
	// 	msg = fmt.Sprintf("It’s always good to see friends.\n\nFetching %s and %s...", keyword("social-skills"), keyword("conversationutils"))
	// }

	// label := "Downloading..."
	// if m.Loaded {
	// 	label = fmt.Sprintf("Downloaded. Exiting in %s seconds...", colorFg(strconv.Itoa(m.Ticks), "79"))
	// }

	// return msg + "\n\n" + label + "\n" + progressbar(80, m.Progress) + "%"
	return fmt.Sprintf(tpl, choices, colorFg(strconv.Itoa(m.Ticks), "79"))
}

func checkbox(label string, checked bool) string {
	if checked {
		return colorFg("[x] "+label, "212")
	}
	return fmt.Sprintf("[ ] %s", label)
}

func progressbar(width int, percent float64) string {
	w := float64(progressBarWidth)

	fullSize := int(math.Round(w * percent))
	var fullCells string
	for i := 0; i < fullSize; i++ {
		fullCells += termenv.String(progressFullChar).Foreground(term.Color(ramp[i])).String()
	}

	emptySize := int(w) - fullSize
	emptyCells := strings.Repeat(progressEmpty, emptySize)

	return fmt.Sprintf("%s%s %3.0f", fullCells, emptyCells, math.Round(percent*100))
}

// Utils

// Color a string's foreground with the given value.
func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

// Return a function that will colorize the foreground of a given string.
func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}

// Color a string's foreground and background with the given value.
func makeFgBgStyle(fg, bg string) func(string) string {
	return termenv.Style{}.
		Foreground(term.Color(fg)).
		Background(term.Color(bg)).
		Styled
}

// Generate a blend of colors.
func makeRamp(colorA, colorB string, steps float64) (s []string) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, colorToHex(c))
	}
	return
}

// Convert a colorful.Color to a hexadecimal format compatible with termenv.
func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and
// 1.
func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}
