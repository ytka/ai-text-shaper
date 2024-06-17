package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type updateSpinnerMsg string
type runMsg string

type Model struct {
	spinner    spinner.Model
	statusText string
	quitting   bool
	err        error
	runner     func(m Model) tea.Cmd
}

func New(text string, runner func(m Model) tea.Cmd) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Model{
		statusText: text,
		spinner:    s,
		runner:     runner,
	}
}

func makeRunMsg() tea.Cmd {
	return func() tea.Msg {
		return runMsg("running")
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, makeRunMsg())
}

func (m Model) UpdateMsg(msg string) tea.Msg {
	return updateSpinnerMsg(msg)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case runMsg:
		m.spinner, _ = m.spinner.Update(tea.Quit)
		return m, m.runner(m)
	case updateSpinnerMsg:
		m.statusText = string(msg)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case error:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m Model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("\n\n%s %s\n\n", m.spinner.View(), m.statusText)
	if m.quitting {
		return str + "\n"
	}
	return str
}
