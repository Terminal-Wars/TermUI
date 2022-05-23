package TermUI

import (
	"fmt"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type Window struct {
	Conn 			*xgb.Conn
	Screen 			*xproto.ScreenInfo
	Window 			xproto.Window
	Width 			uint16
	Height 			uint16

	uiEventChan		chan Event
}

func defaultFlags(flags []uint32) ([]uint32) {
	var flag0, flag1 uint32
	if(len(flags) >= 1) {flag0 = flags[0]}
	if(len(flags) >= 2) {flag0 = flags[1]}
	return []uint32{flag0, flag1 | 
	xproto.EventMaskStructureNotify | 
	xproto.EventMaskExposure | 
	xproto.EventMaskPointerMotion |
	xproto.EventMaskButtonPress |
	xproto.EventMaskButtonRelease,
	}
}

// Function for creating a new window, with default options in place.
func NewWindow(Width, Height uint16, Flags []uint32) (Window Window, Error error) {
	flags := defaultFlags(Flags)
	Window, Error = NewWindowComplex(Width, Height, 0, xproto.CwBackPixel | xproto.CwEventMask, flags)
	return
}

// The same function as above, but it calls NewRawWindowComplex
// [[ CURRENTLY BROKEN, IT CAUSES A CRASH UNLESS YOU REMOVE THE FLAG THAT CAUSES WINDOWS TO BE UNDECORATED ]]
func NewRawWindow(Width, Height uint16, Flags []uint32) (Window Window, Error error) {
	flags := defaultFlags(Flags)
	Window, Error = NewWindowComplex(Width, Height, 0, xproto.CwBackPixel | xproto.CwEventMask | xproto.CwOverrideRedirect, flags)
	return
}

// Function for creating a new window without any default options in place
func NewWindowComplex(Width, Height, BorderWidth uint16, Mask uint32, Flags []uint32) (win Window, err error) {

	// establish a new connection
	X, err := xgb.NewConn()
	if err != nil {
		return
	}

	// establish a new connection
	setup := xproto.Setup(X)
	screen := setup.DefaultScreen(X)

	// and a new window ID
	windowID, err := xproto.NewWindowId(X)
	if err != nil {return}

	// create a window using that ID
	err = xproto.CreateWindowChecked(X, screen.RootDepth, windowID, screen.Root,
		0, 0, Width, Height, BorderWidth,
		xproto.WindowClassInputOutput, screen.RootVisual, 
		Mask,
		Flags).Check()
	if(err != nil) {return}

	// make that window appear on screen.
	err = xproto.MapWindowChecked(X, windowID).Check()
	if err != nil {return}

	// make a new Window object with what we've gotten here
	win = Window{X, screen, windowID,Width,Height, 		make(chan Event)}

	return
}

// functions for getting percentages of the window

func (win *Window) PercentOfWidth(per uint16) (int16) {
	return int16(float32(win.Width)*(float32(per)/100))
}

func (win *Window) PercentOfHeight(per uint16) (int16) {
	fmt.Println()
	return int16(float32(win.Height)*(float32(per)/100))
}