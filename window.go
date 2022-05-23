package TermUI

import (
	"fmt"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

var Atoms map[string]xproto.Atom

type Window struct {
	Conn 			*xgb.Conn
	Screen 			*xproto.ScreenInfo
	Window 			xproto.Window
	Width 			uint16
	Height 			uint16
	X 				int16
	Y  				int16

	uiEventChan		chan Event
}

func init() {
	Atoms = make(map[string]xproto.Atom, 50)
}

// function for returning a []uint32 with the default flags 
func defaultFlags(flags []uint32, size int) ([]uint32) {
	var numOfDefaults int = 5
	var newSize int = size+numOfDefaults
	// If we have a second flag defined, return an array 
	// based on the amount of flags we have in it.
	if(len(flags) >= 2) {
		newSize = size+numOfDefaults+xgb.PopCount(int(flags[1]))
	}
	newFlags := make([]uint32, newSize)

	var defaults uint32 = xproto.EventMaskStructureNotify | 
	xproto.EventMaskExposure | 
	xproto.EventMaskPointerMotion |
	xproto.EventMaskButtonPress |
	xproto.EventMaskButtonRelease

	if(len(flags) >= 1) {newFlags[0] = flags[0]}
	if(len(flags) >= 2) {
		newFlags[1] = flags[1] | defaults
	} else {
		newFlags[1] = defaults
	}
	// Any null flags should be 1, actually.
	for i, v := range newFlags {
		if(v == 0) {
			newFlags[i] = 1
		}
	}

	return newFlags
}

// Function for creating a new window, with default options in place.
func NewWindow(Width, Height uint16, Flags []uint32) (Window Window, Error error) {
	flags := defaultFlags(Flags,2)
	Window, Error = NewWindowComplex(Width, Height, 0, xproto.CwBackPixel | xproto.CwEventMask, flags)
	return
}

func (win *Window) NewChildWindow(Width, Height uint16, Flags []uint32) (Window Window, Error error) {
	flags := defaultFlags(Flags,2)
	Window, Error = win.NewChildWindowComplex(Width, Height, 0, xproto.CwBackPixel | xproto.CwEventMask, flags)
	return
}

// The same function as above, but the window is made undecorated.
func NewUndecoratedWindow(Width, Height uint16, Flags []uint32) (Window Window, Error error) {
	flags := defaultFlags(Flags,2)
	Window, Error = NewWindowComplex(Width, Height, 0, xproto.CwBackPixel | xproto.CwEventMask, flags)
	
	data := []uint{2, 0, 0, 0, 0}
	ChangeProp(&Window, Window.Window, 32, "_MOTIF_WM_HINTS", "_MOTIF_WM_HINTS",data...)

	return
}

func (win *Window) NewUndecoratedChildWindow(Width, Height uint16, Flags []uint32) (Window Window, Error error) {
	flags := defaultFlags(Flags,2)
	Window, Error = win.NewChildWindowComplex(Width, Height, 0, xproto.CwBackPixel | xproto.CwEventMask, flags)
	
	data := []uint{2, 0, 0, 0, 0}
	ChangeProp(&Window, Window.Window, 32, "_MOTIF_WM_HINTS", "_MOTIF_WM_HINTS",data...)

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
	win = Window{
		Conn: X,
		Screen: screen,
		Window: windowID,
		Width: Width,
		Height: Height,
		uiEventChan: make(chan Event),
	}

	return
}

// The same function, but for child windows.
func (win *Window) NewChildWindowComplex(Width, Height, BorderWidth uint16, Mask uint32, Flags []uint32) (cwin Window, err error) {

	// new window ID
	windowID, err := xproto.NewWindowId(win.Conn)
	if err != nil {return}

	// create a window using that ID
	err = xproto.CreateWindowChecked(win.Conn, win.Screen.RootDepth, windowID, win.Screen.Root,
		0, 0, Width, Height, BorderWidth,
		xproto.WindowClassInputOutput, win.Screen.RootVisual, 
		Mask,
		Flags).Check()
	if(err != nil) {return}

	// make that window appear on screen.
	err = xproto.MapWindowChecked(win.Conn, windowID).Check()
	if err != nil {return}

	// make a new Window object with what we've gotten here
	cwin = Window{
		Conn: win.Conn,
		Screen: win.Screen,
		Window: windowID,
		Width: Width,
		Height: Height,
		uiEventChan: make(chan Event),
	}

	return
}


// copy/pasted and modified functions from xgbutil, a mostly undocumented library

// ChangeProperty abstracts the semi-nastiness of xgb.ChangeProperty.
func ChangeProp(win *Window, win_ xproto.Window, format byte, prop string,
	typ string, data ...uint) error {

	buf := make([]byte, len(data)*4)
	for i, datum := range data {
		xgb.Put32(buf[(i*4):], uint32(datum))
	}

	propAtom, err := Atom(win, prop, false)
	if err != nil {
		return err
	}

	typAtom, err := Atom(win, typ, false)
	if err != nil {
		return err
	}

	return xproto.ChangePropertyChecked(win.Conn, xproto.PropModeReplace, win_,
		propAtom, typAtom, format,
		uint32(len(buf)/(int(format)/8)), buf).Check()
}

func Atom(win *Window, name string, onlyIfExists bool) (xproto.Atom, error) {
	// Check the cache first
	if aid, ok := atomGet(name); ok {
		return aid, nil
	}

	reply, err := xproto.InternAtom(win.Conn, onlyIfExists,
		uint16(len(name)), name).Reply()
	if err != nil {
		return 0, fmt.Errorf("Atom: Error interning atom '%s': %s", name, err)
	}

	// If we're here, it means we didn't have this atom cached. So cache it!
	cacheAtom(name, reply.Atom)

	return reply.Atom, nil
}

func atomGet(name string) (xproto.Atom, bool) {
	// leftover code from xkbutil, commenting it out
	// might be important idk
	//xu.AtomsLck.RLock()
	//defer xu.AtomsLck.RUnlock()

	aid, ok := Atoms[name]
	return aid, ok
}

func cacheAtom(name string, aid xproto.Atom) {
	//xu.AtomsLck.Lock()
	//defer xu.AtomsLck.Unlock()

	Atoms[name] = aid
}

// functions for getting percentages of the window

func (win *Window) PercentOfWidth(per uint16) (int16) {
	return int16(float32(win.Width)*(float32(per)/100))
}

func (win *Window) PercentOfHeight(per uint16) (int16) {
	fmt.Println()
	return int16(float32(win.Height)*(float32(per)/100))
}

func (win *Window) UpdateSavedLocation(ev xgb.Event) {
	win.X = ev.(xproto.ConfigureNotifyEvent).X
	win.Y = ev.(xproto.ConfigureNotifyEvent).Y
}

// goroutine for a switch that listens for all the default shit we usually want.
func (win *Window) DefaultListeners(ev xgb.Event) {
	switch ev.(type) {
		case xproto.ExposeEvent: 			win.DrawUIElements()
		case xproto.MotionNotifyEvent:  	win.UpdateMouseCoords(ev)
		case xproto.ButtonPressEvent: 		go win.CheckMousePress(ev) 
		case xproto.ButtonReleaseEvent: 	go win.CheckMouseRelease(ev) 
		case xproto.KeyPressEvent: 			win.CheckKeyPress(ev)
		case xproto.KeyReleaseEvent: 		win.CheckKeyRelease(ev) 
		case xproto.ConfigureNotifyEvent: 	win.UpdateSavedLocation(ev)
		case xproto.DestroyNotifyEvent: return
	}
}