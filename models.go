/*
存放自己定义的数据类型
*/

package main

import (
	"context"
	"github.com/rivo/tview"
	"net"
)

type System struct {
	CTX         context.Context    // 控制所有协程退出的上下文
	Cancel      context.CancelFunc // 退出函数
	APP         *tview.Application
	messageBox  *tview.TextView // 输出信息的地方
	Pages       *tview.Pages
	userConn    *userConnection // 存储用户连接信息
	packageChan chan string
}

type userConnection struct {
	host string
	port string
	nick string
	conn net.Conn
}
