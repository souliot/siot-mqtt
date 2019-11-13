package v5

import (
	"bytes"

	util "github.com/souliot/siot-mqtt/util"
)

type PubCompProperties struct {
	ReasonString string
	UserProperty map[string][]interface{}
}

// 下面是 具体协议编码解码
type PubComp struct {
	FixedHeader       *FixedHeader
	PacketIdentifier  uint16
	ReasonCode        util.ReasonCode
	PubCompProperties *PubCompProperties
}

func (m *PubComp) Encode(buf *bytes.Buffer) (err error) {
	err = m.FixedHeader.Encode(buf)
	err = util.SetUint16(m.PacketIdentifier, buf)
	if m.ReasonCode == 0 && m.PubCompProperties == nil {
		return
	}
	err = util.SetUint8(uint8(m.ReasonCode), buf)

	var cp Properties = m.PubCompProperties
	err = Encode(&cp, buf)
	return
}

func (m *PubComp) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header

	m.PacketIdentifier = util.GetUint16(b, &p)

	if header.RemainingLength == 2 {
		return
	}

	m.ReasonCode = util.ReasonCode(util.GetUint8(b, &p))

	if header.RemainingLength < 4 {
		return
	}

	var properties Properties
	properties = &PubCompProperties{}
	Decode(&properties, b, &p)
	m.PubCompProperties = properties.(*PubCompProperties)

	return

}
