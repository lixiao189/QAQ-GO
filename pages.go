/*
用来生成所有页面的代码
*/

package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"net"
	"time"
)

func loginFormPage() (name string,
	form *tview.Form,
	resize bool,
	visible bool) { // 初始化登陆页面
	name = "loginForm"
	resize = true
	visible = true
	form = tview.NewForm().
		AddInputField("地址: ", system.userConn.host, 60, nil, nil).
		AddInputField("端口: ", system.userConn.port, 60, nil, nil).
		AddInputField("昵称: ", system.userConn.nick, 60, nil, nil).
		AddButton("连接", func() {
			// 更新连接信息
			system.userConn.host = form.GetFormItem(0).(*tview.InputField).GetText()
			system.userConn.port = form.GetFormItem(1).(*tview.InputField).GetText()
			system.userConn.nick = form.GetFormItem(2).(*tview.InputField).GetText()

			var err error
			system.userConn.conn, err = net.Dial("tcp", system.userConn.host+":"+system.userConn.port)
			if err == nil {
				system.Pages.SwitchToPage("chat")
				go messageServices() // 启动后台接收信息服务
			} else {
				system.Pages.SwitchToPage("connectError")
			}
		}).
		AddButton("退出", func() {
			system.APP.Stop()
		})
	form.SetBorderPadding(2, 0, 10, 0)
	form.SetBorder(false)
	form.SetFieldBackgroundColor(tcell.ColorBlack)

	return
}

func connectErrorPage() (name string, modal *tview.Modal, resize bool, visible bool) {
	name = "connectError"
	resize = false
	visible = false

	modal = tview.NewModal().
		SetText("Connect Error").
		AddButtons([]string{"Reconnect", "Quit"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonIndex == 1 {
				system.APP.Stop()
			} else {
				system.Pages.SwitchToPage("loginForm")
			}
		})
	modal.SetBackgroundColor(tcell.ColorBlack)

	return
}

func aboutPage() (name string, modal *tview.Modal, resize bool, visible bool) {
	name = "about"
	resize = false
	visible = false

	modal = tview.NewModal().SetText("About\nA QAQ TUI client which is written in GO\nWritter: Node Sans").
		AddButtons([]string{"Yes"}).
		SetDoneFunc(func(buttonIndex int, buttonLable string) {
			system.Pages.SwitchToPage("chat")
		})
	modal.SetBackgroundColor(tcell.ColorBlack)

	return
}

func chatPage() (name string,
	rootFlex *tview.Flex,
	resize bool,
	visible bool) {
	name = "chat"
	resize = true
	visible = false

	var primitives []tview.Primitive
	rootFlex = tview.NewFlex() // 根容器
	aboutButton := tview.NewButton("关于")
	quitButton := tview.NewButton("退出")
	disconnectButton := tview.NewButton("断开连接")
	messageBox := tview.NewTextView() // 输出窗口
	inputField := tview.NewInputField()
	system.messageBox = messageBox

	primitives = append(primitives, aboutButton)
	primitives = append(primitives, quitButton)
	primitives = append(primitives, disconnectButton)
	primitives = append(primitives, messageBox)
	primitives = append(primitives, inputField)
	currentPrimitiveIndex := 4

	aboutButton.SetSelectedFunc(func() {
		system.Pages.SwitchToPage("about")
	})

	quitButton.SetSelectedFunc(func() {
		_ = system.userConn.conn.Close()
		system.Cancel() // 关闭所有后台协程
		system.APP.Stop()
	})

	disconnectButton.SetSelectedFunc(func() {
		_ = system.userConn.conn.Close()
		system.Cancel()
		system.Pages.SwitchToPage("loginForm")
	})

	inputField.SetLabel(" >> ").
		SetPlaceholderTextColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorBlack).
		SetLabelColor(tcell.ColorBlack).
		SetFieldBackgroundColor(tcell.ColorWhite).
		SetBackgroundColor(tcell.ColorWhite)

	messageBox.SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			system.APP.Draw()
		})

	rootFlex.SetDirection(tview.FlexRow)
	rootFlex.SetBorderPadding(0, 1, 0, 0)
	rootFlex.AddItem(tview.NewFlex().
		AddItem(aboutButton, 0, 1, false).
		AddItem(quitButton, 0, 1, false).
		AddItem(disconnectButton, 0, 1, false),
		1, 0, false)
	rootFlex.AddItem(messageBox, 0, 1, false)
	rootFlex.AddItem(inputField, 1, 0, true)
	rootFlex.SetBorderPadding(0, 1, 0, 0)

	// 设置按键绑定
	rootFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab { // TAB 键切换选中状态
			if currentPrimitiveIndex == len(primitives)-1 {
				currentPrimitiveIndex = 0
			} else {
				currentPrimitiveIndex++
			}
			system.APP.SetFocus(primitives[currentPrimitiveIndex])
		}
		if event.Key() == tcell.KeyEnter && inputField.HasFocus() && inputField.GetText() != "" { // 发送
			// 在 messageBox 中显示个人消息内容
			_, _ = fmt.Fprintln(system.messageBox, system.userConn.nick, time.Now().Format("2006-01-02 15:04:05"))
			_, _ = fmt.Fprintln(system.messageBox, inputField.GetText())
			_, _ = fmt.Fprint(system.messageBox, "\n")
			system.messageBox.ScrollToEnd() // 发送信息后自动滚到底部
			err := sendToServer(inputField.GetText())
			if err != nil {
				system.Pages.SwitchToPage("connectError")
				_ = system.userConn.conn.Close()
				system.Cancel()
			}
			inputField.SetText("")
		}

		return event
	})

	return
}
