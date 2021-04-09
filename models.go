package main

import (
	"context"
	"github.com/rivo/tview"
)

type System struct {
	CTX    context.Context    // 控制所有协程退出的上下文
	Cancel context.CancelFunc // 退出函数
	APP    *tview.Application
	Pages  *tview.Pages
}
