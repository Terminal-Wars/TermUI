package TermUI

import (
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

var MousePos struct {
	X 		int16
	Y 		int16
}

func (win *Window) UpdateMouseCoords(ev xgb.Event) {
	MousePos.X = ev.(xproto.MotionNotifyEvent).EventX
	MousePos.Y = ev.(xproto.MotionNotifyEvent).EventY
}

func (win *Window) CheckMouseHover(ev xgb.Event) {
	v := checkMouseOver(MousePos.X, MousePos.Y)
	if(v == nil) {
		if(lastSelectedUIEvent != nil) {
			v = lastSelectedUIEvent
		} else {
			return
		}
	} else {
		lastSelectedUIEvent = v
	}
	lastSelectedUIEvent = v
	win.changeState(v,0)
	win.DrawUIElement(v)
	evn := UIHoverEvent{v.Name, v}
	win.uiEventChan <- evn
}

func (win *Window) CheckMousePress(ev xgb.Event) {
	v := checkMouseOver(MousePos.X, MousePos.Y)
	if(v == nil) {
		if(lastSelectedUIEvent != nil) {
			v = lastSelectedUIEvent
		} else {
			return
		}
	} else {
		lastSelectedUIEvent = v
	}
	lastSelectedUIEvent = v
	win.changeState(v,1)
	win.DrawUIElement(v)
	evn := UIPressEvent{v.Name, v}
	win.uiEventChan <- evn
}

func (win *Window) CheckMouseRelease(ev xgb.Event) {
	v := checkMouseOver(MousePos.X, MousePos.Y)
	if(v == nil) {
		if(lastSelectedUIEvent != nil) {
			v = lastSelectedUIEvent
		} else {
			return
		}
	} else {
		lastSelectedUIEvent = v
	}
	win.changeState(v,2)
	win.DrawUIElement(v)
	evn := UIReleaseEvent{v.Name, v}
	win.uiEventChan <- evn
}

func checkMouseOver(x int16, y int16) (*UIEvent) {
	for i, v := range UIEvents {
		if(uint8(i) == UIEventNum) {break}
		if(x >= v.X && y >= v.Y && x <= v.X+int16(v.Width) && y <= v.Y+int16(v.Height)) {
			return UIEvents[i]
		}
	}
	return nil
}