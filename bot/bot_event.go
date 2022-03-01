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
	_ = SendWSMsg(data)
}

//群消息
func handleGroupMessage(cli *client.QQClient, event *message.GroupMessage) {
	go func() {
		msg := MiraiMsgToRawMsg2(cli, event.GroupCode, event.Elements)
		go func() {
			AddLogItem(cli.Uin, event.GroupCode, event.Sender.Uin, ACCEPT, ACCEPT_GROUP, msg)
		}()
		log.Info("收到群聊消息")
		if WsCon == nil {
			log.Infof("未找到server")
			return
		}
		go func() {
			intn := rand.Intn(100)
			if intn > 80 {
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
		_ = SendWSMsg(data)
	}()
}

//临时消息
func handleTempMessage(cli *client.QQClient, event *client.TempMessageEvent) {
	msg := MiraiMsgToRawMsg(cli, event.Message.Elements)
	go func() {
		AddLogItem(cli.Uin, event.Message.GroupCode, event.Message.Sender.Uin, ACCEPT, ACCEPT_TEMP, msg)
	}()
	if WsCon == nil {
		return
	}
	var data ws_data.GMCWSData
	data.BotId = cli.Uin
	data.GroupId = event.Message.GroupCode
	data.UserId = event.Message.Sender.Uin
	data.MsgType = ws_data.GMC_TEMP_MESSAGE
	data.Message = msg
	_ = SendWSMsg(data)
}

//有人入群
func handleMemberJoinGroup(cli *client.QQClient, event *client.MemberJoinGroupEvent) {
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
	_ = SendWSMsg(data)
}

//有人离开
func handleMemberLeaveGroup(cli *client.QQClient, event *client.MemberLeaveGroupEvent) {
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
	_ = SendWSMsg(data)
}

//入群申请
func handleUserJoinGroupRequest(cli *client.QQClient, event *client.UserJoinGroupRequest) {
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
	_ = SendWSMsg(data)
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
	_ = SendWSMsg(data)
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
		AddLogItem(botId, groupID, 0, SEND, SEND_GROUP, "风控消息:"+msg)
	}()
	var data ws_data.GMCWSData
	data.MsgType = ws_data.ERROR_MSG
	data.BotId = botId
	data.GroupId = groupID
	data.MessageId = msgId
	data.Message = msg
	_ = SendWSMsg(data)
}

func SendWSMsg(data ws_data.GMCWSData) error {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
	}()
	marshal, _ := json.Marshal(data)
	_, err := WsCon.Write(marshal)
	return err
}
