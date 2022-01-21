package bot

const (
	Ipad         = "IPad"
	AndroidPhone = "安卓手机"
	AndroidWatch = "安卓手表"
	MacOS        = "MacOS"
)

func GetProtocol(protocolType int32) string {
	switch protocolType {
	case 0:
		return Ipad
	case 1:
		return AndroidPhone
	case 2:
		return AndroidWatch
	case 3:
		return MacOS
	}
	return Ipad
}
