package protocol

import "bytes"

type UnSubackProperties struct {
	ReasonString string
	UserProperty map[string][]interface{}
}

type UnSubackPayload struct {
	ReasonCodes []ReasonCode
}

func (m *UnSubackPayload) Encode(buf *bytes.Buffer) (err error) {
	if m == nil {
		return
	}
	for _, v := range m.ReasonCodes {
		err = setUint8(uint8(v), buf)
	}
	return
}

func (m *UnSubackPayload) Decode(b []byte, p *int) {
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
type UnSuback struct {
	FixedHeader        *FixedHeader
	PacketIdentifier   uint16
	UnSubackProperties *UnSubackProperties
	UnSubackPayload    *UnSubackPayload
}

func (m *UnSuback) Encode(buf *bytes.Buffer) (err error) {
	err = m.FixedHeader.Encode(buf)
	err = setUint16(m.PacketIdentifier, buf)

	var cp Properties = m.UnSubackProperties
	err = Encode(&cp, buf)

	err = m.UnSubackPayload.Encode(buf)
	return
}

func (m *UnSuback) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header

	m.PacketIdentifier = getUint16(b, &p)

	var properties Properties
	properties = &UnSubackProperties{}
	Decode(&properties, b, &p)
	m.UnSubackProperties = properties.(*UnSubackProperties)

	sp := &UnSubackPayload{}
	sp.Decode(b, &p)
	m.UnSubackPayload = sp

	return

}
