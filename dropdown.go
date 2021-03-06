package TermUI

import (
	"errors"
	"fmt"
	"github.com/jezek/xgb/xproto"
)

var Dropdowns []*Window

func (win *Window) Dropdown(X, Y int16, Strings []string) (error) {
	// Go through each of the options and find out which is longest.
	var width uint16 = 0
	var height uint16 = uint16(len(Strings))
	for _, v := range Strings {
		if(uint16(len(v)) > width) {
			width = uint16(len(v))
		}
	}
	// Create a new window based on all the other values
	dropdown, err := win.NewUndecoratedChildWindow(width*8,height*16,
			[]uint32{
				0xffffffff,
			})
	if(err != nil) {return errors.New("Couldn't create error: "+err.Error())}
	Dropdowns = append(Dropdowns,&dropdown)

	if err != nil {
		fmt.Errorf("Couldn't move window: %s", err)
	}	

	// Then draw each of the text strings in order.
	for i, v := range Strings {
		dropdown.DrawText(v,4,int16(12*i)+12,12,0x000000,0xffffff)
	}

	return nil
}

func (win *Window) ClearDropdowns() {
	for i, v := range Dropdowns {
		xproto.DestroyWindowChecked(v.Conn, v.Window)
		if(err != nil) {
			fmt.Printf("Couldn't destroy dropdown %d: %s", i, err)
		}
	}
}