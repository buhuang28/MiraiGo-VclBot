package bot

import (
	"fmt"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"github.com/ying32/govcl/vcl/types/colors"
	"github.com/ying32/govcl/vcl/win"
	"strconv"
	"time"
)

func (f *TLogForm) OnFormCreate(sender vcl.IObject) {
	f.SetCaption("日志")
	f.SetHeight(400)
	f.SetWidth(600)
	f.ScreenCenter()
	f.EnabledMaximize(false)
	f.SetBorderStyle(types.BsSingle)
	f.SetShowInTaskBar(types.StAlways)
	f.SetDoubleBuffered(true)
	f.LogListView = vcl.NewListView(f)
	f.LogListView.SetParent(f)
	f.LogListView.SetWidth(600)
	f.LogListView.SetHeight(370)
	f.LogListView.SetLeft(0)
	f.LogListView.SetTop(0)
	f.LogListView.SetRowSelect(true)
	f.LogListView.SetReadOnly(true)
	f.LogListView.SetViewStyle(types.VsReport)
	f.LogListView.SetGridLines(true)

	addCol := func(name string, width int32) {
		col := f.LogListView.Columns().Add()
		col.SetCaption(name)
		col.SetWidth(width)
		col.SetAlignment(types.TaCenter)
	}
	addCol("QQ", 100)
	addCol("类型", 50)
	addCol("时间", 120)
	addCol("消息内容", 380)

	f.LogListView.SetOnAdvancedCustomDrawItem(func(sender *vcl.TListView, item *vcl.TListItem, state types.TCustomDrawState, Stage types.TCustomDrawStage, defaultDraw *bool) {
		canvas := sender.Canvas()
		font := canvas.Font()
		if item.SubItems().Strings(0) == "发送" {
			font.SetColor(colors.ClGreen)
		} else {
			font.SetColor(colors.ClBlue)
		}
	})

	f.ClearText = vcl.NewLabel(f)
	f.ClearText.SetParent(f)
	f.ClearText.SetTop(f.LogListView.Top() + f.LogListView.Height() + 5)
	f.ClearText.SetLeft(20)
	f.ClearText.SetCaption("清除日志")
	f.ClearText.SetOnClick(func(sender vcl.IObject) {
		LogForm.LogListView.Clear()
	})

	f.RollCheck = vcl.NewCheckBox(f)
	f.RollCheck.SetParent(f)
	f.RollCheck.SetTop(f.ClearText.Top() - 3)
	f.RollCheck.SetLeft(f.ClearText.Width() + f.ClearText.Left() + 50)
	f.RollCheck.SetCaption("日志滚动")
}

func (f *TLogForm) GetSubItemRect(hwndLV types.HWND, iItem, iSubItem int32) (ret types.TRect) {
	win.ListView_GetSubItemRect(hwndLV, iItem, iSubItem, win.LVIR_LABEL, &ret)
	return
}

func AddLogItem(botId, groupId, userId int64, acceptOrSend string, messageType int32, message string) {
	botIdStr := strconv.FormatInt(botId, 10)
	messageTimeStr := time.Now().Format("15:04:05")
	logMessage := ""
	switch messageType {
	case ACCEPT_PRIVATE:
		logMessage = fmt.Sprintf("QQ(%v):%v", groupId, userId, message)
	case ACCEPT_GROUP:
		logMessage = fmt.Sprintf("群(%v) QQ(%v):%v", groupId, userId, message)
	case ACCEPT_TEMP:
		logMessage = fmt.Sprintf("群(%v) QQ(%v)临时会话:%v", groupId, userId, message)
	case ACCEPT_GROUP_REQUEST:
		logMessage = fmt.Sprintf("QQ(%v)申请加群(%v)", userId, groupId)
	case ACCPET_GROUP_INVITED:
		logMessage = fmt.Sprintf("机器人被QQ(%v)邀请加入群(%v)", userId, groupId)
	case ACCEPT_KICK:
		logMessage = fmt.Sprintf("QQ(%v)被踢出群(%v)", userId, groupId)
	case ACCEPT_MEMBER_INSERT:
		logMessage = fmt.Sprintf("QQ(%v)加入群(%v)", userId, groupId)
	case ACCEPT_MEMBER_DECREASE:
		logMessage = fmt.Sprintf("QQ(%v)离开群(%v)", userId, groupId)
	case SEND_PRIVATE:
		logMessage = fmt.Sprintf("QQ(%v):", userId, groupId)
	case SEND_GROUP:
		logMessage = fmt.Sprintf("群(%v):%v", groupId, message)
	case SEND_TEMP:
		logMessage = fmt.Sprintf("群(%v) QQ(%v)临时会话:%v", groupId, userId, message)
	}
	vcl.ThreadSync(func() {
		LogForm.LogListView.Items().BeginUpdate()
		item := LogForm.LogListView.Items().Add()
		item.SetCaption(botIdStr)
		subItem := item.SubItems()
		subItem.Add(acceptOrSend)
		subItem.Add(messageTimeStr)
		subItem.Add(logMessage)
		LogForm.LogListView.Items().EndUpdate()
		if LogForm.RollCheck.Checked() {
			item.MakeVisible(true)
		}
	})
}
