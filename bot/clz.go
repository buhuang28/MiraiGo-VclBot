package bot

import (
	"github.com/Mrs4s/MiraiGo/message"
	"io"
)

// 自定义类型

type MyVideoElement struct {
	message.ShortVideoElement
	CoverUrl       string        // 仅用于发送时日志展示
	UploadingCover io.ReadSeeker // 待上传的封面 发送时需要
	UploadingVideo io.ReadSeeker // 待上传的视频 发送时需要
}

type LocalImageElement struct {
	Url      string
	Stream   io.ReadSeeker
	Tp       string // 类型 flash/show
	EffectId int32  // show的特效id，范围40000-40005
}

func (m *LocalImageElement) Type() message.ElementType {
	return message.Image
}

type PokeElement struct {
	Target int64
}

func (g *PokeElement) Type() message.ElementType {
	return message.At
}
