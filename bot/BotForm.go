package bot

import (
	"github.com/ying32/govcl/vcl"
	"sync"
)

var (
	BotForm     *TBotForm
	TempBotData []TTempItem
	//QQ--头像索引
	BotAvatarMap = make(map[int64]int32)
	TempBotLock  sync.Mutex
)

type TTempItem struct {
	IconIndex int32
	QQ        string
	NickName  string
	Status    string
	Protocol  string
	Auto      string
	Note      string
}

type TBotForm struct {
	*vcl.TForm
	BotListView  *vcl.TListView
	Icons        *vcl.TImageList
	TempIco      *vcl.TIcon
	NoSelectMenu *vcl.TPopupMenu
	SelectedMenu *vcl.TPopupMenu
	Avatar       []*vcl.TPngImage
	TForm1Fields
}
