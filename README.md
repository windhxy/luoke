# luoke

一个使用 Go 编写的定时按键小工具。  
支持配置按键数字与触发间隔，主要面向 Windows（通过 `user32.dll` 的 `SendInput` 发送扫描码按键事件）。

## 使用方式

```bash
go run . -key 1 -interval 3s
```

参数说明：

- `-key`：要按下的数字键，仅支持 `0-9`，默认 `1`
- `-interval`：按键间隔，Go duration 格式，默认 `3s`

按 `Ctrl + C` 退出。

## 常见问题

- Windows 若出现 `SendInput 调用失败: The parameter is incorrect.`，通常是 `INPUT` 结构体大小与 WinAPI 不一致导致。当前实现已按 WinAPI 对齐处理。
- 对于部分全屏游戏，虚拟键码（VK）注入可能不被识别。当前实现已改为优先发送扫描码（Scan Code）事件，兼容性更好。
