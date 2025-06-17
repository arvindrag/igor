package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ollama/ollama/api"
)

type model struct {
	viewport DynViewPort
	textarea textarea.Model
	err      error
	oclient  OllamaClient
	cmdchan  *chan tea.Msg
}

func buildUserTextArea() textarea.Model {
	textArea := textarea.New()
	textArea.Placeholder = "Send a message..."
	textArea.Focus()

	textArea.Prompt = "｜ "
	textArea.CharLimit = 280

	textArea.SetWidth(30)
	textArea.SetHeight(2)

	// Remove cursor line styling
	textArea.FocusedStyle.CursorLine = lipgloss.NewStyle()

	textArea.ShowLineNumbers = false
	textArea.KeyMap.InsertNewline.SetEnabled(false)
	return textArea
}

func BuildChatUIModel(cmdchan *chan tea.Msg) model {
	textarea := buildUserTextArea()
	return model{
		textarea: textarea,
		viewport: BuildDynViewPort("Ask me anything!", textarea.Height()),
		err:      nil,
		oclient:  BuildOllamaClient("llama3.2"),
		cmdchan:  cmdchan,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.viewport.Init(), textarea.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)
	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			m.viewport, vpCmd = m.viewport.Update(speakerMsg{"user", m.textarea.Value(), true})
			go m.oclient.Generate(m.textarea.Value(), func(resp api.GenerateResponse) error {
				*m.cmdchan <- speakerMsg{"jenny", resp.Response, resp.Done}
				return nil
			}, true)
			m.textarea.Reset()
		}
	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}
	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		m.textarea.View(),
	)
}
