package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list     list.Model
	selected item
	width    int
	height   int
}

func initialModel() model {
	items := []list.Item{
		item{title: "Something", desc: "special"},
		item{title: "Buy", desc: "me a coffee"},
	}

	m := model{
		list: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	m.list.Title = "What should we buy?"
	return m
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.list.SetSize(m.width, m.height)
	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			m.selected = m.list.SelectedItem().(item)
			return m, tea.Quit
		}
	}

	// update list
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, cmd
}

func (m model) View() string {
	// The header
	s := m.list.View()

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	if m, ok := m.(model); ok && m.selected.title != "" {
		fmt.Println(m.selected.title)
	}
}
