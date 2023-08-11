package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	list     list.Model
	selected gitLog
	height   int
	width    int
	message  string
}

func initialModel() model {
	items := []list.Item{}

	m := model{
		list: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	m.list.Title = "What is the first good commit?"
	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, getGitLogs, getGitBranch)
}

func getGitBranch() tea.Msg {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()

	if err != nil {
		return errMsg{err}
	}
	branchName := string(out)
	branchName = strings.Replace(branchName, "\n", "", 1)

	return gitBranchMsg{branchName}
}

type gitBranchMsg struct {
	branchName string
}

type sayHello string

func getGitLogs() tea.Msg {
	out, err := exec.Command("git", "log", "--format='%H||%an||%as||%s").Output()
	if err != nil {
		return errMsg{err}
	}

	// sending message into the program :)
	p.Send(sayHello("hi baby :)"))

	gitLogString := string(out)
	x := strings.Split(gitLogString, "\n")
	items := []gitLog{}
	for _, line := range x {
		cols := strings.Split(line, "||")
		if len(cols) == 4 {
			items = append(items, gitLog{
				githash: cols[0],
				author:  cols[1],
				date:    cols[2],
				message: cols[3],
			})
		}
	}

	return gitLogMsg{logs: items}
}

type gitLog struct {
	githash string
	author  string
	date    string
	message string
}

func (i gitLog) Title() string       { return i.author + " " + i.message }
func (i gitLog) Description() string { return i.date + " " + i.githash }
func (i gitLog) FilterValue() string { return i.author + " " + i.message }
func (i gitLog) Item() gitLog        { return i }

type gitLogMsg struct {
	logs []gitLog
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case gitBranchMsg:
		m.list.Title = msg.branchName + ": " + m.list.Title
	case errMsg:
		// todo: do something with errors

	case sayHello:
		m.message = string(msg)
	case gitLogMsg:
		items := m.list.Items()
		for _, item := range msg.logs {
			items = append(items, item)
		}
		m.list.SetItems(items)

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
			m.selected = m.list.SelectedItem().(gitLog)
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

var p *tea.Program

func main() {
	p = tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	if m, ok := m.(model); ok && m.selected.message != "" {
		fmt.Println(m.selected.message)
	}
}
