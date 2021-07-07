package v3

import (
	"bytes"

	util "github.com/souliot/siot-mqtt/util"
)

// 下面是 具体协议编码解码
type PubAck struct {
	FixedHeader      *FixedHeader
	PacketIdentifier uint16
}

func (m *PubAck) Encode(buf *bytes.Buffer) (err error) {
	m.FixedHeader.RemainingLength = 2
	err = m.FixedHeader.Encode(buf)
	err = util.SetUint16(m.PacketIdentifier, buf)
	return
}

func (m *PubAck) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header

	m.PacketIdentifier = util.GetUint16(b, &p)

	return

}
