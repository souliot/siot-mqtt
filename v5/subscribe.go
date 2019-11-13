package v5

import (
	"bytes"

	util "github.com/souliot/siot-mqtt/util"
)

type SubscribeProperties struct {
	SubscriptionIdentifier uint32
	UserProperty           map[string][]interface{}
}

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

// 下面是 具体协议编码解码
type Subscribe struct {
	FixedHeader         *FixedHeader
	PacketIdentifier    uint16
	SubscribeProperties *SubscribeProperties
	SubscribePayload    *SubscribePayload
}

func (m *Subscribe) Encode(buf *bytes.Buffer) (err error) {
	err = m.FixedHeader.Encode(buf)
	err = util.SetUint16(m.PacketIdentifier, buf)

	var cp Properties = m.SubscribeProperties
	err = Encode(&cp, buf)

	err = m.SubscribePayload.Encode(buf)
	return
}

func (m *Subscribe) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header

	m.PacketIdentifier = util.GetUint16(b, &p)

	var properties Properties
	properties = &SubscribeProperties{}
	Decode(&properties, b, &p)
	m.SubscribeProperties = properties.(*SubscribeProperties)

	sp := &SubscribePayload{}
	sp.Decode(b, &p)
	m.SubscribePayload = sp

	return

}
