package climenu

import (
	"github.com/scrouthtv/termios"
)

const (
	KEY_up int = 1 + iota
	KEY_down
	KEY_left
	KEY_right
	KEY_escape
	KEY_enter
	KEY_ctrlC
)

func (m *Menu) getInput() int {
	defer func() {
		_ = recover()
	}()

	term, err := termios.Open()
	if err != nil {
		panic(err)
	}
	defer term.Close()

	term.SetRaw(true)
	keys, _ := term.Read()

	k := keys[0]
	switch {
	case k.Type == termios.KeySpecial && k.Mod == 0 && k.Value == termios.SpecialArrowUp:
		return KEY_up
	case k.Type == termios.KeySpecial && k.Mod == 0 && k.Value == termios.SpecialArrowDown:
		return KEY_down
	case k.Type == termios.KeySpecial && k.Mod == 0 && k.Value == termios.SpecialArrowLeft:
		return KEY_left
	case k.Type == termios.KeySpecial && k.Mod == 0 && k.Value == termios.SpecialArrowRight:
		return KEY_right
	case k.Type == termios.KeySpecial && k.Mod == 0 && k.Value == termios.SpecialEnter:
		return KEY_enter
	case k.Type == termios.KeySpecial && k.Mod == 0 && k.Value == termios.SpecialEscape:
		return KEY_escape
	case k.Type == 0 && k.Mod == termios.ModCtrl && k.Value == 'c':
		return KEY_ctrlC
	}
	return 0
}
