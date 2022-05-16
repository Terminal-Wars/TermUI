package TermUI

import (
	"fmt"

	//"github.com/jezek/xgb/xproto"
)

func (win *Window) Label(Name string, ID int16, Width uint16, Height uint16, X int16, Y int16, State int8) {
	if(win.NoMoreUIEvents()) {
		fmt.Println("No more UI events allowed. Refusing to make a textbox.")
		return
	}
	if(LabelNum >= 64) {
		fmt.Println("No more textboxes allowed. Refusing to make a textbox.")
		return
	}
	// Create a textbox.
	ev := win.NewUIEventWithPresetState(Name,ID,Width,Height,X,Y,LabelType,LabelNum,State)
	UIElements.Labels[LabelNum] = ev
	LabelNum++
}

func WordWrap(text string, wrapAt int) (rows []string) {
	for i := 0; i < len(text); i+=wrapAt {
		// If we're near the end of the text...
		if(len(text) < i+wrapAt) {
			// Just append the remainder
			rows = append(rows, text[i:len(text)])
		} else {
			// Otherwise, append the chunk.
			rows = append(rows, text[i:i+wrapAt])
		}
	}
	return
}