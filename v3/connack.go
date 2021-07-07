package v3

import (
	"bytes"

	util "github.com/souliot/siot-mqtt/util"
)

type ConnectAcknowledgeFlags struct {
	SessionPresentFlag bool
}

func (m *ConnectAcknowledgeFlags) Encode(buf *bytes.Buffer) (err error) {
	if m == nil {
		return
	}
	val := util.BoolToByte(m.SessionPresentFlag)
	err = buf.WriteByte(val)
	return
}

func (m *ConnectAcknowledgeFlags) Decode(b []byte, p *int) {
	if len(b) <= *p {
		return
	}
	byte1 := b[*p]
	*p += 1
	*m = ConnectAcknowledgeFlags{
		SessionPresentFlag: byte1&0x01 > 0,
	}
	return
}

// 下面是 具体协议编码解码
type ConnAck struct {
	FixedHeader             *FixedHeader
	ConnectAcknowledgeFlags *ConnectAcknowledgeFlags
	ReasonCode              util.ReasonCode
}

func (m *ConnAck) Encode(buf *bytes.Buffer) (err error) {
	bt := new(bytes.Buffer)
	err = m.ConnectAcknowledgeFlags.Encode(bt)
	err = util.SetUint8(uint8(m.ReasonCode), bt)

	m.FixedHeader.RemainingLength = uint32(bt.Len())
	err = m.FixedHeader.Encode(buf)
	buf.Write(bt.Bytes())
	return
}

func (m *ConnAck) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header

	connectAcknowledgeFlags := &ConnectAcknowledgeFlags{}
	connectAcknowledgeFlags.Decode(b, &p)
	m.ConnectAcknowledgeFlags = connectAcknowledgeFlags

	m.ReasonCode = util.ReasonCode(util.GetUint8(b, &p))
	return

}
