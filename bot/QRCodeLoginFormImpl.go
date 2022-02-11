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

	f.GetQRCodeButton.SetOnClick(func(sender vcl.IObject) {
		qrCodeUrlBytes := GetQRCodeUrl(QRCodeLoginForm.ProtocolCheck.Items().IndexOf(QRCodeLoginForm.ProtocolCheck.Text()))
		QRCodeLoginForm.Image.Picture().LoadFromBytes(qrCodeUrlBytes)
		go func() {
			if qrCodeBot.Online.Load() {
				return
			}
			thisSig := tempLoginSig
			tempLoginSig = []byte("")
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
					break
				}

				log.Infof("扫码登录成功")
				originCli, ok := Clients.Load(qrCodeBot.Uin)
				if ok {
					originCli.Release()
				}
				botLock.Lock()
				index := GetBotIndex(qrCodeBot.Uin)

				var botData TTempItem
				botData.IconIndex = index
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
				SetBotAvatar(qrCodeBot.Uin, index)
				AddTempBotData(botData)
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

func GetQRCodeUrl(clientProtocol int32) []byte {
	if qrCodeBot != nil {
		qrCodeBot.Release()
	}
	qrCodeBot = client.NewClientEmpty()
	tempDeviceInfo = device.GetDevice(time.Now().Unix(), clientProtocol)
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
