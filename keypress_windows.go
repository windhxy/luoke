//go:build windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	inputKeyboard  = 1
	keyeventfKeyup = 0x0002
	vk0            = 0x30
)

type keybdInput struct {
	wVk         uint16
	wScan       uint16
	dwFlags     uint32
	time        uint32
	dwExtraInfo uintptr
}

type input struct {
	rType uint32
	ki    keybdInput
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
		{
			rType: inputKeyboard,
			ki: keybdInput{
				wVk: vk,
			},
		},
		{
			rType: inputKeyboard,
			ki: keybdInput{
				wVk:     vk,
				dwFlags: keyeventfKeyup,
			},
		},
	}

	ret, _, callErr := sendInputProc.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(input{}),
	)
	if ret == 0 {
		if callErr != syscall.Errno(0) {
			return fmt.Errorf("SendInput 调用失败: %w", callErr)
		}
		return fmt.Errorf("SendInput 返回 0")
	}

	return nil
}
