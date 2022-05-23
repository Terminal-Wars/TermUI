package main

import (
	"fmt"
	"log"

	"github.com/Terminal-Wars/TermUI"
	"github.com/jezek/xgb/xproto"
)

// Example program
func main() {
	win, err := TermUI.NewWindow(320,200,
		[]uint32{
			0xffffffff,
			xproto.EventMaskKeyPress |
			xproto.EventMaskKeyRelease,
		})
	if(err != nil) {log.Fatalln(err)}

	win.Square(0,320,200,0,0,0xff0000)
	win.DrawUISquares()

	go func() {
		for {
			ev := win.WaitForUIEvent()
			switch ev.(type) {
				case TermUI.UIPressEvent:
					win.ClearDropdowns()
					win.Dropdown(0,0,[]string{"Option 1", "Option 2", "Option 3"})
			}
		}
	}()
	for {
		ev, xerr := win.Conn.WaitForEvent()
		if xerr != nil {fmt.Printf("Error: %s\n", xerr)}
		if ev == nil && xerr == nil {return}
		win.DefaultListeners(ev);
	}
}