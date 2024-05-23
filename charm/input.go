package charm

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Input struct {
	text         string
	placeholder  string
	initialValue string
	hidden       bool
	required     bool
}

var (
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00"))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")). // Red color
			Bold(true)

	inputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
)

// Function to create and run the bubbletea model
func GetSensitiveInput(placeholder string, defaultValue string) (string, error) {
	initialModel := model{
		textInput: textinput.New(),
	}
	initialModel.textInput.Placeholder = placeholder
	initialModel.textInput.SetValue(defaultValue)
	initialModel.textInput.EchoMode = textinput.EchoPassword // Hides the input
	initialModel.textInput.Focus()

	p := tea.NewProgram(initialModel)
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	return finalModel.(model).textInput.Value(), nil
}

// Bubbletea model
type model struct {
	textInput textinput.Model
	err       error
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.textInput.Value() == "" {
				m.err = fmt.Errorf("input is required")
				return m, nil
			}
			return m, tea.Quit
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var errMsg string
	if m.err != nil {
		errMsg = ErrorStyle.Render(m.err.Error()) + "\n\n"
	}

	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		titleStyle.Render("Enter your token:"),
		inputStyle.Render(m.textInput.View()),
		errMsg,
	)
}
