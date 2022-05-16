package TermUI

import (
	//"fmt"
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

func (win *Window) checkKeyRelease(ev xgb.Event) {
	if(lastSelectedUIEvent == nil) {return}
	_ = DecodeKey(ev.(xproto.KeyReleaseEvent).Detail)
}

func (win *Window) checkKeyPress(ev xgb.Event) {
	if(lastSelectedUIEvent == nil) {return}
	key := DecodeKey(ev.(xproto.KeyPressEvent).Detail)
	// What was the last thing we selected?
	switch(lastSelectedUIEvent.Type) {
		case 1: // textbox
			if(len(key) == 0 && len(lastSelectedUIEvent.Name) >= 1) { // If the length of the character is zero, it's a backspace
				lastSelectedUIEvent.Name = string(lastSelectedUIEvent.Name[0:len(lastSelectedUIEvent.Name)-1])
			} else {
				lastSelectedUIEvent.Name += key
			}
	}
	win.DrawUIElements()
}