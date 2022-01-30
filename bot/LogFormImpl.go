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
	f.LogListView.SetParent(f)
	f.LogListView.SetViewStyle(types.VsReport)
	f.LogListView.SetOwnerData(true)
	f.LogListView.SetGridLines(true)
	f.LogListView.SetReadOnly(true)
	f.LogListView.SetRowSelect(true)
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
		if len(LogItems) == 0 {
			return
		}
		canvas := sender.Canvas()
		boundRect := item.DisplayRect(types.DrBounds)
		//当前状态，鼠标选中的那行显示的颜色
		if state.In(types.CdsFocused) {
			canvas.Brush().SetColor(colors.ClAqua)
		} else {
			canvas.Brush().SetColor(sender.Color())
		}

		canvas.FillRect(boundRect)
		data := LogItems[item.Index()]
		drawFlags := types.NewSet(types.TfCenter, types.TfSingleLine, types.TfVerticalCenter)
		var i int32
		font := canvas.Font()
		switch data.MessageType {
		case "接收":
			font.SetColor(colors.ClBlue)
		case "发送":
			font.SetColor(colors.ClGreen)
		}
		for i = 0; i < sender.Columns().Count(); i++ {
			r := f.GetSubItemRect(sender.Handle(), item.Index(), i)
			switch i {
			case 0:
				canvas.TextRect2(&r, data.BotId, drawFlags)
			case 1:
				canvas.TextRect2(&r, data.MessageType, drawFlags)
			case 2:
				canvas.TextRect2(&r, data.MessageTime, drawFlags)
			case 3:
				canvas.TextRect2(&r, data.Message, drawFlags)
			}
		}
	})

	f.ClearText = vcl.NewLabel(f)
	f.ClearText.SetParent(f)
	f.ClearText.SetTop(f.LogListView.Top() + f.LogListView.Height() + 5)
	f.ClearText.SetLeft(20)
	f.ClearText.SetCaption("清除日志")
	f.ClearText.SetOnClick(func(sender vcl.IObject) {
		LogItems = LogItems[:0]
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

func AddLogItem(botId, groupId, userId int64, acceptOrSend, messageType int32, message string) {
	var logItem LogItem
	switch acceptOrSend {
	case ACCEPT:
		logItem.MessageType = "接收"
	case SEND:
		logItem.MessageType = "发送"
	}
	logItem.BotId = strconv.FormatInt(botId, 10)
	logItem.MessageTime = time.Now().Format("15:04:05")

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
	logItem.Message = logMessage
	LogItems = append(LogItems, logItem)
	itemIndex := int32(len(LogItems))
	LogForm.LogListView.Items().SetCount(itemIndex)
}
