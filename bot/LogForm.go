package bot

import "github.com/ying32/govcl/vcl"

var (
	LogForm *TLogForm
)

type TLogForm struct {
	*vcl.TForm
	LogListView *vcl.TListView
	RollCheck   *vcl.TCheckBox
}
