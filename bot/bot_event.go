package bot

import (
	"MiraiGo-VclBot/util"
	"MiraiGo-VclBot/ws_data"
	"encoding/json"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

const (
	MessageIgnore = 0
	MessageBlock  = 1
)

var (
	REQUEST_ACCEPT int64 = 1
	REQUEST_REJECT int64 = -1
)

//注册事件
func Serve(cli *client.QQClient) {
	cli.OnPrivateMessage(handlePrivateMessage)
	cli.OnGroupMessage(handleGroupMessage)
	cli.OnTempMessage(handleTempMessage)
	//群成员++  (有人入群)
	cli.OnGroupMemberJoined(handleMemberJoinGroup)
	//群成员--  (有人跑路)
	cli.OnGroupMemberLeaved(handleMemberLeaveGroup)
	//入群请求
	cli.OnUserWantJoinGroup(handleUserJoinGroupRequest)
	//机器人被邀请入群
	cli.OnGroupInvited(handleGroupInvitedRequest)

	//cli.OnJoinGroup(handleJoinGroup)  //机器人入群
	//cli.OnLeaveGroup(handleLeaveGroup) //机器人退群
	//cli.OnNewFriendRequest(handleNewFriendRequest)
	//cli.OnGroupMessageRecalled(handleGroupMessageRecalled)
	//cli.OnFriendMessageRecalled(handleFriendMessageRecalled)
	//cli.OnNewFriendAdded(handleNewFriendAdded)
	//cli.OnReceivedOfflineFile(handleOfflineFile)
	//cli.OnGroupMuted(handleGroupMute)
	//cli.OnGroupMemberPermissionChanged(handleMemberPermissionChanged)
}

//私聊消息
func handlePrivateMessage(cli *client.QQClient, event *message.PrivateMessage) {
	//bot.WSWLock.Lock()
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
		//bot.WSWLock.Unlock()
	}()
	msg := MiraiMsgToRawMsg(cli, event.Elements)
	go func() {
		AddLogItem(cli.Uin, 0, event.Sender.Uin, ACCEPT, ACCEPT_PRIVATE, msg)
	}()
	log.Info("收到", event.Sender.Uin, "私聊消息:", msg)
	if WsCon == nil {
		return
	}
	cli.MarkPrivateMessageReaded(event.Sender.Uin, int64(event.Time))
	var data ws_data.GMCWSData
	data.BotId = cli.Uin
	data.UserId = event.Sender.Uin
	data.MsgType = ws_data.GMC_PRIVATE_MESSAGE
	data.Message = msg
	marshal, _ := json.Marshal(data)
	//err := bot.WsCon.WriteMessage(websocket.TextMessage, marshal)
	WsCon.Write(marshal)
	//if err != nil {
	//	log.Info("handlePrivateMessage出错", err)
	//}
}

//群消息
func handleGroupMessage(cli *client.QQClient, event *message.GroupMessage) {
	go func() {
		defer func() {
			e := recover()
			if e != nil {
				util.PrintStackTrace(e)
			}
			//bot.WSWLock.Unlock()
		}()
		//bot.WSWLock.Lock()
		msg := MiraiMsgToRawMsg2(cli, event.GroupCode, event.Elements)
		go func() {
			AddLogItem(cli.Uin, event.GroupCode, event.Sender.Uin, ACCEPT, ACCEPT_GROUP, msg)
		}()
		log.Info("收到群聊消息:", msg)
		if WsCon == nil {
			log.Infof("未找到server")
			return
		}
		go func() {
			intn := rand.Intn(100)
			if intn > 90 {
				cli.MarkGroupMessageReaded(event.GroupCode, int64(event.Id))
			}
		}()
		var data ws_data.GMCWSData
		data.BotId = cli.Uin
		data.GroupId = event.GroupCode
		data.UserId = event.Sender.Uin
		data.MsgType = ws_data.GMC_GROUP_MESSAGE
		data.MessageId = int64(event.Id)
		data.InternalId = event.InternalId
		data.Message = msg
		marshal, _ := json.Marshal(data)
		//e := bot.WsCon.WriteMessage(websocket.TextMessage, marshal)
		//if e != nil {
		//log.Info("handleGroupMessage错误:", e)
		//}
		WsCon.Write(marshal)
	}()
}

//临时消息
func handleTempMessage(cli *client.QQClient, event *client.TempMessageEvent) {
	//bot.WSWLock.Lock()
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
		//bot.WSWLock.Unlock()
	}()
	//cli.MarkPrivateMessageReaded(event.Message.Sender.Uin, time.Now().Unix())
	msg := MiraiMsgToRawMsg(cli, event.Message.Elements)
	go func() {
		AddLogItem(cli.Uin, event.Message.GroupCode, event.Message.Sender.Uin, ACCEPT, ACCEPT_TEMP, msg)
	}()
	//log.Info("收到", event.Message.GroupCode, "群,", event.Message.Sender.Uin, "的临时私聊消息:", msg)
	if WsCon == nil {
		return
	}
	var data ws_data.GMCWSData
	data.BotId = cli.Uin
	data.GroupId = event.Message.GroupCode
	data.UserId = event.Message.Sender.Uin
	data.MsgType = ws_data.GMC_TEMP_MESSAGE
	data.Message = msg
	marshal, _ := json.Marshal(data)
	//err := bot.WsCon.WriteMessage(websocket.TextMessage, marshal)
	//if err != nil {
	//log.Info("handleTempMessage出错:", err)
	//}
	WsCon.Write(marshal)
}

//有人入群
func handleMemberJoinGroup(cli *client.QQClient, event *client.MemberJoinGroupEvent) {
	//bot.WSWLock.Lock()
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
		//bot.WSWLock.Unlock()
	}()
	log.Info("收到入群信息")
	go func() {
		AddLogItem(cli.Uin, event.Group.Code, event.Member.Uin, ACCEPT, ACCEPT_MEMBER_INSERT, "")
	}()
	if WsCon == nil {
		return
	}
	var data ws_data.GMCWSData
	data.GroupId = event.Group.Code
	data.UserId = event.Member.Uin
	data.NickName = event.Member.Nickname
	data.BotId = cli.Uin
	data.MsgType = ws_data.GMC_MEMBER_ADD
	marshal, _ := json.Marshal(data)
	//err := bot.WsCon.WriteMessage(websocket.TextMessage, marshal)
	//if err != nil {
	//log.Info("handleMemberJoinGroup出错", err)
	//}
	WsCon.Write(marshal)
}

//有人离开
func handleMemberLeaveGroup(cli *client.QQClient, event *client.MemberLeaveGroupEvent) {
	//bot.WSWLock.Lock()
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
		//bot.WSWLock.Unlock()
	}()
	log.Info("收到有人退群")
	go func() {
		AddLogItem(cli.Uin, event.Group.Code, event.Member.Uin, ACCEPT, ACCEPT_MEMBER_DECREASE, "")
	}()
	if WsCon == nil {
		return
	}
	var data ws_data.GMCWSData
	data.GroupId = event.Group.Code
	data.UserId = event.Member.Uin
	data.BotId = cli.Uin
	data.MsgType = ws_data.GMC_MEMBER_LEAVE
	marshal, _ := json.Marshal(data)
	//err := bot.WsCon.WriteMessage(websocket.TextMessage, marshal)
	//if err != nil {
	//	log.Info("handleMemberLeaveGroup出错", err)
	//}
	WsCon.Write(marshal)
}

//入群申请
func handleUserJoinGroupRequest(cli *client.QQClient, event *client.UserJoinGroupRequest) {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
	}()
	log.Info("有人申请入群")
	go func() {
		AddLogItem(cli.Uin, event.GroupCode, event.RequesterUin, ACCEPT, ACCEPT_GROUP_REQUEST, "")
	}()
	if WsCon == nil {
		return
	}
	var data ws_data.GMCWSData
	data.MsgType = ws_data.GMC_GROUP_REQUEST
	data.BotId = cli.Uin
	data.GroupId = event.GroupCode
	data.UserId = event.RequesterUin
	data.NickName = event.RequesterNick
	data.RequestId = event.RequestId
	data.Message = event.Message
	marshal, _ := json.Marshal(data)
	//bot.WSWLock.Lock()
	WsCon.Write(marshal)
	//err := bot.WsCon.WriteMessage(websocket.TextMessage, marshal)
	//if err != nil {
	//	log.Info("handleUserJoinGroupRequest出错", err)
	//}
	//bot.WSWLock.Unlock()
	ch := make(chan ws_data.GMCWSData, 1)
	ws_data.ChanMapLock.Lock()
	ws_data.ChanMap[event.RequestId] = ch
	ws_data.ChanMapLock.Unlock()
	go func() {
		select {
		case r := <-ch:
			switch r.GroupRequest {
			case REQUEST_ACCEPT:
				event.Accept()
			case REQUEST_REJECT:
				event.Reject(false, "")
			}
		case <-time.After(time.Second * 30):
			delete(ws_data.ChanMap, event.RequestId)
			return
		}
	}()
}

//机器人被邀请入群
func handleGroupInvitedRequest(cli *client.QQClient, event *client.GroupInvitedRequest) {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
	}()
	log.Info("收到机器人被邀请入群")
	go func() {
		AddLogItem(cli.Uin, event.GroupCode, event.InvitorUin, ACCEPT, ACCPET_GROUP_INVITED, "")
	}()
	if WsCon == nil {
		return
	}
	var data ws_data.GMCWSData
	data.MsgType = ws_data.GMC_BOT_INVITED
	data.BotId = cli.Uin
	data.GroupId = event.GroupCode
	data.NickName = event.InvitorNick
	data.InvitorId = event.InvitorUin
	data.RequestId = event.RequestId
	marshal, _ := json.Marshal(data)
	//bot.WSWLock.Lock()
	WsCon.Write(marshal)
	//err := bot.WsCon.WriteMessage(websocket.TextMessage, marshal)
	//if err != nil {
	//	log.Info("handleGroupInvitedRequest出错", err)
	//}
	//bot.WSWLock.Unlock()
	ch := make(chan ws_data.GMCWSData, 1)
	ws_data.ChanMapLock.Lock()
	ws_data.ChanMap[event.RequestId] = ch
	ws_data.ChanMapLock.Unlock()
	go func() {
		select {
		case r := <-ch:
			switch r.GroupRequest {
			case REQUEST_ACCEPT:
				cli.SolveGroupJoinRequest(event, true, false, "")
			case REQUEST_REJECT:
				cli.SolveGroupJoinRequest(event, false, false, "")
			}
			return
		case <-time.After(time.Second * 30):
			delete(ws_data.ChanMap, event.RequestId)
			return
		}
	}()
}

//发送失败的消息返回给server
func handleErrorMsg(botId int64, groupID, msgId int64, msg string) {
	go func() {
		defer func() {
			e := recover()
			if e != nil {
				util.PrintStackTrace(e)
			}
			//bot.WSWLock.Unlock()
		}()
		AddLogItem(botId, groupID, 0, SEND, SEND_GROUP, "风控消息:"+msg)
	}()
	var data ws_data.GMCWSData
	data.MsgType = ws_data.ERROR_MSG
	data.BotId = botId
	data.GroupId = groupID
	data.MessageId = msgId
	data.Message = msg
	marshal, _ := json.Marshal(data)
	WsCon.Write(marshal)
}
