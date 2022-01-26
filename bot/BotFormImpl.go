package bot

import (
	"MiraiGo-VclBot/util"
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	log "github.com/sirupsen/logrus"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"github.com/ying32/govcl/vcl/types/colors"
	"github.com/ying32/govcl/vcl/win"
	"io/ioutil"
	"os/exec"
	"path"
	"strconv"
	"time"
)

type TForm1Fields struct {
	subItemHit win.TLVHitTestInfo
}

func (f *TBotForm) OnFormCreate(sender vcl.IObject) {
	f.SetCaption("机器人列表")
	f.SetDoubleBuffered(true)
	f.SetHeight(400)
	f.SetWidth(700)
	f.ScreenCenter()
	f.SetBorderStyle(types.BsSingle)
	f.EnabledMaximize(false)

	f.Icons = vcl.NewImageList(f)
	f.Icons.SetHeight(20)
	f.Icons.SetWidth(20)

	f.TempIco = vcl.NewIcon()

	f.BotListView = vcl.NewListView(f)
	f.BotListView.SetAlign(types.AlClient)
	f.BotListView.SetParent(f)
	f.BotListView.SetViewStyle(types.VsReport)
	f.BotListView.SetReadOnly(true)
	f.BotListView.SetOwnerData(true)
	f.BotListView.SetGridLines(true)
	f.BotListView.SetReadOnly(true)
	f.BotListView.SetRowSelect(true)
	f.BotListView.SetSmallImages(f.Icons)

	addCol := func(name string, width int32) {
		col := f.BotListView.Columns().Add()
		col.SetCaption(name)
		col.SetWidth(width)
		col.SetAlignment(types.TaCenter)
	}
	addCol("头像", 50)
	addCol("QQ", 100)
	addCol("昵称", 100)
	addCol("状态", 50)
	addCol("协议", 100)
	addCol("自动", 50)
	addCol("备注", 250)

	item := vcl.NewMenuItem(f.BotListView)
	item.SetCaption("扫码登录")
	item.SetOnClick(func(sender vcl.IObject) {
		qrCodeUrlBytes := GetQRCodeUrl(QRCodeLoginForm.ProtocolCheck.Items().IndexOf(QRCodeLoginForm.ProtocolCheck.Text()))
		QRCodeLoginForm.Image.Picture().LoadFromBytes(qrCodeUrlBytes)
		QRCodeLoginForm.Show()
		go func() {
			if qrCodeBot.Online.Load() {
				return
			}
			thisSig := tempLoginSig
			for i := 0; i < 100; i++ {
				queryQRCodeStatusResp, err := qrCodeBot.QueryQRCodeStatus(thisSig)
				if err != nil {
					log.Info("failed to query qrcode status:", err)
					break
				}
				if queryQRCodeStatusResp.State != client.QRCodeConfirmed {
					time.Sleep(time.Second * 3)
					continue
				}
				loginResp, err := qrCodeBot.QRCodeLogin(queryQRCodeStatusResp.LoginInfo)
				if err != nil || !loginResp.Success {
					vcl.ShowMessage(fmt.Sprintf("扫码登录失败: %+v", err))
					log.Info("扫码登录失败:", err)
					break
				}

				log.Infof("扫码登录成功")
				originCli, ok := Clients.Load(qrCodeBot.Uin)
				if ok {
					originCli.Release()
				}
				botLock.Lock()
				index, ok := botIndexMap[qrCodeBot.Uin]
				if !ok {
					index = botIndexStart
					botIndexStart++
				}
				botIndexMap[qrCodeBot.Uin] = index
				var botData TTempItem
				botData.IconIndex = int32(index)
				botData.NickName = qrCodeBot.Nickname
				botData.QQ = strconv.FormatInt(qrCodeBot.Uin, 10)
				botData.Protocol = QRCodeLoginForm.ProtocolCheck.Text()
				botData.Status = "在线"
				botData.Note = "登录成功"
				if QRCodeLoginForm.AutoLogin.Checked() {
					botData.Auto = "√"
				} else {
					botData.Auto = "X"
				}
				avatarUrl := AvatarUrlPre + strconv.FormatInt(qrCodeBot.Uin, 10)
				bytes, err := util.GetBytes(avatarUrl)
				if err == nil {
					pic := vcl.NewPicture()
					pic.LoadFromBytes(bytes)
					BotForm.Icons.AddSliced(pic.Bitmap(), 1, 1)
					pic.Free()
					BotForm.BotListView.SetStateImages(BotForm.Icons)
				}
				TempBotData = append(TempBotData, botData)
				BotForm.BotListView.Items().SetCount(int32(len(TempBotData))) //   必须主动的设置Virtual List的行数

				var qqInfo QQInfo
				qqInfo.StoreLoginInfo(qrCodeBot.Uin, [16]byte{}, qrCodeBot.GenToken(), int32(tempDeviceInfo.Protocol), QRCodeLoginForm.AutoLogin.Checked())
				Clients.Store(qrCodeBot.Uin, qrCodeBot)
				go AfterLogin(qrCodeBot, int32(tempDeviceInfo.Protocol))
				devicePath := path.Join("device", fmt.Sprintf("device-%d.json", qrCodeBot.Uin))
				_ = ioutil.WriteFile(devicePath, tempDeviceInfo.ToJson(), 0644)
				qrCodeBot = nil
				break
			}
			vcl.ThreadSync(func() {
				QRCodeLoginForm.Hide()
			})
		}()
	})

	item2 := vcl.NewMenuItem(f.BotListView)
	item2.SetCaption("密码登录")
	item2.SetOnClick(func(sender vcl.IObject) {
		PWLoginForm.Show()
	})

	item3 := vcl.NewMenuItem(f.BotListView)
	item3.SetCaption("显示日志")
	item3.SetOnClick(func(sender vcl.IObject) {
		LogForm.Show()
	})
	f.NoSelectMenu = vcl.NewPopupMenu(f.BotListView)
	f.NoSelectMenu.Items().Add(item)
	f.NoSelectMenu.Items().Add(item2)
	f.NoSelectMenu.Items().Add(item3)

	githubItem := vcl.NewMenuItem(f.BotListView)
	githubItem.SetCaption("项目开源地址:https://github.com/buhuang28/MiraiGo-VclBot")
	githubItem.SetOnClick(func(sender vcl.IObject) {
		exec.Command("cmd", "/c", "start", "https://github.com/buhuang28/MiraiGo-VclBot").Start()
	})
	f.NoSelectMenu.Items().Add(githubItem)

	item4 := vcl.NewMenuItem(f.BotListView)
	item4.SetCaption("登录")
	item4.SetOnClick(func(sender vcl.IObject) {
		selected := f.BotListView.Selected()
		if selected == nil {
			fmt.Println("未选中QQ")
			return
		}
		//Login
		LogForm.Show()
	})

	item5 := vcl.NewMenuItem(f.BotListView)
	item5.SetCaption("下线")
	item5.SetOnClick(func(sender vcl.IObject) {
		LogForm.Show()
	})

	item6 := vcl.NewMenuItem(f.BotListView)
	item6.SetCaption("重登")
	item6.SetOnClick(func(sender vcl.IObject) {
		LogForm.Show()
	})
	item7 := vcl.NewMenuItem(f.BotListView)
	item7.SetCaption("删除")
	item7.SetOnClick(func(sender vcl.IObject) {
		LogForm.Show()
	})
	f.SelectedMenu = vcl.NewPopupMenu(f.BotListView)
	f.SelectedMenu.Items().Add(item4)
	f.SelectedMenu.Items().Add(item5)
	f.SelectedMenu.Items().Add(item6)
	f.SelectedMenu.Items().Add(item7)

	f.BotListView.SetOnMouseDown(func(sender vcl.IObject, button types.TMouseButton, shift types.TShiftState, x, y int32) {
		if button == types.MbRight {
			p := f.BotListView.ScreenToClient(vcl.Mouse.CursorPos())
			f.subItemHit.Pt.X = p.X
			f.subItemHit.Pt.Y = p.Y
			win.ListView_SubItemHitTest(f.BotListView.Handle(), &f.subItemHit)
			p = vcl.Mouse.CursorPos()
			if f.subItemHit.IItem != -1 {
				f.SelectedMenu.Popup(p.X, p.Y)
			} else {
				f.NoSelectMenu.Popup(p.X, p.Y)
			}
		}
	})

	f.BotListView.SetOnAdvancedCustomDrawSubItem(func(sender *vcl.TListView, item *vcl.TListItem, subItem int32, state types.TCustomDrawState, stage types.TCustomDrawStage, defaultDraw *bool) {
		if len(TempBotData) == 0 {
			return
		}
		canvas := sender.Canvas()
		boundRect := item.DisplayRect(types.DrBounds)
		//当前状态，鼠标选中的那行显示的颜色
		data := TempBotData[item.Index()]
		drawFlags := types.NewSet(types.TfCenter, types.TfSingleLine, types.TfVerticalCenter)
		var i int32
		font := canvas.Font()
		if state.In(types.CdsFocused) {
			canvas.Brush().SetColor(colors.ClBisque)
		} else {
			canvas.Brush().SetColor(sender.Color())
		}
		canvas.FillRect(boundRect)
		font.SetColor(colors.ClBlack)

		for i = 0; i < sender.Columns().Count(); i++ {
			r := f.GetSubItemRect(sender.Handle(), item.Index(), i)
			switch i {
			case 0:
				var hw int32 = 20
				f.Icons.GetIcon(data.IconIndex, f.TempIco)
				if !f.TempIco.Empty() {
					canvas.Draw(r.Right/2-hw/2, r.Top+(r.Bottom-r.Top-hw)/2, f.TempIco)
				}
			case 1:
				canvas.TextRect2(&r, data.QQ, drawFlags)
			case 2:
				canvas.TextRect2(&r, data.NickName, drawFlags)
			case 3:
				canvas.TextRect2(&r, data.Status, drawFlags)
			case 4:
				canvas.TextRect2(&r, data.Protocol, drawFlags)
			case 5:
				canvas.TextRect2(&r, data.Auto, drawFlags)
			case 6:
				canvas.TextRect2(&r, data.Note, drawFlags)
			}
		}
	})
	f.BotListView.SetOnClick(func(sender vcl.IObject) {
		if len(TempBotData) == 5 {
			return
		}
		var t TTempItem
		t.Status = "123"
		t.QQ = "123"
		t.Protocol = "123"
		t.Note = "2135"
		t.NickName = "5412"
		TempBotData = append(TempBotData, t)
		BotForm.BotListView.Items().SetCount(int32(len(TempBotData))) //   必须主动的设置Virtual List的行数
	})
}

func (f *TBotForm) GetSubItemRect(hwndLV types.HWND, iItem, iSubItem int32) (ret types.TRect) {
	win.ListView_GetSubItemRect(hwndLV, iItem, iSubItem, win.LVIR_LABEL, &ret)
	return
}

//
//func (f *TBotForm) OnBotListViewAdvancedCustomDrawSubItem(sender *vcl.TListView, item *vcl.TListItem, subItem int32, state types.TCustomDrawState, stage types.TCustomDrawStage, defaultDraw *bool)  {
//
//}

func (f *TBotForm) OnFormDestroy(sender vcl.IObject) {
	if f.TempIco != nil {
		f.TempIco.Free()
	}
}
