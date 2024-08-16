package skeleton

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

// header is a helper for rendering the header of the terminal.
type header struct {
	viewport *viewport.Model

	lockTabs   bool
	currentTab int

	keyMap *keyMap

	headers             []commonHeader
	currentSpecialStyle int

	// properties are hold the properties of the header
	properties *headerProperties
}

type headerProperties struct {
	borderColor        string
	leftTabPadding     int
	rightTabPadding    int
	titleStyleActive   lipgloss.Style
	titleStyleInactive lipgloss.Style
	titleStyleDisabled lipgloss.Style
}

func defaultHeaderProperties() *headerProperties {
	borderColor := "39"
	leftPadding := 2
	rightPadding := 2
	return &headerProperties{
		borderColor:     borderColor,
		leftTabPadding:  leftPadding,
		rightTabPadding: rightPadding,
		titleStyleActive: func() lipgloss.Style {
			b := lipgloss.DoubleBorder()
			b.Right = "├"
			b.Left = "┤"
			return lipgloss.NewStyle().BorderStyle(b).
				PaddingLeft(leftPadding).PaddingRight(rightPadding).
				BorderForeground(lipgloss.Color("205"))
		}(),
		titleStyleInactive: func() lipgloss.Style {
			b := lipgloss.RoundedBorder()
			b.Right = "├"
			b.Left = "┤"
			return lipgloss.NewStyle().BorderStyle(b).
				PaddingLeft(leftPadding).PaddingRight(rightPadding).
				BorderForeground(lipgloss.Color("255"))
		}(),
		titleStyleDisabled: func() lipgloss.Style {
			b := lipgloss.RoundedBorder()
			b.Right = "├"
			b.Left = "┤"
			return lipgloss.NewStyle().BorderStyle(b).
				PaddingLeft(leftPadding).PaddingRight(rightPadding).
				BorderForeground(lipgloss.Color("240")).Foreground(lipgloss.Color("240"))
		}(),
	}
}

type commonHeader struct {
	header string
}

// newHeader returns a new header.
func newHeader() *header {
	return &header{
		properties: defaultHeaderProperties(),
		viewport:   newTerminalViewport(),
		currentTab: 0,
		keyMap:     newKeyMap(),
	}
}

func (h *header) Init() tea.Cmd {
	return nil
}

func (h *header) Update(msg tea.Msg) (*header, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, h.keyMap.SwitchTabLeft):
			if !h.GetLockTabs() {
				h.currentTab = max(h.currentTab-1, 0)
			}
		case key.Matches(msg, h.keyMap.SwitchTabRight):
			if !h.GetLockTabs() {
				h.currentTab = min(h.currentTab+1, len(h.headers)-1)
			}
		}
	}

	return h, nil
}

// View renders the header.
func (h *header) View() string {
	var titleLen int
	for _, title := range h.headers {
		titleLen += len(title.header)
		titleLen += h.properties.leftTabPadding + h.properties.rightTabPadding
		titleLen += 2 // for the border between titles
	}

	var renderedTitles []string
	renderedTitles = append(renderedTitles, "")
	for i, title := range h.headers {
		if h.GetLockTabs() {
			if i == 0 {
				renderedTitles = append(renderedTitles, h.properties.titleStyleActive.Render(title.header))
			} else {
				renderedTitles = append(renderedTitles, h.properties.titleStyleDisabled.Render(title.header))
			}
		} else {
			if i == h.currentTab {
				renderedTitles = append(renderedTitles, h.properties.titleStyleActive.Render(title.header))
			} else {
				renderedTitles = append(renderedTitles, h.properties.titleStyleInactive.Render(title.header))
			}
		}
	}

	leftCorner := lipgloss.JoinVertical(lipgloss.Top, "╭", "│")
	rightCorner := lipgloss.JoinVertical(lipgloss.Top, "╮", "│")
	leftCorner = lipgloss.NewStyle().Foreground(lipgloss.Color(h.properties.borderColor)).Render(leftCorner)
	rightCorner = lipgloss.NewStyle().Foreground(lipgloss.Color(h.properties.borderColor)).Render(rightCorner)

	line := strings.Repeat("─", h.viewport.Width-(titleLen+2))
	line = lipgloss.NewStyle().Foreground(lipgloss.Color(h.properties.borderColor)).Render(line)

	return lipgloss.JoinHorizontal(lipgloss.Bottom, leftCorner, lipgloss.JoinHorizontal(lipgloss.Center, append(renderedTitles, line)...), rightCorner)
}

// SetLeftPadding sets the left padding of the header.
func (h *header) SetLeftPadding(padding int) {
	h.properties.leftTabPadding = padding
	h.properties.titleStyleActive = h.properties.titleStyleActive.PaddingLeft(padding)
	h.properties.titleStyleInactive = h.properties.titleStyleInactive.PaddingLeft(padding)
	h.properties.titleStyleDisabled = h.properties.titleStyleDisabled.PaddingLeft(padding)
}

// SetRightPadding sets the right padding of the header.
func (h *header) SetRightPadding(padding int) {
	h.properties.rightTabPadding = padding
	h.properties.titleStyleActive = h.properties.titleStyleActive.PaddingRight(padding)
	h.properties.titleStyleInactive = h.properties.titleStyleInactive.PaddingRight(padding)
	h.properties.titleStyleDisabled = h.properties.titleStyleDisabled.PaddingRight(padding)
}

// SetInactiveTabTextColor sets the idle tab color of the header.
func (h *header) SetInactiveTabTextColor(color string) {
	h.properties.titleStyleInactive = h.properties.titleStyleInactive.Foreground(lipgloss.Color(color))
}

// SetInactiveTabBorderColor sets the idle tab border color of the header.
func (h *header) SetInactiveTabBorderColor(color string) {
	h.properties.titleStyleInactive = h.properties.titleStyleInactive.BorderForeground(lipgloss.Color(color))
}

// SetActiveTabTextColor sets the active tab color of the header.
func (h *header) SetActiveTabTextColor(color string) {
	h.properties.titleStyleActive = h.properties.titleStyleActive.Foreground(lipgloss.Color(color))
}

// SetActiveTabBorderColor sets the active tab border color of the header.
func (h *header) SetActiveTabBorderColor(color string) {
	h.properties.titleStyleActive = h.properties.titleStyleActive.BorderForeground(lipgloss.Color(color))
}

// SetBorderColor sets the border color of the header.
func (h *header) SetBorderColor(color string) {
	h.properties.borderColor = color
}

func (h *header) SetCurrentTab(tab int) {
	h.currentTab = tab
}

func (h *header) SetLockTabs(lock bool) {
	h.lockTabs = lock
}

func (h *header) GetLockTabs() bool {
	return h.lockTabs
}

func (h *header) GetCurrentTab() int {
	return h.currentTab
}

func (h *header) AddCommonHeader(header string) {
	h.headers = append(h.headers, commonHeader{
		header: header,
	})
}
