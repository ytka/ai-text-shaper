package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/errors"
)

// Confirm prompts the user with a message and waits for a confirmation.
func Confirm(message string) (bool, error) {
	m := newConfirmModel(message)

	fM, err := tea.NewProgram(m).Run()
	if err != nil {
		return false, errors.Wrap(err, "failed to run the confirmation program")
	}
	cm, ok := fM.(confirmModel)
	if !ok {
		return false, errors.New("failed to assert type confirmModel")
	}
	return cm.confirmed, nil
}

type confirmModel struct {
	userInput textinput.Model
	confirmed bool
}

func newConfirmModel(initialPrompt string) confirmModel {
	i := textinput.New()
	i.Prompt = initialPrompt
	i.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	i.CursorEnd()
	i.Focus()
	return confirmModel{userInput: i}
}

func (m confirmModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.Type { //nolint:exhaustive
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

func (m confirmModel) View() string {
	return "\n" + m.userInput.View()
}
