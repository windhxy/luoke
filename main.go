package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	key := flag.String("key", "1", "要按下的数字键，仅支持 0-9")
	interval := flag.Duration("interval", 3*time.Second, "按键间隔，例如 3s、500ms")
	flag.Parse()

	digit, err := validateDigitKey(*key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "参数错误: %v\n", err)
		os.Exit(1)
	}
	if *interval <= 0 {
		fmt.Fprintln(os.Stderr, "参数错误: -interval 必须大于 0")
		os.Exit(1)
	}

	fmt.Printf("开始定时按键: key=%c interval=%s\n", digit, *interval)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

	for {
		select {
		case <-sigCh:
			fmt.Println("收到退出信号，程序结束。")
			return
		case <-ticker.C:
			if err := pressNumberKey(digit); err != nil {
				fmt.Fprintf(os.Stderr, "按键失败: %v\n", err)
			}
		}
	}
}

func validateDigitKey(key string) (rune, error) {
	if len(key) != 1 || key[0] < '0' || key[0] > '9' {
		return 0, fmt.Errorf("-key 仅支持单个数字字符 0-9，当前为 %q", key)
	}
	return rune(key[0]), nil
}
