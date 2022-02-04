package bot

import (
	"MiraiGo-VclBot/device"
	"github.com/Mrs4s/MiraiGo/client"
	log "github.com/sirupsen/logrus"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
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
	//f.Image.Picture().LoadFromFile("C:\\images\\2.png")
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
