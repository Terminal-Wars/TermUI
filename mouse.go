package TermUI

import (

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

var MousePos struct {
	X 		int16
	Y 		int16
}

func (win *Window) CheckMouseMove(ev xgb.Event) {
	MousePos.X = ev.(xproto.MotionNotifyEvent).EventX
	MousePos.Y = ev.(xproto.MotionNotifyEvent).EventY
	for i, v := range UIEvents {
		if(int16(i) == UIEventNum) {break}
		if(MousePos.X >= v.X && MousePos.Y >= v.Y && MousePos.X <= v.X+v.Width && MousePos.Y <= v.Y+v.Height) {
			ev := UIHoverEvent{v.ID, v.Name}
			win.uiEventChan <- ev
		}
	}
}