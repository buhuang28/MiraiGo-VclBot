package bot

import (
	log "github.com/sirupsen/logrus"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"strings"
)

func (f *TSMSForm) OnFormCreate(sender vcl.IObject) {
	f.SetCaption("短信")
	f.SetDoubleBuffered(true)
	f.SetHeight(120)
	f.SetWidth(240)
	f.ScreenCenter()
	f.SetBorderStyle(types.BsSingle)
	f.EnabledMaximize(false)
	f.SetShowInTaskBar(types.StAlways)

	f.TipLabel = vcl.NewLabel(f)
	f.TipLabel.SetParent(f)
	f.TipLabel.SetLeft(10)
	f.TipLabel.SetTop(20)

	f.SMSLabel = vcl.NewLabel(f)
	f.SMSLabel.SetParent(f)
	f.SMSLabel.SetTop(f.TipLabel.Top() + 25)
	f.SMSLabel.SetLeft(20)
	f.SMSLabel.SetCaption("验证码:")

	f.SMSCode = vcl.NewEdit(f)
	f.SMSCode.SetParent(f)
	f.SMSCode.SetTop(f.SMSLabel.Top() - 3)
	f.SMSCode.SetLeft(f.SMSLabel.Left() + f.SMSLabel.Width()/2 + 15)
	f.SMSCode.SetWidth(85)

	f.SubmitButton = vcl.NewButton(f)
	f.SubmitButton.SetParent(f)
	f.SubmitButton.SetTop(f.SMSCode.Top() + 30)
	f.SubmitButton.SetWidth(90)
	f.SubmitButton.SetCaption("提交")
	f.SubmitButton.SetLeft(f.Width()/2 - f.SubmitButton.Width()/2)

	f.SubmitButton.SetOnClick(func(sender vcl.IObject) {
		if TempCaptchaQQ == 0 {
			return
		}
		cli, _ := Clients.Load(TempCaptchaQQ)
		code := strings.TrimSpace(f.SMSCode.Text())
		if code == "" {
			vcl.ShowMessage("验证码不可为空")
			return
		}
		rsp, err := cli.SubmitSMS(code)
		f.Hide()
		f.SMSLabel.SetCaption("")
		f.SMSCode.Clear()
		defer func() {
			TempCaptchaQQ = 0
		}()
		if err != nil || !rsp.Success {
			UpdateBotItem(TempCaptchaQQ, "", OFFLINE, "", "", err.Error())
			log.Info("短信验证码提交后出错:", err)
			cli.Disconnect()
			return
		}
		if rsp.Success {
			UpdateBotItem(TempCaptchaQQ, "", ONLINE, "", "", LOGIN_SUCCESS)
			go AfterLogin(cli, -1)
		}
	})

}
