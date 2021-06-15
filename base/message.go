package base

import (
	"bytes"
	"errors"

	util "github.com/souliot/siot-mqtt/util"
	v3 "github.com/souliot/siot-mqtt/v3"
	v5 "github.com/souliot/siot-mqtt/v5"
)

// Message is the interface that all MQTT messages implement.
type Message interface {
	Encode(buf *bytes.Buffer) (err error)

	Decode(b []byte)
}

func GetMessageType(b []byte) util.MessageType {
	p := 0
	fh := &v5.FixedHeader{}
	fh.Decode(b, &p)
	return fh.MsgType
}

func GetProtocolLevel(b []byte) uint8 {
	p := 0
	fh := &v5.FixedHeader{}
	fh.Decode(b, &p)

	util.GetString(b, &p)
	ProtocolLevel := util.GetUint8(b, &p)
	return ProtocolLevel
}

func NewMessage(b []byte, protocolLevel uint8, msgType util.MessageType) (msg Message, err error) {
	if protocolLevel == 4 {
		switch msgType {
		case util.MsgConnect:
			msg = &v3.Connect{}
		case util.MsgConnAck:
			msg = &v3.ConnAck{}
		case util.MsgPublish:
			msg = &v3.Publish{}
		case util.MsgPubAck:
			msg = &v3.PubAck{}
		case util.MsgPubRec:
			msg = &v3.PubRec{}
		case util.MsgPubRel:
			msg = &v3.PubRel{}
		case util.MsgPubComp:
			msg = &v3.PubComp{}
			// case util.MsgSubscribe:
			// 	msg = &v3.Subscribe{}
			// case util.MsgSubAck:
			// 	msg = &v3.SubAck{}
			// case util.MsgUnsubscribe:
			// 	msg = &v3.Unsubscribe{}
			// case util.MsgUnsubAck:
			// 	msg = &v3.UnsubAck{}
			// case util.MsgPingReq:
			// 	msg = &v3.PingReq{}
			// case util.MsgPingResp:
			// 	msg = &v3.PingResp{}
			// case util.MsgDisconnect:
			// 	msg = &v3.Disconnect{}
			// case util.MsgAuth:
			// 	msg = &v3.Auth{}
		}
	} else if protocolLevel == 5 {
		switch msgType {
		case util.MsgConnect:
			msg = &v5.Connect{}
		case util.MsgConnAck:
			msg = &v5.ConnAck{}
		case util.MsgPublish:
			msg = &v5.Publish{}
		case util.MsgPubAck:
			msg = &v5.PubAck{}
		case util.MsgPubRec:
			msg = &v5.PubRec{}
		case util.MsgPubRel:
			msg = &v5.PubRel{}
		case util.MsgPubComp:
			msg = &v5.PubComp{}
		case util.MsgSubscribe:
			msg = &v5.Subscribe{}
		case util.MsgSubAck:
			msg = &v5.SubAck{}
		case util.MsgUnsubscribe:
			msg = &v5.Unsubscribe{}
		case util.MsgUnsubAck:
			msg = &v5.UnsubAck{}
		case util.MsgPingReq:
			msg = &v5.PingReq{}
		case util.MsgPingResp:
			msg = &v5.PingResp{}
		case util.MsgDisconnect:
			msg = &v5.Disconnect{}
		case util.MsgAuth:
			msg = &v5.Auth{}
		}
	} else {
		err = errors.New("Not Support Protocol Level")
	}
	msg.Decode(b)
	return
}
