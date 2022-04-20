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
	ID 			int16
	State		int16 // 0: regular, 1: hover, 2: active, 3: disabled
}

var UIEventNum int16 = 0
const UIEventMax = 256
var UIEvents []*UIEvent = make([]*UIEvent, UIEventMax)

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

var ButtonNum int16 = 0

// Check if we have any UI events left
func (win *Window) NoMoreUIEvents() (bool) {
	return (UIEventNum == UIEventMax-1)
}

// Make a new UI event
func (win *Window) NewUIEvent(Name string, ID int16, Width uint16, Height uint16, X int16, Y int16) (*UIEvent) {
	ev := UIEvent{Name, Width,Height,X,Y,ID,0}
	UIEvents[UIEventNum] = &ev
	UIEventNum++
	return &ev
}

// Draw all UI elements.
func (win *Window) DrawUIElements() {
	for i := 0; int16(i) < ButtonNum; i++ {
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
}

// Wait for a UI event
func (win *Window) WaitForUIEvent() (Event)  {
	return <-win.uiEventChan
}
