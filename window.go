package TermUI

import (
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type Window struct {
	Conn 		*xgb.Conn
	Screen 		*xproto.ScreenInfo
	Window 		xproto.Window
}

// Function for creating a new window, with default options in place.
func NewWindow(Width uint16, Height uint16, Flags []uint32) (Window Window, Errors [3]error) {
	Window, Errors = NewWindowComplex(Width, Height, 0, xproto.CwBackPixel | xproto.CwEventMask, Flags)
	return
}

func NewWindowComplex(Width uint16, Height uint16, BorderWidth uint16, Mask uint32, Flags []uint32) (win Window, errors [3]error) {

	// establish a new connection
	X, err := xgb.NewConn()
	if err != nil {
		errors[0] = err
	}

	// establish a new connection
	setup := xproto.Setup(X)
	screen := setup.DefaultScreen(X)

	// and a new window ID
	windowID, err := xproto.NewWindowId(X)
	if err != nil {errors[1] = err}

	// create a window using that ID
	xproto.CreateWindow(X, screen.RootDepth, windowID, screen.Root,
		0, 0, Width, Height, BorderWidth,
		xproto.WindowClassInputOutput, screen.RootVisual, 
		Mask,
		Flags)

	// make that window appear on screen.
	err = xproto.MapWindowChecked(X, windowID).Check()
	if err != nil {errors[2] = err}

	// make a new Window object with what we've gotten here
	win = Window{X, screen, windowID}

	return
}