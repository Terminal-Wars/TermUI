package TermUI

import (
	"fmt"

	"github.com/jezek/xgb/xproto"
)

// Generic event
type Event interface {
	String() string
}

// Any event that needs to be drawn and checked on
type UIEvent struct {
	Name 		string
	Width 		uint16
	Height  	uint16
	X 			int16
	Y 			int16
	ID 			int16 // user supplied id
	State		int8  // dependent on the element, see their respective comments
	Type		int8  // 0: button, 1: textbox,
	id 			uint8 // library supplied id
}

var UIEventNum uint8 = 0
var UIEvents []*UIEvent = make([]*UIEvent, 255)

var lastSelectedUIEvent *UIEvent // the last ui event we selected.

// Types of elements
var UIElements struct {
	Buttons		[16]*UIEvent
	Textboxes	[16]*UIEvent
}

// What returns when we hover over something
type UIHoverEvent struct {
	Name 	string
	Event *UIEvent
}
func (v UIHoverEvent) String() string {
	return fmt.Sprintf("UIHoverEvent {Name: %s, Event: %p}", v.Name, v.Event)
}

// What returns when we click something
type UIPressEvent struct {
	Name 	string
	Event *UIEvent
}
func (v UIPressEvent) String() string {
	return fmt.Sprintf("UIPressEvent {Name: %s, Event: %p}", v.Name, v.Event)
}

// What returns when we release our click on something
type UIReleaseEvent struct {
	Name 	string
	Event *UIEvent
}
func (v UIReleaseEvent) String() string {
	return fmt.Sprintf("UIReleaseEvent {Name: %s, Event: %p}", v.Name, v.Event)
}

// Check if we have any UI events left
func (win *Window) NoMoreUIEvents() (bool) {
	return (UIEventNum == 254)
}

// Make a new UI event
func (win *Window) NewUIEvent(Name string, ID int16, Width uint16, Height uint16, X int16, Y int16, eventType int8, id uint8) (*UIEvent) {
	ev := UIEvent{Name, Width,Height,X,Y,ID,0,eventType,id}
	UIEvents[UIEventNum] = &ev
	UIEventNum++
	return &ev
}

// Change the state of a UI element, based on it's type.
func (win *Window) changeState(v *UIEvent, trigger int8) {
	// TRIGGERS (not states): 0 hover, 1 press, 2 release
	if(v.Type != 1) { // any time we change the state of something other then the textbox...
		for i := 0; uint8(i) < TextboxNum; i++ {
			textbox := UIElements.Textboxes[i]
			textbox.State = 0
			win.DrawUITextbox(uint8(i))
		}
	}	
	switch(v.Type) {
		case 0: // button
			if(v.State != 3) {
				switch(trigger) {
					case 1:
						v.State = 2
					case 2:
						v.State = 0
				}
			}
		case 1: // textbox
			switch(trigger) {
				case 1:
					v.State = 1
			}
	}
}

// ========
// BUTTONS
// ========
// STATES: 0 regular, 1 hover, 2 active, 3 disabled

var ButtonNum uint8 = 0

// Draw a UI button
func (win *Window) DrawUIButton(i uint8) {
	btn := UIElements.Buttons[i]
	textLength := int16(len(btn.Name))
	btnX := (btn.X+int16(btn.Width/2))-(textLength*3)
	btnY := (btn.Y+int16(btn.Height/2))+6
	if(btn.State == 2) {
		win.BaseSunken(btn.Width,btn.Height,btn.X,btn.Y)
		btnY++
		btnX++
	} else {
		win.BaseRaised(btn.Width,btn.Height,btn.X,btn.Y)
	}
	win.DrawText(btn.Name,btnX,btnY,12,0x000000,0xafb5b5)
}

// ========
// TEXTBOX
// ========
// STATES: 0 inactive, 1 listening

var TextboxNum uint8 = 0

// Draw a UI textbox
func (win *Window) DrawUITextbox(i uint8) {
	txt := UIElements.Textboxes[i]
	txtY := (txt.Y+int16(txt.Height/2))+6
	win.TextboxBase(txt.Width,txt.Height,txt.X,txt.Y)
	// If it's active, draw a little marker
	if (txt.State == 1) {
		txtX_ := txt.X+int16(len(txt.Name))*7
		txtY_ := txt.Y+4
		win.Square(1, txt.Height-7, txtX_, txtY_, 0x000000)
	}
	// TODO: word wrapping
	win.DrawText(txt.Name,txt.X+6,txtY,12,0x000000,0xffffff)
}

// ========
// GENERAL/MISC 
// ========

// Draw one UI Element
func (win *Window) DrawUIElement(ev *UIEvent) {
	switch(ev.Type) {
		case 0: win.DrawUIButton(ev.id)
		case 1: win.DrawUITextbox(ev.id)
	}
}

// Draw all UI elements.
func (win *Window) DrawUIElements() {
	for i := 0; uint8(i) < ButtonNum; i++ {
		win.DrawUIButton(uint8(i))
	}
	for i := 0; uint8(i) < TextboxNum; i++ {
		win.DrawUITextbox(uint8(i))
	}
}

// Wait for a UI event
func (win *Window) WaitForUIEvent() (Event)  {
	return <-win.uiEventChan
}

// Draw a simple square 
func (win *Window) Square(Width uint16, Height uint16, X int16, Y int16, color uint32) {
	// Drawing context
	draw := xproto.Drawable(win.Window)

	// Create a gcontext for each setting
	gcontext, err := xproto.NewGcontextId(win.Conn)
	if(err != nil) {
		fmt.Println("error creating context for square call",err)
		return
	}

	// gcontext settings
	mask := uint32(xproto.GcForeground)
	values := []uint32{color}
	xproto.CreateGC(win.Conn, gcontext, draw, mask, values)

	// basic rectangle
	rectangle := []xproto.Rectangle{{X: X, Y: Y, Width: Width, Height: Height}}
	xproto.PolyFillRectangle(win.Conn, draw, gcontext, rectangle)
} 
