package v5

import (
	"bytes"

	util "github.com/souliot/siot-mqtt/util"
)

type UnsubAckProperties struct {
	ReasonString string
	UserProperty map[string][]interface{}
}

type UnsubAckPayload struct {
	ReasonCodes []util.ReasonCode
}

func (m *UnsubAckPayload) Encode(buf *bytes.Buffer) (err error) {
	if m == nil {
		return
	}
	for _, v := range m.ReasonCodes {
		err = util.SetUint8(uint8(v), buf)
	}
	return
}

func (m *UnsubAckPayload) Decode(b []byte, p *int) {
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
type UnsubAck struct {
	FixedHeader        *FixedHeader
	PacketIdentifier   uint16
	UnsubAckProperties *UnsubAckProperties
	UnsubAckPayload    *UnsubAckPayload
}

func (m *UnsubAck) Encode(buf *bytes.Buffer) (err error) {
	bt := new(bytes.Buffer)
	err = util.SetUint16(m.PacketIdentifier, bt)

	var cp Properties
	if m.UnsubAckProperties != nil {
		cp = m.UnsubAckProperties
	} else {
		cp = new(UnsubAckProperties)
	}
	err = Encode(&cp, bt)

	err = m.UnsubAckPayload.Encode(bt)

	m.FixedHeader.RemainingLength = uint32(bt.Len())
	err = m.FixedHeader.Encode(buf)
	buf.Write(bt.Bytes())
	return
}

func (m *UnsubAck) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header

	m.PacketIdentifier = util.GetUint16(b, &p)

	var properties Properties
	properties = &UnsubAckProperties{}
	Decode(&properties, b, &p)
	m.UnsubAckProperties = properties.(*UnsubAckProperties)

	sp := &UnsubAckPayload{}
	sp.Decode(b, &p)
	m.UnsubAckPayload = sp

	return

}
