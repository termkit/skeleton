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
	Viewport *viewport.Model

	header   *Header
	lockTabs bool

	KeyMap *KeyMap

	currentTab int
	Pages      []tea.Model
}

var lockTabs bool

func SetLockTabs(lock bool) {
	lockTabs = lock
}

func GetLockTabs() bool {
	return lockTabs
}

func (a *Alpha) AddPage(title Title, page tea.Model) {
	a.header.AddCommonHeader(title.Title)
	a.Pages = append(a.Pages, page)
}

type Title struct {
	Title string
	Style TitleStyle
}

type TitleStyle struct {
	Active   lipgloss.Style
	Inactive lipgloss.Style
}

var (
	onceSkeletonAlpha sync.Once
	skeletonAlpha     *Alpha
)

// NewSkeletonAlpha returns a new Alpha.
func NewSkeletonAlpha() *Alpha {
	onceSkeletonAlpha.Do(func() {
		skeletonAlpha = &Alpha{
			Viewport: newTerminalViewport(),
			header:   newHeader(),
			KeyMap:   NewKeyMap(),
		}
	})
	return skeletonAlpha
}

type SwitchTab struct {
	Tab int
}

func (a *Alpha) SetCurrentTab(tab int) {
	a.currentTab = tab
	a.header.SetCurrentTab(tab)
}

func (a *Alpha) Init() tea.Cmd {
	self := func() tea.Msg {
		return SwitchTab{}
	}

	inits := make([]tea.Cmd, len(a.Pages)+1) // +1 for self
	for i := range a.Pages {
		inits[i] = a.Pages[i].Init()
	}

	inits[len(a.Pages)] = self

	return tea.Batch(inits...)
}

func (a *Alpha) Update(msg tea.Msg) (*Alpha, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	a.currentTab = a.header.GetCurrentTab()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.Viewport.Width = msg.Width
		a.Viewport.Height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, a.KeyMap.Quit):
			return a, tea.Quit
		case key.Matches(msg, a.KeyMap.SwitchTabLeft):
			if !GetLockTabs() {
				a.currentTab = max(a.currentTab-1, 0)
			}
		case key.Matches(msg, a.KeyMap.SwitchTabRight):
			if !GetLockTabs() {
				a.currentTab = min(a.currentTab+1, len(a.Pages)-1)
			}
		}
	}

	a.header, cmd = a.header.Update(msg)
	cmds = append(cmds, cmd)

	a.Pages[a.currentTab], cmd = a.Pages[a.currentTab].Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a *Alpha) View() string {
	base := lipgloss.NewStyle().
		BorderForeground(lipgloss.Color("39")).
		Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderTop(false).
		Width(a.Viewport.Width - 2)

	return lipgloss.JoinVertical(lipgloss.Top, a.header.View(), base.Render(a.Pages[a.currentTab].View()))
}
