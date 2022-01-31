package bot

import (
	log "github.com/sirupsen/logrus"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"strings"
)

func (f *TCaptchaForm) OnFormCreate(sender vcl.IObject) {
	f.SetCaption("图形验证码")
	f.SetDoubleBuffered(true)
	f.SetHeight(145)
	f.SetWidth(250)
	f.ScreenCenter()
	f.SetBorderStyle(types.BsSingle)
	f.EnabledMaximize(false)
	f.SetShowInTaskBar(types.StAlways)

	f.Captcha = vcl.NewImage(f)
	f.Captcha.SetParent(f)
	f.Captcha.SetTop(10)
	f.Captcha.SetWidth(200)
	f.Captcha.SetHeight(50)
	f.Captcha.SetLeft(f.Width()/2 - 75)

	f.CodeLabel = vcl.NewLabel(f)
	f.CodeLabel.SetParent(f)
	f.CodeLabel.SetTop(f.Captcha.Top() + f.Captcha.Height() + 20)
	f.CodeLabel.SetLeft(50)
	f.CodeLabel.SetCaption("验证码:")

	f.Code = vcl.NewEdit(f)
	f.Code.SetParent(f)
	f.Code.SetWidth(90)
	f.Code.SetTop(f.CodeLabel.Top() - 3)
	f.Code.SetLeft(f.CodeLabel.Left() + f.CodeLabel.Width()/2 + 15)

	f.SubmitButton = vcl.NewButton(f)
	f.SubmitButton.SetParent(f)
	f.SubmitButton.SetWidth(100)
	f.SubmitButton.SetLeft(f.Width()/2 - f.SubmitButton.Width()/2)
	f.SubmitButton.SetTop(f.Code.Top() + 35)
	f.SubmitButton.SetCaption("提交")

	f.SubmitButton.SetOnClick(func(sender vcl.IObject) {
		if TempCaptchaQQ == 0 {
			return
		}
		code := strings.TrimSpace(f.Code.Text())
		if code == "" {
			vcl.ShowMessage("验证码不可为空")
			return
		}
		cli, _ := Clients.Load(TempCaptchaQQ)
		rsp, err := cli.SubmitCaptcha(code, TempCaptchSign)
		f.Hide()
		f.Captcha.Hide()
		f.Code.Clear()
		TempCaptchaQQ = 0
		//index := botIndexMap[TempCaptchaQQ]
		index := GetBotIndex(TempCaptchaQQ)
		if err != nil || !rsp.Success {
			TempBotData[index].Status = "离线"
			TempBotData[index].Note = "登录失败"
			log.Info("图片验证码提交后出错:", err)
			cli.Disconnect()
			return
		}
		if rsp.Success {
			TempBotData[index].Status = "在线"
			TempBotData[index].Note = "登录成功"
			go AfterLogin(cli, -1)
		}
	})
}
