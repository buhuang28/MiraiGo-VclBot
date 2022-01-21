package bot

import "github.com/ying32/govcl/vcl"

var (
	DeviceVerifyForm *TDeviceVerifyForm
)

type TDeviceVerifyForm struct {
	*vcl.TForm
	TipLabel *vcl.TLabel
	QRCode   *vcl.TImage
}
