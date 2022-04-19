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
			xproto.EventMaskKeyRelease |
			xproto.EventMaskStructureNotify | 
			xproto.EventMaskExposure,
		})

	// Check errors (there are three returned at once, so they're in an array)
	for _, e := range err {
		if(e != nil) {
			fmt.Println(e)
			return
		}
	}

	// Event loop
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
		case xproto.ExposeEvent:
			win.Button(64,16,160,100)
		case xproto.KeyPressEvent:
			// See https://pkg.go.dev/github.com/jezek/xgb/xproto#KeyPressEvent
			// for documentation about a key press event.
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