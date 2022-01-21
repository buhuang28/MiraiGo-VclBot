package bot

import (
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

func (f *TQRCodeLoginForm) OnFormCreate(sender vcl.IObject) {
	f.SetCaption("扫码")
	f.EnabledMaximize(false)
	f.SetBorderStyle(types.BsSingle)
	f.SetHeight(180)
	f.SetWidth(280)
	f.ScreenCenter()
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
	f.ProtocolCheck.SetLeft(160)
	f.ProtocolCheck.SetTop(30)
	f.ProtocolCheck.SetItemIndex(0)
	f.AutoLogin = vcl.NewCheckBox(f)
	f.AutoLogin.SetParent(f)
	f.AutoLogin.SetCaption("自动登录")
	f.AutoLogin.SetTop(80)
	f.AutoLogin.SetLeft(180)
}
