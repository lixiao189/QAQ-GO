package main

import "encoding/base64"

func sendToServer(message string) error {
	message = base64.StdEncoding.EncodeToString([]byte(message))
	packageData := "{msg&;send&;" + message + "}"
	_, err := system.userConn.conn.Write([]byte(packageData))
	if err != nil {
		return err
	}
	return nil
}

func getHistoryMessage() {
	// 查询历史消息
	_, err := system.userConn.conn.Write([]byte("{msg&;list}"))
	if err != nil {
		system.Pages.SwitchToPage("connectError")
		_ = system.userConn.conn.Close()
		return
	}
}
