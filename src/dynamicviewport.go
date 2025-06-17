package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type DynViewPort struct {
	viewport    viewport.Model
	msgrenderer *glamour.TermRenderer
	spinner     spinner.Model
	messages    []string
	yoffset     int
	speaker     string
	latest      string
	spin        bool
}

func BuildDynViewPort(initialMsg string, yoffset int) (m DynViewPort) {
	viewport := viewport.New(32, 5)

	msgrenderer, _ := glamour.NewTermRenderer(glamour.WithAutoStyle())

	spin := spinner.New()
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#EE8080"))
	spin.Spinner = spinner.MiniDot

	return DynViewPort{
		viewport:    viewport,
		msgrenderer: msgrenderer,
		spinner:     spin,
		messages:    []string{initialMsg},
		yoffset:     yoffset,
		speaker:     "",
		latest:      "",
		spin:        true,
	}
}

func (m DynViewPort) Init() tea.Cmd {
	return m.spinner.Tick
}

type speakerMsg struct {
	speaker string
	msg     string
	done    bool
}

func (m DynViewPort) mdrender(msg string) string {
	md, _ := m.msgrenderer.Render(msg)
	trimmed := strings.Trim(md, "\n")
	return trimmed
}
func (m DynViewPort) renderContent() DynViewPort {
	var content []string
	var spin = ""
	if m.spin {
		spin = m.spinner.View()
	}
	if m.latest != "" {
		content = append(m.messages, spin+m.mdrender(fmt.Sprintf("**%s**: %s", m.speaker, m.latest)))
	}
	m.viewport.SetContent(
		lipgloss.NewStyle().Width(
			m.viewport.Width).Render(
			strings.Join(content, "\n")))
	return m
}
func (m DynViewPort) Update(msg tea.Msg) (DynViewPort, tea.Cmd) {
	var (
		vpCmd tea.Cmd
		spCmd tea.Cmd
	)
	// update viewport
	m.viewport, vpCmd = m.viewport.Update(msg)
	// update spinner
	m.spinner, spCmd = m.spinner.Update(msg)
	// other message types
	switch msg := msg.(type) {
	// dynamic to resizing
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - m.yoffset - lipgloss.Height(gap)
		m.viewport.GotoBottom()
	case speakerMsg:
		if msg.speaker == m.speaker {
			m.latest = m.latest + msg.msg
			m.spin = !msg.done
		} else {
			if m.latest != "" {
				md := m.mdrender(fmt.Sprintf("**%s**: %s", m.speaker, m.latest))
				m.messages = append(m.messages, md)
			}
			m.speaker = msg.speaker
			m.latest = msg.msg
			m.spin = true
		}
		m = m.renderContent()
		m.viewport.GotoBottom()
	}
	return m, tea.Batch(vpCmd, spCmd)
}

func (m DynViewPort) View() string {
	m = m.renderContent()
	return m.viewport.View()
}
