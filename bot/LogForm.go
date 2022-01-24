package bot

import "github.com/ying32/govcl/vcl"

var (
	LogForm  *TLogForm
	LogItems []LogItem
)

type LogItem struct {
	BotId       string
	MessageType string
	MessageTime string
	Message     string
}

type TLogForm struct {
	*vcl.TForm
	LogListView *vcl.TListView
	RollCheck   *vcl.TCheckBox
	ClearText   *vcl.TLabel
}
