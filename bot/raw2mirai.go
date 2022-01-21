package bot

import (
	"encoding/xml"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	log "github.com/sirupsen/logrus"
	"html"
	"regexp"
	"strings"
)

type Node struct {
	XMLName xml.Name
	Attr    []xml.Attr `xml:",any,attr"`
}

var re = regexp.MustCompile("<[\\s\\S]+?/>")

func RawMsgToMiraiMsg(cli *client.QQClient, str string) []message.IMessageElement {
	containReply := false
	var node Node
	textList := re.Split(str, -1)
	codeList := re.FindAllString(str, -1)
	elemList := make([]message.IMessageElement, 0)
	for len(textList) > 0 || len(codeList) > 0 {
		if len(textList) > 0 && strings.HasPrefix(str, textList[0]) {
			text := textList[0]
			textList = textList[1:]
			str = str[len(text):]
			if text != "" {
				elemList = append(elemList, message.NewText(text))
			}
		}
		if len(codeList) > 0 && strings.HasPrefix(str, codeList[0]) {
			code := codeList[0]
			codeList = codeList[1:]
			str = str[len(code):]
			err := xml.Unmarshal([]byte(code), &node)
			if err != nil {
				elemList = append(elemList, message.NewText(code))
				continue
			}
			attrMap := make(map[string]string)
			for _, attr := range node.Attr {
				attrMap[attr.Name.Local] = html.UnescapeString(attr.Value)
			}
			switch node.XMLName.Local {
			case "at":
				elemList = append(elemList, ProtoAtToMiraiAt(attrMap))
			case "poke":
				elemList = append(elemList, ProtoPokeToMiraiPoke(attrMap))
			case "img", "image":
				elemList = append(elemList, ProtoImageToMiraiImage(attrMap)) // TODO 为了兼容我的旧代码偷偷加的
			case "face":
				elemList = append(elemList, ProtoFaceToMiraiFace(attrMap))
			case "share":
				elemList = append(elemList, ProtoShareToMiraiShare(attrMap))
			case "voice":
				elemList = append(elemList, ProtoVoiceToMiraiVoice(attrMap))
			case "record":
				elemList = append(elemList, ProtoVoiceToMiraiVoice(attrMap))
			case "text":
				elemList = append(elemList, ProtoTextToMiraiText(attrMap))
			case "light_app":
				elemList = append(elemList, ProtoLightAppToMiraiLightApp(attrMap))
			case "service":
				elemList = append(elemList, ProtoServiceToMiraiService(attrMap))
			case "reply":
				if replyElement := ProtoReplyToMiraiReply(attrMap); replyElement != nil && !containReply {
					containReply = true
					elemList = append([]message.IMessageElement{replyElement}, elemList...)
				}
			case "sleep":
				ProtoSleep(attrMap)
			case "tts":
				elemList = append(elemList, ProtoTtsToMiraiTts(cli, attrMap))
			case "video":
				elemList = append(elemList, ProtoVideoToMiraiVideo(cli, attrMap))
			case "gift":
				//elemList = append(elemList, ProtoGiftToMiraiGift(cli, attrMap))
			default:
				log.Warnf("不支持的类型 %s", code)
				elemList = append(elemList, message.NewText(code))
			}
		}
	}
	return elemList
}
