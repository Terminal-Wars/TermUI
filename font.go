package TermUI

import (
	"fmt"
	"unicode/utf16"

	"github.com/jezek/xgb/xproto"
)

// "After all, writing text is way more comfortable using Xft"
// (refuses to elaborate on how to fucking use xft with this library)

// don't @ me about this code it's all taken from the "shapes.go" example on
// https://github.com/jezek/xgb

func (win *Window) DrawText(Text string, X int16, Y int16, Size int16, FGColor uint32, BGColor uint32) (error) {
	// new connection
	font, err := xproto.NewFontId(win.Conn)
	if err != nil {return err}

	// open font (todo: convert Ubuntu or something into a font that X11 can use)
	fontname := fmt.Sprintf("-*-fixed-*-*-*-*-%v-*-*-*-*-*-*-*",Size)
	err = xproto.OpenFontChecked(win.Conn, font, uint16(len(fontname)), fontname).Check()
	if err != nil {return err}

	// new graphics context and drawable
	gcontext, err := xproto.NewGcontextId(win.Conn)
	if err != nil {return err}
	drawable := xproto.Drawable(win.Window)
	xproto.CreateGC(win.Conn, gcontext, drawable, uint32(xproto.GcForeground) | uint32(xproto.GcBackground), []uint32{FGColor, BGColor})

	buttonText := convertStringToChar2b(Text)

	xproto.CloseFont(win.Conn, font)

	// draw the actual text
	xproto.ImageText16(win.Conn, byte(len(Text)), drawable, gcontext, X, Y, buttonText)

	return nil
}

// i copy and pasted these next two functions out of anger i'm not even gonna try and comment them

func convertStringToChar2b(s string) []xproto.Char2b {
	var chars []xproto.Char2b
	var p []uint16

	for _, r := range []rune(s) {
		p = utf16.Encode([]rune{r})
		if len(p) == 1 {
			chars = append(chars, convertUint16ToChar2b(p[0]))
		} else {
			// If the utf16 representation is larger than 2 bytes
			// we can not use it and insert a blank instead:
			chars = append(chars, xproto.Char2b{Byte1: 0, Byte2: 32})
		}
	}

	return chars
}
func convertUint16ToChar2b(u uint16) xproto.Char2b {
	return xproto.Char2b{
		Byte1: byte((u & 0xff00) >> 8),
		Byte2: byte((u & 0x00ff)),
	}
}