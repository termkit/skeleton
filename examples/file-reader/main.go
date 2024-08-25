package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/termkit/skeleton"
)

type mainModel struct {
	skeleton *skeleton.Skeleton
}

func (m *mainModel) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		tea.SetWindowTitle("File Reader"),
		m.skeleton.Init(),
	)
}

func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.skeleton, cmd = m.skeleton.Update(msg)
	return m, cmd
}

func (m *mainModel) View() string {
	return m.skeleton.View()
}

func main() {
	s := skeleton.NewSkeleton()
	s.SetPagePosition(lipgloss.Left)

	s.SetBorderColor("#ff0055")
	s.SetActiveTabBorderColor("#00aaff")

	s.AddPage("explorer", "Explorer", newExplorer(s))

	m := &mainModel{skeleton: s}

	if err := tea.NewProgram(m).Start(); err != nil {
		panic(err)
	}
}
