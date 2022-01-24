package bot

var (
	ACCEPT                 int32 = 1
	SEND                   int32 = 2
	ACCEPT_PRIVATE         int32 = 1 //私聊消息
	ACCEPT_GROUP           int32 = 2 //群聊消息
	ACCEPT_TEMP            int32 = 3 //临时消息
	ACCEPT_GROUP_REQUEST   int32 = 4 //有人申请入群
	ACCPET_GROUP_INVITED   int32 = 5 //机器人被邀请入群
	ACCEPT_KICK            int32 = 6 //踢出消息  (群员被踢)
	ACCEPT_MEMBER_INSERT   int32 = 7 //群员增加
	ACCEPT_MEMBER_DECREASE int32 = 8 //群员减少

	SEND_PRIVATE int32 = 9
	SEND_GROUP   int32 = 10
	SEND_TEMP    int32 = 11
)
