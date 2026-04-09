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
	keyeventfScan  = 0x0008
	mapvkVkToVsc   = 0
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
	// Padding ensures union data starts at pointer-aligned offset, matching WinAPI INPUT.
	// On 64-bit this adds 4 bytes (8-4), on 32-bit no padding is needed (4-4=0).
	_     [unsafe.Sizeof(uintptr(0)) - unsafe.Sizeof(uint32(0))]byte
	// INPUT is a C union (MOUSEINPUT/KEYBDINPUT/HARDWAREINPUT).
	// mouseInput is used because it is the largest union member across supported archs.
	data  [unsafe.Sizeof(mouseInput{})]byte
}

var (
	user32DLL       = syscall.NewLazyDLL("user32.dll")
	sendInputProc   = user32DLL.NewProc("SendInput")
	mapVirtualKeyW  = user32DLL.NewProc("MapVirtualKeyW")
)

func pressNumberKey(digit rune) error {
	if digit < '0' || digit > '9' {
		return fmt.Errorf("仅支持数字键 0-9，当前为 %q", string(digit))
	}
	vk := uint16(vk0 + (digit - '0'))
	scanCode, err := lookupDigitScanCode(vk, digit)
	if err != nil {
		return fmt.Errorf("无法获取按键扫描码: %q: %w", string(digit), err)
	}

	inputs := []input{
		newKeyboardInput(vk, scanCode, keyeventfScan),
		newKeyboardInput(vk, scanCode, keyeventfScan|keyeventfKeyUp),
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

func newKeyboardInput(vk uint16, scanCode uint16, flags uint32) input {
	in := input{rType: inputKeyboard}
	// Safe: data is sized to the largest INPUT union member (mouseInput), larger than keybdInput.
	ki := (*keybdInput)(unsafe.Pointer(&in.data[0]))
	ki.wVk = vk
	ki.wScan = scanCode
	ki.dwFlags = flags
	return in
}

func lookupDigitScanCode(vk uint16, digit rune) (uint16, error) {
	if err := mapVirtualKeyW.Find(); err != nil {
		return 0, fmt.Errorf("MapVirtualKeyW 不可用: %w", err)
	}
	sc, _, _ := mapVirtualKeyW.Call(uintptr(vk), uintptr(mapvkVkToVsc))
	if sc != 0 {
		return uint16(sc), nil
	}
	// Fallback to standard keyboard top-row digit scan codes.
	switch digit {
	case '1':
		return 0x02, nil
	case '2':
		return 0x03, nil
	case '3':
		return 0x04, nil
	case '4':
		return 0x05, nil
	case '5':
		return 0x06, nil
	case '6':
		return 0x07, nil
	case '7':
		return 0x08, nil
	case '8':
		return 0x09, nil
	case '9':
		return 0x0A, nil
	case '0':
		return 0x0B, nil
	default:
		return 0, fmt.Errorf("不支持的数字键: %q", string(digit))
	}
}
