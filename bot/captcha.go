package bot

import (
	"MiraiGo-VclBot/util"
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	log "github.com/sirupsen/logrus"
	"github.com/ying32/govcl/vcl"
	"strconv"
	"strings"
	"time"
)

//type WaitingCaptcha struct {
//	Prom    *promise.Promise
//}

//go:generate go run github.com/a8m/syncmap -o "gen_captcha_map.go" -pkg bot -name CaptchaMap "map[int64]*WaitingCaptcha"

// TODO sync
//var WaitingCaptchas CaptchaMap

var (
	TempCaptchaQQ  int64
	TempCaptchSign []byte
	//CaptchaPath = "./captcha.jpg"
)

func ProcessLoginRsp(cli *client.QQClient, rsp *client.LoginResponse) bool {
	if rsp.Success {
		//index := GetBotIndex(cli.Uin)
		//TempBotData[index].NickName = cli.Nickname
		//TempBotData[index].Status = "在线"
		//TempBotData[index].Note = "登录成功"
		UpdateBotItem(cli.Uin, cli.Nickname, ONLINE, "", "", LOGIN_SUCCESS)
		return true
	}
	if rsp.Error == client.SMSOrVerifyNeededError {
		rsp.Error = client.SMSNeededError
	}

	if strings.Contains(rsp.ErrorMessage, "冻结") {
		UpdateBotItem(cli.Uin, "", "冻结", "", "", "账号被冻结")
		Clients.Delete(cli.Uin)
		return false
	}

	if TempCaptchaQQ != 0 {
		//vcl.ShowMessage("还有正在处理登录验证的机器人:" + strconv.FormatInt(TempCaptchaQQ, 10) + "错误问题:" + rsp.ErrorMessage)
		return false
	}
	TempCaptchaQQ = cli.Uin
	switch rsp.Error {
	case client.SliderNeededError:
		log.Info("遇到滑块验证码")
		go func() {
			vcl.ThreadSync(func() {
				BotSlideForm.VerifyUrl.SetText(rsp.VerifyUrl)
				code := util.CreateQRCode(rsp.VerifyUrl, 150)
				BotSlideForm.VerifyQRCode.Show()
				BotSlideForm.VerifyQRCode.Picture().LoadFromFile(code)
				BotSlideForm.Show()
			})
		}()
		return false
	case client.NeedCaptcha:
		log.Info("遇到图形验证码")
		go func() {
			TempCaptchSign = rsp.CaptchaSign
			vcl.ThreadSync(func() {
				CaptchaForm.Captcha.Show()
				CaptchaForm.Captcha.Picture().LoadFromBytes(rsp.CaptchaImage)
				CaptchaForm.Show()
			})
		}()
		return false
	case client.SMSNeededError:
		log.Info("遇到短信验证码")
		if cli.RequestSMS() {
			go func() {
				vcl.ThreadSync(func() {
					SMSForm.TipLabel.SetCaption(fmt.Sprintf("已经向手机号为:%s 发送短信验证码", rsp.SMSPhone))
					SMSForm.Show()
				})
			}()
			return false
		} else {
			go func() {
				UpdateBotItem(cli.Uin, "", "", "", "", "手机号"+rsp.SMSPhone+"请求短信验证码错误，可能是太频繁")
				vcl.ThreadSync(func() {
					vcl.ShowMessage("手机号" + rsp.SMSPhone + "请求短信验证码错误，可能是太频繁")
				})
			}()
			return false
		}
	case client.UnsafeDeviceError:
		log.Info("遇到设备锁，需要手机QQ扫码验证,请在3分钟内完成验证")
		go func() {
			qrCode := util.CreateQRCode(rsp.VerifyUrl, 150)
			vcl.ThreadSync(func() {
				DeviceVerifyForm.QRCode.Picture().LoadFromFile(qrCode)
				DeviceVerifyForm.QRCode.Show()
				DeviceVerifyForm.Show()
			})
		}()
		var i int32
		for i = 0; i < 30; i++ {
			cli.Disconnect()
			time.Sleep(5 * time.Second)
			resp, err := cli.Login()
			if err != nil || !resp.Success {
				continue
			} else {
				UpdateBotItem(TempCaptchaQQ, "", ONLINE, "", "", LOGIN_SUCCESS)
				TempCaptchaQQ = 0
				go func() {
					vcl.ThreadSync(func() {
						DeviceVerifyForm.QRCode.Hide()
						DeviceVerifyForm.Hide()
					})
				}()
				go AfterLogin(cli, -1)
				return true
			}
		}
		TempCaptchaQQ = 0
		go func() {
			//TempBotData[index].Status = "离线"
			//TempBotData[index].Note = "设备扫码验证失败，登录失败"
			UpdateBotItem(TempCaptchaQQ, "", OFFLINE, "", "", QRCODE_CAPTCHA_ERROR)
			vcl.ThreadSync(func() {
				DeviceVerifyForm.QRCode.Hide()
				DeviceVerifyForm.Hide()
			})
		}()
		log.Info("设备扫码验证失败")
		cli.Disconnect()
		return false
	case client.OtherLoginError, client.UnknownLoginError:
		Clients.Delete(cli.Uin)
		if strings.Contains(rsp.ErrorMessage, "冻结") {
			UpdateBotItem(cli.Uin, "", "冻结", "", "", "账号被冻结")
			return false
		}
		//log.Errorf(rsp.ErrorMessage)
		go func() {
			vcl.ThreadSync(func() {
				vcl.ShowMessage(strconv.FormatInt(cli.Uin, 10) + "登录失败:" + rsp.ErrorMessage)
			})
		}()
		log.Info(cli.Uin, "遇到登录错误:", rsp.ErrorMessage)
		return false
	}
	log.Info("process login error")
	return false
}
