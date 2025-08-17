package main

import (
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	if m.(model).err != nil {
		return "", m.(model).err
	}

	return m.(model).textarea.Value(), nil
}

type (
	errMsg error
)

type model struct {
	textarea textarea.Model
	err      error
	prompt   string
}

func initialModel(prompt string) model {
	ti := textarea.New()
	ti.Placeholder = "Enter your response here..."
	ti.Focus()

	return model{
		textarea: ti,
		err:      nil,
		prompt:   prompt,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return lipgloss.Place(
		100, 100,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.prompt,
			m.textarea.View(),
		),
	)
}
