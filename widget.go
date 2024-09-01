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

	// widgetLength is hold the length of the widget
	widgetLength int

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

// GetBorderColor returns the border color of the Widget.
func (w *widget) GetBorderColor() string {
	return w.properties.borderColor
}

// SetWidgetBorderColor sets the border color of the Widget.
func (w *widget) SetWidgetBorderColor(color string) *widget {
	w.properties.widgetStyle = w.properties.widgetStyle.BorderForeground(lipgloss.Color(color))
	return w
}

// GetWidgetBorderColor returns the border color of the Widget.
func (w *widget) GetWidgetBorderColor() string {
	return w.properties.widgetStyle.BorderForeground().String()
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

// DeleteAllWidgets deletes all the widgets.
func (w *widget) DeleteAllWidgets() {
	w.widgets = nil
	w.calculateWidgetLength()
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
		return <-w.updateChan
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

	w.calculateWidgetLength()
}

func (w *widget) updateWidgetContent(key, value string) {
	x := w.GetWidget(key)
	if x != nil {
		x.Value = value
	}

	w.calculateWidgetLength()
}

func (w *widget) deleteWidget(key string) {
	for i, widget := range w.widgets {
		if widget.Key == key {
			w.widgets = append(w.widgets[:i], w.widgets[i+1:]...)
			break
		}
	}

	w.calculateWidgetLength()
}

type WidgetSizeMsg struct {
	NotEnoughToHandleWidgets bool
}

func (w *widget) SendIsTerminalSizeEnough(isEnough bool) {
	go func() {
		w.updateChan <- WidgetSizeMsg{
			NotEnoughToHandleWidgets: isEnough,
		}
	}()
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

		w.calculateWidgetLength()
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

// calculateWidgetLength calculates the length of the widgets.
func (w *widget) calculateWidgetLength() {
	var widgetLen int
	for _, widget := range w.widgets {
		widgetLen += len(widget.Value)
		widgetLen += w.properties.leftTabPadding + w.properties.rightTabPadding
		widgetLen += 2 // for the border between widgets
	}

	requiredLineCount := w.viewport.Width - (widgetLen + 2)
	if requiredLineCount < 0 {
		w.SendIsTerminalSizeEnough(false)
	} else {
		w.SendIsTerminalSizeEnough(true)
	}

	w.widgetLength = widgetLen
}

func (w *widget) View() string {
	if !w.termReady {
		return "setting up terminal..."
	}

	requiredLineCount := w.viewport.Width - (w.widgetLength + 2)

	if requiredLineCount < 0 {
		return ""
	}

	line := strings.Repeat("─", requiredLineCount)
	line = lipgloss.NewStyle().Foreground(lipgloss.Color(w.properties.borderColor)).Render(line)

	var renderedWidgets = make([]string, len(w.widgets))
	for i, wgt := range w.widgets {
		renderedWidgets[i] = w.properties.widgetStyle.Render(wgt.Value)
	}

	leftCorner := lipgloss.JoinVertical(lipgloss.Top, "│", "╰")
	rightCorner := lipgloss.JoinVertical(lipgloss.Top, "│", "╯")
	leftCorner = lipgloss.NewStyle().Foreground(lipgloss.Color(w.properties.borderColor)).Render(leftCorner)
	rightCorner = lipgloss.NewStyle().Foreground(lipgloss.Color(w.properties.borderColor)).Render(rightCorner)

	var bottom []string
	bottom = append(bottom, line)
	bottom = append(bottom, renderedWidgets...)

	position := lipgloss.Center
	if len(w.widgets) > 0 {
		position = lipgloss.Top
	}

	return lipgloss.JoinHorizontal(position, leftCorner, lipgloss.JoinHorizontal(lipgloss.Center, bottom...), rightCorner)
}
