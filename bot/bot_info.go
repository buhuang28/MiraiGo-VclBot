package bot

import (
	"MiraiGo-VclBot/device"
	"MiraiGo-VclBot/util"
	"encoding/json"
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

func (q *QQInfo) StoreLoginInfo(qq int64, pw [16]byte, token []byte, clientProtocol int32) bool {
	q.QQ = qq
	if q.QQ == 0 {
		return false
	}
	q.PassWord = pw
	q.Token = token
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
		success := CreateBotImplMd5(q.QQ, q.PassWord, q.QQ, q.ClientProtocol, q.AutoLogin)
		if !success {
			time.Sleep(time.Second * 2)
			success = CreateBotImplMd5(q.QQ, q.PassWord, q.QQ, q.ClientProtocol, q.AutoLogin)
		}
		if success {
			time.Sleep(time.Second)
			return success
		}
	}
	var botClient = client.NewClientEmpty()
	deviceInfo := device.GetDevice(q.QQ, q.ClientProtocol)
	botClient.UseDevice(deviceInfo)
	err := botClient.TokenLogin(q.Token)
	if err != nil {
		time.Sleep(time.Second * 2)
		err = botClient.TokenLogin(q.Token)
	}
	if err != nil {
		return false
	} else {
		time.Sleep(time.Second)
		Clients.Store(q.QQ, botClient)
		AfterLogin(botClient, q.ClientProtocol)
		return true
	}
}
