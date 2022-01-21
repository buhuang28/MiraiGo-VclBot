package bot

import "github.com/ying32/govcl/vcl"

func (f *TLogForm) OnFormCreate(sender vcl.IObject) {
	f.SetHeight(400)
	f.SetWidth(400)
	f.ScreenCenter()

}
