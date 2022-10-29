package conn

import (
	"encoding/json"
	"src/lyr"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const ABOUT = "1.鸣谢\n\t1.fxa为了该程序在 Windows 平台的运行干翻了三节正课\n\t2.yjx贡献了设计哲学 Less is more\n3.所有回答我在各个平台上提问的老哥\n4.gqx, mcy, llw对于本程序在逻辑上的改进\n\n2.教程\n\t1.AC(Awesome Chat) 中有两种重要的数据 一种是 UID 即你的用户 ID, 一种是 GID 即该群组的ID\n\t2.AC 一切皆群 传统意义上的私聊就是两个人的群 你拥有创建和加入群的权力\n3.在“添加群组”功能中直接输入对方 UID 以创建你们两个的群聊, 别人可以通过“添加群组”输入 GID 加入该群\n\n3.其他\n\t1.新增了用户机制 群机制\n\n完"

func (conn *Conn) MainWindow() {
	conn.Window.Resize(fyne.NewSize(800, 600))

	T1 := widget.NewLabel("聊天区")
	conn.MainL = widget.NewList(
		func() int {
			return len(conn.ChatHistory)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("tmp")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(conn.ChatHistory[i])
		})

	tmpE1 := widget.NewEntry()
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			dialog.NewForm("键入对方UID 群名(空格隔开)", "确定", "取消", []*widget.FormItem{
				widget.NewFormItem("键入", tmpE1),
			}, func(b bool) {
				if b {
					conn.LastCommand = NewGroup
					newMsg := Msg{Command: NewGroup, UID: conn.UID, Text: tmpE1.Text, Zone: conn.Zone}
					msg, _ := json.Marshal(&newMsg)
					conn.Sender.Write(msg)
				}
			}, conn.Window).Show()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.SearchIcon(), func() {
			dialog.NewForm("键入GID", "确定", "取消", []*widget.FormItem{
				widget.NewFormItem("键入", tmpE1),
			}, func(b bool) {
				if b {
					conn.LastCommand = AddGroup
					newMsg := Msg{Command: AddGroup, UID: conn.UID, Text: tmpE1.Text, Zone: conn.Zone}
					msg, _ := json.Marshal(&newMsg)
					conn.Sender.Write(msg)
				}
			}, conn.Window).Show()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.AccountIcon(), func() {
			dialog.NewInformation("账户信息", "UID "+lyr.NumToString(conn.UID)+"\n用户名 "+conn.UserName, conn.Window).Show()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.InfoIcon(), func() {
			dialog.NewInformation("关于", ABOUT, conn.Window).Show()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ListIcon(), func() {
			var addedGroupsString string
			for _, x := range conn.AddedGroups {
				addedGroupsString += "组名: "
				addedGroupsString += x.GroupName
				addedGroupsString += "组员UID: "
				for _, y := range x.Users {
					addedGroupsString += lyr.NumToString(y)
					addedGroupsString += " "
				}
				addedGroupsString += "组ID(GID): "
				addedGroupsString += x.GID
				addedGroupsString += "\n"
			}
			dialog.NewInformation("已经加入的群聊", addedGroupsString, conn.Window).Show()
		}),
	)

	layerToolbar := container.NewBorder(toolbar, nil, nil, nil, widget.NewLabel("新建群组|添加群组|账户信息|关于 Awesome Chat (教程)|查看已经加入的群组"))
	conn.MainE = widget.NewEntry()
	conn.MainE.Validator = conn.MonitorEnter
	conn.MainE.SetPlaceHolder("键入信息 回车发送")
	layerL1 := container.NewGridWrap(fyne.Size{Width: 730, Height: 500}, conn.MainL)
	layerL1 = container.NewHBox(T1, layerL1)
	layer := container.NewVBox(layerL1, conn.MainE, layerToolbar)
	conn.Window.SetContent(layer)
	layer.Show()
}
