package bot

import "github.com/ying32/govcl/vcl"

var (
	PWLoginForm *TPWLoginForm
)

type TPWLoginForm struct {
	*vcl.TForm
	QQ            *vcl.TEdit
	QQLabel       *vcl.TLabel
	PW            *vcl.TEdit
	PWLabel       *vcl.TLabel
	ProtocolText  *vcl.TLabel
	ProtocolCheck *vcl.TComboBox
	AutoLogin     *vcl.TCheckBox
	LoginButton   *vcl.TButton
}
