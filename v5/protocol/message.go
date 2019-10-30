package protocol

import (
	"bytes"
)

type MessageType uint8

const (
	MsgConnect     = MessageType(iota + 1) //客户端到服务端	客户端请求连接服务端
	MsgConnAck                             //服务端到客户端	连接报文确认
	MsgPublish                             //两个方向都允许	发布消息
	MsgPubAck                              //两个方向都允许	QoS 1 消息发布收到确认
	MsgPubRec                              //两个方向都允许	发布收到(保证交付第一步)
	MsgPubRel                              //两个方向都允许	发布释放(保证交付第二步)
	MsgPubComp                             //两个方向都允许	QoS 2 消息发布完成(保证交付第三步)
	MsgSubscribe                           //客户端到服务端	客户端订阅请求
	MsgSubAck                              //服务端到客户端	订阅请求报文确认
	MsgUnsubscribe                         //客户端到服务端	客户端取消订阅请求
	MsgUnsubAck                            //服务端到客户端	取消订阅报文确认
	MsgPingReq                             //客户端到服务端	心跳请求
	MsgPingResp                            //服务端到客户端	心跳响应
	MsgDisconnect                          //两个方向都允许	断开连接通知
	MsgAuth                                //两个方向都允许	认证信息交换
)

var MessageType_name = map[MessageType]string{
	1:  "MsgConnect",
	2:  "MsgConnAck",
	3:  "MsgPublish",
	4:  "MsgPubAck",
	5:  "MsgPubRec",
	6:  "MsgPubRel",
	7:  "MsgPubComp",
	8:  "MsgSubscribe",
	9:  "MsgSubAck",
	10: "MsgUnsubscribe",
	11: "MsgUnsubAck",
	12: "MsgPingReq",
	13: "MsgPingResp",
	14: "MsgDisconnect",
	15: "MsgAuth",
}

type QosLevel uint8

const (
	QosAtMostOnce = QosLevel(iota)
	QosAtLeastOnce
	QosExactlyOnce
)

var QosLevel_name = map[QosLevel]string{
	0: "QosAtMostOnce",
	1: "QosAtLeastOnce",
	2: "QosExactlyOnce",
}

type ReasonCode uint8

var ReasonCodeDesc = map[ReasonCode]string{
	0:   "成功",
	1:   "授权的QoS 1",
	2:   "授权的QoS 2",
	4:   "包含遗嘱的断开",
	16:  "无匹配订阅",
	17:  "订阅不存在",
	24:  "继续认证",
	25:  "重新认证",
	128: "未指明的错误",
	129: "无效报文",
	130: "协议错误",
	131: "实现错误",
	132: "协议版本不支持",
	133: "客户标识符无效",
	134: "用户名密码错误",
	135: "未授权",
	136: "服务端不可用",
	137: "服务端正忙",
	138: "禁止",
	139: "服务端关闭中",
	140: "无效的认证方法",
	141: "保活超时",
	142: "会话被接管",
	143: "主题过滤器无效",
	144: "主题名无效",
	145: "报文标识符已被占用",
	146: "报文标识符无效",
	147: "接收超出最大数量",
	148: "主题别名无效",
	149: "报文过长",
	150: "消息太过频繁",
	151: "超出配额",
	152: "管理行为",
	153: "载荷格式无效",
	154: "不支持保留",
	155: "不支持的QoS等级",
	156: "(临时)使用其他服务端",
	157: "服务端已(永久)移动",
	158: "不支持共享订阅",
	159: "超出连接速率限制",
	160: "最大连接时间",
	161: "不支持订阅标识符",
	162: "不支持通配符订阅",
}

// Message is the interface that all MQTT messages implement.
type Message interface {
	Encode(buf *bytes.Buffer) (err error)

	Decode(b []byte)
}
