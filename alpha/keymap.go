package alpha

import (
	teakey "github.com/charmbracelet/bubbles/key"
	"sync"
)

type KeyMap struct {
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
	keyMap     *KeyMap
)

func NewKeyMap() *KeyMap {
	onceKeyMap.Do(func() {
		keyMap = &KeyMap{
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
	return keyMap
}

// --------------------------------------------

func (k *KeyMap) SetKeyNextTab(keybinding teakey.Binding) {
	k.SwitchTabRight = keybinding
}

func (k *KeyMap) SetKeyPrevTab(keybinding teakey.Binding) {
	k.SwitchTabLeft = keybinding
}

func (k *KeyMap) SetKeyQuit(keybinding teakey.Binding) {
	k.Quit = keybinding
}

func (k *KeyMap) GetKeyNextTab() teakey.Binding {
	return k.SwitchTabRight
}

func (k *KeyMap) GetKeyPrevTab() teakey.Binding {
	return k.SwitchTabLeft
}

func (k *KeyMap) GetKeyQuit() teakey.Binding {
	return k.Quit
}
