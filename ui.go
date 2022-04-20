package TermUI

import (
	"fmt"
)

// Generic event
type Event interface {
	String() string
}

// Types of elements
var UIElements struct {
	Buttons		[16]Button
}

// Any UI Event
type UIEvent struct {
	Name 	string
	Width 	int16
	Height  int16
	X 		int16
	Y 		int16
	ID 		int16
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
var ButtonMax int16 = 16

var UIEventNum int16 = 0
const UIEventMax = 64
var UIEvents []UIEvent = make([]UIEvent, UIEventMax)

// Check if we have any UI events left
func (win *Window) NoMoreUIEvents() (bool) {
	return (UIEventNum == UIEventMax-1)
}

// Check if we have any buttons left
func (win *Window) NoMoreButtons() (bool) {
	return (ButtonNum == ButtonMax-1)
}

// Make a new UI event
func (win *Window) NewUIEvent(Name string, Width uint16, Height uint16, X int16, Y int16) {
	UIEvents[UIEventNum] = UIEvent{Name, int16(Width),int16(Height),X,Y,UIEventNum}
	UIEventNum++
}

// Draw all UI elements.
func (win *Window) DrawUIElements() {
	for i := 0; int16(i) < ButtonNum; i++ {
		btn := UIElements.Buttons[i]
		win.Base(btn.Width,btn.Height,btn.X,btn.Y)
	}
}

// Wait for a UI event
func (win *Window) WaitForUIEvent() (Event)  {
	return <-win.uiEventChan
}
