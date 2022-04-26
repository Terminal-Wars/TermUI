package TermUI

// #include "/tmp/keymap.h"
import "C"

import (
	"fmt"
	"bytes"
	"unicode/utf8"
	"unicode/utf16"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

var HoldingShift int8 	= -1
var HoldingAlt int8 	= -1
var HoldingAltGr int8 	= -1
var HoldingCtrl int8 	= -1

func (win *Window) CheckKeyRelease(ev xgb.Event) {
	_ = DecodeKey(ev.(xproto.KeyReleaseEvent).Detail)
}

func (win *Window) CheckKeyPress(ev xgb.Event) {
	_ = DecodeKey(ev.(xproto.KeyPressEvent).Detail)
}

func DecodeKey(kc xproto.Keycode) (string) {
	// get the utf-16 keycode that matches to our key but subtract it by fucking 8.
	// (this is stupidly consistent and i don't fucking know why)
	kcim := uint16(C.plain_map[kc-8]) // KeyCodeInMap

	// convert it to utf8 (using golang's crackhead methods)
	ret := &bytes.Buffer{}
	bytebuf := make([]byte, 4)
	uint16s := make([]uint16, 1)

	switch(kcim) {
		case 63232: HoldingShift *= -1
		case 63233: HoldingAlt *= -1
		case 63234: HoldingCtrl *= -1
	}
	fmt.Println(kcim)

	// the first two bytes in the utf16 value aren't actaully useful to us and only exist to throw
	// off the final letter. golang's way of working with bits is messy (and i couldn't even get it
	// working right) so we'll just subtract accordingly
	// todo: find out how to get the last two bits because it's probably faster and more predictable 
	switch {
		case (kcim >= 0xf000 && kcim <= 0xf0ff): uint16s[0] = kcim-0xf000
		case (kcim >= 0xf100 && kcim <= 0xf1ff): uint16s[0] = kcim-0xf100
		case (kcim >= 0xf200 && kcim <= 0xf2ff): uint16s[0] = kcim-0xf200
		case (kcim >= 0xf300 && kcim <= 0xf3ff): uint16s[0] = kcim-0xf300
		case (kcim >= 0xf400 && kcim <= 0xf4ff): uint16s[0] = kcim-0xf400
		case (kcim >= 0xf500 && kcim <= 0xf5ff): uint16s[0] = kcim-0xf500
		case (kcim >= 0xf600 && kcim <= 0xf6ff): uint16s[0] = kcim-0xf600
		case (kcim >= 0xf700 && kcim <= 0xf7ff): uint16s[0] = kcim-0xf700
		case (kcim >= 0xf800 && kcim <= 0xf8ff): uint16s[0] = kcim-0xf800
		case (kcim >= 0xf900 && kcim <= 0xf9ff): uint16s[0] = kcim-0xf900
		case (kcim >= 0xfa00 && kcim <= 0xfaff): uint16s[0] = kcim-0xfa00
		case (kcim >= 0xfb00 && kcim <= 0xfbff): uint16s[0] = kcim-0xfb00
		case (kcim >= 0xfc00 && kcim <= 0xfcff): uint16s[0] = kcim-0xfc00
		case (kcim >= 0xfd00 && kcim <= 0xfdff): uint16s[0] = kcim-0xf600
		case (kcim >= 0xfe00 && kcim <= 0xfeff): uint16s[0] = kcim-0xfe00
		case (kcim >= 0xff00 && kcim <= 0xffff): uint16s[0] = kcim-0xff00
	}

	 // subtract whatever we have by 64256 btw
	foo := utf16.Decode(uint16s)
	_ = utf8.EncodeRune(bytebuf, foo[0])

	// it gives us the correct key
	ret.Write(bytebuf)

	return ret.String()
}
