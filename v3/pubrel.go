package v3

import (
	"bytes"

	util "github.com/souliot/siot-mqtt/util"
)

// 下面是 具体协议编码解码
type PubRel struct {
	FixedHeader      *FixedHeader
	PacketIdentifier uint16
}

func (m *PubRel) Encode(buf *bytes.Buffer) (err error) {
	err = m.FixedHeader.Encode(buf)
	err = util.SetUint16(m.PacketIdentifier, buf)

	return
}

func (m *PubRel) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header

	m.PacketIdentifier = util.GetUint16(b, &p)

	return

}
