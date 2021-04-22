/*
所有的后台服务
*/

package main

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func messageServices() {
	base64Nick := base64.StdEncoding.EncodeToString([]byte(system.userConn.nick))
	_, _ = system.userConn.conn.Write([]byte("{user&;named&;" + base64Nick + "}")) // 设置用户名
	getHistoryMessage()                                                            // 查询历史记录

	isStarted := false
	var packageData []byte
	for {
		select {
		case <-system.CTX.Done(): // 停止接收消息
			return
		default:
			serverData := make([]byte, 128)
			n, err := system.userConn.conn.Read(serverData)
			if err != nil {
				system.Pages.SwitchToPage("connectError")
				_ = system.userConn.conn.Close()
				return
			}
			for i := 0; i < n; i++ {
				if serverData[i] == '{' {
					isStarted = true
					continue
				}
				if serverData[i] == '}' {
					isStarted = false
					system.packageChan <- string(packageData)
					packageData = nil // 清空所有包的数据
					continue
				}

				if isStarted {
					packageData = append(packageData, serverData[i])
				}
			}
		}
	}
}

func handlePackgeData() {
	for {
		select {
		case <-system.CTX.Done():
			return
		default:
			packageData := <-system.packageChan
			args := strings.Split(packageData, "&;")
			n := len(args)
			for len(args) < 10 { // 填充一些元素防止爆炸
				args = append(args, "")
			}
			if args[0] == "msg" { // 有推送的新消息来
				name := []byte("lazy")
				msg := []byte("A error mesasage")

				name, _ = base64.StdEncoding.DecodeString(args[1])
				msg, _ = base64.StdEncoding.DecodeString(args[3])

				_, _ = fmt.Fprintln(system.messageBox, string(name), args[2])
				_, _ = fmt.Fprintln(system.messageBox, string(msg))
				_, _ = fmt.Fprint(system.messageBox, "\n")
			}
			if args[0] == "msghistory" { // 接收到历史记录
				for i := 1; i < n; i += 3 {
					if args[i] != "" {
						name := []byte("lazy")
						name, _ = base64.StdEncoding.DecodeString(args[i])
						_, _ = fmt.Fprint(system.messageBox, string(name)+" ")
					}
					if i+1 < n && args[i+1] != "" {
						_, _ = fmt.Fprint(system.messageBox, args[i+1]+"\n")
					}
					if i+2 < n && args[i+2] != "" {
						msg := []byte("A error message")
						msg, _ = base64.StdEncoding.DecodeString(args[i+2])
						_, _ = fmt.Fprintln(system.messageBox, string(msg))
						_, _ = fmt.Fprint(system.messageBox, "\n")
					}
				}
				system.messageBox.ScrollToEnd()
			}
		}
	}
}
