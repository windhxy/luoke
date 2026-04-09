//go:build windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	inputKeyboard  = 1
	keyeventfKeyUp = 0x0002
	vk0            = 0x30
)

type keybdInput struct {
	wVk         uint16
	wScan       uint16
	dwFlags     uint32
	time        uint32
	dwExtraInfo uintptr
}

type mouseInput struct {
	dx          int32
	dy          int32
	mouseData   uint32
	dwFlags     uint32
	time        uint32
	dwExtraInfo uintptr
}

type input struct {
	rType uint32
	// Keep INPUT.type followed by pointer-sized alignment, matching WinAPI layout.
	// On 64-bit this becomes 4 bytes (8-4), on 32-bit it becomes 0 bytes.
	_     [unsafe.Sizeof(uintptr(0)) - unsafe.Sizeof(uint32(0))]byte
	// INPUT is a C union (MOUSEINPUT/KEYBDINPUT/HARDWAREINPUT).
	// mouseInput is used because it is the largest union member across supported archs.
	data  [unsafe.Sizeof(mouseInput{})]byte
}

var (
	user32DLL     = syscall.NewLazyDLL("user32.dll")
	sendInputProc = user32DLL.NewProc("SendInput")
)

func pressNumberKey(digit rune) error {
	if digit < '0' || digit > '9' {
		return fmt.Errorf("仅支持数字键 0-9，当前为 %q", string(digit))
	}
	vk := uint16(vk0 + (digit - '0'))

	inputs := []input{
		newKeyboardInput(vk, 0),
		newKeyboardInput(vk, keyeventfKeyUp),
	}

	ret, _, callErr := sendInputProc.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(inputs[0]),
	)
	if ret == 0 {
		if callErr != syscall.Errno(0) {
			return fmt.Errorf("SendInput 调用失败: %w", callErr)
		}
		return fmt.Errorf("SendInput 返回 0")
	}

	return nil
}

func newKeyboardInput(vk uint16, flags uint32) input {
	in := input{rType: inputKeyboard}
	ki := (*keybdInput)(unsafe.Pointer(&in.data[0]))
	ki.wVk = vk
	ki.dwFlags = flags
	return in
}
