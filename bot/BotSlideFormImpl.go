package bot

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"strings"
)

func (f *TBotSlideForm) OnFormCreate(sender vcl.IObject) {
	f.SetCaption("滑块验证码")
	f.SetDoubleBuffered(true)
	f.SetHeight(270)
	f.SetWidth(400)
	f.ScreenCenter()
	f.SetBorderStyle(types.BsSingle)
	f.EnabledMaximize(false)
	f.SetShowInTaskBar(types.StAlways)

	f.VerifyUrlLabel = vcl.NewLabel(f)
	f.VerifyUrlLabel.SetParent(f)
	f.VerifyUrlLabel.SetTop(15)
	f.VerifyUrlLabel.SetLeft(15)
	f.VerifyUrlLabel.SetCaption("验证链接:")

	f.VerifyUrl = vcl.NewEdit(f)
	f.VerifyUrl.SetParent(f)
	f.VerifyUrl.SetTop(f.VerifyUrlLabel.Top() - 3)
	f.VerifyUrl.SetLeft(f.VerifyUrlLabel.Left() + f.VerifyUrlLabel.Width() + 10)
	f.VerifyUrl.SetWidth(300)

	f.VerifyQRCodeLabel = vcl.NewLabel(f)
	f.VerifyQRCodeLabel.SetParent(f)
	f.VerifyQRCodeLabel.SetLeft(f.VerifyUrlLabel.Left())
	f.VerifyQRCodeLabel.SetTop(f.VerifyUrlLabel.Top() + 40)
	f.VerifyQRCodeLabel.SetCaption("验证二维码:")

	f.VerifyQRCode = vcl.NewImage(f)
	f.VerifyQRCode.SetParent(f)
	f.VerifyQRCode.SetTop(f.VerifyQRCodeLabel.Top() - 3)
	f.VerifyQRCode.SetLeft(f.VerifyQRCodeLabel.Left() + f.VerifyQRCodeLabel.Width() + 10)
	f.VerifyQRCode.SetWidth(200)
	f.VerifyQRCode.SetHeight(200)

	f.TicketLabel = vcl.NewLabel(f)
	f.TicketLabel.SetParent(f)
	f.TicketLabel.SetLeft(f.VerifyUrlLabel.Left())
	f.TicketLabel.SetTop(f.VerifyQRCode.Top() + 160)
	f.TicketLabel.SetCaption("ticket:")

	f.Ticket = vcl.NewEdit(f)
	f.Ticket.SetParent(f)
	f.Ticket.SetTop(f.TicketLabel.Top() - 3)
	f.Ticket.SetLeft(f.TicketLabel.Left() + f.TicketLabel.Width() + 10)
	f.Ticket.SetWidth(300)

	f.SubmitButton = vcl.NewButton(f)
	f.SubmitButton.SetParent(f)
	f.SubmitButton.SetTop(f.Ticket.Top() + 30)
	f.SubmitButton.SetWidth(100)
	f.SubmitButton.SetLeft(f.Width()/2 - 50)
	f.SubmitButton.SetCaption("提交")

	f.SubmitButton.SetOnClick(func(sender vcl.IObject) {
		go func() {
			if TempCaptchaQQ == 0 {
				return
			}
			cli, _ := Clients.Load(TempCaptchaQQ)
			ticket := strings.TrimSpace(f.Ticket.Text())
			fmt.Println("提交是ticket:", ticket)
			rsp, err := cli.SubmitTicket(ticket)
			f.Hide()
			f.Ticket.Clear()
			f.VerifyQRCode.Hide()
			TempCaptchaQQ = 0
			if rsp.Success {
				UpdateBotItem(cli.Uin, "", ONLINE, "", "", LOGIN_SUCCESS)
				go AfterLogin(cli, -1)
				return
			}
			if err != nil || !rsp.Success {
				UpdateBotItem(cli.Uin, "", OFFLINE, "", "", rsp.ErrorMessage)
				ProcessLoginRsp(cli, rsp)
				log.Info("滑块提交后出错:", err)
				log.Info("滑块提交后出错2:", rsp.ErrorMessage)
				log.Info("滑块提交后出错3:", rsp.Success)
				cli.Disconnect()
				return
			}
		}()
	})
}
