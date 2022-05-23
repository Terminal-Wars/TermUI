package main

import (
	"fmt"
	"log"

	"github.com/Terminal-Wars/TermUI"
	//"github.com/jezek/xgb/xproto"
)

// Example program
func main() {
	// Create a 320x200 window that listens for nothing
	win, err := TermUI.NewWindow(320,200,
		[]uint32{
			0xffffffff,
		})
	if(err != nil) {log.Fatalln(err)}

	// We don't need to check for any TermUI events.
	// Only X events, and only to close the program properly.
	for {
		ev, xerr := win.Conn.WaitForEvent()
		if xerr != nil {fmt.Printf("Error: %s\n", xerr)}
		if ev == nil && xerr == nil {return}
	}
}