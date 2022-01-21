package bot

import (
	"MiraiGo-VclBot/util"
	"MiraiGo-VclBot/ws_data"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/Mrs4s/MiraiGo/utils"
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

const MAX_TEXT_LENGTH = 80

var (
	GroupMsgRecordMap  = make(map[int64][]MessageRecord)
	GroupMsgRecordLock sync.Mutex
)

// 预处理私聊消息，上传图片，MiraiGo更新后删除
func preProcessPrivateSendingMessage(cli *client.QQClient, target int64, m *message.SendingMessage) {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
	}()
	newElements := make([]message.IMessageElement, 0, len(m.Elements))
	for _, element := range m.Elements {
		if _, ok := element.(*PokeElement); ok {
			cli.SendFriendPoke(target)
			continue
		}
		if i, ok := element.(*LocalImageElement); ok {
			img, err := cli.UploadPrivateImage(target, i.Stream)
			if err != nil {
				log.Errorf("failed to upload private image, %+v", err)
				continue
			}
			if i.Tp == "flash" {
				img.Flash = true
			}
			newElements = append(newElements, img)
			continue
		}
		if i, ok := element.(*message.VoiceElement); ok {
			gm, err := cli.UploadPrivatePtt(target, bytes.NewReader(i.Data))
			if err != nil {
				log.Errorf("failed to upload private ptt, %+v", err)
				continue
			}
			newElements = append(newElements, gm)
			continue
		}
		if i, ok := element.(*MyVideoElement); ok {
			gm, err := cli.UploadGroupShortVideo(target, i.UploadingVideo, i.UploadingCover)
			if err != nil {
				log.Errorf("failed to upload private video, %+v", err)
				continue
			}
			newElements = append(newElements, gm)
			continue
		}
		newElements = append(newElements, element)
	}
	m.Elements = newElements
}

// 预处理群消息，上传图片/语音，MiraiGo更新后删除
func preProcessGroupSendingMessage(cli *client.QQClient, groupCode int64, m *message.SendingMessage) {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
	}()
	newElements := make([]message.IMessageElement, 0, len(m.Elements))
	for _, element := range m.Elements {
		if i, ok := element.(*message.TextElement); ok {
			for _, text := range utils.ChunkString(i.Content, MAX_TEXT_LENGTH) {
				if text != "" {
					newElements = append(newElements, message.NewText(text))
				}
			}
			continue
		}
		if i, ok := element.(*LocalImageElement); ok {
			img, err := cli.UploadGroupImage(groupCode, i.Stream)
			if err != nil {
				log.Errorf("failed to upload group image, %+v", err)
				continue
			}
			if i.Tp == "flash" {
				img.Flash = true
			} else if i.Tp == "show" {
				img.EffectID = i.EffectId
			}
			newElements = append(newElements, img)
			continue
		}
		if i, ok := element.(*message.VoiceElement); ok {
			gm, err := cli.UploadGroupPtt(groupCode, bytes.NewReader(i.Data))
			if err != nil {
				log.Errorf("failed to upload group ptt, %+v", err)
				continue
			}
			newElements = append(newElements, gm)
			continue
		}
		if i, ok := element.(*message.AtElement); ok && i.Target != 0 {
			i.Display = "@" + func() string {
				mem := cli.FindGroup(groupCode).FindMember(i.Target)
				if mem != nil {
					return mem.DisplayName()
				}
				return strconv.FormatInt(i.Target, 10)
			}()
			newElements = append(newElements, i)
			continue
		}
		if i, ok := element.(*PokeElement); ok {
			cli.SendGroupPoke(groupCode, i.Target)
			continue
		}
		if i, ok := element.(*MyVideoElement); ok {
			gm, err := cli.UploadGroupShortVideo(groupCode, i.UploadingVideo, i.UploadingCover)
			if err != nil {
				log.Errorf("failed to upload group video, %+v", err)
				continue
			}
			newElements = append(newElements, gm)
			continue
		}
		//if i, ok := element.(*GiftElement); ok {
		//	cli.SendGroupGift(uint64(groupCode), uint64(i.Target), i.GiftId)
		//	continue
		//}
		newElements = append(newElements, element)
	}
	m.Elements = newElements
}

//往QQ群发送消息
func BuHuangSendGroupMsg(cli *client.QQClient, miraiMsg []message.IMessageElement, msgId, groupId int64) int64 {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
	}()
	if cli.FindGroup(groupId) == nil {
		return 0
	}
	sendingMessage := &message.SendingMessage{Elements: miraiMsg}
	preProcessGroupSendingMessage(cli, groupId, sendingMessage)
	if len(sendingMessage.Elements) == 0 {
		return -1
	}
	var record MessageRecord
	nt := time.Now().Unix()
	ret := cli.SendGroupMessage(groupId, sendingMessage, false)
	nt2 := time.Now().Unix()
	if ret.Id == -1 {
		log.Info("发送失败使用时间:", nt2-nt)
	} else {
		log.Info("发送成功使用时间:", nt2-nt)
	}
	if ret.Id == -1 {
		ret = cli.SendGroupMessage(groupId, sendingMessage, false)
		log.Info(cli.Uin, "消息第二次重发返回id:", ret.Id)
	}
	if ret.Id == -1 {
		ret = cli.SendGroupMessage(groupId, sendingMessage, false)
		log.Info(cli.Uin, "消息第三次重发返回id:", ret.Id)
	}
	if ret.Id != -1 {
		GroupMsgRecordLock.Lock()
		defer func() {
			GroupMsgRecordLock.Unlock()
		}()
	}
	record.EventId = ret.Id
	record.InternalId = ret.InternalId
	record.GroupCode = ret.GroupCode
	recordList := GroupMsgRecordMap[msgId]
	if len(recordList) == 0 {
		var newRecordList []MessageRecord
		recordList = newRecordList
	}
	recordList = append(recordList, record)
	GroupMsgRecordMap[msgId] = recordList
	return int64(ret.Id)
}

//发送私聊、临时消息
func BuHuangSendPrivateMsg(cli *client.QQClient, miraiMsg []message.IMessageElement, userId, groupId int64) int64 {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
	}()
	sendingMessage := &message.SendingMessage{Elements: miraiMsg}
	if userId != 0 { // 私聊+临时
		preProcessPrivateSendingMessage(cli, userId, sendingMessage)
	}
	if groupId != 0 && userId != 0 { // 临时
		ret := cli.SendGroupTempMessage(groupId, userId, sendingMessage)
		return int64(ret.Id)
	}
	preProcessPrivateSendingMessage(cli, userId, sendingMessage)
	ret := cli.SendPrivateMessage(userId, sendingMessage)
	return int64(ret.Id)
}

//撤回群消息 -- 这里的msgId，如果internalId是0,那么是自己发送群消息的ID，否则就是接收群消息的时候event的id
func BuBuhuangWithDrawMsg(cli *client.QQClient, groupId, msgId int64, internalId int32) {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
	}()
	if msgId > 0 && internalId != 0 && groupId > 0 {
		_ = cli.RecallGroupMessage(groupId, int32(msgId), internalId)
	}

	GroupMsgRecordLock.Lock()
	defer func() {
		e := recover()
		if e != nil {
			fmt.Println(e)
		}
		GroupMsgRecordLock.Unlock()
	}()
	recordList := GroupMsgRecordMap[msgId]
	for _, v := range recordList {
		_ = cli.RecallGroupMessage(v.GroupCode, v.EventId, v.InternalId)
	}
}

//往websocket-server传递上线消息
func BuhuangBotOnline(botId int64) {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
	}()
	if WsCon != nil {
		var data ws_data.GMCWSData
		data.BotId = botId
		data.MsgType = ws_data.GMC_ONLINE
		//WSWLock.Lock()
		defer func() {
			e := recover()
			if e != nil {
				util.PrintStackTrace(e)
			}
			//WSWLock.Unlock()
		}()
		marshal, _ := json.Marshal(data)
		//WsCon.WriteMessage(websocket.TextMessage, marshal)
		WsCon.Write(marshal)
	}
}

//往websocket-server传递下线消息
func BuhuangBotOffline(botId int64) {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
	}()
	if WsCon != nil {
		var data ws_data.GMCWSData
		data.BotId = botId
		data.MsgType = ws_data.GMC_OFFLINE
		marshal, _ := json.Marshal(data)
		//WsCon.WriteMessage(websocket.TextMessage, marshal)
		WsCon.Write(marshal)
	}
}

//禁言
func BuhuangBanGroupMember(cli *client.QQClient, groupId, userId, banTime int64) {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
	}()
	if group := cli.FindGroup(groupId); group != nil {
		if member := group.FindMember(userId); member != nil {
			if err := member.Mute(uint32(banTime)); err != nil {
				return
			}
			return
		} else {
			log.Infof("禁言信息:%d在%d群里找不到%d", cli.Uin, group, userId)
		}
	} else {
		log.Infof("禁言信息:%d找不到QQ群:%d", cli.Uin, group)
	}
	return
}

//拒绝入群  ---- 不需要单独的Api，用channel + select 处理
func BuhuangRejectAddGroupRequest(cli *client.QQClient, groupId, userId int64) {
}

//踢出
func BuhuangKickGroupMember(cli *client.QQClient, groupId, userId int64) {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
	}()
	if group := cli.FindGroup(groupId); group != nil {
		if member := group.FindMember(userId); member != nil {
			if err := member.Kick("", false); err != nil {
				return
			}
			return
		} else {
			log.Infof("踢出信息:%d在%d群里找不到%d", cli.Uin, group, userId)
		}
	} else {
		log.Infof("踢出信息:%d找不到QQ群:%d", cli.Uin, group)
	}
	return
}

//获取QQ群列表
func BuhuangGetGroupList(cli *client.QQClient) []ws_data.GMCGroup {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
	}()
	var groupInfoList []ws_data.GMCGroup
	groupList := cli.GroupList
	for _, v := range groupList {
		var group ws_data.GMCGroup
		group.GroupId = v.Code
		group.GroupName = v.Name
		groupInfoList = append(groupInfoList, group)
	}
	return groupInfoList
}

//获取特定群成员
func BuhuangGetGroupMemberList(cli *client.QQClient, groupId int64) []ws_data.GMCMember {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
	}()
	var memberList []ws_data.GMCMember

	group := cli.FindGroup(groupId)
	if group == nil {
		return nil
	}
	members, err := cli.GetGroupMembers(group)
	if err != nil || len(members) == 0 {
		log.Infof("获取%d群群成员失败", groupId)
		return nil
	}
	for _, v := range members {
		var gmcMember ws_data.GMCMember
		gmcMember.QQ = v.Uin
		gmcMember.Level = v.Level
		gmcMember.Permission = int64(v.Permission)
		memberList = append(memberList, gmcMember)
	}
	return memberList
}

//获取全部群、成员
func BuhuangGetAllGroupListAndMemberList(cli *client.QQClient) ws_data.GMCAllGroupMember {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
	}()
	var allGroupMember ws_data.GMCAllGroupMember
	data := make(map[int64][]ws_data.GMCMember)
	groupList := cli.GroupList
	for _, v := range groupList {
		var gmcMemberList []ws_data.GMCMember
		for _, v2 := range v.Members {
			var member ws_data.GMCMember
			member.QQ = v2.Uin
			member.Level = v2.Level
			member.Permission = int64(v2.Permission)
			gmcMemberList = append(gmcMemberList, member)
		}
		if gmcMemberList != nil && len(gmcMemberList) > 0 {
			data[v.Code] = gmcMemberList
		}
	}
	allGroupMember.Data = data
	return allGroupMember
}

//上传群文件
func BuhuangUploadGroupFile(cli *client.QQClient, groupId int64, fileName, filePath string) {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
	}()

	//url := cli.GetGroupFileUrl(fromGroup, filePath, int32(busId))
	//downName := "C:\\data\\" + fileName
	//err := util.DownloadFile(url, downName, 0, nil)
	//if err == nil {
	system, _ := cli.GetGroupFileSystem(groupId)
	system.UploadFile(filePath, fileName, "/")
	//}
}
