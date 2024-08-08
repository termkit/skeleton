package alpha

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"sync"
)

// Header is a helper for rendering the Header of the terminal.
type Header struct {
	Viewport *viewport.Model

	KeyMap *KeyMap

	currentTab int

	commonHeaders       []commonHeader
	currentSpecialStyle int
}

var (
	TitleStyleActive = func() lipgloss.Style {
		b := lipgloss.DoubleBorder()
		b.Right = "├"
		b.Left = "┤"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 2).BorderForeground(lipgloss.Color("205"))
	}()

	TitleStyleInactive = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		b.Left = "┤"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 2).BorderForeground(lipgloss.Color("255"))
	}()

	TitleStyleDisabled = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		b.Left = "┤"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 2).BorderForeground(lipgloss.Color("240")).Foreground(lipgloss.Color("240"))
	}()
)

type commonHeader struct {
	header    string
	rawHeader string

	inactiveStyle lipgloss.Style
	activeStyle   lipgloss.Style
}

// Define sync.Once and newHeader should return same instance
var (
	onceHeader sync.Once
	header     *Header
)

// newHeader returns a new Header.
func newHeader() *Header {
	onceHeader.Do(func() {
		header = &Header{
			Viewport:   newTerminalViewport(),
			currentTab: 0,
			KeyMap:     NewKeyMap(),
		}
	})
	return header
}

func (h *Header) SetCurrentTab(tab int) {
	h.currentTab = tab
}

func (h *Header) SetLockTabs(lock bool) {
	h.SetLockTabs(lock)
}

func (h *Header) GetLockTabs() bool {
	return h.GetLockTabs()
}

func (h *Header) GetCurrentTab() int {
	return h.currentTab
}

func (h *Header) AddCommonHeader(header string) {
	h.commonHeaders = append(h.commonHeaders, commonHeader{
		header:    header,
		rawHeader: header,
		//inactiveStyle: inactiveStyle,
		//activeStyle:   activeStyle,
	})
}

type UpdateMsg struct {
	Msg               string
	UpdatingComponent string
}

func (h *Header) Init() tea.Cmd {
	return nil
}

func (h *Header) Update(msg tea.Msg) (*Header, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, h.KeyMap.SwitchTabLeft):
			if !h.GetLockTabs() {
				h.currentTab = max(h.currentTab-1, 0)
			}
		case key.Matches(msg, h.KeyMap.SwitchTabRight):
			if !h.GetLockTabs() {
				h.currentTab = min(h.currentTab+1, len(h.commonHeaders)-1)
			}
		}
	}

	return h, nil
}

// View renders the Header.
func (h *Header) View() string {
	var titleLen int
	for _, title := range h.commonHeaders {
		titleLen += len(title.rawHeader)
		titleLen += TitleStyleActive.GetPaddingLeft() + TitleStyleActive.GetPaddingRight()
		titleLen += 2 // for the border between titles
	}

	var renderedTitles []string
	renderedTitles = append(renderedTitles, "")
	for i, title := range h.commonHeaders {
		if h.GetLockTabs() {
			if i == 0 {
				renderedTitles = append(renderedTitles, TitleStyleActive.Render(title.header))
			} else {
				renderedTitles = append(renderedTitles, TitleStyleDisabled.Render(title.header))
			}
		} else {
			if i == h.currentTab {
				renderedTitles = append(renderedTitles, TitleStyleActive.Render(title.header))
			} else {
				renderedTitles = append(renderedTitles, TitleStyleInactive.Render(title.header))
			}
		}
	}

	leftCorner := lipgloss.JoinVertical(lipgloss.Top, "╭", "│")
	rightCorner := lipgloss.JoinVertical(lipgloss.Top, "╮", "│")
	leftCorner = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Render(leftCorner)
	rightCorner = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Render(rightCorner)

	line := strings.Repeat("─", h.Viewport.Width-(titleLen+2))
	line = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Render(line)

	return lipgloss.JoinHorizontal(lipgloss.Bottom, leftCorner, lipgloss.JoinHorizontal(lipgloss.Center, append(renderedTitles, line)...), rightCorner)
}
