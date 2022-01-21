package bot

import "github.com/ying32/govcl/vcl"

var (
	CaptchaForm *TCaptchaForm
)

type TCaptchaForm struct {
	*vcl.TForm
	Captcha      *vcl.TImage
	CodeLabel    *vcl.TLabel
	Code         *vcl.TEdit
	SubmitButton *vcl.TButton
}
