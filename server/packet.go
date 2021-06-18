package server

import (
	"bytes"

	"github.com/souliot/fetcp"
	logs "github.com/souliot/siot-log"
	"github.com/souliot/siot-mqtt/base"
	util "github.com/souliot/siot-mqtt/util"
)

type Packet struct {
	MsgType       util.MessageType
	ProtocolLevel uint8
	Message       base.Message
}

func (p *Packet) Serialize() []byte {
	buf := new(bytes.Buffer)
	err := p.Message.Encode(buf)
	if err != nil {
		logs.Error("Packet Serialize err:", err)
		return nil
	}
	logs.Info(buf.Bytes())
	return buf.Bytes()
}

func NewPacket(b []byte, c *fetcp.Conn) (p *Packet, err error) {
	p = &Packet{}
	msgType := base.GetMessageType(b)
	logs.Info(msgType)
	var protocolLevel uint8
	if msgType == util.MsgConnect {
		protocolLevel = base.GetProtocolLevel(b)
	} else {
		protocolLevel = c.GetExtraData().(*ExtraData).ProtocolLevel
	}
	msg, err := base.NewMessage(b, protocolLevel, msgType)
	if err != nil {
		logs.Error("New Packet Message err:", err)
		return
	}
	p.ProtocolLevel = protocolLevel
	p.MsgType = msgType
	p.Message = msg
	return
}
