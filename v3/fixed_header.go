package v3

import (
	"bytes"

	util "github.com/souliot/siot-mqtt/util"
)

type FixedHeader struct {
	DupFlag, Retain bool
	QosLevel        util.QosLevel
	MsgType         util.MessageType
	RemainingLength uint32
}

func (m *FixedHeader) Encode(buf *bytes.Buffer) (err error) {
	if m == nil {
		return
	}
	val := byte(m.MsgType) << 4

	switch util.MessageType_name[m.MsgType] {
	case "MsgPublish":
		val |= (util.BoolToByte(m.DupFlag) << 3)
		val |= byte(m.QosLevel) << 1
		val |= util.BoolToByte(m.Retain)
	case "MsgPubRel", "MsgSubscribe", "MsgUnsubscribe":
		val |= byte(1) << 1
	}
	err = buf.WriteByte(val)

	if err != nil {
		return
	}

	err = util.EncodeLength(m.RemainingLength, buf)

	return
}

func (m *FixedHeader) Decode(b []byte, p *int) {
	if len(b) <= *p {
		return
	}
	byte1 := b[*p]
	*p += 1

	*m = FixedHeader{
		DupFlag:         byte1&0x08 > 0,
		QosLevel:        util.QosLevel(byte1 & 0x06 >> 1),
		Retain:          byte1&0x01 > 0,
		MsgType:         util.MessageType(byte1 & 0xF0 >> 4),
		RemainingLength: util.DecodeLength(b, p),
	}

	return
}
