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
		vp = &viewport.Model{Width: 80, Height: 24}
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

func (s *Skeleton) GetTerminalViewportWidth() int {
	return vp.Width
}

func (s *Skeleton) GetTerminalViewportHeight() int {
	return vp.Height
}
