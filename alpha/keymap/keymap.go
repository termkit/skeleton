package keymap

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
	SwitchTabRight = "ctrl+right"
	SwitchTabLeft  = "ctrl+left"
	Quit           = "ctrl+c"
)

var (
	once sync.Once
	k    *KeyMap
)

func NewKeyMap() *KeyMap {
	once.Do(func() {
		k = &KeyMap{
			SwitchTabRight: teakey.NewBinding(
				teakey.WithKeys(SwitchTabRight),
			),
			SwitchTabLeft: teakey.NewBinding(
				teakey.WithKeys(SwitchTabLeft),
			),
			Quit: teakey.NewBinding(
				teakey.WithKeys(Quit),
			),
		}
	})
	return k
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
