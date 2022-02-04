# MiraiGo-VclBot
基于MiraiGo做出来的多Q机器人......框架？

开发文档？开发文档还在鸽，咕咕咕。
Event事件文档自己开个websocket-server试试，目前咕了，下次一定补上。
Api文档目前也咕了大部分，就写了一小部分Api，下次一定补上


## API

##### 发送私聊

```json
{
    "msg_type":1,
    "bot_id":123,
    "message":"消息内容",
    "user_id":123,
}
```

|   参数   |  类型  | 必填 |      说明      |
| :------: | :----: | :--: | :------------: |
| msg_type |  int   |  是  |  消息类型：1   |
|  bot_id  |  int   |  是  |   机器人QQ号   |
| message  | string |  是  |    消息内容    |
| user_id  |  int   |  是  | 消息接收人QQ号 |



##### 发送群消息

```json
{
    "msg_type":2,
    "bot_id":123,
    "message":"我是内容",
    "group_id":123,
    "message_id":2
}
```

|    参数    |  类型  | 必填 |        说明         |
| :--------: | :----: | :--: | :-----------------: |
|  msg_type  |  int   |  是  |     消息类型：2     |
|   bot_id   |  int   |  是  |     机器人QQ号      |
|  message   | string |  是  |      消息内容       |
| message_id |  int   |  是  | 消息id,可以用于撤回 |
|  grouo_id  |  int   |  是  |     发送群群号      |



##### 撤回群消息

```json
{
    "msg_type":4,
    "bot_id":123,
    "group_id":0,
    "user_id":0,
    "message_id":4,
    "internal_id":0
}
```

|    参数     | 类型 | 必填 |                       说明                       |
| :---------: | :--: | :--: | :----------------------------------------------: |
|  msg_type   | int  |  是  |                   消息类型：4                    |
|   bot_id    | int  |  是  |                    机器人QQ号                    |
|  group_id   | int  |  否  |    撤回自己的消息则不填，否则填要撤回的群群号    |
|   user_id   | int  |  否  |     撤回自己的消息则不填，否则填要撤回的QQ号     |
| message_id  | int  |  是  |                      消息id                      |
| internal_id | int  |  否  | 撤回自己的则不填，否则按照event的internal_id传入 |



##### 	获取群列表

```json
{
    "msg_type":8,
    "bot_id":123
}
```

|   参数   | 类型 | 必填 |    说明     |
| :------: | :--: | :--: | :---------: |
| msg_type | int  |  是  | 消息类型：8 |
|  bot_id  | int  |  是  | 机器人QQ号  |



##### 踢出群成员

```json
{
    "msg_type":9,
    "bot_id":123,
	"group_id":123,
	"user_id":123
}
```

|   参数   | 类型 | 必填 |    说明     |
| :------: | :--: | :--: | :---------: |
| msg_type | int  |  是  | 消息类型：9 |
|  bot_id  | int  |  是  | 机器人QQ号  |
| group_id | int  |  是  |    群号     |
| user_id  | int  |  是  |  成员QQ号   |



##### 禁言群成员

```json
{
    "msg_type":10,
    "bot_id":123,
	"group_id":123,
	"user_id":123,
    "time":180
}
```

|   参数   | 类型 | 必填 |     说明     |
| :------: | :--: | :--: | :----------: |
| msg_type | int  |  是  | 消息类型:10  |
|  bot_id  | int  |  是  |  机器人QQ号  |
| group_id | int  |  是  |     群号     |
| user_id  | int  |  是  |   成员QQ号   |
|   time   | int  |  是  | 禁言时间(秒) |
