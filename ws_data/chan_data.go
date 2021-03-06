package ws_data

import (
	log "github.com/sirupsen/logrus"
	"sync"
)

var (
	ChanMap     = make(map[int64]chan GMCWSData)
	ChanMapLock sync.Mutex
)

//处理入群请求、机器人被邀请进群的回调结果
func HandleCallBackEvent(data GMCWSData) {
	ChanMapLock.Lock()
	defer func() {
		e := recover()
		if e != nil {
			log.Info(e)
		}
		ChanMapLock.Unlock()
	}()
	ch := ChanMap[data.RequestId]
	ch <- data
}
