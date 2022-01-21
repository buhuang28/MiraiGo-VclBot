package bot

import "github.com/ying32/govcl/vcl"

var (
	QRCodeLoginForm *TQRCodeLoginForm
)

type TQRCodeLoginForm struct {
	*vcl.TForm
	Image         *vcl.TImage
	ProtocolText  *vcl.TLabel
	ProtocolCheck *vcl.TComboBox
	AutoLogin     *vcl.TCheckBox
}
