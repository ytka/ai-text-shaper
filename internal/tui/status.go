package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	updateStatusMsg string
	quitMsg         struct{}
)

type StatusUI struct {
	program *tea.Program
	model   statusModel
}

func NewStatusUI(initialMessage string) *StatusUI {
	model := newStatusModel(initialMessage)
	return &StatusUI{
		program: tea.NewProgram(model),
		model:   model,
	}
}

func (s *StatusUI) Run() error {
	_, err := s.program.Run()
	if err != nil {
		return fmt.Errorf("failed to run program: %w", err)
	}
	return nil
}

func (s *StatusUI) Quit() {
	s.program.Send(quitMsg{})
}

func (s *StatusUI) UpdateStatusText(statusText string) {
	s.program.Send(updateStatusMsg(statusText))
}

type statusModel struct {
	spinner    spinner.Model
	statusText string
	quitting   bool
	err        error
}

func newStatusModel(text string) statusModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return statusModel{
		statusText: text,
		spinner:    s,
	}
}

func (m statusModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m statusModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case quitMsg:
		m.spinner.Update(tea.Quit())
		m.quitting = true
		return m, tea.Quit
	case updateStatusMsg:
		m.statusText = string(msg)
		return m, nil
	case spinner.TickMsg:
		if m.quitting {
			return m, nil
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case error:
		m.err = msg
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m statusModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("model error: %v", m.err) // fixed: use fmt.Sprintf instead of fmt.Errorf for simple strings in View method
	}
	str := fmt.Sprintf("%s %s", m.spinner.View(), m.statusText)
	if m.quitting {
		return ""
	}
	return str
}
