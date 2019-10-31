package protocol

import "bytes"

type UnSubscribeProperties struct {
	UserProperty map[string][]interface{}
}

type UnSubscribeTopic struct {
	TopicFilter string
}

type UnSubscribePayload struct {
	UnSubscribeTopics []*UnSubscribeTopic
}

func (m *UnSubscribePayload) Encode(buf *bytes.Buffer) (err error) {
	if m == nil {
		return
	}
	for _, v := range m.UnSubscribeTopics {
		err = setString(v.TopicFilter, buf)
	}
	return
}

func (m *UnSubscribePayload) Decode(b []byte, p *int) {
	if len(b) <= *p {
		return
	}
	usts := make([]*UnSubscribeTopic, 0)
	for *p < len(b) {
		ust := &UnSubscribeTopic{}
		ust.TopicFilter = getString(b, p)

		usts = append(usts, ust)
	}
	m.UnSubscribeTopics = usts

	return
}

// 下面是 具体协议编码解码
type UnSubscribe struct {
	FixedHeader           *FixedHeader
	PacketIdentifier      uint16
	UnSubscribeProperties *UnSubscribeProperties
	UnSubscribePayload    *UnSubscribePayload
}

func (m *UnSubscribe) Encode(buf *bytes.Buffer) (err error) {
	err = m.FixedHeader.Encode(buf)
	err = setUint16(m.PacketIdentifier, buf)

	var cp Properties = m.UnSubscribeProperties
	err = Encode(&cp, buf)

	err = m.UnSubscribePayload.Encode(buf)
	return
}

func (m *UnSubscribe) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header

	m.PacketIdentifier = getUint16(b, &p)

	var properties Properties
	properties = &UnSubscribeProperties{}
	Decode(&properties, b, &p)
	m.UnSubscribeProperties = properties.(*UnSubscribeProperties)

	sp := &UnSubscribePayload{}
	sp.Decode(b, &p)
	m.UnSubscribePayload = sp

	return

}
