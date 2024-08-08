package alpha

import (
	"github.com/charmbracelet/bubbles/viewport"
	"sync"
)

// --------------------------------------------

var (
	onceViewport sync.Once
	vp           *viewport.Model
)

func newTerminalViewport() *viewport.Model {
	onceViewport.Do(func() {
		vp = &viewport.Model{Width: 80, Height: 24}
	})
	return vp
}

// --------------------------------------------

func GetTerminalViewport() *viewport.Model {
	return vp
}

func SetTerminalViewportWidth(width int) {
	vp.Width = width
}

func SetTerminalViewportHeight(height int) {
	vp.Height = height
}

func GetTerminalViewportWidth() int {
	return vp.Width
}

func GetTerminalViewportHeight() int {
	return vp.Height
}
