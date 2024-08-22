package skeleton

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

// widget is a helper for rendering the widget of the terminal.
type widget struct {
	termReady bool

	viewport *viewport.Model

	// widgets are hold the widgets
	widgets []*commonWidget

	// properties are hold the properties of the widget
	properties *widgetProperties

	// updateChan is hold the update channel
	updateChan chan any
}

// newWidget returns a new Widget.
func newWidget() *widget {
	return &widget{
		properties: defaultWidgetProperties(),
		viewport:   newTerminalViewport(),
		updateChan: make(chan any),
	}
}

type commonWidget struct {
	Key   string // Key is the name of the Value
	Value string // Value is the content of the Value
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

func (w *widget) AddWidget(key string, value string) {
	go func() {
		w.updateChan <- AddNewWidget{
			Key:   key,
			Value: value,
		}
	}()
}

// GetWidget returns the Value by the given key.
func (w *widget) GetWidget(key string) *commonWidget {
	for _, widget := range w.widgets {
		if widget.Key == key {
			return widget
		}
	}

	return nil
}

// UpdateWidgetValue updates the Value content by the given key.
func (w *widget) UpdateWidgetValue(key string, value string) {
	go func() {
		w.updateChan <- UpdateWidgetContent{
			Key:   key,
			Value: value,
		}
	}()
}

// DeleteWidget deletes the Value by the given key.
func (w *widget) DeleteWidget(key string) {
	go func() {
		w.updateChan <- DeleteWidget{
			Key: key,
		}
	}()
}

type AddNewWidget struct {
	Key   string
	Value string
}

type UpdateWidgetContent struct {
	Key   string
	Value string
}

type DeleteWidget struct {
	Key string
}

func (w *widget) Listen() tea.Cmd {
	return func() tea.Msg {
		select {
		case o := <-w.updateChan:
			return o
		default:
			return <-w.updateChan
		}
	}
}

func (w *widget) addNewWidget(key, value string) {
	// skip if key already exists
	if w.GetWidget(key) != nil {
		return
	}

	w.widgets = append(w.widgets, &commonWidget{
		Key:   key,
		Value: value,
	})
}

func (w *widget) updateWidgetContent(key, value string) {
	x := w.GetWidget(key)
	if x != nil {
		x.Value = value
	}
}

func (w *widget) deleteWidget(key string) {
	for i, widget := range w.widgets {
		if widget.Key == key {
			w.widgets = append(w.widgets[:i], w.widgets[i+1:]...)
			break
		}
	}
}

func (w *widget) Init() tea.Cmd {
	return tea.Batch(w.Listen())
}

func (w *widget) Update(msg tea.Msg) (*widget, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !w.termReady {
			if msg.Width > 0 && msg.Height > 0 {
				w.termReady = true
			}
		}
		w.viewport.Width = msg.Width
		w.viewport.Height = msg.Height

	case AddNewWidget:
		w.addNewWidget(msg.Key, msg.Value)

	case UpdateWidgetContent:
		w.updateWidgetContent(msg.Key, msg.Value)

	case DeleteWidget:
		w.deleteWidget(msg.Key)
	}

	cmds = append(cmds, w.Listen())

	return w, tea.Batch(cmds...)
}

func (w *widget) View() string {
	if !w.termReady {
		return "setting up terminal..."
	}

	var widgetLen int
	for _, widget := range w.widgets {
		widgetLen += len(widget.Value)
		widgetLen += w.properties.leftTabPadding + w.properties.rightTabPadding
		widgetLen += 2 // for the border between widgets
	}

	var renderedWidgets []string
	for _, wgt := range w.widgets {
		renderedWidgets = append(renderedWidgets, w.properties.widgetStyle.Render(wgt.Value))
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
