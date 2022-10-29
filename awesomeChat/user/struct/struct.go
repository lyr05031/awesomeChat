package conn

import (
	"encoding/json"
	"net"
	"os"
	"src/lyr"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/flopp/go-findfont"
)

type Conn struct {
	//81, 70, 215, 158
	IP            [4]byte
	PORT          int
	Sender        *net.UDPConn
	Err           error
	ChatHistory   []string
	MainLastInput string
	MainE         *widget.Entry
	MainL         *widget.List
	Application   fyne.App
	Window        fyne.Window
	Locker        bool
	User
}

type User struct {
	UID         int
	AddedGroups []Group
	UserName    string
	Zone        string
	LastCommand int
}

type Group struct {
	GID       string
	GroupName string
	Users     []int
}

type Msg struct {
	Command int
	UID     int
	Zone    string
	Text    string
}

const (
	LogIn = iota
	SignUp
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

func (newConn *Conn) Start() {
	newConn.Sender, newConn.Err = net.DialUDP("udp", nil, &net.UDPAddr{IP: net.IPv4(newConn.IP[0], newConn.IP[1], newConn.IP[2], newConn.IP[3]), Port: newConn.PORT})
	lyr.IfErr(newConn.Err)
	fontPath, err := findfont.Find("Arial Unicode.ttf")
	if err != nil {
		fontPath, _ = findfont.Find("simkai.ttf")
	}
	err = os.Setenv("FYNE_FONT", fontPath)
	lyr.IfErr(err)

	newConn.Zone = "public"
	newConn.Application = app.New()
	go newConn.Rev()
	newConn.LogInWindow()

	newMsg := Msg{Command: Delete, UID: newConn.UID, Zone: newConn.Zone, Text: ""}
	msg, _ := json.Marshal(&newMsg)
	newConn.Sender.Write(msg)
	newConn.Locker = false
}
