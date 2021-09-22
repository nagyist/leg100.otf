package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/leg100/go-tfe"
	"github.com/leg100/otf"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type run struct {
	*otf.Run
}

func (r run) Title() string       { return r.ID }
func (r run) Description() string { return string(r.Status) }
func (r run) FilterValue() string { return r.ID }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	runs := []list.Item{
		run{&otf.Run{ID: "run-1", Status: tfe.RunApplied}},
		run{&otf.Run{ID: "run-2", Status: tfe.RunPlanQueued}},
		run{&otf.Run{ID: "run-3", Status: tfe.RunPending}},
	}

	m := model{list: list.NewModel(runs, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Runs"

	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
