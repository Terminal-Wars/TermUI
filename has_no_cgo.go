// +build !cgo

package TermUI

import (
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

const cgoEnabled = false

func init() {
	fmt.Println("built without cgo. some features may be disabled.")
}

// goroutine for a switch that listens for all the default shit we usually want.
func (win *Window) DefaultListeners(ev xgb.Event) {
	switch ev.(type) {
		case xproto.ExposeEvent: 		go win.DrawUIElements()
		case xproto.ButtonPressEvent: 	go win.CheckMousePress(ev) 
		case xproto.ButtonReleaseEvent: go win.CheckMouseRelease(ev) 
		case xproto.DestroyNotifyEvent: return
	}
};