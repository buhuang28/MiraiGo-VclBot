package bot

import (
	"MiraiGo-VclBot/device"
	"MiraiGo-VclBot/util"
	"encoding/json"
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
)

var (
	botLock sync.Mutex
)

func CreateBotImplMd5(uin int64, passwordMd5 [16]byte, deviceRandSeed int64, clientProtocol int32, autoLogin bool) bool {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
	}()
	log.Infof("开始初始化设备信息")
	deviceInfo := device.GetDevice(uin, clientProtocol)
	if deviceRandSeed != 0 && deviceRandSeed != uin {
		deviceInfo = device.GetDevice(deviceRandSeed, clientProtocol)
	}
	log.Infof("设备信息 %+v", string(deviceInfo.ToJson()))
	log.Infof("创建机器人 %+v", uin)

	cli := client.NewClientMd5(uin, passwordMd5)
	cli.UseDevice(deviceInfo)
	log.Infof("初始化日志")
	InitLog(cli)

	index := GetBotIndex(cli.Uin)
	var botData TTempItem
	if index == -1 {
		index = int32(len(TempBotData))
	}
	botData.IconIndex = index
	SetBotAvatar(cli.Uin, index)
	botData.QQ = strconv.FormatInt(cli.Uin, 10)
	botData.Protocol = GetProtocol(clientProtocol)
	botData.Status = LOGIN
	botData.NickName = ""
	if autoLogin {
		botData.Auto = "√"
	} else {
		botData.Auto = "X"
	}
	botData.Note = LOGIN
	AddTempBotData(botData)
	BotForm.BotListView.Items().SetCount(int32(len(TempBotData))) //   必须主动的设置Virtual List的行数
	log.Info(uin, "密码登录中...")
	Clients.Store(uin, cli)
	ok, err := Login(cli)
	if err != nil || !ok {
		log.Info(uin, "密码登录失败")
		return false
	}
	cli.Online.Store(true)
	log.Info(uin, "密码登录成功")
	AfterLogin(cli, clientProtocol)
	return true
}

func AfterLogin(cli *client.QQClient, clientProtocol int32) bool {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
	}()
	Serve(cli)
	log.Infof("刷新好友列表")
	if err := cli.ReloadFriendList(); err != nil {
		log.Info("刷新好友列表失败:", err)
		//util.FatalError(fmt.Errorf("failed to load friend list, err: %+v", err))
	}
	log.Infof("%v共加载 %v 个好友.", cli.Uin, len(cli.FriendList))

	log.Infof("刷新群列表")
	if err := cli.ReloadGroupList(); err != nil {
		log.Info("刷新群列表失败:", err)
		//util.FatalError(fmt.Errorf("failed to load group list, err: %+v", err))
	}
	log.Infof("%v共加载 %v 个群.", cli.Uin, len(cli.GroupList))

	SetRelogin(cli, 30, 20)
	BuhuangBotOnline(cli.Uin)
	go func() {
		var qqInfo QQInfo
		fileByte := util.ReadFileByte(QQINFOPATH + strconv.FormatInt(cli.Uin, 10) + QQINFOSKIN)
		json.Unmarshal(fileByte, &qqInfo)
		getToken := cli.GenToken()
		if clientProtocol == -1 {
			clientProtocol = qqInfo.ClientProtocol
		}
		LoginTokens.Store(cli.Uin, getToken)
		qqInfo.StoreLoginInfo(cli.Uin, qqInfo.PassWord, getToken, clientProtocol, qqInfo.AutoLogin)
		fmt.Println("获取token成功")
	}()
	return true
}
