//go:build !windows

package main

import "fmt"

func pressNumberKey(digit rune) error {
	return fmt.Errorf("当前系统不支持发送按键: %c（仅支持 Windows）", digit)
}
