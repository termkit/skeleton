package skeleton

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Skeleton is a helper for rendering the Skeleton of the terminal.
type Skeleton struct {
	viewport *viewport.Model

	lockTabs   bool
	currentTab int

	header *header
	KeyMap *keyMap

	pages []tea.Model

	// properties are hold the properties of the Skeleton
	properties *skeletonProperties
}

type skeletonProperties struct {
	borderColor string
}

func defaultSkeletonProperties() *skeletonProperties {
	return &skeletonProperties{
		borderColor: "39",
	}
}

// NewSkeleton returns a new Skeleton.
func NewSkeleton() *Skeleton {
	return &Skeleton{
		properties: defaultSkeletonProperties(),
		viewport:   newTerminalViewport(),
		header:     newHeader(),
		KeyMap:     newKeyMap(),
	}
}

// SetBorderColor sets the border color of the Skeleton.
func (s *Skeleton) SetBorderColor(color string) *Skeleton {
	s.header.SetBorderColor(color)
	s.properties.borderColor = color
	return s
}

// SetInactiveTabTextColor sets the idle tab color of the Skeleton.
func (s *Skeleton) SetInactiveTabTextColor(color string) *Skeleton {
	s.header.SetInactiveTabTextColor(color)
	return s
}

// SetInactiveTabBorderColor sets the idle tab border color of the Skeleton.
func (s *Skeleton) SetInactiveTabBorderColor(color string) *Skeleton {
	s.header.SetInactiveTabBorderColor(color)
	return s
}

// SetActiveTabTextColor sets the active tab color of the Skeleton.
func (s *Skeleton) SetActiveTabTextColor(color string) *Skeleton {
	s.header.SetActiveTabTextColor(color)
	return s
}

// SetActiveTabBorderColor sets the active tab border color of the Skeleton.
func (s *Skeleton) SetActiveTabBorderColor(color string) *Skeleton {
	s.header.SetActiveTabBorderColor(color)
	return s
}

// SetTabLeftPadding sets the left padding of the Skeleton.
func (s *Skeleton) SetTabLeftPadding(padding int) *Skeleton {
	s.header.SetLeftPadding(padding)
	return s
}

// SetTabRightPadding sets the right padding of the Skeleton.
func (s *Skeleton) SetTabRightPadding(padding int) *Skeleton {
	s.header.SetRightPadding(padding)
	return s
}

func (s *Skeleton) SetLockTabs(lock bool) *Skeleton {
	s.header.SetLockTabs(lock)
	s.lockTabs = lock
	return s
}

func (s *Skeleton) GetLockTabs() bool {
	return s.lockTabs
}

func (s *Skeleton) AddPage(pageName string, page tea.Model) *Skeleton {
	s.header.AddCommonHeader(pageName)
	s.pages = append(s.pages, page)
	return s
}

func (s *Skeleton) SetCurrentTab(tab int) *Skeleton {
	s.currentTab = tab
	s.header.SetCurrentTab(tab)
	return s
}

func (s *Skeleton) Init() tea.Cmd {
	self := func() tea.Msg {
		return nil
	}

	inits := make([]tea.Cmd, len(s.pages)+1) // +1 for self
	for i := range s.pages {
		inits[i] = s.pages[i].Init()
	}

	inits[len(s.pages)] = self

	return tea.Batch(inits...)
}

func (s *Skeleton) Update(msg tea.Msg) (*Skeleton, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	s.currentTab = s.header.GetCurrentTab()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.viewport.Width = msg.Width
		s.viewport.Height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.KeyMap.Quit):
			return s, tea.Quit
		case key.Matches(msg, s.KeyMap.SwitchTabLeft):
			if !s.GetLockTabs() {
				s.currentTab = max(s.currentTab-1, 0)
			}
		case key.Matches(msg, s.KeyMap.SwitchTabRight):
			if !s.GetLockTabs() {
				s.currentTab = min(s.currentTab+1, len(s.pages)-1)
			}
		}
	}

	s.header, cmd = s.header.Update(msg)
	cmds = append(cmds, cmd)

	s.pages[s.currentTab], cmd = s.pages[s.currentTab].Update(msg)
	cmds = append(cmds, cmd)

	return s, tea.Batch(cmds...)
}

func (s *Skeleton) View() string {
	base := lipgloss.NewStyle().
		BorderForeground(lipgloss.Color(s.properties.borderColor)).
		Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderTop(false).
		Width(s.viewport.Width - 2)

	return lipgloss.JoinVertical(lipgloss.Top, s.header.View(), base.Render(s.pages[s.currentTab].View()))
}
