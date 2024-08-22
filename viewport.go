package skeleton

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
		vp = &viewport.Model{Width: 80, Height: 24} // Question: Is it best to use 80x24 as default?
	})
	return vp
}

// --------------------------------------------

func (s *Skeleton) GetTerminalViewport() *viewport.Model {
	return vp
}

func (s *Skeleton) SetTerminalViewportWidth(width int) {
	vp.Width = width
}

func (s *Skeleton) SetTerminalViewportHeight(height int) {
	vp.Height = height
}

func (s *Skeleton) GetTerminalWidth() int {
	return vp.Width
}

func (s *Skeleton) GetTerminalHeight() int {
	return vp.Height
}
