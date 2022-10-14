package cli

import (
  "fmt"
  "os"

  tea "github.com/charmbracelet/bubbletea"
)

type model struct {
  choices []string  // items on the list
  cursor int // which item our cursor is pointing at
  selected map[int]struct{} // which items are selected
}

func initialModel() model {

  return model{
    choices: []string{"Manchester United", "Arsenal", "Tottenham", "Manchesthair Shitty :DDDD"},
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
  for i, choice := range m.choices {
    cursor := " " //no cursor
    if m.cursor == i {
      cursor = ">"
    }

    if  _, ok := m.selected[i]; ok {
      checked = m.choices[i]
    }

    s += fmt.Sprintf("%s is checked", checked)

    s += fmt.Sprintf("%s %s\n", cursor, choice)
  }

  if checked != " " {
    s = fmt.Sprintf("%s HAS BEEN SELECTED", checked)
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
