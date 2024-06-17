package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type model struct {
	userInput textinput.Model
	confirmed bool
}

func newModel(initialValue string) (m model) {
	i := textinput.New()
	i.Prompt = initialValue
	i.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	i.CursorEnd()
	i.Focus()
	m.userInput = i
	return
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.Type {
		case tea.KeyCtrlC, tea.KeyEscape:
			return m, tea.Quit
		case tea.KeyEnter:
			v := strings.ToLower(m.userInput.Value())
			if v == "y" || v == "yes" {
				m.confirmed = true
				return m, tea.Quit
			}
			if v == "n" || v == "no" {
				return m, tea.Quit
			}
			m.userInput.Reset()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.userInput, cmd = m.userInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf("\n%s", m.userInput.View())
}

func Confirm(message string) (bool, error) {
	m := newModel(message)

	fM, err := tea.NewProgram(m).Run()
	if err != nil {
		return false, err
	}
	return fM.(model).confirmed, nil
}
