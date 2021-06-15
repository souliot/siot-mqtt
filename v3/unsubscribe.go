package v3

import (
	"bytes"

	util "github.com/souliot/siot-mqtt/util"
)

type UnsubscribeTopic struct {
	TopicFilter string
}

type UnsubscribePayload struct {
	UnsubscribeTopics []*UnsubscribeTopic
}

func (m *UnsubscribePayload) Encode(buf *bytes.Buffer) (err error) {
	if m == nil {
		return
	}
	for _, v := range m.UnsubscribeTopics {
		err = util.SetString(v.TopicFilter, buf)
	}
	return
}

func (m *UnsubscribePayload) Decode(b []byte, p *int) {
	if len(b) <= *p {
		return
	}
	usts := make([]*UnsubscribeTopic, 0)
	for *p < len(b) {
		ust := &UnsubscribeTopic{}
		ust.TopicFilter = util.GetString(b, p)

		usts = append(usts, ust)
	}
	m.UnsubscribeTopics = usts

	return
}

// 下面是 具体协议编码解码
type Unsubscribe struct {
	FixedHeader        *FixedHeader
	PacketIdentifier   uint16
	UnsubscribePayload *UnsubscribePayload
}

func (m *Unsubscribe) Encode(buf *bytes.Buffer) (err error) {
	err = m.FixedHeader.Encode(buf)
	err = util.SetUint16(m.PacketIdentifier, buf)

	err = m.UnsubscribePayload.Encode(buf)
	return
}

func (m *Unsubscribe) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header

	m.PacketIdentifier = util.GetUint16(b, &p)

	sp := &UnsubscribePayload{}
	sp.Decode(b, &p)
	m.UnsubscribePayload = sp

	return

}
