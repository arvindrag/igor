package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

const gap = "\n\n"

func main() {
	chatUI := tea.NewProgram(BuildChatUIModel(), tea.WithAltScreen())

	if _, err := chatUI.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

// import (
// 	"fmt"
// 	"os"
// 	"strings"

// 	"github.com/charmbracelet/bubbles/spinner"
// 	"github.com/charmbracelet/bubbles/textarea"
// 	"github.com/charmbracelet/bubbles/viewport"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"
// )

// type model struct {
// 	viewport viewport.Model
// 	textarea textarea.Model
// 	spinner  spinner.Model
// 	messages []string
// 	width    int
// 	height   int
// 	input    string
// 	msgChan  chan string
// }

// type tickMsg struct{}
// type incomingMsg string

// func initialModel() model {
// 	ta := textarea.New()
// 	ta.Placeholder = "Type something..."
// 	ta.Focus()
// 	ta.ShowLineNumbers = false
// 	ta.Prompt = "┃ "

// 	sp := spinner.New()
// 	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#EE8080"))
// 	sp.Spinner = spinner.MiniDot

// 	return model{
// 		textarea: ta,
// 		spinner:  sp,
// 		msgChan:  make(chan string, 10),
// 	}
// }

// func (m model) Init() tea.Cmd {
// 	return tea.Batch(
// 		tea.EnterAltScreen,
// 		m.spinner.Tick,
// 		// m.listenForMessages(),
// 	)
// }

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	var cmds []tea.Cmd

// 	switch msg := msg.(type) {
// 	case tea.WindowSizeMsg:
// 		m.width = msg.Width
// 		m.height = msg.Height

// 		m.viewport = viewport.New(msg.Width, msg.Height-m.textarea.Height())
// 		m.viewport.SetContent(strings.Join(m.messages, "\n"))

// 	case tea.KeyMsg:
// 		switch msg.Type {
// 		case tea.KeyCtrlC, tea.KeyEsc:
// 			return m, tea.Quit
// 		case tea.KeyEnter:
// 			userMsg := m.textarea.Value()
// 			if strings.TrimSpace(userMsg) != "" {
// 				m.messages = append(m.messages, "You: "+userMsg)
// 				m.viewport.SetContent(strings.Join(m.messages, "\n"))
// 				m.textarea.Reset()
// 			}
// 		}
// 	case spinner.TickMsg:
// 		var cmd tea.Cmd
// 		m.spinner, cmd = m.spinner.Update(msg)
// 		cmds = append(cmds, cmd)
// 		m.viewport.SetContent(strings.Join(m.messages, "\n") + fmt.Sprintf("\n%s Thinking...", m.spinner.View()))
// 	case tea.Msg:
// 		select {
// 		case newMsg := <-m.msgChan:
// 			m.messages = append(m.messages, "Bot: "+string(newMsg))
// 			m.viewport.SetContent(strings.Join(m.messages, "\n") + fmt.Sprintf("\n%s Thinking...", m.spinner.View()))
// 		default:
// 		}
// 	}

// 	// Update textarea and viewport
// 	var taCmd, vpCmd tea.Cmd
// 	m.textarea, taCmd = m.textarea.Update(msg)
// 	m.viewport, vpCmd = m.viewport.Update(msg)

// 	cmds = append(cmds, taCmd, vpCmd)

// 	// Spinner should keep ticking
// 	cmds = append(cmds, m.spinner.Tick)

// 	return m, tea.Batch(cmds...)
// }

// func (m model) View() string {
// 	content := m.viewport.View()
// 	return lip(content + "\n\n" + m.textarea.View())
// }

// func lip(s string) string {
// 	return strings.TrimSuffix(s, "\n")
// }

// func main() {
// 	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
// 	if err := p.Start(); err != nil {
// 		fmt.Println("Error running program:", err)
// 		os.Exit(1)
// 	}
// }
