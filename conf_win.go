//go:build windows

package climenu

import (
	"os"

	"golang.org/x/sys/windows"
)

// 初始化终端，使终端可以接受\033码等控制字符
func init() {
	// https://learn.microsoft.com/zh-cn/windows/console/console-virtual-terminal-sequences
	var outMode uint32
	out := windows.Handle(os.Stdout.Fd())
	if err := windows.GetConsoleMode(out, &outMode); err != nil {
		return
	}
	if outMode&windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING == 0 {
		outMode |= windows.ENABLE_PROCESSED_OUTPUT | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		windows.SetConsoleMode(out, outMode)
	}
}
