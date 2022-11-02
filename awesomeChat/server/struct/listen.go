/*
监听器
*/

package fetch

import (
	"encoding/json"
	"fmt"
	"net"
	"src/lyr"
	"strconv"
	"strings"
)

var num int
var tmpHeader *net.UDPAddr
var data [4096]byte
var jsonText *Msg
var jsonMsg []byte

func (server *Server) Listener() {
	/*
	   监听器本体
	*/

	for {
		num, tmpHeader, _ = server.Listen.ReadFromUDP(data[:])
		json.Unmarshal(data[:num], &jsonText)

		switch jsonText.Command {
		case LogIn:
			server.logIn()
		case SignIn:
			server.signUp()
		case Delete:
			server.delete()
		case NewGroup:
			server.newGroup()
		case AddGroup:
			server.addGroup()
		case Send:
			println("asdf")
			server.send()
		}
	}
}

// 处理消息
func (server *Server) send() {
	var newMsg Msg
	for _, x := range server.AllGroups[jsonText.Zone].Users {
		if server.AllUsers[x].State {
			newMsg = Msg{Command: Receive, UID: jsonText.UID, Text: "<" + jsonText.Zone + ">" + server.AllUsers[jsonText.UID].UserName + " :" + jsonText.Text}
			jsonMsg, _ = json.Marshal(newMsg)
			server.Listen.WriteToUDP(jsonMsg, server.AllUsers[x].Header)
		}
	}
}

// 用户登陆
func (server *Server) logIn() {
	// 获取发过来的账户密码
	uidAndPassword := strings.Split(jsonText.Text, " ")
	// 获取uid
	uid, _ := strconv.Atoi(uidAndPassword[0])
	// 无聊的判断
	if content, ok := server.AllUsers[uid]; ok {
		if content.Password == uidAndPassword[1] {
			if server.AllUsers[uid].State {
				jsonText = &Msg{Command: LogIns}
			} else {
				server.AllUsers[uid].State = true
				server.AllUsers[uid].Header = tmpHeader
				jsonText = &Msg{Command: Sucess, Text: lyr.NumToString(server.AllUsers[uid].UID) + " " + server.AllUsers[uid].UserName}
				jsonMsg, _ = json.Marshal(jsonText)
				server.Listen.WriteToUDP(jsonMsg, server.AllUsers[uid].Header)
				server.SendAddedGroup(uid)
				fmt.Println("连接用户")
				server.ShowInfo()
			}
		} else {
			jsonText = &Msg{Command: Fail}
		}
	} else {
		jsonText = &Msg{Command: Fail}
	}
	jsonMsg, _ = json.Marshal(jsonText)
	server.Listen.WriteToUDP(jsonMsg, tmpHeader)
}

func (server *Server) signUp() {
	/*
		注册新用户
	*/

	// 初始化新用户
	userNameAndPassword := strings.Split(jsonText.Text, " ")
	newUser := new(User)
	server.LastUID++
	newUser.State = true
	newUser.UID = server.LastUID
	newUser.Header = tmpHeader
	newUser.Zone = ""
	newUser.UserName = userNameAndPassword[0]
	newUser.Password = userNameAndPassword[1]
	// 写入本地文件
	server.UserFile.WriteString(lyr.NumToString(newUser.UID) + " " + newUser.UserName + " " + userNameAndPassword[1] + "\n")
	server.AllUsers[newUser.UID] = newUser
	jsonText = &Msg{Command: Sucess, Text: lyr.NumToString(newUser.UID) + " " + newUser.UserName}
	jsonMsg, _ = json.Marshal(jsonText)
	server.Listen.WriteToUDP(jsonMsg, tmpHeader)
	fmt.Println("连接用户")
	server.ShowInfo()
}

func (server *Server) delete() {
	/*
		用户登出
	*/

	// 将状态改为 false 表示不在线
	server.AllUsers[jsonText.UID].State = false
	fmt.Println("删除用户")
	server.ShowInfo()
}

func (server *Server) newGroup() {
	/*
		用户创建新群组
	*/

	var newMsg *Msg
	// 判断用户输入是否合法
	splitText := strings.Split(jsonText.Text, " ")
	if len(splitText) == 2 {
		tarUID, err := strconv.Atoi(splitText[0])
		if err == nil {
			if tarUserContent, ok := server.AllUsers[tarUID]; ok {
				// 读取 GID 的数字部分
				tmp, _ := strconv.Atoi(server.LastGID[1:])
				// 给 GID + 1
				server.LastGID = "g" + lyr.NumToString(tmp+1)
				// 读取组名
				groupName := splitText[1]
				newGroup := &Group{GID: server.LastGID, GroupName: groupName, Users: []int{jsonText.UID, tarUID}}
				// 赋值给各个部分
				server.AllGroups[server.LastGID] = newGroup
				tarUserContent.AddedGroups = append(tarUserContent.AddedGroups, *newGroup)
				server.AllUsers[jsonText.UID].AddedGroups = append(server.AllUsers[jsonText.UID].AddedGroups, *newGroup)
				server.GroupFile.WriteString(server.LastGID + " " + newGroup.GroupName + " " + (lyr.NumToString(jsonText.UID) + "|" + lyr.NumToString(tarUID)) + "\n")
				newMsg = &Msg{Command: Sucess, Text: groupName}
				server.SendAddedGroup(jsonText.UID)
				if server.AllUsers[tarUID].State {
					server.SendAddedGroup(tarUID)
				}
				fmt.Println("新群组")
				server.ShowInfo()
			} else {
				newMsg = &Msg{Command: Fail}
			}
		} else {
			newMsg = &Msg{Command: Fail}
		}
	} else {
		newMsg = &Msg{Command: Fail}
	}
	jsonMsg, _ := json.Marshal(&newMsg)
	server.Listen.WriteToUDP(jsonMsg, tmpHeader)
}

// 添加群组
func (server *Server) addGroup() {
	var newMsg Msg
	// 获取要添加的 GID 名字
	tarGID := jsonText.Text
	if _, ok := server.AllGroups[tarGID]; ok {
		for _, x := range server.AllGroups[tarGID].Users {
			// 避免重复添加
			if x == jsonText.UID {
				newMsg = Msg{Command: Fail}
				jsonMsg, _ = json.Marshal(newMsg)
				server.Listen.WriteToUDP(jsonMsg, tmpHeader)
				return
			}
		}
		// 添加到群组
		server.AllUsers[jsonText.UID].AddedGroups = append(server.AllUsers[jsonText.UID].AddedGroups, *server.AllGroups[tarGID])
		server.AllGroups[tarGID].Users = append(server.AllGroups[tarGID].Users, jsonText.UID)
		newMsg = Msg{Command: Sucess}

		// 写入本地文件
		server.GroupFile.WriteString("APPEND " + tarGID + " " + lyr.NumToString(jsonText.UID) + "\n")

		for _, x := range server.AllGroups[tarGID].Users {
			for num, y := range server.AllUsers[x].AddedGroups {
				if y.GID == tarGID {
					server.AllUsers[x].AddedGroups[num].Users = append(server.AllUsers[x].AddedGroups[num].Users, jsonText.UID)
				}
			}
		}

		for _, x := range server.AllGroups[tarGID].Users {
			// 广播给组中的在线用户
			if server.AllUsers[x].State {
				server.SendAddedGroup(x)
			}
		}

		server.ShowInfo()
	} else {
		newMsg = Msg{Command: Fail}
	}
	jsonMsg, _ = json.Marshal(newMsg)
	server.Listen.WriteToUDP(jsonMsg, tmpHeader)
}
