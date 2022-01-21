package ws_data

type GMCWSData struct {
	MsgType        int64             `json:"msg_type"`
	BotId          int64             `json:"bot_id"`
	Message        string            `json:"message"`
	GroupId        int64             `json:"group_id"`
	UserId         int64             `json:"user_id"`
	MessageId      int64             `json:"message_id"`
	InternalId     int32             `json:"internal_id"`
	MemberList     []GMCMember       `json:"member_list,omitempty"`
	GroupList      []GMCGroup        `json:"group_list,omitempty"`
	ManageGroup    []int64           `json:"manage_group,omitempty"`
	AllGroupMember GMCAllGroupMember `json:"all_group_member,omitempty"`
	Time           int64             `json:"time,omitempty"`
	BusId          int64             `json:"bus_id,omitempty"`
	FileId         string            `json:"file_id,omitempty"`
	FileFromGroup  int64             `json:"file_from_group,omitempty"`
	FilePath       string            `json:"file,omitempty"`
	NickName       string            `json:"nick_name,omitempty"`
	RequestId      int64             `json:"request_id,omitempty"`
	GroupRequest   int64             `json:"group_request,omitempty"`
	InvitorId      int64             `json:"invitor_id,omitempty"`
	InvitorName    string            `json:"invitor_name,omitempty"`
}

type GMCMember struct {
	QQ         int64  `json:"qq"`
	Permission int64  `json:"permission"`
	NickName   string `json:"nick_name,omitempty"`
	Level      uint16 `json:"lv"`
}

type GMCAllGroupMember struct {
	Data map[int64][]GMCMember `json:"data"`
}

type GMCGroup struct {
	GroupId   int64  `json:"group_id"`
	GroupName string `json:"group_name,omitempty"`
}

const (
	GMC_PRIVATE_MESSAGE  = 1
	GMC_GROUP_MESSAGE    = 2
	GMC_TEMP_MESSAGE     = 3
	GMC_WITHDRAW_MESSAGE = 4
	GMC_ONLINE           = 5
	GMC_OFFLINE          = 6
	//全部群成员
	GMC_ALLGROUPMEMBER = 7
	//群列表
	GMC_GROUP_LIST = 8
	//踢出
	GMC_KICK = 9
	//禁言
	GMC_BAN = 10
	//入群请求
	GMC_GROUP_REQUEST = 11
	//可能是机器人被邀请入群也可能是别人被邀请入群
	GMC_BOT_INVITED = 12
	//成员++
	GMC_MEMBER_ADD = 13
	//群员退群
	GMC_MEMBER_LEAVE = 14
	//群文件
	GMC_GROUP_FILE = 15
)
