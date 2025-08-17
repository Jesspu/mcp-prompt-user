package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F8")).
			Padding(1, 0, 0, 0)

	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2)

	instructionsStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Padding(1, 0)
)

func RunPrompt(prompt string) (string, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return "", err
	}
	defer tty.Close()

	p := tea.NewProgram(initialModel(prompt), tea.WithInput(tty), tea.WithOutput(tty), tea.WithAltScreen())

	m, err := p.Run()
	if err != nil {
		return "", err
	}

	if finalModel, ok := m.(model); ok {
		if finalModel.err != nil {
			return "", finalModel.err
		}
		if !finalModel.submitted {
			return "", fmt.Errorf("prompt cancelled by user")
		}
		return finalModel.textarea.Value(), nil
	}

	return "", fmt.Errorf("unexpected error during prompt")
}

type (
	errMsg error
)

type model struct {
	textarea  textarea.Model
	err       error
	prompt    string
	submitted bool
	width     int
	height    int
}

func initialModel(prompt string) model {
	ti := textarea.New()
	ti.Placeholder = "Enter your response here..."
	ti.Focus()
	ti.SetWidth(50)
	ti.SetHeight(5)

	return model{
		textarea:  ti,
		err:       nil,
		prompt:    prompt,
		submitted: false,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.submitted = true
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	title := titleStyle.Render("User Input Required")
	prompt := promptStyle.Render(m.prompt)
	instructions := instructionsStyle.Render("Press Enter to submit, Esc to cancel")

	ui := lipgloss.JoinVertical(lipgloss.Left,
		title,
		prompt,
		m.textarea.View(),
		instructions,
	)

	borderedUI := containerStyle.Render(ui)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		borderedUI,
	)
}
