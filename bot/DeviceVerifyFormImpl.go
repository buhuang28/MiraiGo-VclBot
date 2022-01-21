package bot

import (
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

func (f *TDeviceVerifyForm) OnFormCreate(sender vcl.IObject) {
	f.SetCaption("扫码验证")
	f.SetDoubleBuffered(true)
	f.SetHeight(200)
	f.SetWidth(200)
	f.ScreenCenter()
	f.SetBorderStyle(types.BsSingle)
	f.EnabledMaximize(false)
	f.SetShowInTaskBar(types.StAlways)

	f.TipLabel = vcl.NewLabel(f)
	f.TipLabel.SetParent(f)
	f.TipLabel.SetCaption("请在3分钟内使用手机QQ扫码验证")
	f.TipLabel.SetTop(10)
	f.TipLabel.SetLeft(10)

	f.QRCode = vcl.NewImage(f)
	f.QRCode.SetParent(f)
	f.QRCode.SetWidth(150)
	f.QRCode.SetHeight(150)
	f.QRCode.SetLeft((f.Width() - f.QRCode.Width()) / 2)
	f.QRCode.SetTop(f.TipLabel.Top() + 30)
}
