package protocol

import "bytes"

type FixedHeader struct {
	DupFlag, Retain bool
	QosLevel        QosLevel
	MsgType         MessageType
	RemainingLength uint32
}

func (m *FixedHeader) Encode(buf *bytes.Buffer) (err error) {
	if m == nil {
		return
	}
	val := byte(m.MsgType) << 4

	switch MessageType_name[m.MsgType] {
	case "MsgPublish":
		val |= (boolToByte(m.DupFlag) << 3)
		val |= byte(m.QosLevel) << 1
		val |= boolToByte(m.Retain)
	case "MsgPubRel", "MsgSubscribe", "MsgUnsubscribe":
		val |= byte(1) << 1
	}
	err = buf.WriteByte(val)

	if err != nil {
		return
	}

	err = encodeLength(m.RemainingLength, buf)

	return
}

func (m *FixedHeader) Decode(b []byte, p *int) {
	if len(b) == 0 {
		return
	}
	byte1 := b[*p]
	*p += 1

	*m = FixedHeader{
		DupFlag:         byte1&0x08 > 0,
		QosLevel:        QosLevel(byte1 & 0x06 >> 1),
		Retain:          byte1&0x01 > 0,
		MsgType:         MessageType(byte1 & 0xF0 >> 4),
		RemainingLength: decodeLength(b, p),
	}

	return
}
