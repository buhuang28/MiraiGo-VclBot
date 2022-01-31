package bot

import (
	"MiraiGo-VclBot/device"
	"MiraiGo-VclBot/util"
	"encoding/json"
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	log "github.com/sirupsen/logrus"
	"github.com/ying32/govcl/vcl"
	"strconv"
	"sync"
)

var (
	//botIndexMap       = make(map[int64]int)
	//botIndexStart int = 0
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

	botLock.Lock()
	index := GetBotIndex(cli.Uin)
	var botData TTempItem
	if index == -1 {
		index = int32(len(TempBotData))
	}
	botData.IconIndex = index
	avatarUrl := AvatarUrlPre + strconv.FormatInt(cli.Uin, 10)
	bytes, err2 := util.GetBytes(avatarUrl)
	if err2 != nil {
		fmt.Println(err2)
	} else {
		pic := vcl.NewPicture()
		pic.LoadFromBytes(bytes)
		BotForm.Icons.AddSliced(pic.Bitmap(), 1, 1)
		pic.Free()
		SetBotAvatarIndex(cli.Uin, index)
	}
	BotForm.BotListView.SetStateImages(BotForm.Icons)
	botData.QQ = strconv.FormatInt(cli.Uin, 10)
	botData.Protocol = GetProtocol(clientProtocol)
	botData.Status = "登录中"
	botData.NickName = ""
	if autoLogin {
		botData.Auto = "√"
	} else {
		botData.Auto = "X"
	}
	botData.Note = "登录中"
	TempBotLock.Lock()
	TempBotData = append(TempBotData, botData)
	TempBotLock.Unlock()
	BotForm.BotListView.Items().SetCount(int32(len(TempBotData))) //   必须主动的设置Virtual List的行数
	botLock.Unlock()

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
	//for {
	//	time.Sleep(5 * time.Second)
	//	if cli.Online.Load() {
	//		break
	//	}
	//	log.Warnf("%+v机器人不在线，可能在等待输入验证码，或出错了。如果出错请重启。", cli.Uin)
	//}
	Serve(cli)
	log.Infof("刷新好友列表")
	if err := cli.ReloadFriendList(); err != nil {
		log.Info("刷新好友列表失败:", err)
		//util.FatalError(fmt.Errorf("failed to load friend list, err: %+v", err))
	}
	log.Infof("共加载 %v 个好友.", len(cli.FriendList))

	log.Infof("刷新群列表")
	if err := cli.ReloadGroupList(); err != nil {
		log.Info("刷新群列表失败:", err)
		//util.FatalError(fmt.Errorf("failed to load group list, err: %+v", err))
	}
	log.Infof("共加载 %v 个群.", len(cli.GroupList))

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
		qqInfo.StoreLoginInfo(cli.Uin, qqInfo.PassWord, getToken, clientProtocol, qqInfo.AutoLogin)
		fmt.Println("获取token成功")
	}()
	return true
}
