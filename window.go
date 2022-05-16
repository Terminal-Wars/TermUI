package TermUI

import (
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

type Window struct {
	Conn 			*xgb.Conn
	Screen 			*xproto.ScreenInfo
	Window 			xproto.Window

	uiEventChan		chan Event
}

// Function for creating a new window, with default options in place.
func NewWindow(Width uint16, Height uint16, Flags []uint32) (Window Window, Errors [3]error) {
	flags := []uint32{Flags[0], Flags[1] | 
	xproto.EventMaskStructureNotify | 
	xproto.EventMaskExposure | 
	xproto.EventMaskPointerMotion |
	xproto.EventMaskButtonPress |
	xproto.EventMaskButtonRelease,
	}
	Window, Errors = NewWindowComplex(Width, Height, 0, xproto.CwBackPixel | xproto.CwEventMask, flags)
	return
}

// Function for creating a new window witohut any default options in place
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
	win = Window{X, screen, windowID,make(chan Event)}

	return
}

// goroutine for a switch that listens for all the default shit we usually want.
func (win *Window) DefaultListeners(ev xgb.Event) {
	switch ev.(type) {
		case xproto.ExposeEvent: 		go win.DrawUIElements()
		case xproto.ButtonPressEvent: 	go win.CheckMousePress(ev) 
		case xproto.ButtonReleaseEvent: go win.CheckMouseRelease(ev) 
		case xproto.KeyPressEvent: 		win.CheckKeyPress(ev)
		case xproto.KeyReleaseEvent: 	win.CheckKeyRelease(ev) 
		case xproto.DestroyNotifyEvent: return
	}
};