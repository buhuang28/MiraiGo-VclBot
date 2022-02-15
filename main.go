package main

import (
	"MiraiGo-VclBot/bot"
	"MiraiGo-VclBot/bot_log"
	_ "github.com/ying32/govcl/pkgs/winappres"
	"github.com/ying32/govcl/vcl"
)

func init() {
	go bot_log.InitLog()
	go bot.WSDailCall()
	go bot.HandleWSMsg()
}

func main() {
	vcl.Application.Initialize()
	vcl.Application.CreateForm(&bot.BotForm)
	vcl.Application.CreateForm(&bot.PWLoginForm)
	vcl.Application.CreateForm(&bot.QRCodeLoginForm)
	vcl.Application.CreateForm(&bot.LogForm)
	vcl.Application.CreateForm(&bot.BotSlideForm)
	vcl.Application.CreateForm(&bot.CaptchaForm)
	vcl.Application.CreateForm(&bot.SMSForm)
	vcl.Application.CreateForm(&bot.DeviceVerifyForm)

	bot.PWLoginForm.Hide()
	bot.BotSlideForm.Hide()
	bot.QRCodeLoginForm.Hide()
	bot.LogForm.Hide()
	bot.CaptchaForm.Hide()
	bot.SMSForm.Hide()
	bot.DeviceVerifyForm.Hide()
	go func() {
		bot.AutoLogin()
	}()
	//bot.PWLoginForm.Show()
	//bot.BotSlideForm.Show()
	//bot.QRCodeLoginForm.Show()
	//bot.LogForm.Show()
	//bot.CaptchaForm.Show()
	//bot.SMSForm.Show()
	//bot.DeviceVerifyForm.Show()
	vcl.Application.Run()

}
