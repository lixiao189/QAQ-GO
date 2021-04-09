package main

import "github.com/rivo/tview"

func loginFormPage() (name string,
	primitive *tview.Form,
	resize bool,
	visible bool) { // 初始化登陆页面
	name = "loginForm"
	resize = true
	visible = true
	primitive = tview.NewForm().
		AddInputField("地址", "127.0.0.1", 60, nil, nil).
		AddInputField("端口", "8080", 60, nil, nil).
		AddInputField("昵称", "", 60, nil, nil).
		AddButton("连接", nil).
		AddButton("退出", func() {
			system.APP.Stop()
		})
	primitive.SetBorder(false)
	return
}
