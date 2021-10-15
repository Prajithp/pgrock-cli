package client

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"github.com/mattn/go-runewidth"
)

type Term struct {
	Screen tcell.Screen
	Status chan string
	Width  int
	Height int
}

func NewScreen() (*Term, error) {
	encoding.Register()

	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err = s.Init(); err != nil {
		return nil, err
	}

	// s.SetStyle(
	// 	tcell.StyleDefault.Foreground(
	// 		tcell.ColorBlack,
	// 	).Background(
	// 		tcell.ColorWhite,
	// 	),
	// )
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	w, h := s.Size()
	s.Clear()

	term := &Term{
		Screen: s,
		Status: make(chan string),
		Width:  w,
		Height: h,
	}

	return term, nil
}

func (t *Term) Run(a *Agent) {
	quit := make(chan struct{})
	go func() {
		for {
			switch ev := t.Screen.PollEvent().(type) {
			case *tcell.EventResize:
				t.Screen.Sync()
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					close(quit)
					return
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-quit:
				t.Screen.Fini()
				a.CloseEvent <- 1
				break
			case res := <-t.Status:
				t.connStatusUpdate(res)
			}
		}
	}()
}

func (t *Term) SetStatus(status string) {
	t.Status <- status
}

func (t *Term) connStatusUpdate(status string) {
	lines := strings.Split(status, "\n")

	style := tcell.StyleDefault

	w := 0
	h := 1

	for _, line := range lines {
		t.emitStr(w, h, style, line)
		h += 1
	}
	t.Screen.Sync()
}

func (t *Term) emitStr(x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		t.Screen.SetContent(x, y, c, comb, style)
		x += w
	}
}
