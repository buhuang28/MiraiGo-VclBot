package bot

import "github.com/ying32/govcl/vcl"

var (
	BotForm     *TBotForm
	TempBotData []TTempItem
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
}
