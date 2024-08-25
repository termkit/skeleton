package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/termkit/skeleton"
	"os"
	"strings"
)

type fileReader struct {
	skeleton *skeleton.Skeleton
	viewport viewport.Model
	fileName string
}

func (m *fileReader) Init() tea.Cmd {
	return nil
}

func (m *fileReader) CalculatePercent() {
	percent := m.viewport.ScrollPercent() * 100
	m.skeleton.UpdateWidgetValue("percent", fmt.Sprintf("%s | %.2f%%", m.fileName, percent))
}

func (m *fileReader) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case skeleton.IAMActivePage:
		m.skeleton.DeleteAllWidgets()
		m.CalculatePercent()
	case tea.WindowSizeMsg:
		m.viewport.Height = msg.Height - 10
		m.viewport.Width = msg.Width - 2
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+w":
			m.skeleton.DeletePage(m.fileName)
		default:
			m.CalculatePercent()
		}
	}

	m.viewport, _ = m.viewport.Update(msg)

	return m, tea.Batch(cmds...)
}

func (m *fileReader) View() string {
	var b strings.Builder
	b.WriteString(m.viewport.View())
	b.WriteString("\n\n\n")

	helperWindow := lipgloss.NewStyle().UnsetBorderStyle().Foreground(lipgloss.Color("#00ffff"))
	helper := helperWindow.Render(fmt.Sprintf("%s | %s Switch Tab - %s close tab",
		m.skeleton.KeyMap.SwitchTabLeft.Keys(),
		m.skeleton.KeyMap.SwitchTabRight.Keys(), "ctrl+w"))

	b.WriteString(helper)
	return b.String()
}

func newFileReader(skeleton *skeleton.Skeleton, fileName string, filePath string) *fileReader {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("could not load file:", err)
		os.Exit(1)
	}

	vp := viewport.Model{Width: skeleton.GetTerminalWidth() - 2, Height: skeleton.GetTerminalHeight() - 10}
	vp.SetContent(string(content))

	return &fileReader{
		skeleton: skeleton,
		viewport: vp,
		fileName: fileName,
	}
}
