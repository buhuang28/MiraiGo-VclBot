package bot

import (
	"MiraiGo-VclBot/util"
	"runtime/debug"
	"time"

	"github.com/Mrs4s/MiraiGo/client"
	log "github.com/sirupsen/logrus"
)

//go:generate go run github.com/a8m/syncmap -o "gen_client_map.go" -pkg bot -name ClientMap "map[int64]*client.QQClient"
//go:generate go run github.com/a8m/syncmap -o "gen_token_map.go" -pkg bot -name TokenMap "map[int64][]byte"
var (
	Clients      ClientMap
	LoginTokens  TokenMap
	AvatarUrlPre = "http://q2.qlogo.cn/headimg_dl?spec=640&dst_uin="
)

func InitLog(cli *client.QQClient) {
	cli.OnLog(func(c *client.QQClient, e *client.LogEvent) {
		switch e.Type {
		case "INFO":
			log.Info("MiraiGo -> " + e.Message)
		case "ERROR":
			log.Error("MiraiGo -> " + e.Message)
			log.Debugf("%+v", string(debug.Stack()))
		case "DEBUG":
			log.Debug("MiraiGo -> " + e.Message)
		}
	})
	cli.OnServerUpdated(func(bot *client.QQClient, e *client.ServerUpdatedEvent) bool {
		log.Infof("收到服务器地址更新通知, 将在下一次重连时应用. ")
		return true // 如果是 false 表示不应用
	})
}

func Login(cli *client.QQClient) (bool, error) {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
	}()
	cli.AllowSlider = true
	rsp, err := cli.Login()
	if err != nil {
		return false, err
	}
	ok := ProcessLoginRsp(cli, rsp)
	return ok, nil
}

func SetRelogin(cli *client.QQClient, retryInterval int, retryCount int) {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
	}()
	LoginTokens.Store(cli.Uin, cli.GenToken())
	cli.OnDisconnected(func(bot *client.QQClient, e *client.ClientDisconnectedEvent) {
		if bot.Online.Load() {
			return
		}
		bot.Disconnect()
		var times = 1

		for IsClientExist(bot.Uin) {
			if bot.Online.Load() {
				log.Warn("Bot已登录")
				return
			}
			if times > retryCount {
				break
			}
			log.Warnf("Bot已离线 (%v)，将在 %v 秒后尝试重连. 重连次数：%v",
				e.Message, retryInterval, times)
			times++
			time.Sleep(time.Second * time.Duration(retryInterval))

			if token, ok := LoginTokens.Load(bot.Uin); ok {
				// 尝试token登录
				if err := bot.TokenLogin(token); err != nil {
					log.Errorf("failed to relogin with token, try to login with password, %+v", err)
					bot.Disconnect()
				} else {
					LoginTokens.Store(bot.Uin, bot.GenToken())
					log.Info("succeed to relogin with token")
					return
				}
			}

			time.Sleep(time.Second)

			// 尝试密码登录
			ok, err := Login(bot)

			if err != nil {
				log.Errorf("重连失败: %v", err)
				bot.Disconnect()
				continue
			}
			if ok {
				LoginTokens.Store(bot.Uin, bot.GenToken())
				log.Info("重连成功")
				return
			}
		}
		log.Errorf("failed to reconnect: 重连次数达到设置的上限值, %+v", cli.Uin)
		ReleaseClient(cli)
	})
}

// ReleaseClient 断开连接并释放资源
func ReleaseClient(cli *client.QQClient) {
	defer func() {
		e := recover()
		if e != nil {
			util.PrintStackTrace(e)
		}
	}()
	BuhuangBotOffline(cli.Uin)
	cli.Release()
	Clients.Delete(cli.Uin) // 必须先删Clients，影响IsClientExist
	LoginTokens.Delete(cli.Uin)
}

func IsClientExist(uin int64) bool {
	_, ok := Clients.Load(uin)
	return ok
}
