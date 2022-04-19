package TermUI

import (
	"fmt"

	"github.com/jezek/xgb/xproto"
)

func Base(win *Window, Width uint16, Height uint16, X int16, Y int16) {

	conn := win.Conn
	drawable := win.Drawable
	fmt.Println(conn)

	// Graphics context
	foreground, err := xproto.NewGcontextId(conn)
	if(err != nil) {
		fmt.Println("error creating base context:", err)
		return
	}

	// Create that graphics context
	mask := uint32(xproto.GcLineWidth)
	values := []uint32{10}
	xproto.CreateGC(conn, foreground, drawable, mask, values)

	rectangle := []xproto.Rectangle{
		{X: X, Y: Y, Width: Width, Height: Height},
	}
	
	// Create the rectangle
	xproto.PolyRectangle(conn, drawable, foreground, rectangle)
	fmt.Println("j")
} 