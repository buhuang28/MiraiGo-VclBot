package bot

import (
	"crypto/md5"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"strconv"
	"strings"
)

func (f *TPWLoginForm) OnFormCreate(sender vcl.IObject) {
	f.SetCaption("密码")
	f.EnabledMaximize(false)
	f.SetBorderStyle(types.BsSingle)
	f.SetHeight(180)
	f.SetWidth(200)
	f.ScreenCenter()
	f.SetShowInTaskBar(types.StAlways)

	//f.SetColor()
	f.QQLabel = vcl.NewLabel(f)
	f.QQLabel.SetParent(f)
	f.QQLabel.SetCaption("Q Q:")
	f.QQLabel.SetLeft(20)
	f.QQLabel.SetTop(15)
	f.QQ = vcl.NewEdit(f)
	f.QQ.SetParent(f)
	f.QQ.SetWidth(120)
	f.QQ.SetLeft(f.QQLabel.Left() + 30)
	f.QQ.SetTop(f.QQLabel.Top() - 3)
	f.PWLabel = vcl.NewLabel(f)
	f.PWLabel.SetParent(f)
	f.PWLabel.SetCaption("密码:")
	f.PWLabel.SetLeft(20)
	f.PWLabel.SetTop(f.QQLabel.Top() + f.QQLabel.Height() + 20)
	f.PW = vcl.NewEdit(f)
	f.PW.SetParent(f)
	f.PW.SetWidth(120)
	f.PW.SetLeft(f.PWLabel.Left() + 30)
	f.PW.SetTop(f.PWLabel.Top() - 3)
	f.PW.SetPasswordChar('*')

	f.ProtocolText = vcl.NewLabel(f)
	f.ProtocolText.SetParent(f)
	f.ProtocolText.SetLeft(20)
	f.ProtocolText.SetTop(f.PWLabel.Top() + f.PW.Height() + 10)
	f.ProtocolText.SetCaption("登录协议:")

	f.ProtocolCheck = vcl.NewComboBox(f)
	f.ProtocolCheck.SetParent(f)
	f.ProtocolCheck.SetTop(f.ProtocolText.Top() - 5)
	f.ProtocolCheck.SetLeft(f.ProtocolText.Left() + f.ProtocolText.Width())
	f.ProtocolCheck.Items().Add(Ipad)
	f.ProtocolCheck.Items().Add(AndroidPhone)
	f.ProtocolCheck.Items().Add(AndroidWatch)
	f.ProtocolCheck.Items().Add(MacOS)
	f.ProtocolCheck.Items().Add(QiDian)
	f.ProtocolCheck.SetItemIndex(0)
	f.ProtocolCheck.SetWidth(85)

	f.AutoLogin = vcl.NewCheckBox(f)
	f.AutoLogin.SetParent(f)
	f.AutoLogin.SetTop(f.ProtocolCheck.Top() + 35)
	f.AutoLogin.SetCaption("自动登录")
	f.AutoLogin.SetLeft(f.Width()/2 - f.AutoLogin.Width()/3)

	f.LoginButton = vcl.NewButton(f)
	f.LoginButton.SetParent(f)
	f.LoginButton.SetTop(f.AutoLogin.Top() + 30)
	f.LoginButton.SetCaption("登录")
	f.LoginButton.SetLeft(f.Width()/2 - f.LoginButton.Width()/2)
	f.LoginButton.SetOnClick(func(sender vcl.IObject) {
		QQ := strings.TrimSpace(f.QQ.Text())
		if QQ == "" {
			vcl.ShowMessage("请输入QQ号")
			return
		}
		QQInt, err := strconv.ParseInt(QQ, 10, 64)
		if err != nil {
			vcl.ShowMessage("请输入正确的QQ号")
			return
		}
		PW := strings.TrimSpace(f.PW.Text())
		if PW == "" {
			vcl.ShowMessage("请输入密码")
			return
		}
		PWMd5 := md5.Sum([]byte(PW))
		loginProtocol := f.ProtocolCheck.Items().IndexOf(f.ProtocolCheck.Text())
		var qqInfo QQInfo
		qqInfo.StoreLoginInfo(QQInt, PWMd5, nil, loginProtocol, f.AutoLogin.Checked())
		f.Hide()
		f.QQ.Clear()
		f.PW.Clear()
		f.AutoLogin.SetChecked(false)
		go CreateBotImplMd5(QQInt, PWMd5, QQInt, loginProtocol, f.AutoLogin.Checked())
	})
}
