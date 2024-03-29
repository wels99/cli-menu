package climenu

import (
	"github.com/scrouthtv/termios"
)

const (
	KEY_ignore = iota
	KEY_up
	KEY_down
	KEY_left
	KEY_right
	KEY_escape
	KEY_enter
	KEY_ctrlC
	KEY_backspace
	KEY_pageup
	KEY_pagedown
	KEY_filterstring
)

func (m *Menu) getInput() (int, string) {
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
		return KEY_up, ""
	case k.Type == termios.KeySpecial && k.Mod == 0 && k.Value == termios.SpecialArrowDown:
		return KEY_down, ""
	case k.Type == termios.KeySpecial && k.Mod == 0 && k.Value == termios.SpecialArrowLeft:
		return KEY_left, ""
	case k.Type == termios.KeySpecial && k.Mod == 0 && k.Value == termios.SpecialArrowRight:
		return KEY_right, ""
	case k.Type == termios.KeySpecial && k.Mod == 0 && k.Value == termios.SpecialEnter:
		return KEY_enter, ""
	case k.Type == termios.KeySpecial && k.Mod == 0 && k.Value == termios.SpecialEscape:
		return KEY_escape, ""
	case k.Type == termios.KeySpecial && k.Mod == 0 && k.Value == termios.SpecialPgUp:
		return KEY_pageup, ""
	case k.Type == termios.KeySpecial && k.Mod == 0 && k.Value == termios.SpecialPgDown:
		return KEY_pagedown, ""
	case k.Type == 0 && k.Mod == termios.ModCtrl && k.Value == 'c':
		return KEY_ctrlC, ""
	case k.Type == termios.KeySpecial && k.Mod == termios.ModCtrl && k.Value == termios.SpecialBackspace:
		return KEY_backspace, ""
	case k.Type == 0 && k.Mod == 0 && (k.Value >= '!' && k.Value <= '~'):
		return KEY_filterstring, string(k.Value)
	}
	return KEY_ignore, ""
}
