package alpha

import (
	"github.com/charmbracelet/bubbles/viewport"
	"sync"
)

type Spirit struct {
	lockTabs bool
}

var (
	onceSpirit  sync.Once
	modelSpirit *Spirit
)

// newSpirit returns a new Spirit.
func newSpirit() *Spirit {
	onceSpirit.Do(func() {
		modelSpirit = &Spirit{}
	})
	return modelSpirit
}

func (s *Spirit) SetLockTabs(lock bool) {
	s.lockTabs = lock
}

func (s *Spirit) GetLockTabs() bool {
	return s.lockTabs
}

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
