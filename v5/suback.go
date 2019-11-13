package v5

import (
	"bytes"

	util "github.com/souliot/siot-mqtt/util"
)

type SubAckProperties struct {
	ReasonString string
	UserProperty map[string][]interface{}
}

type SubAckPayload struct {
	ReasonCodes []util.ReasonCode
}

func (m *SubAckPayload) Encode(buf *bytes.Buffer) (err error) {
	if m == nil {
		return
	}
	for _, v := range m.ReasonCodes {
		err = util.SetUint8(uint8(v), buf)
	}
	return
}

func (m *SubAckPayload) Decode(b []byte, p *int) {
	if len(b) <= *p {
		return
	}
	rcs := make([]util.ReasonCode, 0)
	for *p < len(b) {
		rc := util.ReasonCode(util.GetUint8(b, p))

		rcs = append(rcs, rc)
	}
	m.ReasonCodes = rcs

	return
}

// 下面是 具体协议编码解码
type SubAck struct {
	FixedHeader      *FixedHeader
	PacketIdentifier uint16
	SubAckProperties *SubAckProperties
	SubAckPayload    *SubAckPayload
}

func (m *SubAck) Encode(buf *bytes.Buffer) (err error) {
	err = m.FixedHeader.Encode(buf)
	err = util.SetUint16(m.PacketIdentifier, buf)

	var cp Properties = m.SubAckProperties
	err = Encode(&cp, buf)

	err = m.SubAckPayload.Encode(buf)
	return
}

func (m *SubAck) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header

	m.PacketIdentifier = util.GetUint16(b, &p)

	var properties Properties
	properties = &SubAckProperties{}
	Decode(&properties, b, &p)
	m.SubAckProperties = properties.(*SubAckProperties)

	sp := &SubAckPayload{}
	sp.Decode(b, &p)
	m.SubAckPayload = sp

	return

}
