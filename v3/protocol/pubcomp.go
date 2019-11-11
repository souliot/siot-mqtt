package protocol

import "bytes"

type PubcompProperties struct {
	ReasonString string
	UserProperty map[string][]interface{}
}

// 下面是 具体协议编码解码
type Pubcomp struct {
	FixedHeader       *FixedHeader
	PacketIdentifier  uint16
	ReasonCode        ReasonCode
	PubcompProperties *PubcompProperties
}

func (m *Pubcomp) Encode(buf *bytes.Buffer) (err error) {
	err = m.FixedHeader.Encode(buf)
	err = setUint16(m.PacketIdentifier, buf)
	if m.ReasonCode == 0 && m.PubcompProperties == nil {
		return
	}
	err = setUint8(uint8(m.ReasonCode), buf)

	var cp Properties = m.PubcompProperties
	err = Encode(&cp, buf)
	return
}

func (m *Pubcomp) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header

	m.PacketIdentifier = getUint16(b, &p)

	if header.RemainingLength == 2 {
		return
	}

	m.ReasonCode = ReasonCode(getUint8(b, &p))

	if header.RemainingLength < 4 {
		return
	}

	var properties Properties
	properties = &PubcompProperties{}
	Decode(&properties, b, &p)
	m.PubcompProperties = properties.(*PubcompProperties)

	return

}
