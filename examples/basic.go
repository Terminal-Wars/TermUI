package main

import (
	"fmt"

	"github.com/Terminal-Wars/TermUI"
	"github.com/jezek/xgb/xproto"
)

// Example program
func main() {
	// Create a 320x200 window that listens for key presses and releases
	win, err := TermUI.NewWindow(320,200,
		[]uint32{
			0xffffffff,
			xproto.EventMaskKeyPress |
			xproto.EventMaskKeyRelease,
		})

	// Check errors (there are three returned at once, so they're in an array)
	for _, e := range err {
		if(e != nil) {
			fmt.Println(e)
			return
		}
	}

	// Now we create any elements we want
	win.Button("My Cool Button",
		0, // ui event id
		uint16(win.PercentOfWidth(30)), 	// Width
		uint16(win.PercentOfHeight(13)), 	// Height
		win.PercentOfWidth(69),				// X
		win.PercentOfHeight(85),				// Y
	)

	win.Textbox("my cool textbox",
		0,
		uint16(win.PercentOfWidth(68)),
		uint16(win.PercentOfHeight(13)),
		win.PercentOfWidth(1),
		win.PercentOfHeight(85),
	)

	win.Label("AAAAAAAAAAAAAAAAAAAAAAAAASFESDHJSDFGAKLSDJQIWEROYFSDLKGHKLXCVJVKXCJLVBHKLJXCVKLHHLKJXCBLHJKXCVLHJLHKJXCBLHJKXCVLBHXCVHJKJHKLXCBLJKHXCVLJKLHJKXCXCVHBLXCVKJXLCKJVBXLKCVJXLKCVJ",
		0,
		uint16(win.PercentOfWidth(98)),
		uint16(win.PercentOfHeight(74)),
		win.PercentOfWidth(1),
		win.PercentOfHeight(1),
	)

	// We have two event loops
	// One for checking for UI events, which shouldn't hang.
	go func() {
		for {
			ev := win.WaitForUIEvent()
			fmt.Println(ev)
			switch ev.(type) {
				case TermUI.UIReleaseEvent:
					fmt.Println(ev.(TermUI.UIReleaseEvent))
			}
		}
	}()
	// And one for checking for X events, which SHOULD hang.
	for {
		ev, xerr := win.Conn.WaitForEvent()
		if xerr != nil {fmt.Printf("Error: %s\n", xerr)}

		// (in some WMs this happens when you close the program)
		if ev == nil && xerr == nil {return}

		// All the default listeners required to make the UI work; you are
		// free to put your own function in place of this to cut
		// anything you aren't listening on but this suffices for most
		// cases
		win.DefaultListeners(ev);

		// making it a goroutine too seems weird on paper but
		// without it, the cpu usage is usually 95%. using one brings 
		// that down to like 0.2% ¯\_(ツ)_/¯

		// that said, if you do put your own switch or another switch here,
		// it should be in an inline goroutine
	}
}