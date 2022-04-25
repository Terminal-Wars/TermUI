package TermUI

import (
	"fmt"
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
	State		int8  // 0: regular, 1: hover, 2: active, 3: disabled
	Type		int8  // 0: button,  
	id 			uint8 // library supplied id
}

var UIEventNum uint8 = 0
var UIEvents []*UIEvent = make([]*UIEvent, 255)

// Types of elements
var UIElements struct {
	Buttons	[16]*UIEvent
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

// BUTTONS
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

// Draw one UI Element
func (win *Window) DrawUIElement(ev *UIEvent) {
	switch(ev.Type) {
		case 0: win.DrawUIButton(ev.id)
	}
}

// Draw all UI elements.
func (win *Window) DrawUIElements() {
	for i := 0; uint8(i) < ButtonNum; i++ {
		win.DrawUIButton(uint8(i))
	}
}

// Wait for a UI event
func (win *Window) WaitForUIEvent() (Event)  {
	return <-win.uiEventChan
}
