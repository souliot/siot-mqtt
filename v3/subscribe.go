package v3

import (
	"bytes"

	util "github.com/souliot/siot-mqtt/util"
)

type SubscriptionOptions struct {
	QosLevel                   util.QosLevel
	NoLocal, RetainAsPublished bool
	RetainHandling             uint8
	Reserved                   uint8
}

func (m *SubscriptionOptions) Encode(buf *bytes.Buffer) (err error) {
	if m == nil {
		return
	}
	val := byte(m.Reserved) << 6
	val |= byte(m.RetainHandling) << 4
	val |= util.BoolToByte(m.RetainAsPublished) << 3
	val |= util.BoolToByte(m.NoLocal) << 2
	val |= byte(m.QosLevel)

	err = buf.WriteByte(val)
	return
}

func (m *SubscriptionOptions) Decode(b []byte, p *int) {
	if len(b) <= *p {
		return
	}
	byte1 := b[*p]
	*p += 1
	*m = SubscriptionOptions{
		QosLevel:          util.QosLevel(byte1 & 0x03),
		NoLocal:           byte1&0x04 > 0,
		RetainAsPublished: byte1&0x08 > 0,
		RetainHandling:    uint8(byte1 & 0x30 >> 4),
		Reserved:          uint8(byte1 & 0xc0 >> 6),
	}
	return
}

type SubscribeTopic struct {
	TopicFilter         string
	SubscriptionOptions *SubscriptionOptions
}

type SubscribePayload struct {
	SubscribeTopics []*SubscribeTopic
}

func (m *SubscribePayload) Encode(buf *bytes.Buffer) (err error) {
	if m == nil {
		return
	}
	for _, v := range m.SubscribeTopics {
		err = util.SetString(v.TopicFilter, buf)
		err = v.SubscriptionOptions.Encode(buf)
	}
	return
}

func (m *SubscribePayload) Decode(b []byte, p *int) {
	if len(b) <= *p {
		return
	}
	sts := make([]*SubscribeTopic, 0)
	for *p < len(b) {
		st := &SubscribeTopic{}
		st.TopicFilter = util.GetString(b, p)

		so := &SubscriptionOptions{}
		so.Decode(b, p)
		st.SubscriptionOptions = so

		sts = append(sts, st)
	}
	m.SubscribeTopics = sts

	return
}

func (m *SubscribePayload) Merger(s *SubscribePayload) {
	for _, v := range s.SubscribeTopics {
		if !m.HasSubscribeTopic(v) {
			m.SubscribeTopics = append(m.SubscribeTopics, v)
		}
	}
	return
}

func (m *SubscribePayload) Remove(s *UnsubscribePayload) {
	for _, v := range s.UnsubscribeTopics {
		m.RemoveUnsubscribeTopic(v)
	}
	return
}

func (m *SubscribePayload) HasSubscribeTopic(s *SubscribeTopic) bool {
	for _, v := range m.SubscribeTopics {
		if v.TopicFilter == s.TopicFilter && *(v.SubscriptionOptions) == *(s.SubscriptionOptions) {
			return true
		}
	}
	return false
}

func (m *SubscribePayload) HasPublish(s *Publish) bool {
	for _, v := range m.SubscribeTopics {
		if v.TopicFilter == s.TopicName && v.SubscriptionOptions.QosLevel == s.FixedHeader.QosLevel {
			return true
		}
	}
	return false
}

func (m *SubscribePayload) RemoveUnsubscribeTopic(s *UnsubscribeTopic) {
	for i, v := range m.SubscribeTopics {
		if v.TopicFilter == s.TopicFilter {
			m.SubscribeTopics = append(m.SubscribeTopics[:i], m.SubscribeTopics[i+1:]...)
			m.RemoveUnsubscribeTopic(s)
			return
		}
	}
	return
}

// 下面是 具体协议编码解码
type Subscribe struct {
	FixedHeader      *FixedHeader
	PacketIdentifier uint16
	SubscribePayload *SubscribePayload
}

func (m *Subscribe) Encode(buf *bytes.Buffer) (err error) {
	err = m.FixedHeader.Encode(buf)
	err = util.SetUint16(m.PacketIdentifier, buf)

	err = m.SubscribePayload.Encode(buf)
	return
}

func (m *Subscribe) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header

	m.PacketIdentifier = util.GetUint16(b, &p)

	sp := &SubscribePayload{}
	sp.Decode(b, &p)
	m.SubscribePayload = sp

	return

}
