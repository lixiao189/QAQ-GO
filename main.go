package main

import (
	"context"
	"github.com/rivo/tview"
)

var system System

func main() {
	system.CTX, system.Cancel = context.WithCancel(context.Background())
	system.APP = tview.NewApplication()
	system.Pages = tview.NewPages()
	initPages() // 初始化所有页面

	system.APP.SetRoot(system.Pages, true).SetFocus(system.Pages) // 将 pages 组件设置成根组件
	_ = system.APP.Run()                                          // 启动应用
}
