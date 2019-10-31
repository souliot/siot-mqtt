package protocol

import "bytes"

type SubackProperties struct {
	ReasonString string
	UserProperty map[string][]interface{}
}

type SubackPayload struct {
	ReasonCodes []ReasonCode
}

func (m *SubackPayload) Encode(buf *bytes.Buffer) (err error) {
	if m == nil {
		return
	}
	for _, v := range m.ReasonCodes {
		err = setUint8(uint8(v), buf)
	}
	return
}

func (m *SubackPayload) Decode(b []byte, p *int) {
	if len(b) <= *p {
		return
	}
	rcs := make([]ReasonCode, 0)
	for *p < len(b) {
		rc := ReasonCode(getUint8(b, p))

		rcs = append(rcs, rc)
	}
	m.ReasonCodes = rcs

	return
}

// 下面是 具体协议编码解码
type Suback struct {
	FixedHeader      *FixedHeader
	PacketIdentifier uint16
	SubackProperties *SubackProperties
	SubackPayload    *SubackPayload
}

func (m *Suback) Encode(buf *bytes.Buffer) (err error) {
	err = m.FixedHeader.Encode(buf)
	err = setUint16(m.PacketIdentifier, buf)

	var cp Properties = m.SubackProperties
	err = Encode(&cp, buf)

	err = m.SubackPayload.Encode(buf)
	return
}

func (m *Suback) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header

	m.PacketIdentifier = getUint16(b, &p)

	var properties Properties
	properties = &SubackProperties{}
	Decode(&properties, b, &p)
	m.SubackProperties = properties.(*SubackProperties)

	sp := &SubackPayload{}
	sp.Decode(b, &p)
	m.SubackPayload = sp

	return

}
