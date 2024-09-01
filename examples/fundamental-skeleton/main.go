package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/termkit/skeleton"

	tea "github.com/charmbracelet/bubbletea"
)

// -----------------------------------------------------------------------------
// Tiny Model
// The Tiny Model is a sub-model for the tabs. It's a simple model that just shows the title of the tab.

// tinyModel is a sub-model for the tabs
type tinyModel struct {
	skeleton *skeleton.Skeleton
	title    string
}

// newTinyModel returns a new tinyModel
func newTinyModel(skeleton *skeleton.Skeleton, title string) *tinyModel {
	return &tinyModel{
		skeleton: skeleton,
		title:    title,
	}
}

func (m tinyModel) Init() tea.Cmd {
	return nil
}
func (m tinyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
func (m tinyModel) View() string {
	verticalCenter := m.skeleton.GetTerminalHeight()/2 - 3 // 3 for vertical length of the title section
	requiredNewLines := strings.Repeat("\n", verticalCenter)

	return fmt.Sprintf("%s%s | %d x %d", requiredNewLines, m.title, m.skeleton.GetTerminalWidth(), m.skeleton.GetTerminalHeight())
}

// -----------------------------------------------------------------------------
// Main Program
func main() {
	s := skeleton.NewSkeleton()

	// Add tabs (pages)
	s.AddPage("first", "First Tab", newTinyModel(s, "First"))
	s.AddPage("second", "Second Tab", newTinyModel(s, "Second"))
	s.AddPage("third", "Third Tab", newTinyModel(s, "Third"))

	// Add a widget to entire screen ( Optional )
	// Battery level is hardcoded. You can use a library to get the battery level of your system.
	s.AddWidget("battery", "Battery %92") // Add a widget to entire screen

	// Add current time ( Optional )
	s.AddWidget("time", time.Now().Format("15:04:05"))

	// Update the time widget every second ( Optional )
	go func() {
		time.Sleep(time.Second)
		for {
			s.UpdateWidgetValue("time", time.Now().Format("15:04:05"))
			time.Sleep(time.Second)
		}
	}()

	p := tea.NewProgram(s)
	if err := p.Start(); err != nil {
		panic(err)
	}
}
