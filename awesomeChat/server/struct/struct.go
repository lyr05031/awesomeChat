package fetch

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Server struct {
	// 172.21.0.15
	IP        [4]byte
	PORT      int
	AllUsers  map[int]*User
	AllGroups map[string]*Group
	Listen    *net.UDPConn
	UserFile  *os.File
	GroupFile *os.File
	LastUID   int
	LastGID   string
}

const (
	LogIn = iota
	SignIn
	Send
	Delete
	Receive
	Sucess
	Fail
	LogIns
	NewGroup
	AddGroup
	LongData
)

type Msg struct {
	Command int
	UID     int
	Zone    string
	Text    string
}

type Group struct {
	GID       string
	GroupName string
	Users     []int
}

type User struct {
	State       bool
	Password    string
	UID         int
	Zone        string
	Header      *net.UDPAddr
	UserName    string
	AddedGroups []Group
	UnSendMsg   []Msg
}

var textByte [4096]byte
var textString []string
var err error

func (server *Server) LoadUser() {
	server.UserFile, err = os.OpenFile("user.txt", os.O_RDWR, 0)
	if err != nil {
		os.Create("user.txt")
		server.UserFile, _ = os.OpenFile("user.txt", os.O_RDWR, 0)
	}
	num, _ := server.UserFile.Read(textByte[:])
	textString = strings.Split(string(textByte[:num]), "\n")
	if len(textString) != 0 {
		for _, x := range textString {
			// tmp[0] -> uid, tmp[1] -> userName, tmp[2] -> password
			tmp := strings.Split(string(x), " ")
			uid, _ := strconv.Atoi(tmp[0])
			server.LastUID = uid
			server.AllUsers[uid] = &User{State: false, UID: uid, UserName: tmp[1], Password: tmp[2]}
		}
	}
}

// 加载群组
func (server *Server) LoadGroup() {
	// 打开文件
	server.GroupFile, err = os.OpenFile("group.txt", os.O_RDWR, 0)
	if err != nil {
		os.Create("group.txt")
		server.UserFile, _ = os.OpenFile("group.txt", os.O_RDWR, 0)
	}
	num, _ = server.GroupFile.Read(textByte[:])
	textString = strings.Split(string(textByte[:num]), "\n")
	if len(textString) != 0 {
		for _, x := range textString {
			tmp := strings.Split(string(x), " ")
			// 如果不是追加
			if tmp[0] != "APPEND" {
				server.LastGID = tmp[0]
				tmpIntArray := []int{}
				for _, y := range strings.Split(tmp[2], "|") {
					tmpInt, _ := strconv.Atoi(y)
					tmpIntArray = append(tmpIntArray, tmpInt)
				}
				for _, y := range strings.Split(tmp[2], "|") {
					tmpInt, _ := strconv.Atoi(y)
					server.AllUsers[tmpInt].AddedGroups = append(server.AllUsers[tmpInt].AddedGroups, Group{server.LastGID, tmp[1], tmpIntArray})
				}
				server.AllGroups[tmp[0]] = &Group{server.LastGID, tmp[1], tmpIntArray}
			} else {
				tmpInt, _ := strconv.Atoi(tmp[2])
				server.AllUsers[tmpInt].AddedGroups = append(server.AllUsers[tmpInt].AddedGroups, *server.AllGroups[tmp[1]])
				server.AllGroups[tmp[1]].Users = append(server.AllGroups[tmp[1]].Users, tmpInt)
				for num, x := range server.AllUsers {
					for num2, y := range x.AddedGroups {
						if y.GID == tmp[1] {
							server.AllUsers[num].AddedGroups[num2].Users = append(server.AllUsers[num].AddedGroups[num2].Users, tmpInt)
						}
					}
				}
			}
		}
	}
}

// 发送最新的群组信息
func (server *Server) SendAddedGroup(uid int) {
	// 需要发送的 Group 信息
	innerJsonMsg := server.AllUsers[uid].AddedGroups
	innerJsonText, _ := json.Marshal(innerJsonMsg)
	// 标准的信息
	jsonMsg := &Msg{Command: LongData, Text: string(innerJsonText)}
	jsonText, _ := json.Marshal(jsonMsg)
	server.Listen.WriteToUDP(jsonText, server.AllUsers[uid].Header)
}

// 打印用户和群组信息
func (server *Server) ShowInfo() {
	fmt.Println("用户\nUID 密码 用户名 加入的群 是否在线")
	for _, content := range server.AllUsers {
		fmt.Println(content.UID, content.Password, content.UserName, content.State, content.AddedGroups)
	}
	fmt.Println()
	fmt.Println("群组\nGID 群名 用户")
	for _, content := range server.AllGroups {
		fmt.Println(content.GID, content.GroupName, content.Users)
	}
	fmt.Println()
}

// RUN THE CODE
// GOD PLEASE
func (server *Server) Start() {
	// 初始化
	server.AllUsers = make(map[int]*User)
	server.AllGroups = make(map[string]*Group)

	// 加载用户和群组
	server.LoadUser()
	server.LoadGroup()

	fmt.Println("初始化")
	server.ShowInfo()

	// 开始监听
	server.Listen, _ = net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(server.IP[0], server.IP[1], server.IP[2], server.IP[3]), Port: server.PORT})
	server.Listener()
}
