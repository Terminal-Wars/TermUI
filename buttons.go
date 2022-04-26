package TermUI

import (
	"fmt"

	"github.com/jezek/xgb/xproto"
)

func (win *Window) Base(Width uint16, Height uint16, X int16, Y int16, colors [5]uint32) {
	// Drawing contexts
	draw := [5]xproto.Drawable{
		xproto.Drawable(win.Window),
		xproto.Drawable(win.Window),
		xproto.Drawable(win.Window),
		xproto.Drawable(win.Window),
		xproto.Drawable(win.Window),
	}

	// Graphics context
	gcontext := [5]xproto.Gcontext{}

	// Create a gcontext for each setting
	for i, v := range colors {
		gcontext_, err := xproto.NewGcontextId(win.Conn)
		if(err != nil) {
			fmt.Println("error creating context ",i,err)
			return
		}
		gcontext[i] = gcontext_

		mask := uint32(xproto.GcForeground)
		values := []uint32{v}

		xproto.CreateGC(win.Conn, gcontext[i], draw[i], mask, values)
	}

	// Now we create the shapes that make up the base of the button.
	rectangles := [][]xproto.Rectangle{
		// bg
		[]xproto.Rectangle{{X: X, Y: Y, Width: Width, Height: Height}},

		// outer shadow
		[]xproto.Rectangle{
			{X: X, Y: Y+int16(Height), Width: Width+1, Height: 1},
			{X: X+int16(Width), Y: Y, Width: 1, Height: Height},
		}, 

		// inner shadow
		[]xproto.Rectangle{
			{X: X+int16(Width-1), Y: Y, Width: 1, Height: Height-1},
			{X: X+2, Y: Y+int16(Height-1), Width: Width-2, Height: 1},
		},

		// outer highlight
		[]xproto.Rectangle{
			{X: X, Y: Y, Width: Width, Height: 1},
			{X: X, Y: Y, Width: 1, Height: Height},
		}, 

		// inner highlight
		[]xproto.Rectangle{
			{X: X+1, Y: Y+1, Width: Width-1, Height: 1},
			{X: X+1, Y: Y+1, Width: 1, Height: Height-1},
		},

	}

	for i, v := range rectangles {
		xproto.PolyFillRectangle(win.Conn, draw[i], gcontext[i], v)
	}
} 

func (win *Window) BaseRaised(Width uint16, Height uint16, X int16, Y int16) {
	colors := [5]uint32{0xafb5b5,0x000000,0x808080,0xffffff,0xd7dfdf}
	win.Base(Width,Height,X,Y,colors)
}

func (win *Window) BaseSunken(Width uint16, Height uint16, X int16, Y int16) {
	colors := [5]uint32{0xafb5b5,0xffffff,0xd7dfdf,0x000000,0x808080}
	win.Base(Width,Height,X,Y,colors)
}

func (win *Window) Button(Name string, ID int16, Width uint16, Height uint16, X int16, Y int16) {
	if(win.NoMoreUIEvents()) {
		fmt.Println("No more UI events allowed. Refusing to make a button.")
		return
	}
	if(ButtonNum >= 16) {
		fmt.Println("No more buttons allowed. Refusing to make a button.")
		return
	}
	// Create a button.
	ev := win.NewUIEvent(Name,ID,Width,Height,X,Y,0,ButtonNum)
	UIElements.Buttons[ButtonNum] = ev
	ButtonNum++
}