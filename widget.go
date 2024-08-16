package skeleton

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

// widget is a helper for rendering the widget of the terminal.
type widget struct {
	viewport *viewport.Model

	widgets []commonWidget

	// properties are hold the properties of the widget
	properties *widgetProperties
}

type commonWidget struct {
	widget string
}

type widgetProperties struct {
	borderColor     string
	leftTabPadding  int
	rightTabPadding int
	widgetStyle     lipgloss.Style
}

func defaultWidgetProperties() *widgetProperties {
	borderColor := "39"
	leftPadding := 2
	rightPadding := 2
	return &widgetProperties{
		borderColor:     borderColor,
		leftTabPadding:  leftPadding,
		rightTabPadding: rightPadding,
		widgetStyle: func() lipgloss.Style {
			b := lipgloss.RoundedBorder()
			b.Right = "├"
			b.Left = "┤"
			return lipgloss.NewStyle().BorderStyle(b).
				PaddingLeft(leftPadding).PaddingRight(rightPadding).
				BorderForeground(lipgloss.Color("49"))
		}(),
	}
}

// newWidget returns a new Widget.
func newWidget() *widget {
	return &widget{
		properties: defaultWidgetProperties(),
		viewport:   newTerminalViewport(),
	}
}

// SetBorderColor sets the border color of the Widget.
func (w *widget) SetBorderColor(color string) *widget {
	w.properties.borderColor = color
	return w
}

// SetWidgetBorderColor sets the border color of the Widget.
func (w *widget) SetWidgetBorderColor(color string) *widget {
	w.properties.widgetStyle = w.properties.widgetStyle.BorderForeground(lipgloss.Color(color))
	return w
}

// SetLeftPadding sets the left padding of the Widget.
func (w *widget) SetLeftPadding(padding int) *widget {
	w.properties.leftTabPadding = padding
	w.properties.widgetStyle = w.properties.widgetStyle.PaddingLeft(padding)
	return w
}

// SetRightPadding sets the right padding of the Widget.
func (w *widget) SetRightPadding(padding int) *widget {
	w.properties.rightTabPadding = padding
	w.properties.widgetStyle = w.properties.widgetStyle.PaddingRight(padding)
	return w
}

func (w *widget) AddWidget(widget string) *widget {
	w.widgets = append(w.widgets, commonWidget{widget: widget})
	return w
}

func (w *widget) Init() tea.Cmd {
	return nil
}

func (w *widget) Update(msg tea.Msg) (*widget, tea.Cmd) {
	return w, nil
}

func (w *widget) View() string {
	var widgetLen int
	for _, widget := range w.widgets {
		widgetLen += len(widget.widget)
		widgetLen += w.properties.leftTabPadding + w.properties.rightTabPadding
		widgetLen += 2 // for the border between widgets
	}

	var renderedWidgets []string
	for _, wgt := range w.widgets {
		renderedWidgets = append(renderedWidgets, w.properties.widgetStyle.Render(wgt.widget))
	}

	leftCorner := lipgloss.JoinVertical(lipgloss.Top, "│", "╰")
	rightCorner := lipgloss.JoinVertical(lipgloss.Top, "│", "╯")
	leftCorner = lipgloss.NewStyle().Foreground(lipgloss.Color(w.properties.borderColor)).Render(leftCorner)
	rightCorner = lipgloss.NewStyle().Foreground(lipgloss.Color(w.properties.borderColor)).Render(rightCorner)

	line := strings.Repeat("─", w.viewport.Width-(widgetLen+2))
	line = lipgloss.NewStyle().Foreground(lipgloss.Color(w.properties.borderColor)).Render(line)

	var bottom []string
	bottom = append(bottom, line)
	bottom = append(bottom, renderedWidgets...)

	position := lipgloss.Center
	if len(w.widgets) > 0 {
		position = lipgloss.Top
	}

	return lipgloss.JoinHorizontal(position, leftCorner, lipgloss.JoinHorizontal(lipgloss.Center, bottom...), rightCorner)
}
