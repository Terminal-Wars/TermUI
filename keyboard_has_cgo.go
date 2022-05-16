// (i stole this from gioui btw)

//go:build cgo
// +build cgo

package TermUI

import (
	"errors"
	"fmt"
	"log"
	"os"
	"unicode"
	"unicode/utf8"
	"unsafe"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

/*
#cgo linux pkg-config: xkbcommon xkbcommon-x11 x11-xcb xcb-xkb
#cgo linux LDFLAGS: -lX11

#include <xcb/xkb.h>
#include <stdlib.h>
#include <xkbcommon/xkbcommon.h>
#include <xkbcommon/xkbcommon-compose.h>
#include <xkbcommon/xkbcommon-x11.h>
#include <X11/Xlib.h>
#include <X11/Xlib-xcb.h>

*/
import "C"

// ======
// XKB SHIT
// ======

type Context struct {
	Ctx       *C.struct_xkb_context
	keyMap    *C.struct_xkb_keymap
	state     *C.struct_xkb_state
	compTable *C.struct_xkb_compose_table
	compState *C.struct_xkb_compose_state
	utf8Buf   []byte
}

// Modifiers
type Modifiers uint32

const (
	// ModCtrl is the ctrl modifier key.
	ModCtrl Modifiers = 1 << iota
	// ModCommand is the command modifier key
	// found on Apple keyboards.
	ModCommand
	// ModShift is the shift modifier key.
	ModShift
	// ModAlt is the alt modifier key, or the option
	// key on Apple keyboards.
	ModAlt
	// ModSuper is the "logo" modifier key, often
	// represented by a Windows logo.
	ModSuper
)

var (
	_XKB_MOD_NAME_CTRL  = []byte("Control\x00")
	_XKB_MOD_NAME_SHIFT = []byte("Shift\x00")
	_XKB_MOD_NAME_ALT   = []byte("Mod1\x00")
	_XKB_MOD_NAME_LOGO  = []byte("Mod4\x00")

	disp 				*C.Display
	
	KeyContext 			*Context
	err 				error
)

func (x *Context) Destroy() {
	if x.compState != nil {
		C.xkb_compose_state_unref(x.compState)
		x.compState = nil
	}
	if x.compTable != nil {
		C.xkb_compose_table_unref(x.compTable)
		x.compTable = nil
	}
	x.DestroyKeymapState()
	if x.Ctx != nil {
		C.xkb_context_unref(x.Ctx)
		x.Ctx = nil
	}
}

func init() {
	disp = C.XOpenDisplay(nil);
	// new key context
	KeyContext, err = New()
	if(err != nil) {
		log.Fatalln(err)
	}
	err := KeyContext.UpdateKeymap()
	if(err != nil) {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func New() (*Context, error) {
	ctx := &Context{
		Ctx: C.xkb_context_new(C.XKB_CONTEXT_NO_FLAGS),
	}
	if ctx.Ctx == nil {
		return nil, errors.New("newXKB: xkb_context_new failed")
	}
	locale := os.Getenv("LC_ALL")
	if locale == "" {
		locale = os.Getenv("LC_CTYPE")
	}
	if locale == "" {
		locale = os.Getenv("LANG")
	}
	if locale == "" {
		locale = "C"
	}
	cloc := C.CString(locale)
	defer C.free(unsafe.Pointer(cloc))
	ctx.compTable = C.xkb_compose_table_new_from_locale(ctx.Ctx, cloc, C.XKB_COMPOSE_COMPILE_NO_FLAGS)
	if ctx.compTable == nil {
		ctx.Destroy()
		return nil, errors.New("newXKB: xkb_compose_table_new_from_locale failed")
	}
	ctx.compState = C.xkb_compose_state_new(ctx.compTable, C.XKB_COMPOSE_STATE_NO_FLAGS)
	if ctx.compState == nil {
		ctx.Destroy()
		return nil, errors.New("newXKB: xkb_compose_state_new failed")
	}
	return ctx, nil
}

func (x *Context) DestroyKeymapState() {
	if x.state != nil {
		C.xkb_state_unref(x.state)
		x.state = nil
	}
	if x.keyMap != nil {
		C.xkb_keymap_unref(x.keyMap)
		x.keyMap = nil
	}
}

// SetKeymap sets the keymap and state. The context takes ownership of the
// keymap and state and frees them in Destroy.
func (x *Context) SetKeymap(xkbKeyMap, xkbState unsafe.Pointer) {
	x.DestroyKeymapState()
	x.keyMap = (*C.struct_xkb_keymap)(xkbKeyMap)
	x.state = (*C.struct_xkb_state)(xkbState)
}

func (x *Context) UpdateKeymap() error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered panic!\n", r)
		}
	}()
	KeyContext.DestroyKeymapState()
	ctx := (*C.struct_xkb_context)(unsafe.Pointer(KeyContext.Ctx))
	if disp == nil {
		return errors.New("x11: XOpenDisplay failed")
	}
	xcb := C.XGetXCBConnection(disp)
	if xcb == nil {
		return errors.New("x11: XGetXCBConnection failed")
	}
	xkbDevID := C.xkb_x11_get_core_keyboard_device_id(xcb)
	if xkbDevID == -1 {
		var theError *C.xcb_generic_error_t
		var cookie C.xcb_xkb_get_device_info_cookie_t
		cookie = C.xcb_xkb_get_device_info(xcb, 256,
                                0, 0, 0, 0, 0, 0);
		reply := C.xcb_xkb_get_device_info_reply(xcb, cookie, &theError)
		if(reply != nil) {
			msg := int(theError.full_sequence)
			return errors.New(fmt.Sprintf("x11: xkb_x11_get_core_keyboard_device_id failed. %d",msg))
		} else {
			return errors.New("x11: xkb_x11_get_core_keyboard_device_id failed. so did xcb_xkb_get_device_info_reply.")
		}
	}
	keymap := C.xkb_x11_keymap_new_from_device(ctx, xcb, xkbDevID, C.XKB_KEYMAP_COMPILE_NO_FLAGS)
	if keymap == nil {
		return errors.New("x11: xkb_x11_keymap_new_from_device failed")
	}
	state := C.xkb_x11_state_new_from_device(keymap, xcb, xkbDevID)
	if state == nil {
		C.xkb_keymap_unref(keymap)
		return errors.New("x11: xkb_x11_keymap_new_from_device failed")
	}
	KeyContext.SetKeymap(unsafe.Pointer(keymap), unsafe.Pointer(state))
	return nil
}

func (x *Context) Modifiers() Modifiers {
	var mods Modifiers
	if(x.state == nil) {
		return mods
	}
	if C.xkb_state_mod_name_is_active(x.state, (*C.char)(unsafe.Pointer(&_XKB_MOD_NAME_CTRL[0])), C.XKB_STATE_MODS_EFFECTIVE) == 1 {
		mods |= ModCtrl
	}
	if C.xkb_state_mod_name_is_active(x.state, (*C.char)(unsafe.Pointer(&_XKB_MOD_NAME_SHIFT[0])), C.XKB_STATE_MODS_EFFECTIVE) == 1 {
		mods |= ModShift
	}
	if C.xkb_state_mod_name_is_active(x.state, (*C.char)(unsafe.Pointer(&_XKB_MOD_NAME_ALT[0])), C.XKB_STATE_MODS_EFFECTIVE) == 1 {
		mods |= ModAlt
	}
	if C.xkb_state_mod_name_is_active(x.state, (*C.char)(unsafe.Pointer(&_XKB_MOD_NAME_LOGO[0])), C.XKB_STATE_MODS_EFFECTIVE) == 1 {
		mods |= ModSuper
	}
	return mods
}

func DecodeKey(kc xproto.Keycode) (string) {
	err := KeyContext.UpdateKeymap()
	if(err != nil) {
		fmt.Println(err)
	}
	kcim := KeyContext.DispatchKey(uint32(kc)) // KeyCodeInMap
	return kcim
}

func (x *Context) DispatchKey(keyCode uint32) (string) {
	if x.state == nil {
		return "state is null"
	}
	// get the keycode through the library
	kc := C.xkb_keycode_t(keyCode)
	if len(x.utf8Buf) == 0 {
		x.utf8Buf = make([]byte, 1)
	}

	// get the keysym
	sym := C.xkb_state_key_get_one_sym(x.state, kc)

	// "Feed one keysym to the Compose sequence state machine."
	C.xkb_compose_state_feed(x.compState, sym)
	var str []byte
	switch C.xkb_compose_state_get_status(x.compState) {
	case C.XKB_COMPOSE_CANCELLED, C.XKB_COMPOSE_COMPOSING:
		return "cancelled/composing"
	case C.XKB_COMPOSE_COMPOSED:
		size := C.xkb_compose_state_get_utf8(x.compState, (*C.char)(unsafe.Pointer(&x.utf8Buf[0])), C.size_t(len(x.utf8Buf)))
		if int(size) >= len(x.utf8Buf) {
			x.utf8Buf = make([]byte, size+1)
			size = C.xkb_compose_state_get_utf8(x.compState, (*C.char)(unsafe.Pointer(&x.utf8Buf[0])), C.size_t(len(x.utf8Buf)))
		}
		C.xkb_compose_state_reset(x.compState)
		str = x.utf8Buf[:size]
	case C.XKB_COMPOSE_NOTHING:
		mod := x.Modifiers()
		if mod&(ModCtrl|ModAlt|ModSuper) == 0 {
			str = x.charsForKeycode(kc)
		}
	}
	
	// convert some characters to readable ones
	if(len(str) <= 0) {
		for n, v := range str {
			str_, _ := convertKeysym(C.uint(v))
			str[n] = []byte(str_)[0]
		}
	}

	// Report only printable runes.
	var n int
	for n < len(str) {
		r, s := utf8.DecodeRune(str)
		if unicode.IsPrint(r) {
			n += s
		} else {
			copy(str[n:], str[n+s:])
			str = str[:len(str)-s]
		}
	}

	return string(str)
}

func (x *Context) charsForKeycode(keyCode C.xkb_keycode_t) []byte {
	size := C.xkb_state_key_get_utf8(x.state, keyCode, (*C.char)(unsafe.Pointer(&x.utf8Buf[0])), C.size_t(len(x.utf8Buf)))
	if int(size) >= len(x.utf8Buf) {
		x.utf8Buf = make([]byte, size+1)
		size = C.xkb_state_key_get_utf8(x.state, keyCode, (*C.char)(unsafe.Pointer(&x.utf8Buf[0])), C.size_t(len(x.utf8Buf)))
	}
	return x.utf8Buf[:size]
}

func convertKeysym(s C.xkb_keysym_t) (string, bool) {
	if 'a' <= s && s <= 'z' {
		return string(rune(s - 'a' + 'A')), true
	}
	if ' ' < s && s <= '~' {
		return string(rune(s)), true
	}
	var n string
	switch s {
	case C.XKB_KEY_Escape:
		n = NameEscape
	case C.XKB_KEY_Left:
		n = NameLeftArrow
	case C.XKB_KEY_Right:
		n = NameRightArrow
	case C.XKB_KEY_Return:
		n = NameReturn
	case C.XKB_KEY_KP_Enter:
		n = NameEnter
	case C.XKB_KEY_Up:
		n = NameUpArrow
	case C.XKB_KEY_Down:
		n = NameDownArrow
	case C.XKB_KEY_Home:
		n = NameHome
	case C.XKB_KEY_End:
		n = NameEnd
	case C.XKB_KEY_BackSpace:
		n = NameDeleteBackward
	case C.XKB_KEY_Delete:
		n = NameDeleteForward
	case C.XKB_KEY_Page_Up:
		n = NamePageUp
	case C.XKB_KEY_Page_Down:
		n = NamePageDown
	case C.XKB_KEY_F1:
		n = NameF1
	case C.XKB_KEY_F2:
		n = NameF2
	case C.XKB_KEY_F3:
		n = NameF3
	case C.XKB_KEY_F4:
		n = NameF4
	case C.XKB_KEY_F5:
		n = NameF5
	case C.XKB_KEY_F6:
		n = NameF6
	case C.XKB_KEY_F7:
		n = NameF7
	case C.XKB_KEY_F8:
		n = NameF8
	case C.XKB_KEY_F9:
		n = NameF9
	case C.XKB_KEY_F10:
		n = NameF10
	case C.XKB_KEY_F11:
		n = NameF11
	case C.XKB_KEY_F12:
		n = NameF12
	case C.XKB_KEY_Tab, C.XKB_KEY_KP_Tab, C.XKB_KEY_ISO_Left_Tab:
		n = NameTab
	case 0x20, C.XKB_KEY_KP_Space:
		n = NameSpace
	case C.XKB_KEY_Control_L, C.XKB_KEY_Control_R:
		n = NameCtrl
	case C.XKB_KEY_Shift_L, C.XKB_KEY_Shift_R:
		n = NameShift
	case C.XKB_KEY_Alt_L, C.XKB_KEY_Alt_R:
		n = NameAlt
	case C.XKB_KEY_Super_L, C.XKB_KEY_Super_R:
		n = NameSuper
	default:
		return "", false
	}
	return n, true
}

// We have functions to call an unexported function because under cgo these can't be
// run as go routines, but without it they can, so if we have the opprutunity we should 
// run these on seperate threads.

func (win *Window) CheckKeyRelease(ev xgb.Event) {
	win.checkKeyRelease(ev)
}
func (win *Window) CheckKeyPress(ev xgb.Event) {
	win.checkKeyPress(ev)
}