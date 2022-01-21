package bot

import "github.com/ying32/govcl/vcl"

var (
	SMSForm *TSMSForm
)

type TSMSForm struct {
	*vcl.TForm
	TipLabel     *vcl.TLabel
	SMSLabel     *vcl.TLabel
	SMSCode      *vcl.TEdit
	SubmitButton *vcl.TButton
}
