package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/termkit/skeleton"
)

func main() {
	s := skeleton.NewSkeleton()
	s.SetPagePosition(lipgloss.Left)

	s.SetBorderColor("#ff0055")
	s.SetActiveTabBorderColor("#00aaff")

	s.AddPage("explorer", "Explorer", newExplorer(s))

	if err := tea.NewProgram(s).Start(); err != nil {
		panic(err)
	}
}
