package conn

import (
	"encoding/json"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (conn *Conn) LogInWindow() {
	conn.Window = conn.Application.NewWindow("Awesome Chat")
	conn.Window.Resize(fyne.Size{Width: 425, Height: 185})
	logInL1 := widget.NewLabel("键入用户名与密码来登录 新用户请注册")
	logInE1 := widget.NewEntry()
	logInE1.PlaceHolder = "键入UID(注册的话输入用户名)"
	logInE2 := widget.NewEntry()
	logInE2.PlaceHolder = "键入密码"
	logInE2.Password = true
	logInB1 := widget.NewButton("登录", func() {
		newMsg := Msg{Command: LogIn, UID: conn.UID, Text: logInE1.Text + " " + logInE2.Text, Zone: conn.Zone}
		msg, _ := json.Marshal(&newMsg)
		conn.Sender.Write(msg)
		conn.LastCommand = LogIn
	})
	logInB2 := widget.NewButton("注册", func() {
		newMsg := Msg{Command: SignUp, UID: conn.UID, Text: logInE1.Text + " " + logInE2.Text, Zone: conn.Zone}
		msg, _ := json.Marshal(&newMsg)
		conn.Sender.Write(msg)
		conn.LastCommand = SignUp
	})
	buttonLayer := container.NewHBox(logInB1, logInB2)
	logInLayerAll := container.NewVBox(logInL1, logInE1, logInE2, buttonLayer)
	conn.Window.SetContent(logInLayerAll)
	logInLayerAll.Show()
	conn.Window.ShowAndRun()
}
