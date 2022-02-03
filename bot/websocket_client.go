package bot

import (
	"MiraiGo-VclBot/util"
	"MiraiGo-VclBot/ws_data"
	"encoding/json"
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	log "github.com/sirupsen/logrus"
	//"github.com/gorilla/websocket"
	"golang.org/x/net/websocket"
	"sync"
	"time"
)

var (
	WsCon       *websocket.Conn
	WsConSucess bool = false
	WSCallLock  sync.Mutex
	BotLock     sync.Mutex
)

const (
	WSServerAddr   = "ws://127.0.0.1:9801/gmc_event"
	WSClientOrigin = "http://127.0.0.1:9801"
)

func WSDailCall() {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
			go func() {
				log.Info("WsConSucess:", WsConSucess)
				WSDailCall()
			}()
		}
		WSCallLock.Unlock()
	}()
	WSCallLock.Lock()
	//var tempHeader http.Header = make(map[string][]string)
	//tempHeader.Add("origin", WSClientOrigin)
	//WSClientHeader = tempHeader
	var err error
	for {
		if WsCon != nil && WsConSucess {
			return
		}
		fmt.Println("开始连接")
		//WsCon, _, err = websocket.DefaultDialer.Dial(WSServerAddr, WSClientHeader)
		WsCon, err = websocket.Dial(WSServerAddr, "", WSClientOrigin)
		if err != nil || WsCon == nil {
			log.Infof("ws连接出错:", err)
			time.Sleep(time.Second * 2)
			continue
		} else {
			WsConSucess = true
			time.Sleep(time.Second)
			Clients.Range(func(_ int64, cli *client.QQClient) bool {
				if cli.Online.Load() {
					fmt.Println(cli.Uin, "发送上线事件")
					BuhuangBotOnline(cli.Uin)
				}
				return true
			})
			return
		}
	}
}

//处理Websocket-Server的消息，一般负责调用API
func HandleWSMsg() {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
		go func() {
			HandleWSMsg()
		}()
	}()
	for {
		if WsCon == nil && !WsConSucess {
			time.Sleep(time.Second)
			continue
		}
		//WSRLock.Lock()
		//_, message, e := WsCon.ReadMessage()
		//_, message, e := WsCon.Read()
		//WSRLock.Unlock()
		request := make([]byte, 2048)

		readLen, e := WsCon.Read(request)
		if e != nil || WsCon == nil {
			log.Println("出错了：", e)
			time.Sleep(time.Second * 2)
			WsConSucess = false
			go func() {
				log.Println("ws-server掉线，正在重连")
				WSDailCall()
			}()
			continue
		}
		go func() {
			var data ws_data.GMCWSData
			_ = json.Unmarshal(request[:readLen], &data)
			BotLock.Lock()
			cli, ok := Clients.Load(data.BotId)
			BotLock.Unlock()
			if !ok {
				log.Info("加载QQ Cli不存在:", data.BotId)
				return
			}
			miraiMsg := RawMsgToMiraiMsg(cli, data.Message)
			switch data.MsgType {
			case ws_data.GMC_PRIVATE_MESSAGE, ws_data.GMC_TEMP_MESSAGE:
				var logMessageType int32
				if data.GroupId == 0 {
					logMessageType = SEND_PRIVATE
				} else {
					logMessageType = SEND_TEMP
				}
				AddLogItem(data.BotId, data.GroupId, data.UserId, SEND, logMessageType, data.Message)
				BuHuangSendPrivateMsg(cli, miraiMsg, data.UserId, data.GroupId)
			case ws_data.GMC_GROUP_MESSAGE:
				AddLogItem(data.BotId, data.GroupId, data.UserId, SEND, SEND_GROUP, data.Message)
				BuHuangSendGroupMsg(cli, miraiMsg, data.MessageId, data.GroupId)
			case ws_data.GMC_WITHDRAW_MESSAGE:
				BuBuhuangWithDrawMsg(cli, data.GroupId, data.MessageId, data.InternalId)
			case ws_data.GMC_ALLGROUPMEMBER:
				HandleGetAllMember(cli)
			case ws_data.GMC_GROUP_LIST:
				HandleGroupList(cli)
			case ws_data.GMC_KICK:
				BuhuangKickGroupMember(cli, data.GroupId, data.UserId)
			case ws_data.GMC_BAN:
				BuhuangBanGroupMember(cli, data.GroupId, data.UserId, data.Time)
			case ws_data.GMC_GROUP_FILE:
				BuhuangUploadGroupFile(cli, data.GroupId, data.Message, data.FilePath)
			case ws_data.GMC_GROUP_REQUEST, ws_data.GMC_BOT_INVITED:
				ws_data.HandleCallBackEvent(data)
			}
		}()
	}
}
