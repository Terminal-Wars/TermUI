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
	win.Button("EXAMPLE_BUTTON",64,16,160,100)

	// We have two event loops
	// One for checking for UI events...
	go func() {
		for {
			ev := win.WaitForUIEvent()
			switch ev.(type) {
				case TermUI.UIReleaseEvent:
					fmt.Println(ev.(TermUI.UIReleaseEvent))
			}
		}
	}()
	// And one for checking for X events, which SHOULD hang.
	for {
		ev, xerr := win.Conn.WaitForEvent()

		// (in some WMs this happens when you close the program)
		if ev == nil && xerr == nil {
			fmt.Println("Both event and error are nil. Exiting...")
			return
		}

		if xerr != nil {
			fmt.Printf("Error: %s\n", xerr)
		}

		switch ev.(type) {
        case xproto.MotionNotifyEvent:
        	go win.CheckMouseHover(ev) // required for checking object hovers
		case xproto.ButtonReleaseEvent:
			go win.CheckMouseRelease(ev) // required for checking object clicks
		case xproto.ExposeEvent:
			go win.DrawUIElements() // required for drawing any ui elements
		case xproto.KeyPressEvent:
			kpe := ev.(xproto.KeyPressEvent)
			fmt.Printf("Key pressed: %d\n", kpe.Detail)
			// todo: exit on "q" not "24th key"
			if kpe.Detail == 24 {
				return // exit on q
			}
		case xproto.DestroyNotifyEvent:
			return
		}
	}
}