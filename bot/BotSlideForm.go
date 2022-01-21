package bot

import "github.com/ying32/govcl/vcl"

var (
	BotSlideForm *TBotSlideForm
)

type TBotSlideForm struct {
	*vcl.TForm
	VerifyUrlLabel    *vcl.TLabel
	VerifyUrl         *vcl.TEdit
	VerifyQRCodeLabel *vcl.TLabel
	VerifyQRCode      *vcl.TImage
	TicketLabel       *vcl.TLabel
	Ticket            *vcl.TEdit
	SubmitButton      *vcl.TButton
}
