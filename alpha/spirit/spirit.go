package spirit

import (
	"github.com/charmbracelet/bubbles/viewport"
	"sync"
)

type ModelSpirit struct {
	lockTabs bool
}

var (
	once sync.Once
	s    *ModelSpirit
)

// NewSpirit returns a new Spirit.
func NewSpirit() *ModelSpirit {
	once.Do(func() {
		s = &ModelSpirit{}
	})
	return s
}

func (s *ModelSpirit) SetLockTabs(lock bool) {
	s.lockTabs = lock
}

func (s *ModelSpirit) GetLockTabs() bool {
	return s.lockTabs
}

// --------------------------------------------

var (
	onceViewport sync.Once
	vp           *viewport.Model
)

func NewTerminalViewport() *viewport.Model {
	onceViewport.Do(func() {
		vp = &viewport.Model{}
	})
	return vp
}
