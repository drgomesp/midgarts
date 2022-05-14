package window

import (
	"github.com/veandco/go-sdl2/sdl"
)

type KeyState struct {
	window *sdl.Window
	states map[sdl.Keycode]bool
}

func NewKeyState(window *sdl.Window) *KeyState {
	ks := new(KeyState)
	ks.window = window
	ks.states = map[sdl.Keycode]bool{
		sdl.K_RETURN:             false,
		sdl.K_ESCAPE:             false,
		sdl.K_BACKSPACE:          false,
		sdl.K_TAB:                false,
		sdl.K_SPACE:              false,
		sdl.K_EXCLAIM:            false,
		sdl.K_QUOTEDBL:           false,
		sdl.K_HASH:               false,
		sdl.K_PERCENT:            false,
		sdl.K_DOLLAR:             false,
		sdl.K_AMPERSAND:          false,
		sdl.K_QUOTE:              false,
		sdl.K_LEFTPAREN:          false,
		sdl.K_RIGHTPAREN:         false,
		sdl.K_ASTERISK:           false,
		sdl.K_PLUS:               false,
		sdl.K_COMMA:              false,
		sdl.K_MINUS:              false,
		sdl.K_PERIOD:             false,
		sdl.K_SLASH:              false,
		sdl.K_0:                  false,
		sdl.K_1:                  false,
		sdl.K_2:                  false,
		sdl.K_3:                  false,
		sdl.K_4:                  false,
		sdl.K_5:                  false,
		sdl.K_6:                  false,
		sdl.K_7:                  false,
		sdl.K_8:                  false,
		sdl.K_9:                  false,
		sdl.K_COLON:              false,
		sdl.K_SEMICOLON:          false,
		sdl.K_LESS:               false,
		sdl.K_EQUALS:             false,
		sdl.K_GREATER:            false,
		sdl.K_QUESTION:           false,
		sdl.K_AT:                 false,
		sdl.K_LEFTBRACKET:        false,
		sdl.K_BACKSLASH:          false,
		sdl.K_RIGHTBRACKET:       false,
		sdl.K_CARET:              false,
		sdl.K_UNDERSCORE:         false,
		sdl.K_BACKQUOTE:          false,
		sdl.K_a:                  false,
		sdl.K_b:                  false,
		sdl.K_c:                  false,
		sdl.K_d:                  false,
		sdl.K_e:                  false,
		sdl.K_f:                  false,
		sdl.K_g:                  false,
		sdl.K_h:                  false,
		sdl.K_i:                  false,
		sdl.K_j:                  false,
		sdl.K_k:                  false,
		sdl.K_l:                  false,
		sdl.K_m:                  false,
		sdl.K_n:                  false,
		sdl.K_o:                  false,
		sdl.K_p:                  false,
		sdl.K_q:                  false,
		sdl.K_r:                  false,
		sdl.K_s:                  false,
		sdl.K_t:                  false,
		sdl.K_u:                  false,
		sdl.K_v:                  false,
		sdl.K_w:                  false,
		sdl.K_x:                  false,
		sdl.K_y:                  false,
		sdl.K_z:                  false,
		sdl.K_CAPSLOCK:           false,
		sdl.K_F1:                 false,
		sdl.K_F2:                 false,
		sdl.K_F3:                 false,
		sdl.K_F4:                 false,
		sdl.K_F5:                 false,
		sdl.K_F6:                 false,
		sdl.K_F7:                 false,
		sdl.K_F8:                 false,
		sdl.K_F9:                 false,
		sdl.K_F10:                false,
		sdl.K_F11:                false,
		sdl.K_F12:                false,
		sdl.K_PRINTSCREEN:        false,
		sdl.K_SCROLLLOCK:         false,
		sdl.K_PAUSE:              false,
		sdl.K_INSERT:             false,
		sdl.K_HOME:               false,
		sdl.K_PAGEUP:             false,
		sdl.K_DELETE:             false,
		sdl.K_END:                false,
		sdl.K_PAGEDOWN:           false,
		sdl.K_RIGHT:              false,
		sdl.K_LEFT:               false,
		sdl.K_DOWN:               false,
		sdl.K_UP:                 false,
		sdl.K_NUMLOCKCLEAR:       false,
		sdl.K_KP_DIVIDE:          false,
		sdl.K_KP_MULTIPLY:        false,
		sdl.K_KP_MINUS:           false,
		sdl.K_KP_PLUS:            false,
		sdl.K_KP_ENTER:           false,
		sdl.K_KP_1:               false,
		sdl.K_KP_2:               false,
		sdl.K_KP_3:               false,
		sdl.K_KP_4:               false,
		sdl.K_KP_5:               false,
		sdl.K_KP_6:               false,
		sdl.K_KP_7:               false,
		sdl.K_KP_8:               false,
		sdl.K_KP_9:               false,
		sdl.K_KP_0:               false,
		sdl.K_KP_PERIOD:          false,
		sdl.K_APPLICATION:        false,
		sdl.K_POWER:              false,
		sdl.K_KP_EQUALS:          false,
		sdl.K_F13:                false,
		sdl.K_F14:                false,
		sdl.K_F15:                false,
		sdl.K_F16:                false,
		sdl.K_F17:                false,
		sdl.K_F18:                false,
		sdl.K_F19:                false,
		sdl.K_F20:                false,
		sdl.K_F21:                false,
		sdl.K_F22:                false,
		sdl.K_F23:                false,
		sdl.K_F24:                false,
		sdl.K_EXECUTE:            false,
		sdl.K_HELP:               false,
		sdl.K_MENU:               false,
		sdl.K_SELECT:             false,
		sdl.K_STOP:               false,
		sdl.K_AGAIN:              false,
		sdl.K_UNDO:               false,
		sdl.K_CUT:                false,
		sdl.K_COPY:               false,
		sdl.K_PASTE:              false,
		sdl.K_FIND:               false,
		sdl.K_MUTE:               false,
		sdl.K_VOLUMEUP:           false,
		sdl.K_VOLUMEDOWN:         false,
		sdl.K_KP_COMMA:           false,
		sdl.K_KP_EQUALSAS400:     false,
		sdl.K_ALTERASE:           false,
		sdl.K_SYSREQ:             false,
		sdl.K_CANCEL:             false,
		sdl.K_CLEAR:              false,
		sdl.K_PRIOR:              false,
		sdl.K_RETURN2:            false,
		sdl.K_SEPARATOR:          false,
		sdl.K_OUT:                false,
		sdl.K_OPER:               false,
		sdl.K_CLEARAGAIN:         false,
		sdl.K_CRSEL:              false,
		sdl.K_EXSEL:              false,
		sdl.K_KP_00:              false,
		sdl.K_KP_000:             false,
		sdl.K_THOUSANDSSEPARATOR: false,
		sdl.K_DECIMALSEPARATOR:   false,
		sdl.K_CURRENCYUNIT:       false,
		sdl.K_CURRENCYSUBUNIT:    false,
		sdl.K_KP_LEFTPAREN:       false,
		sdl.K_KP_RIGHTPAREN:      false,
		sdl.K_KP_LEFTBRACE:       false,
		sdl.K_KP_RIGHTBRACE:      false,
		sdl.K_KP_TAB:             false,
		sdl.K_KP_BACKSPACE:       false,
		sdl.K_KP_A:               false,
		sdl.K_KP_B:               false,
		sdl.K_KP_C:               false,
		sdl.K_KP_D:               false,
		sdl.K_KP_E:               false,
		sdl.K_KP_F:               false,
		sdl.K_KP_XOR:             false,
		sdl.K_KP_POWER:           false,
		sdl.K_KP_PERCENT:         false,
		sdl.K_KP_LESS:            false,
		sdl.K_KP_GREATER:         false,
		sdl.K_KP_AMPERSAND:       false,
		sdl.K_KP_DBLAMPERSAND:    false,
		sdl.K_KP_VERTICALBAR:     false,
		sdl.K_KP_DBLVERTICALBAR:  false,
		sdl.K_KP_COLON:           false,
		sdl.K_KP_HASH:            false,
		sdl.K_KP_SPACE:           false,
		sdl.K_KP_AT:              false,
		sdl.K_KP_EXCLAM:          false,
		sdl.K_KP_MEMSTORE:        false,
		sdl.K_KP_MEMRECALL:       false,
		sdl.K_KP_MEMCLEAR:        false,
		sdl.K_KP_MEMADD:          false,
		sdl.K_KP_MEMSUBTRACT:     false,
		sdl.K_KP_MEMMULTIPLY:     false,
		sdl.K_KP_MEMDIVIDE:       false,
		sdl.K_KP_PLUSMINUS:       false,
		sdl.K_KP_CLEAR:           false,
		sdl.K_KP_CLEARENTRY:      false,
		sdl.K_KP_BINARY:          false,
		sdl.K_KP_OCTAL:           false,
		sdl.K_KP_DECIMAL:         false,
		sdl.K_KP_HEXADECIMAL:     false,
		sdl.K_LCTRL:              false,
		sdl.K_LSHIFT:             false,
		sdl.K_LALT:               false,
		sdl.K_LGUI:               false,
		sdl.K_RCTRL:              false,
		sdl.K_RSHIFT:             false,
		sdl.K_RALT:               false,
		sdl.K_RGUI:               false,
		sdl.K_MODE:               false,
		sdl.K_AUDIONEXT:          false,
		sdl.K_AUDIOPREV:          false,
		sdl.K_AUDIOSTOP:          false,
		sdl.K_AUDIOPLAY:          false,
		sdl.K_AUDIOMUTE:          false,
		sdl.K_MEDIASELECT:        false,
		sdl.K_WWW:                false,
		sdl.K_MAIL:               false,
		sdl.K_CALCULATOR:         false,
		sdl.K_COMPUTER:           false,
		sdl.K_AC_SEARCH:          false,
		sdl.K_AC_HOME:            false,
		sdl.K_AC_BACK:            false,
		sdl.K_AC_FORWARD:         false,
		sdl.K_AC_STOP:            false,
		sdl.K_AC_REFRESH:         false,
		sdl.K_AC_BOOKMARKS:       false,
		sdl.K_BRIGHTNESSDOWN:     false,
		sdl.K_BRIGHTNESSUP:       false,
		sdl.K_DISPLAYSWITCH:      false,
		sdl.K_KBDILLUMTOGGLE:     false,
		sdl.K_KBDILLUMDOWN:       false,
		sdl.K_KBDILLUMUP:         false,
		sdl.K_EJECT:              false,
		sdl.K_SLEEP:              false,
	}

	// it's not work for OSX. Runtime Error: Terminating app due to uncaught exception
	// 'NSInternalInconsistencyException',
	// reason: 'Modifications to the layout engine must not be performed from a background thread after it has been
	// accessed from the main thread.'
	//go ks.pollEventKeys()

	return ks
}

func (ks *KeyState) Pressed(keyCode sdl.Keycode) bool {
	return ks.states[keyCode]
}

func(ks *KeyState) Update(event sdl.Event) {
	if event == nil {
		return
	}

	switch event := event.(type) {
	case *sdl.KeyboardEvent:
		if event.State == sdl.PRESSED {
			ks.states[event.Keysym.Sym] = true
		} else if event.State == sdl.RELEASED {
			ks.states[event.Keysym.Sym] = false
		}
	}
}

func (ks *KeyState) pollEventKeys() {
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event := event.(type) {
			case *sdl.KeyboardEvent:
				if event.State == sdl.PRESSED {
					ks.states[event.Keysym.Sym] = true
				} else if event.State == sdl.RELEASED {
					ks.states[event.Keysym.Sym] = false
				}
			}
		}
	}
}
