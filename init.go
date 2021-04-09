/*
存放初始化函数
*/

package main

import (
	"context"
	"github.com/rivo/tview"
)

func initSystem() {
	system.CTX, system.Cancel = context.WithCancel(context.Background())
	system.APP = tview.NewApplication()
	system.Pages = tview.NewPages()
	initPages()
	system.APP.SetRoot(system.Pages, true).SetFocus(system.Pages) // 将 pages 组件设置成根组件
}

func initPages() { // 初始化 pages
	system.Pages.AddPage(loginFormPage())
}
