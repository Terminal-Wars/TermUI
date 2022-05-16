package main

import (
	"fmt"
	"log"

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
	if(err != nil) {log.Fatalln(err)}

	// Now we create any elements we want
	// You're advised to call them with go to speed up boot time.
	go win.Button("Clear screen",
		0, // ui event id
		uint16(win.PercentOfWidth(30)), 	// Width
		uint16(win.PercentOfHeight(13)), 	// Height
		win.PercentOfWidth(69),				// X
		win.PercentOfHeight(85),			// Y
	)

	go win.Textbox("",
		0,
		uint16(win.PercentOfWidth(68)),
		uint16(win.PercentOfHeight(13)),
		win.PercentOfWidth(1),
		win.PercentOfHeight(85),
	)

	go win.Label("",
		0,
		uint16(win.PercentOfWidth(98)),
		uint16(win.PercentOfHeight(74)),
		win.PercentOfWidth(1),
		win.PercentOfHeight(1),
		1, // states for labels can be set to have a background and scrollbar.
	)

	// We have two event loops
	// One for checking for UI events, which shouldn't hang.
	go func() {
		for {
			ev := win.WaitForUIEvent()
			// When we get an event, see what type of event it is.
			// For this example, we're interested in the button click, and enter being pressed.
			switch ev.(type) {
				// When enter is pressed...
				case TermUI.UITextboxSubmitEvent:
					// Select the elements we want to work with
					textarea := TermUI.UIElements.GetByID(0, TermUI.LabelType)
					input := TermUI.UIElements.GetByID(0, TermUI.TextboxType)
					text := ev.(TermUI.UITextboxSubmitEvent).Name
					if(len(text) >= 1) {
						textarea.Name += ev.(TermUI.UITextboxSubmitEvent).Name+"\n"
						input.Name = ""
					}
					// Redraw the objects.
					// win.DrawUIElements is tempting but avoid it. Redraw only what you need.
					go win.DrawUIElement(textarea)
					go win.DrawUIElement(input)

				// When the button is clicked...
				case TermUI.UIPressEvent:
					// What event is it?
					// (it is only one, yes, this is just an example of how you'd normally do things)
					switch(ev.(TermUI.UIPressEvent).Event) {
						case TermUI.UIElements.GetByID(0, TermUI.ButtonType):
							textarea := TermUI.UIElements.GetByID(0, TermUI.LabelType)
							textarea.Name = ""
							// In any case, redraw the screen.
							win.DrawUIElement(textarea)
					}
					
					
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
		// free to put your own function in place of or in addition to this to cut
		// anything you aren't listening on but this suffices for most cases
		win.DefaultListeners(ev);
	}
}