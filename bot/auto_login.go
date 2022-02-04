package bot

import (
	"MiraiGo-VclBot/util"
	"encoding/json"
	"fmt"
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
			var botData TTempItem
			botData.IconIndex = int32(len(TempBotData))
			SetBotAvatar(qqInfo.QQ, int32(len(TempBotData)))
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
			//TempBotData = append(TempBotData, botData)
			AddTempBotData(botData)
			BotForm.BotListView.Items().SetCount(int32(len(TempBotData))) //   必须主动的设置Virtual List的行数
			if qqInfo.AutoLogin {
				qqInfo.Login()
				time.Sleep(time.Second * 2)
			}
		}
	}
}
