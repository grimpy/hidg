//+build ignore

// Copyright 2015 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// mouse displays a text box and tests mouse interaction.  As you click
// and drag, boxes are displayed on screen.  Other events are reported in
// the box.  Press ESC twice to exit the program.
package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"

	"github.com/mattn/go-runewidth"
	"github.com/maugsburger/evdev"
)

/*
 * mapping
 *
 * Maps scancodes to USB event IDs
 *
 * This is based on the linux kernel file drivers/input/hid-input.c
 * by inverting the mapping defined in hid_keyboard[]
 */
var mapping []byte = []byte{3, 41, 30, 31, 32, 33, 34, 35, 36, 37, 38,
	39, 45, 46, 42, 43, 20, 26, 8, 21, 23, 28, 24, 12, 18, 19,
	47, 48, 40, 224, 4, 22, 7, 9, 10, 11, 13, 14, 15, 51, 52,
	53, 225, 50, 29, 27, 6, 25, 5, 17, 16, 54, 55, 56, 229, 85,
	226, 44, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 83,
	71, 95, 96, 97, 86, 92, 93, 94, 87, 89, 90, 91, 98, 99, 0,
	148, 100, 68, 69, 135, 146, 147, 138, 136, 139, 140, 88, 228,
	84, 70, 230, 0, 74, 82, 75, 80, 79, 77, 81, 78, 73, 76, 0,
	239, 238, 237, 102, 103, 0, 72, 0, 133, 144, 145, 137, 227,
	231, 101, 243, 121, 118, 122, 119, 124, 116, 125, 244, 123,
	117, 0, 251, 0, 248, 0, 0, 0, 0, 0, 0, 0, 240, 0,
	249, 0, 0, 0, 0, 0, 241, 242, 0, 236, 0, 235, 232, 234,
	233, 0, 0, 0, 0, 0, 0, 250, 0, 0, 247, 245, 246, 182,
	183, 0, 0, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113,
	114}

/*
 * Map modifier events to the modifier bitcode
 */
var kb_mod map[uint16]byte = map[uint16]byte{
	29:  0x01, // --left-ctrl
	97:  0x10, // --right-ctrl
	42:  0x02, // --left-shift
	54:  0x20, // --right-shift
	56:  0x04, // --left-alt
	100: 0x40, // --right-alt
	125: 0x08, // --left-meta
	126: 0x80, // --right-meta
}

var tcelltoev = map[rune]uint16{
	'0':                       evdev.Key0,
	'1':                       evdev.Key1,
	'2':                       evdev.Key2,
	'3':                       evdev.Key3,
	'4':                       evdev.Key4,
	'5':                       evdev.Key5,
	'6':                       evdev.Key6,
	'7':                       evdev.Key7,
	'8':                       evdev.Key8,
	'9':                       evdev.Key9,
	'a':                       evdev.KeyA,
	'b':                       evdev.KeyB,
	'c':                       evdev.KeyC,
	'd':                       evdev.KeyD,
	'e':                       evdev.KeyE,
	'f':                       evdev.KeyF,
	'g':                       evdev.KeyG,
	'h':                       evdev.KeyH,
	'i':                       evdev.KeyI,
	'j':                       evdev.KeyJ,
	'k':                       evdev.KeyK,
	'l':                       evdev.KeyL,
	'm':                       evdev.KeyM,
	'n':                       evdev.KeyN,
	'o':                       evdev.KeyO,
	'p':                       evdev.KeyP,
	'q':                       evdev.KeyQ,
	'r':                       evdev.KeyR,
	's':                       evdev.KeyS,
	't':                       evdev.KeyT,
	'u':                       evdev.KeyU,
	'v':                       evdev.KeyV,
	'w':                       evdev.KeyW,
	'x':                       evdev.KeyX,
	'y':                       evdev.KeyY,
	'z':                       evdev.KeyZ,
	' ':                       evdev.KeySpace,
	'=':                       evdev.KeyEqual,
	',':                       evdev.KeyComma,
	'.':                       evdev.KeyDot,
	'[':                       evdev.KeyLeftBrace,
	']':                       evdev.KeyRightBrace,
	';':                       evdev.KeySemiColon,
	'\'':                      evdev.KeyApostrophe,
	'/':                       evdev.KeySlash,
	'\\':                      evdev.KeyBackSlash,
	rune(tcell.KeyBackspace2): evdev.KeyBackSpace,
	rune(tcell.KeyBackspace):  evdev.KeyBackSpace,
	rune(tcell.KeyEnter):      evdev.KeyEnter,
	rune(tcell.KeyUp):         evdev.KeyUp,
	rune(tcell.KeyDown):       evdev.KeyDown,
	rune(tcell.KeyLeft):       evdev.KeyLeft,
	rune(tcell.KeyRight):      evdev.KeyRight,
	rune(tcell.KeyTAB):        evdev.KeyTab,
	rune(tcell.KeyEsc):        evdev.KeyEscape,
}
var tcellctrltoev = map[rune]uint16{
	rune(tcell.KeyCtrlA): evdev.KeyA,
	rune(tcell.KeyCtrlB): evdev.KeyB,
	rune(tcell.KeyCtrlC): evdev.KeyC,
	rune(tcell.KeyCtrlD): evdev.KeyD,
	rune(tcell.KeyCtrlE): evdev.KeyE,
	rune(tcell.KeyCtrlF): evdev.KeyF,
	rune(tcell.KeyCtrlG): evdev.KeyG,
	rune(tcell.KeyCtrlH): evdev.KeyH,
	rune(tcell.KeyCtrlI): evdev.KeyI,
	rune(tcell.KeyCtrlJ): evdev.KeyJ,
	rune(tcell.KeyCtrlK): evdev.KeyK,
	rune(tcell.KeyCtrlL): evdev.KeyL,
	rune(tcell.KeyCtrlM): evdev.KeyM,
	rune(tcell.KeyCtrlN): evdev.KeyN,
	rune(tcell.KeyCtrlO): evdev.KeyO,
	rune(tcell.KeyCtrlP): evdev.KeyP,
	rune(tcell.KeyCtrlQ): evdev.KeyQ,
	rune(tcell.KeyCtrlR): evdev.KeyR,
	rune(tcell.KeyCtrlS): evdev.KeyS,
	rune(tcell.KeyCtrlT): evdev.KeyT,
	rune(tcell.KeyCtrlU): evdev.KeyU,
	rune(tcell.KeyCtrlV): evdev.KeyV,
	rune(tcell.KeyCtrlW): evdev.KeyW,
	rune(tcell.KeyCtrlX): evdev.KeyX,
	rune(tcell.KeyCtrlY): evdev.KeyY,
	rune(tcell.KeyCtrlZ): evdev.KeyZ,
}
var tcellshiftokey = map[rune]uint16{
	'!': evdev.Key1,
	'@': evdev.Key2,
	'#': evdev.Key3,
	'$': evdev.Key4,
	'%': evdev.Key5,
	'^': evdev.Key6,
	'&': evdev.Key7,
	'*': evdev.Key8,
	'(': evdev.Key9,
	')': evdev.Key0,
	'_': evdev.KeyMinus,
	'+': evdev.KeyEqual,
	'{': evdev.KeyLeftBrace,
	'}': evdev.KeyRightBrace,
	'"': evdev.KeyApostrophe,
	'/': evdev.KeySlash,
	'<': evdev.KeyComma,
	'>': evdev.KeyDot,
	'|': evdev.KeyBackSlash,
}

/*
 * USBHid structure for maintainig state of a USB Hid instance
 */
type UsbHid struct {
	ev     chan evdev.Event
	exit   chan bool
	file   *os.File
	report [8]byte
	keys   int
}

var defStyle tcell.Style

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, r rune) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}
	if y1 != y2 && x1 != x2 {
		// Only add corners if we need to
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		for col := x1 + 1; col < x2; col++ {
			s.SetContent(col, row, r, nil, style)
		}
	}
}

func drawSelect(s tcell.Screen, x1, y1, x2, y2 int, sel bool) {

	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			mainc, combc, style, width := s.GetContent(col, row)
			if style == tcell.StyleDefault {
				style = defStyle
			}
			style = style.Reverse(sel)
			s.SetContent(col, row, mainc, combc, style)
			col += width - 1
		}
	}
}

func Test() {
	path := "/dev/hidg0"

	testEvents := []evdev.Event{
		evdev.Event{Type: evdev.EvKeys, Code: evdev.KeyH, Value: 1},
		evdev.Event{Type: evdev.EvKeys, Code: evdev.KeyH, Value: 0},
		evdev.Event{Type: evdev.EvKeys, Code: evdev.KeyLeftShift, Value: 1},
		evdev.Event{Type: evdev.EvKeys, Code: evdev.KeyE, Value: 1},
		evdev.Event{Type: evdev.EvKeys, Code: evdev.KeyE, Value: 0},
		evdev.Event{Type: evdev.EvKeys, Code: evdev.KeyL, Value: 1},
		evdev.Event{Type: evdev.EvKeys, Code: evdev.KeyL, Value: 0},
		evdev.Event{Type: evdev.EvKeys, Code: evdev.KeyLeftShift, Value: 0},
		evdev.Event{Type: evdev.EvKeys, Code: evdev.KeyL, Value: 1},
		evdev.Event{Type: evdev.EvKeys, Code: evdev.KeyL, Value: 0},
		evdev.Event{Type: evdev.EvKeys, Code: evdev.KeyO, Value: 1},
		evdev.Event{Type: evdev.EvKeys, Code: evdev.KeyO, Value: 0},
		evdev.Event{Type: evdev.EvKeys, Code: evdev.KeyA, Value: 1},
		evdev.Event{Type: evdev.EvKeys, Code: evdev.KeyB, Value: 1},
		evdev.Event{Type: evdev.EvKeys, Code: evdev.KeyC, Value: 1},
		evdev.Event{Type: evdev.EvKeys, Code: evdev.KeyA, Value: 0},
		evdev.Event{Type: evdev.EvKeys, Code: evdev.KeyB, Value: 0},
		evdev.Event{Type: evdev.EvKeys, Code: evdev.KeyC, Value: 0},
	}

	hid, err := Open(path)
	if err != nil {
		fmt.Println("ERROR: Could not open ", path)
		return
	}
	defer hid.Close()

	for _, ev := range testEvents {
		hid.ForwardEvent(ev)
	}
}

func (hid *UsbHid) ForwardEvent(ev evdev.Event) {
	hid.ev <- ev
}

func Open(path string) (*UsbHid, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	hid := new(UsbHid)
	hid.ev = make(chan evdev.Event)
	hid.exit = make(chan bool)
	hid.file = file

	go eventWriter(hid)

	return hid, nil
}

func (hid *UsbHid) Close() {
	hid.exit <- true
}

func (hid *UsbHid) updateReport(ev evdev.Event) {

	if ev.Type == evdev.EvKeys {

		if kb_mod[ev.Code] != 0 {
			// This code is a modifier
			if ev.Value != 0 {
				hid.report[0] |= kb_mod[ev.Code]
			} else {
				hid.report[0] &^= kb_mod[ev.Code]
			}

		} else {
			// This code is a normal key
			if mapping[ev.Code] == 0 {
				fmt.Printf("Warning: No mapping for event code: %d\n", ev.Code)
				return
			}

			keyPos := -1
			for i, c := range hid.report[2 : 2+hid.keys] {
				if c == mapping[ev.Code] {
					keyPos = i
					break
				}
			}

			if keyPos != -1 {
				if ev.Value == 0 {
					// When removing a key from the middle of the byte
					if hid.keys > keyPos {
						hid.report[keyPos+2] = hid.report[hid.keys+1]
					}
					hid.report[hid.keys+1] = 0
					hid.keys--
				}
			} else {
				if hid.keys < len(hid.report)-2 {
					hid.report[2+hid.keys] = mapping[ev.Code]
					hid.keys++
				} else {
					hid.report[len(hid.report)-1] = mapping[ev.Code]
				}
			}
		}
	}
}

func eventWriter(hid *UsbHid) {

	defer hid.file.Close()

	for {
		select {
		case ev := <-hid.ev:
			//fmt.Printf("Got event: Code = %d, Value = %d, mapping = %d\n", ev.Code, ev.Value, mapping[ev.Code])
			hid.updateReport(ev)
			n, _ := hid.file.Write(hid.report[:])
			if n != len(hid.report) {
				fmt.Println("ERROR: Write failed")
				return
			}

			hid.file.Sync()
		case <-hid.exit:
			return
		}
	}
}

// This program just shows simple mouse and keyboard events.  Press ESC twice to
// exit.
func main() {

	hid, err := Open("/dev/hidg0")
	if err != nil {
		fmt.Println("ERROR: Could not open ", "/dev/hidg0")
		return
	}
	defer hid.Close()

	encoding.Register()

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	defStyle = tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.Clear()

	posfmt := "Mouse: %d, %d  "
	btnfmt := "Buttons: %s"
	keyfmt := "Keys: %s"
	white := tcell.StyleDefault.
		Foreground(tcell.ColorWhite).Background(tcell.ColorRed)

	mx, my := -1, -1
	ox, oy := -1, -1
	bx, by := -1, -1
	w, h := s.Size()
	lchar := '*'
	bstr := ""
	lks := ""
	ecnt := 0

	for {
		drawBox(s, 1, 1, 42, 6, white, ' ')
		emitStr(s, 2, 2, white, "Press ESC twice to exit, C to clear.")
		emitStr(s, 2, 3, white, fmt.Sprintf(posfmt, mx, my))
		emitStr(s, 2, 4, white, fmt.Sprintf(btnfmt, bstr))
		emitStr(s, 2, 5, white, fmt.Sprintf(keyfmt, lks))

		s.Show()
		bstr = ""
		ev := s.PollEvent()
		st := tcell.StyleDefault.Background(tcell.ColorRed)
		up := tcell.StyleDefault.
			Background(tcell.ColorBlue).
			Foreground(tcell.ColorBlack)
		w, h = s.Size()

		// always clear any old selection box
		if ox >= 0 && oy >= 0 && bx >= 0 {
			drawSelect(s, ox, oy, bx, by, false)
		}

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
			s.SetContent(w-1, h-1, 'R', nil, st)
		case *tcell.EventKey:
			pressedkey := ev.Rune()
			shifted := false
			var key uint16
			var ok bool
			if ev.Modifiers()&tcell.ModCtrl != 0 {
				modkey := evdev.KeyLeftCtrl
				hid.ForwardEvent(evdev.Event{Type: evdev.EvKeys, Code: uint16(modkey), Value: 1})
				key := tcellctrltoev[pressedkey]
				hid.ForwardEvent(evdev.Event{Type: evdev.EvKeys, Code: uint16(key), Value: 1})
				hid.ForwardEvent(evdev.Event{Type: evdev.EvKeys, Code: uint16(key), Value: 0})
				hid.ForwardEvent(evdev.Event{Type: evdev.EvKeys, Code: uint16(modkey), Value: 0})
			} else {

				if key, ok = tcellshiftokey[pressedkey]; ok {
					hid.ForwardEvent(evdev.Event{Type: evdev.EvKeys, Code: evdev.KeyLeftShift, Value: 1})
					shifted = true
				} else {
					if pressedkey >= 'A' && pressedkey <= 'Z' {
						pressedkey += ('a' - 'A')
						shifted = true
						hid.ForwardEvent(evdev.Event{Type: evdev.EvKeys, Code: evdev.KeyLeftShift, Value: 1})
					}
					key = tcelltoev[pressedkey]
				}
				hid.ForwardEvent(evdev.Event{Type: evdev.EvKeys, Code: uint16(key), Value: 1})
				hid.ForwardEvent(evdev.Event{Type: evdev.EvKeys, Code: uint16(key), Value: 0})
				if shifted {
					hid.ForwardEvent(evdev.Event{Type: evdev.EvKeys, Code: evdev.KeyLeftShift, Value: 0})
				}

			}

			s.SetContent(w-2, h-2, ev.Rune(), nil, st)
			s.SetContent(w-1, h-1, 'K', nil, st)
			if ev.Key() == tcell.KeyEscape {
				ecnt++
				if ecnt > 1 {
					s.Fini()
					os.Exit(0)
				}
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else {
				ecnt = 0
				if ev.Rune() == 'C' || ev.Rune() == 'c' {
					s.Clear()
				}
			}
			lks = ev.Name()
		case *tcell.EventMouse:
			x, y := ev.Position()
			button := ev.Buttons()
			for i := uint(0); i < 8; i++ {
				if int(button)&(1<<i) != 0 {
					bstr += fmt.Sprintf(" Button%d", i+1)
				}
			}
			if button&tcell.WheelUp != 0 {
				bstr += " WheelUp"
			}
			if button&tcell.WheelDown != 0 {
				bstr += " WheelDown"
			}
			if button&tcell.WheelLeft != 0 {
				bstr += " WheelLeft"
			}
			if button&tcell.WheelRight != 0 {
				bstr += " WheelRight"
			}
			// Only buttons, not wheel events
			button &= tcell.ButtonMask(0xff)
			ch := '*'

			if button != tcell.ButtonNone && ox < 0 {
				ox, oy = x, y
			}
			switch ev.Buttons() {
			case tcell.ButtonNone:
				if ox >= 0 {
					bg := tcell.Color((lchar - '0') * 2)
					drawBox(s, ox, oy, x, y,
						up.Background(bg),
						lchar)
					ox, oy = -1, -1
					bx, by = -1, -1
				}
			case tcell.Button1:
				ch = '1'
			case tcell.Button2:
				ch = '2'
			case tcell.Button3:
				ch = '3'
			case tcell.Button4:
				ch = '4'
			case tcell.Button5:
				ch = '5'
			case tcell.Button6:
				ch = '6'
			case tcell.Button7:
				ch = '7'
			case tcell.Button8:
				ch = '8'
			default:
				ch = '*'

			}
			if button != tcell.ButtonNone {
				bx, by = x, y
			}
			lchar = ch
			s.SetContent(w-1, h-1, 'M', nil, st)
			mx, my = x, y
		default:
			s.SetContent(w-1, h-1, 'X', nil, st)
		}

		if ox >= 0 && bx >= 0 {
			drawSelect(s, ox, oy, bx, by, true)
		}
	}
}
