package TermUI

import (
	//"fmt"
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

const (
	// Names for special keys.
	NameLeftArrow      = "←"
	NameRightArrow     = "→"
	NameUpArrow        = "↑"
	NameDownArrow      = "↓"
	NameReturn         = "⏎"
	NameEnter          = "⌤"
	NameEscape         = "⎋"
	NameHome           = "⇱"
	NameEnd            = "⇲"
	NameDeleteBackward = "⌫"
	NameDeleteForward  = "⌦"
	NamePageUp         = "⇞"
	NamePageDown       = "⇟"
	NameTab            = "Tab"
	NameSpace          = "Space"
	NameCtrl           = "Ctrl"
	NameShift          = "Shift"
	NameAlt            = "Alt"
	NameSuper          = "Super"
	NameCommand        = "⌘"
	NameF1             = "F1"
	NameF2             = "F2"
	NameF3             = "F3"
	NameF4             = "F4"
	NameF5             = "F5"
	NameF6             = "F6"
	NameF7             = "F7"
	NameF8             = "F8"
	NameF9             = "F9"
	NameF10            = "F10"
	NameF11            = "F11"
	NameF12            = "F12"
	NameBack           = "Back"
)

func (win *Window) checkKeyRelease(ev xgb.Event) {
	if(lastSelectedUIEvent == nil) {return}
	_ = DecodeKey(ev.(xproto.KeyReleaseEvent).Detail)
}

func (win *Window) checkKeyPress(ev xgb.Event) {
	if(lastSelectedUIEvent == nil) {return}
	key := DecodeKey(ev.(xproto.KeyPressEvent).Detail)
	// if it's a shift, alt, or ctrl key, ignore it.
	if(key == "­") {return}
	// What was the last thing we selected?
	switch(lastSelectedUIEvent.Type) {
		case 1: // textbox
			switch(key) {
				// If the length of the character is zero, it's a backspace ⌤
				case "":
					if(len(lastSelectedUIEvent.Name) >= 1) {
						lastSelectedUIEvent.Name = string(lastSelectedUIEvent.Name[0:len(lastSelectedUIEvent.Name)-1])
					}
				// If it's the character for enter, though...
				case "⌤":
					//
				default:
					lastSelectedUIEvent.Name += key
			}
	}
	win.DrawUIElements()
}