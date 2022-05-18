package TermUI

import (
	"fmt"

	//"github.com/jezek/xgb/xproto"
)

func (win *Window) TextboxBase(Width uint16, Height uint16, X int16, Y int16) {
	colors := [5]uint32{0xffffff,0xffffff,0xe0e0e0,0x7e7e7e,0x000000}
	win.Base(Width,Height,X,Y,colors)
}

func (win *Window) Textbox(Name string, ID int16, Width uint16, Height uint16, X int16, Y int16) (*UIEvent) {
	if(win.NoMoreUIEvents()) {
		fmt.Println("No more UI events allowed. Refusing to make a textbox.")
		return nil
	}
	if(TextboxNum >= 16) {
		fmt.Println("No more textboxes allowed. Refusing to make a textbox.")
		return nil
	}
	// Create a textbox.
	ev := win.NewUIEvent(Name,ID,Width,Height,X,Y,TextboxType,TextboxNum)
	UIElements.Textboxes[TextboxNum] = ev
	TextboxNum++
	return ev
}