package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type chatMsg struct {
	speaker string
	msg     string
}

func buildUserTextArea() textarea.Model {
	textArea := textarea.New()
	textArea.Placeholder = "Send a message..."
	textArea.Focus()

	textArea.Prompt = "| "
	textArea.CharLimit = 280

	textArea.SetWidth(30)
	textArea.SetHeight(2)

	// Remove cursor line styling
	textArea.FocusedStyle.CursorLine = lipgloss.NewStyle()

	textArea.ShowLineNumbers = false
	textArea.KeyMap.InsertNewline.SetEnabled(false)
	return textArea
}

func buildViewPort() viewport.Model {
	viewport := viewport.New(32, 5)
	viewport.SetContent(`Ask me anything!`)
	return viewport
}

func buildSpinner() spinner.Model {
	sp := spinner.New()
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#EE8080"))
	sp.Spinner = spinner.MiniDot
	return sp
}

type model struct {
	viewport   viewport.Model
	textarea   textarea.Model
	spinner    spinner.Model
	messages   []string
	latest     string
	speaker    string
	mdrenderer *glamour.TermRenderer
	msgChan    chan chatMsg
	err        error
}

func BuildChatUIModel() model {
	renderer, _ := glamour.NewTermRenderer(glamour.WithAutoStyle())
	return model{
		textarea:   buildUserTextArea(),
		viewport:   buildViewPort(),
		messages:   []string{},
		latest:     "",
		spinner:    buildSpinner(),
		msgChan:    make(chan chatMsg, 10),
		mdrenderer: renderer,
		speaker:    "???",
		err:        nil,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, m.spinner.Tick)
}
func (m model) renderMessages() {
	messages := append(m.messages, fmt.Sprintf("\n%s Thinking...", m.spinner.View()))
	m.viewport.SetContent(
		lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(messages, "")))
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		spCmd tea.Cmd
	)
	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)

		if len(m.messages) > 0 {
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "")))
		}
		m.viewport.GotoBottom()
	case spinner.TickMsg:
		m.spinner, spCmd = m.spinner.Update(msg)
		messages := append(m.messages, fmt.Sprintf("\n%s Thinking...", m.spinner.View()))
		m.viewport.SetContent(
			lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(messages, "")))

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			m.msgChan <- chatMsg{"You", m.textarea.Value()}
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}
	case tea.Msg:
		select {
		case newChatMsg := <-m.msgChan:
			if newChatMsg.speaker == m.speaker {
				m.latest = m.latest + newChatMsg.msg
				m.speaker = newChatMsg.speaker
			} else {
				md, _ := m.mdrenderer.Render("**" + newChatMsg.speaker + "**: " + newChatMsg.msg)
				trimmed := strings.TrimRight(md, "\n")
				m.messages = append(m.messages, trimmed)
				m.latest = ""
			}
			m.renderMessages()
		default:
		}
	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}
	return m, tea.Batch(tiCmd, vpCmd, spCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		m.textarea.View(),
	)
}
