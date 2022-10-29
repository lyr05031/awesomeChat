package conn

import (
	"encoding/json"
	"src/lyr"
	"strconv"
	"strings"

	"fyne.io/fyne/v2/dialog"
)

// 监听 enter
func (conn *Conn) MonitorEnter(input string) error {
	if conn.MainLastInput == input && conn.MainLastInput != "nil" {
		conn.Send()
		conn.MainLastInput = "nil"
	}
	conn.MainLastInput = input
	return nil
}

// 发送信息
func (conn *Conn) Send() {
	if len(conn.MainE.Text) > 5 {
		if conn.MainE.Text[:5] == "/join" {
			conn.Zone = conn.MainE.Text[6:]
			conn.MainE.Text = ""
			conn.MainE.Refresh()
			conn.ChatHistory = append(conn.ChatHistory, "切换群组->"+conn.Zone)
			conn.MainL.Refresh()
			return
		}
	}
	if conn.Zone != "public" && conn.MainE.Text != "" {
		if len(conn.MainE.Text) >= 4095 {
			conn.ChatHistory = append(conn.ChatHistory, "信息太长 请分开发送")
		}
		newMsg := Msg{Command: Send, UID: conn.UID, Text: conn.MainE.Text, Zone: conn.Zone}
		msg, _ := json.Marshal(&newMsg)
		_, conn.Err = conn.Sender.Write(msg)
		lyr.IfErr(conn.Err)
		conn.MainE.Text = ""
		conn.MainE.Refresh()
	}
}

// 控制 GUI 生死
// 监听消息
func (conn *Conn) Rev() {
	var RevDataByte [4096]byte
	var RevDataUn Msg
	var RevDataGroupInfo []Group
	for {
		// 判断程序是否还在运行
		if !conn.Locker {
			conn.Window.Close()
			conn.Application.Quit()
			return
		}
		// 拆分收到的消息
		n, _, _ := conn.Sender.ReadFromUDP(RevDataByte[:])
		json.Unmarshal(RevDataByte[:n], &RevDataUn)

		// 判断语句
		switch RevDataUn.Command {
		case Receive:
			conn.ChatHistory = append(conn.ChatHistory, string(RevDataUn.Text))
			conn.MainL.Refresh()
			conn.MainL.ScrollToBottom()
		case Sucess:
			switch conn.LastCommand {
			case LogIn:
				fallthrough
			case SignUp:
				conn.UID, _ = strconv.Atoi(strings.Split(RevDataUn.Text, " ")[0])
				conn.UserName = strings.Split(RevDataUn.Text, " ")[1]
				conn.MainWindow()
			case NewGroup:
				dialog.NewInformation("成功", "成功创建群组!", conn.Window).Show()
			}
		case Fail:
			switch conn.LastCommand {
			case LogIn:
				fallthrough
			case SignUp:
				dialog.NewInformation("失败", "UID或密码错误", conn.Window).Show()
			case NewGroup:
				dialog.NewInformation("失败", "请检查输入的信息是否正确", conn.Window).Show()
			case AddGroup:
				dialog.NewInformation("失败", "请检查输入的信息是否正确或是否已经在该组中", conn.Window).Show()
			}
		case LogIns:
			dialog.NewInformation("失败", "AC 暂不支持多终端登录同一用户", conn.Window).Show()
		case LongData:
			json.Unmarshal([]byte(RevDataUn.Text), &RevDataGroupInfo)
			conn.AddedGroups = RevDataGroupInfo
		}
	}
}
