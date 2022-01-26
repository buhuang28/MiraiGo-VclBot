package bot

import (
	"MiraiGo-VclBot/util"
	"encoding/json"
	"fmt"
	"github.com/ying32/govcl/vcl"
	"strconv"
	"strings"
	"time"
)

func AutoLogin() {
	dir := util.ReadDir(QQINFOROOTPATH)
	for _, v := range dir {
		if dir == nil {
			return
		}
		if strings.Contains(v.Name(), QQINFOSKIN) {
			var qqInfo QQInfo
			fileByte := util.ReadFileByte(QQINFOPATH + v.Name())
			err := json.Unmarshal(fileByte, &qqInfo)
			if err != nil {
				fmt.Println("反序列化失败:", err)
				continue
			}
			if qqInfo.AutoLogin {
				qqInfo.Login()
				time.Sleep(time.Second * 2)
			} else {
				botLock.Lock()
				index, ok := botIndexMap[qqInfo.QQ]
				if !ok {
					index = botIndexStart
					botIndexStart++
				}
				botIndexMap[qqInfo.QQ] = index
				var botData TTempItem
				botData.IconIndex = int32(index)
				avatarUrl := AvatarUrlPre + strconv.FormatInt(qqInfo.QQ, 10)
				bytes, err2 := util.GetBytes(avatarUrl)
				if err2 != nil {
					fmt.Println(err2)
				} else {
					pic := vcl.NewPicture()
					pic.LoadFromBytes(bytes)
					BotForm.Icons.AddSliced(pic.Bitmap(), 1, 1)
					pic.Free()
				}
				BotForm.BotListView.SetStateImages(BotForm.Icons)
				botData.QQ = strconv.FormatInt(qqInfo.QQ, 10)
				botData.Protocol = GetProtocol(qqInfo.ClientProtocol)
				botData.Status = "离线"
				botData.NickName = ""
				if qqInfo.AutoLogin {
					botData.Auto = "√"
				} else {
					botData.Auto = "X"
				}
				botData.Note = "离线"
				TempBotData = append(TempBotData, botData)
				BotForm.BotListView.Items().SetCount(int32(len(TempBotData))) //   必须主动的设置Virtual List的行数
				botLock.Unlock()
			}
		}
	}
}
