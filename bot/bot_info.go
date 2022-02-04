package bot

import (
	"MiraiGo-VclBot/device"
	"MiraiGo-VclBot/util"
	"encoding/json"
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

var (
	QQINFOROOTPATH = "C:\\data"
	QQINFOPATH     = "C:\\data\\"
	QQINFOSKIN     = ".info"
)

type QQInfo struct {
	QQ       int64    `json:"qq"`
	PassWord [16]byte `json:"pass_word"`
	Token    []byte   `json:"token"`
	//对应的随机种子文件
	SeedFile       string `json:"seed_file"`
	ClientProtocol int32  `json:"client_protocol"`
	AutoLogin      bool   `json:"auto_login"`
}

func (q *QQInfo) StoreLoginInfo(qq int64, pw [16]byte, token []byte, clientProtocol int32, autoLogin bool) bool {
	q.QQ = qq
	if q.QQ == 0 {
		return false
	}
	q.PassWord = pw
	q.Token = token
	q.AutoLogin = autoLogin
	q.ClientProtocol = clientProtocol
	fileName := QQINFOPATH + strconv.FormatInt(q.QQ, 10) + QQINFOSKIN
	marshal, err := json.Marshal(q)
	if err != nil {
		return false
	}
	return util.WriteFile(fileName, marshal)
}

func (q *QQInfo) Login() bool {
	log.Info("开始登录", q.QQ)

	if q.QQ != 0 && q.PassWord != [16]byte{} {
		fmt.Println("开始密码登录")
		success := CreateBotImplMd5(q.QQ, q.PassWord, q.QQ, q.ClientProtocol, q.AutoLogin)
		if !success {
			time.Sleep(time.Second * 3)
			success = CreateBotImplMd5(q.QQ, q.PassWord, q.QQ, q.ClientProtocol, q.AutoLogin)
		}
		if success {
			return success
		}
	}

	if q.Token != nil {
		var botClient = client.NewClientEmpty()
		deviceInfo := device.GetDevice(q.QQ, q.ClientProtocol)
		botClient.UseDevice(deviceInfo)
		err := botClient.TokenLogin(q.Token)
		fmt.Println("开始使用token登录:", err)
		index := GetBotIndex(q.QQ)
		var botData TTempItem
		botData.IconIndex = index
		SetBotAvatar(q.QQ, index)
		botData.QQ = strconv.FormatInt(q.QQ, 10)
		botData.Protocol = GetProtocol(q.ClientProtocol)
		botData.Status = "登录中"
		botData.NickName = ""
		if q.AutoLogin {
			botData.Auto = "√"
		} else {
			botData.Auto = "X"
		}
		botData.Note = "登录中"
		if index == -1 {
			AddTempBotData(botData)
		}
		BotForm.BotListView.Items().SetCount(int32(len(TempBotData))) //   必须主动的设置Virtual List的行数
		if err == nil {
			TempBotData[index].NickName = botClient.Nickname
			TempBotData[index].Status = "在线"
			TempBotData[index].Note = "在线"
			InitLog(botClient)
			log.Infof("初始化日志")
			Clients.Store(q.QQ, botClient)
			AfterLogin(botClient, q.ClientProtocol)
			return true
		}
	}

	return false
}

func GetBotIndex(botId int64) int32 {
	TempBotLock.Lock()
	defer TempBotLock.Unlock()
	botIdStr := strconv.FormatInt(botId, 10)
	for k, v := range TempBotData {
		if v.QQ == botIdStr {
			return int32(k)
		}
	}
	return -1
}
