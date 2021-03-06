package bot

import (
	"MiraiGo-VclBot/device"
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	log "github.com/sirupsen/logrus"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
	"time"
)

func (f *TQRCodeLoginForm) OnFormCreate(sender vcl.IObject) {
	f.SetCaption("扫码")
	f.SetHeight(180)
	f.SetWidth(280)
	f.ScreenCenter()
	f.EnabledMaximize(false)
	f.SetBorderStyle(types.BsSingle)
	f.SetShowInTaskBar(types.StAlways)
	f.SetDoubleBuffered(true)
	f.Image = vcl.NewImage(f)
	f.Image.SetTop(10)
	f.Image.SetLeft(10)
	f.Image.SetHeight(180)
	f.Image.SetWidth(180)
	f.Image.Update()
	f.Image.SetParent(f)
	f.ProtocolText = vcl.NewLabel(f)
	f.ProtocolText.SetCaption("登录协议:")
	f.ProtocolText.SetParent(f)
	f.ProtocolText.SetLeft(180)
	f.ProtocolText.SetTop(10)
	f.ProtocolCheck = vcl.NewComboBox(f)
	f.ProtocolCheck.SetParent(f)
	f.ProtocolCheck.Items().Add(Ipad)
	f.ProtocolCheck.Items().Add(AndroidPhone)
	f.ProtocolCheck.Items().Add(AndroidWatch)
	f.ProtocolCheck.Items().Add(MacOS)
	f.ProtocolCheck.Items().Add(QiDian)
	f.ProtocolCheck.SetLeft(160)
	f.ProtocolCheck.SetTop(30)
	f.ProtocolCheck.SetItemIndex(0)
	f.AutoLogin = vcl.NewCheckBox(f)
	f.AutoLogin.SetParent(f)
	f.AutoLogin.SetCaption("自动登录")
	f.AutoLogin.SetTop(80)
	f.AutoLogin.SetLeft(180)
	f.GetQRCodeButton = vcl.NewButton(f)
	f.GetQRCodeButton.SetParent(f)
	f.GetQRCodeButton.SetCaption("获取登录二维码")
	f.GetQRCodeButton.SetHeight(25)
	f.GetQRCodeButton.SetWidth(120)
	f.GetQRCodeButton.SetLeft(18)
	f.GetQRCodeButton.SetTop(145)

	f.SeedLabel = vcl.NewLabel(f)
	f.SeedLabel.SetParent(f)
	f.SeedLabel.SetCaption("种子:")
	f.SeedLabel.SetTop(f.AutoLogin.Top() + 30)
	f.SeedLabel.SetLeft(f.AutoLogin.Left() + 15)

	f.Seed = vcl.NewEdit(f)
	f.Seed.SetParent(f)
	f.Seed.SetTop(f.SeedLabel.Top() + 20)
	f.Seed.SetWidth(80)
	f.Seed.SetLeft(f.SeedLabel.Left() - 20)

	f.GetQRCodeButton.SetOnClick(func(sender vcl.IObject) {
		var seed int64 = 0
		seedStr := strings.TrimSpace(f.Seed.Text())
		parseInt, err := strconv.ParseInt(seedStr, 10, 64)
		if err != nil {
			vcl.ShowMessage("种子错误，必须纯数字")
			return
		}
		if parseInt != 0 {
			seed = parseInt
		}
		qrCodeUrlBytes := GetQRCodeUrl(seed, QRCodeLoginForm.ProtocolCheck.Items().IndexOf(QRCodeLoginForm.ProtocolCheck.Text()))
		QRCodeLoginForm.Image.Picture().LoadFromBytes(qrCodeUrlBytes)
		QRCodeLoginForm.Image.Show()
		go func() {
			if qrCodeBot.Online.Load() {
				log.Info(qrCodeBot.Uin, "已在线")
				return
			}
			thisSig := tempLoginSig
			for i := 0; i < 100; i++ {
				queryQRCodeStatusResp, err := qrCodeBot.QueryQRCodeStatus(thisSig)
				if err != nil {
					log.Info("failed to query qrcode status:", err)
					break
				}
				if queryQRCodeStatusResp.State != client.QRCodeConfirmed {
					time.Sleep(time.Second * 3)
					continue
				}
				loginResp, err := qrCodeBot.QRCodeLogin(queryQRCodeStatusResp.LoginInfo)
				if err != nil || !loginResp.Success {
					UpdateBotItem(qrCodeBot.Uin, qrCodeBot.Nickname, OFFLINE, "", "", loginResp.ErrorMessage)
					log.Info("扫码登录失败:", err)
					log.Infof("扫码登录失败:%v", loginResp)
					break
				}
				log.Infof("扫码登录成功")
				originCli, ok := Clients.Load(qrCodeBot.Uin)
				if ok {
					originCli.Release()
				}
				var botData TTempItem

				index := GetBotIndex(qrCodeBot.Uin)
				if index == -1 {
					botData.IconIndex = int32(len(TempBotData))
				} else {
					botData.IconIndex = index
				}
				SetBotAvatar(qrCodeBot.Uin, botData.IconIndex)
				//这里index一般是-1
				botData.NickName = qrCodeBot.Nickname
				botData.QQ = strconv.FormatInt(qrCodeBot.Uin, 10)
				botData.Protocol = QRCodeLoginForm.ProtocolCheck.Text()
				botData.Status = ONLINE
				botData.Note = LOGIN_SUCCESS
				if QRCodeLoginForm.AutoLogin.Checked() {
					botData.Auto = "√"
				} else {
					botData.Auto = "X"
				}
				AddTempBotData(botData)
				//index = GetBotIndex(qrCodeBot.Uin)
				//SetBotAvatar(qrCodeBot.Uin, index)
				var qqInfo QQInfo
				qqInfo.StoreLoginInfo(qrCodeBot.Uin, [16]byte{}, qrCodeBot.GenToken(), int32(tempDeviceInfo.Protocol), QRCodeLoginForm.AutoLogin.Checked())
				Clients.Store(qrCodeBot.Uin, qrCodeBot)
				go AfterLogin(qrCodeBot, int32(tempDeviceInfo.Protocol))
				devicePath := path.Join("device", fmt.Sprintf("device-%d.json", qrCodeBot.Uin))
				_ = ioutil.WriteFile(devicePath, tempDeviceInfo.ToJson(), 0644)
				qrCodeBot = nil
				break
			}
			vcl.ThreadSync(func() {
				QRCodeLoginForm.Image.Hide()
				QRCodeLoginForm.Hide()
			})
		}()
	})
}

var (
	tempDeviceInfo *client.DeviceInfo
	qrCodeBot      *client.QQClient
	tempLoginSig   []byte
)

func GetQRCodeUrl(seed int64, clientProtocol int32) []byte {
	if qrCodeBot != nil {
		qrCodeBot.Release()
	}
	qrCodeBot = client.NewClientEmpty()
	if seed != 0 {
		tempDeviceInfo = device.GetDevice(seed, clientProtocol)
	} else {
		tempDeviceInfo = device.GetDevice(time.Now().Unix(), clientProtocol)
	}
	qrCodeBot.UseDevice(tempDeviceInfo)
	log.Infof("初始化日志")
	InitLog(qrCodeBot)
	fetchQRCodeResp, err := qrCodeBot.FetchQRCode()
	if err != nil {
		vcl.ShowMessage("获取二维码失败")
		log.Info("获取二维码失败:", err)
		return nil
	}
	//QRCodeLoginForm.Image.Picture().LoadFromBytes(fetchQRCodeResp.ImageData)
	tempLoginSig = fetchQRCodeResp.Sig
	return fetchQRCodeResp.ImageData
}
