package main

import (
	"fmt"
	"log"

	"github.com/Terminal-Wars/TermUI"
	"github.com/jezek/xgb/xproto"
)

// Example program
func main() {
	win, err := TermUI.NewWindow(320,200,
		[]uint32{
			0xffffffff
		})
	if(err != nil) {log.Fatalln(err)}

	go func() {
		for {
			ev := win.WaitForUIEvent()
			/*switch ev.(type) {
			}*/
		}
	}()
	for {
		ev, xerr := win.Conn.WaitForEvent()
		if xerr != nil {fmt.Printf("Error: %s\n", xerr)}
		if ev == nil && xerr == nil {return}
		win.DefaultListeners(ev);
	}
}