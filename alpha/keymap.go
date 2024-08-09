package alpha

import (
	teakey "github.com/charmbracelet/bubbles/key"
	"sync"
)

type keyMap struct {
	SwitchTabRight teakey.Binding
	SwitchTabLeft  teakey.Binding
	Quit           teakey.Binding
}

const (
	keymapSwitchTabRight = "ctrl+right"
	keymapSwitchTabLeft  = "ctrl+left"
	keymapQuit           = "ctrl+c"
)

var (
	onceKeyMap sync.Once
	varKeyMap  *keyMap
)

func newKeyMap() *keyMap {
	onceKeyMap.Do(func() {
		varKeyMap = &keyMap{
			SwitchTabRight: teakey.NewBinding(
				teakey.WithKeys(keymapSwitchTabRight),
			),
			SwitchTabLeft: teakey.NewBinding(
				teakey.WithKeys(keymapSwitchTabLeft),
			),
			Quit: teakey.NewBinding(
				teakey.WithKeys(keymapQuit),
			),
		}
	})
	return varKeyMap
}

// --------------------------------------------

func (k *keyMap) SetKeyNextTab(keybinding teakey.Binding) {
	k.SwitchTabRight = keybinding
}

func (k *keyMap) SetKeyPrevTab(keybinding teakey.Binding) {
	k.SwitchTabLeft = keybinding
}

func (k *keyMap) SetKeyQuit(keybinding teakey.Binding) {
	k.Quit = keybinding
}

func (k *keyMap) GetKeyNextTab() teakey.Binding {
	return k.SwitchTabRight
}

func (k *keyMap) GetKeyPrevTab() teakey.Binding {
	return k.SwitchTabLeft
}

func (k *keyMap) GetKeyQuit() teakey.Binding {
	return k.Quit
}
