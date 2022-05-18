package TermUI

import (
	"fmt"
	"strings"

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
	Type		int8  // 0: button, 1: textbox, 2: label
	id 			uint8 // library supplied id
}

const ( // names for the types above
	ButtonType	int8 = 0
	TextboxType int8 = 1
	LabelType 	int8 = 2
)

var UIEventNum uint8 = 0
var UIEvents []*UIEvent = make([]*UIEvent, 255)

var lastSelectedUIEvent *UIEvent // the last ui event we selected.

// initialization function
func init() {

}

// Types of elements
type UIElementStruct struct {
	Buttons		[]*UIEvent
	Textboxes	[]*UIEvent
	Labels  	[]*UIEvent
}
var UIElements UIElementStruct = UIElementStruct{
	make([]*UIEvent, 16),
	make([]*UIEvent, 16),
	make([]*UIEvent, 64),
}

func (v UIElementStruct) GetByID(ID uint8, FindType int8) *UIEvent {
	// Which type should we be searching?
	var searchThru []*UIEvent
	switch(FindType) {
		case ButtonType: searchThru = UIElements.Buttons
		case TextboxType: searchThru = UIElements.Textboxes
		case LabelType: searchThru = UIElements.Labels
	}
	for _, v := range searchThru {
		if(v.id == ID) {
			return v
		}
	}
	return nil
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

// What returns when we hit enter on a textbox
type UITextboxSubmitEvent struct {
	Name 	string
	Event *UIEvent
}
func (v UITextboxSubmitEvent) String() string {
	return fmt.Sprintf("UITextboxSubmitEvent {Name: %s, Event: %p}", v.Name, v.Event)
}

// Check if we have any UI events left
func (win *Window) NoMoreUIEvents() (bool) {
	return (UIEventNum == 254)
}

// Make a new UI event
func (win *Window) NewUIEvent(Name string, ID int16, Width uint16, Height uint16, X int16, Y int16, eventType int8, id uint8) (*UIEvent) {
	return win.NewUIEventWithPresetState(Name,ID,Width,Height,X,Y,eventType,id,0)
}

func (win *Window) NewUIEventWithPresetState(Name string, ID int16, Width uint16, Height uint16, X int16, Y int16, eventType int8, id uint8, defaultState int8) (*UIEvent) {
	ev := UIEvent{Name, Width,Height,X,Y,ID,defaultState,eventType,id}
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

func (win *Window) DrawUIButtons() {
	for i := 0; uint8(i) < ButtonNum; i++ {
		go win.DrawUIButton(uint8(i))
	}
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
	// If it's active (but we're not out of bounds), draw a little marker
	if (txt.State == 1 && len(txt.Name) < int(txt.Width/6-1)) {
		txtX_ := txt.X+int16((len(txt.Name)*6)+6)
		txtY_ := txt.Y+4
		win.Square(1, txt.Height-7, txtX_, txtY_, []uint32{0x000000})
	}
	name := txt.Name
	if(len(txt.Name) > int(txt.Width/6-1)) {
		name = txt.Name[len(txt.Name)-int(txt.Width/6-1):len(txt.Name)]
	}
	win.DrawText(name,txt.X+6,txtY,12,0x000000,0xffffff)
}

func (win *Window) DrawUITextboxes() {
	for i := 0; uint8(i) < TextboxNum; i++ {
		go win.DrawUITextbox(uint8(i))
	}
}

// ========
// LABEL
// ========
// STATES: 0 no background, 1 gfx. anything above is the offset used for scrolling vertically

var LabelNum uint8 = 0

// Draw a UI label
func (win *Window) DrawUILabel(i uint8) {
	txt := UIElements.Labels[i]
	txtY := (txt.Y+13)
	maxWidth := int(txt.Width/6-2)
	maxHeight := int(txt.Height/12)
	var rows []string
	// If the text goes beyond the textbox's width
	if(len(txt.Name) > maxWidth) {
		rows = WordWrap(txt.Name, maxWidth)
	} else {
		rows = append(rows, txt.Name) 
	}
	// We want to split the rows again based on newlines.
	var rows_ []string
	for _, v := range rows {
		splitString := strings.Split(v,"\n")
		for _, w := range splitString {
			rows_ = append(rows_, w)
		}
	}
	// Anything beyond should only be considered if we have the gfx state on.
	if(txt.State >= 1) {
		// Draw the base either way
		win.TextboxBase(txt.Width,txt.Height,txt.X,txt.Y)
		// Is the resulting array out of bounds?
		// If so...
		if(len(rows_) > maxHeight) {
			// another clone yay
			var rows__ []string
			// push the text to the side a bit.
			rows__ = WordWrap(txt.Name, maxWidth-2)
			// ...shrink the array down.
			rows__ = rows_[len(rows_)-maxHeight:len(rows_)]
			// put in a checkerboard pattern in place of where a scrollbar would go
			go win.Square(8, uint16(txt.Height)-8, int16(txt.Width-8), txt.Y+4, []uint32{0xdedede})
			rows_ = rows__
		}
	}
	// Then go through the array and draw it.
	for i, v := range rows_ {
		win.DrawText(v,txt.X+6,txtY+(int16(i)*13),12,0x000000,0xffffff)
	}
}

func (win *Window) DrawUILabels() {
	for i := 0; uint8(i) < LabelNum; i++ {
		go win.DrawUILabel(uint8(i))
	}
}

// ========
// GENERAL/MISC 
// ========

// Draw one UI Element
func (win *Window) DrawUIElement(ev *UIEvent) {
	switch(ev.Type) {
		case ButtonType: 	go win.DrawUIButton(ev.id)
		case TextboxType:	go win.DrawUITextbox(ev.id)
		case LabelType: 	go win.DrawUILabel(ev.id)
	}
}

// Draw all UI elements.
func (win *Window) DrawUIElements() {
	go win.DrawUIButtons()
	go win.DrawUITextboxes()
	go win.DrawUILabels()
}

// Wait for a UI event
func (win *Window) WaitForUIEvent() (Event)  {
	return <-win.uiEventChan
}

// Draw a simple square 
func (win *Window) Square(Width uint16, Height uint16, X int16, Y int16, Colors []uint32) {
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
	values := Colors
	xproto.CreateGC(win.Conn, gcontext, draw, mask, values)

	// basic rectangle
	rectangle := []xproto.Rectangle{{X: X, Y: Y, Width: Width, Height: Height}}
	xproto.PolyFillRectangle(win.Conn, draw, gcontext, rectangle)

	xproto.FreeGC(win.Conn,gcontext)
} 
