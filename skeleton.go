package skeleton

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

// Skeleton is a helper for rendering the Skeleton of the terminal.
type Skeleton struct {
	termReady bool

	viewport *viewport.Model

	lockTabs   bool
	currentTab int

	header *header
	widget *widget
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
		widget:     newWidget(),
		KeyMap:     newKeyMap(),
	}
}

// SetBorderColor sets the border color of the Skeleton.
func (s *Skeleton) SetBorderColor(color string) *Skeleton {
	s.header.SetBorderColor(color)
	s.widget.SetBorderColor(color)
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

// SetWidgetBorderColor sets the border color of the Widget.
func (s *Skeleton) SetWidgetBorderColor(color string) *Skeleton {
	s.widget.SetWidgetBorderColor(color)
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

// SetWidgetLeftPadding sets the left padding of the Skeleton.
func (s *Skeleton) SetWidgetLeftPadding(padding int) *Skeleton {
	s.widget.SetLeftPadding(padding)
	return s
}

// SetWidgetRightPadding sets the right padding of the Skeleton.
func (s *Skeleton) SetWidgetRightPadding(padding int) *Skeleton {
	s.widget.SetRightPadding(padding)
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

// AddWidget adds a new widget to the Skeleton.
func (s *Skeleton) AddWidget(key string, value string) *Skeleton {
	s.widget.AddWidget(key, value)
	return s
}

// UpdateWidgetValue updates the Value content by the given key.
func (s *Skeleton) UpdateWidgetValue(key string, value string) *Skeleton {
	s.widget.UpdateWidgetValue(key, value)
	return s
}

// DeleteWidget deletes the Value by the given key.
func (s *Skeleton) DeleteWidget(key string) *Skeleton {
	s.widget.DeleteWidget(key)
	return s
}

func (s *Skeleton) SetCurrentTab(tab int) *Skeleton {
	s.currentTab = tab
	s.header.SetCurrentTab(tab)
	return s
}

func (s *Skeleton) Init() tea.Cmd {
	if len(s.pages) == 0 {
		panic("skeleton: no pages added, please add at least one page")
	}

	self := func() tea.Msg {
		return nil
	}

	inits := make([]tea.Cmd, len(s.pages)+3) // +3 ( for self, header, Value)
	for i := range s.pages {
		inits[i] = s.pages[i].Init()
	}

	// and init self, header and Value
	inits[len(s.pages)] = self
	inits[len(s.pages)+1] = s.header.Init()
	inits[len(s.pages)+2] = s.widget.Init()

	return tea.Batch(inits...)
}

func (s *Skeleton) Update(msg tea.Msg) (*Skeleton, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	s.currentTab = s.header.GetCurrentTab()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !s.termReady {
			if msg.Width > 0 && msg.Height > 0 {
				s.termReady = true
			}
		}
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

	s.widget, cmd = s.widget.Update(msg)
	cmds = append(cmds, cmd)

	s.pages[s.currentTab], cmd = s.pages[s.currentTab].Update(msg)
	cmds = append(cmds, cmd)

	return s, tea.Batch(cmds...)
}

func (s *Skeleton) View() string {
	if !s.termReady {
		return "setting up terminal..."
	}

	base := lipgloss.NewStyle().
		BorderForeground(lipgloss.Color(s.properties.borderColor)).
		Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderTop(false).BorderBottom(false).
		Width(s.viewport.Width - 2)

	body := s.pages[s.currentTab].View()

	bodyHeight := s.viewport.Height - 5 // 6 is the header height and Value height
	if len(s.widget.widgets) > 0 {
		bodyHeight -= 1
	}
	if lipgloss.Height(body) < bodyHeight {
		body += strings.Repeat("\n", bodyHeight-lipgloss.Height(body))
	}

	return lipgloss.JoinVertical(lipgloss.Top, s.header.View(), base.Render(body), s.widget.View())
}
