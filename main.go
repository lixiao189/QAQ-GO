package main

var system System

func main() {
	initSystem() // 初始化系统

	_ = system.APP.Run() // 启动应用
}
