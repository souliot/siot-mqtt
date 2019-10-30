package protocol

import "bytes"

type PubackProperties struct {
	ReasonString string
	UserProperty map[string][]interface{}
}

// 下面是 具体协议编码解码
type Puback struct {
	FixedHeader      *FixedHeader
	PacketIdentifier uint16
	ReasonCode       ReasonCode
	PubackProperties *PubackProperties
}

func (m *Puback) Encode(buf *bytes.Buffer) (err error) {
	err = m.FixedHeader.Encode(buf)
	err = setUint16(m.PacketIdentifier, buf)
	if m.ReasonCode == 0 && m.PubackProperties != nil {
		err = setUint8(uint8(m.ReasonCode), buf)

		var cp Properties = m.PubackProperties
		err = Encode(&cp, buf)
	}
	return
}

func (m *Puback) Decode(b []byte) {
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
	properties = &PubackProperties{}
	Decode(&properties, b, &p)
	m.PubackProperties = properties.(*PubackProperties)

	return

}
