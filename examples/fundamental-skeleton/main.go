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

func (m tinyModel) Init() tea.Cmd { return nil }
func (m tinyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
func (m tinyModel) View() string {
	verticalCenter := m.skeleton.GetTerminalHeight()/2 - 3 // 3 for vertical length of the title section
	requiredNewLines := strings.Repeat("\n", verticalCenter)

	return fmt.Sprintf("%s%s | %d x %d", requiredNewLines, m.title, m.skeleton.GetTerminalWidth(), m.skeleton.GetTerminalHeight())
}

// -----------------------------------------------------------------------------
// Main Model
// The Main Model is the main model for the program. It contains the skeleton and the tab models.

type mainModel struct {
	skeleton *skeleton.Skeleton
}

func (m *mainModel) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		tea.SetWindowTitle("Basic Tab Example"),
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
	skel := skeleton.NewSkeleton()

	// Add tabs (pages)
	skel.AddPage("first", "First Tab", newTinyModel(skel, "First"))
	skel.AddPage("second", "Second Tab", newTinyModel(skel, "Second"))
	skel.AddPage("third", "Third Tab", newTinyModel(skel, "Third"))

	// Add a widget to entire screen
	// Battery level is hardcoded. You can use a library to get the battery level of your system.
	skel.AddWidget("battery", "Battery %92") // Add a widget to entire screen

	// Add current time
	skel.AddWidget("time", time.Now().Format("15:04:05"))

	// Update the time widget every second
	go func() {
		time.Sleep(time.Second)
		for {
			skel.UpdateWidgetValue("time", time.Now().Format("15:04:05"))
			time.Sleep(time.Second)
		}
	}()

	model := &mainModel{
		skeleton: skel,
	}

	p := tea.NewProgram(model)
	if err := p.Start(); err != nil {
		panic(err)
	}
}
