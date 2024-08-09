package alpha

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sync"
)

// Alpha is a helper for rendering the Alpha of the terminal.
type Alpha struct {
	viewport *viewport.Model

	lockTabs   bool
	currentTab int

	header *Header
	KeyMap *keyMap

	pages []tea.Model
}

var (
	onceSkeletonAlpha sync.Once
	skeletonAlpha     *Alpha
)

// SkeletonAlpha returns a new Alpha.
func SkeletonAlpha() *Alpha {
	onceSkeletonAlpha.Do(func() {
		skeletonAlpha = &Alpha{
			viewport: newTerminalViewport(),
			header:   newHeader(),
			KeyMap:   newKeyMap(),
		}
	})
	return skeletonAlpha
}

func (a *Alpha) Init() tea.Cmd {
	self := func() tea.Msg {
		return nil
	}

	inits := make([]tea.Cmd, len(a.pages)+1) // +1 for self
	for i := range a.pages {
		inits[i] = a.pages[i].Init()
	}

	inits[len(a.pages)] = self

	return tea.Batch(inits...)
}

func (a *Alpha) Update(msg tea.Msg) (*Alpha, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	a.currentTab = a.header.GetCurrentTab()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.viewport.Width = msg.Width
		a.viewport.Height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, a.KeyMap.Quit):
			return a, tea.Quit
		case key.Matches(msg, a.KeyMap.SwitchTabLeft):
			if !a.GetLockTabs() {
				a.currentTab = max(a.currentTab-1, 0)
			}
		case key.Matches(msg, a.KeyMap.SwitchTabRight):
			if !a.GetLockTabs() {
				a.currentTab = min(a.currentTab+1, len(a.pages)-1)
			}
		}
	}

	a.header, cmd = a.header.Update(msg)
	cmds = append(cmds, cmd)

	a.pages[a.currentTab], cmd = a.pages[a.currentTab].Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a *Alpha) View() string {
	base := lipgloss.NewStyle().
		BorderForeground(lipgloss.Color("39")).
		Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderTop(false).
		Width(a.viewport.Width - 2)

	return lipgloss.JoinVertical(lipgloss.Top, a.header.View(), base.Render(a.pages[a.currentTab].View()))
}

func (a *Alpha) SetLockTabs(lock bool) {
	a.header.SetLockTabs(lock)
	a.lockTabs = lock
}

func (a *Alpha) GetLockTabs() bool {
	return a.lockTabs
}

func (a *Alpha) AddPage(pageName string, page tea.Model) {
	a.header.AddCommonHeader(pageName)
	a.pages = append(a.pages, page)
}

func (a *Alpha) SetCurrentTab(tab int) {
	a.currentTab = tab
	a.header.SetCurrentTab(tab)
}
