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
	win.Button("My Cool Button",0,94,23,160,100)
	win.Button("B",1,32,32,0,0)

	// We have two event loops
	// One for checking for UI events...

	/*
	go func() {
		for {
			ev := win.WaitForUIEvent()
			switch ev.(type) {
				case TermUI.UIReleaseEvent:
					fmt.Println(ev.(TermUI.UIReleaseEvent))
			}
		}
	}()
	*/
	// And one for checking for X events, which SHOULD hang.
	for {
		ev, xerr := win.Conn.WaitForEvent()
		if xerr != nil {fmt.Printf("Error: %s\n", xerr)}

		// (in some WMs this happens when you close the program)
		if ev == nil && xerr == nil {return}

		// making the switch a goroutine too seems weird on paper but
		// without it, the cpu usage is usually 95%. using one brings 
		// that down by like 94.5%, it's fucking insane considering this
		// is being deployed to wasm which is single threaded but ¯\_(ツ)_/¯
		// (actually i think v86 uses web-workers to do threads so this maybe
		// makes sense)
		go func() {
			switch ev.(type) {
			case xproto.ExposeEvent:
				go win.DrawUIElements()
			/*
	        // required for checking object hovers;
	        // you almost never need this but it's here if you want
	        case xproto.MotionNotifyEvent:
	        	go win.CheckMouseHover(ev) 
	        */
			case xproto.ButtonPressEvent:
				go win.CheckMousePress(ev) // required for checking object clicks
			case xproto.ButtonReleaseEvent:
				go win.CheckMouseRelease(ev) // required for checking object releases
			case xproto.KeyPressEvent:
				kpe := ev.(xproto.KeyPressEvent)
				// todo: exit on "q" not "24th key"
				if kpe.Detail == 24 {
					return // exit on q
				}
			case xproto.DestroyNotifyEvent:
				return
			}
		}();
	}
}