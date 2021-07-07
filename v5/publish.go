package v5

import (
	"bytes"

	util "github.com/souliot/siot-mqtt/util"
)

type PublishProperties struct {
	PayloadFormatIndicator uint8
	MessageExpiryInterval  uint32
	TopicAlias             uint16
	ResponseTopic          string
	CorrelationData        []byte
	UserProperty           map[string][]interface{}
	SubscriptionIdentifier uint32
	ContentType            string
}

type Publish struct {
	FixedHeader       *FixedHeader
	TopicName         string
	PacketIdentifier  uint16
	PublishProperties *PublishProperties
	Payload           []byte
}

func (m *Publish) Encode(buf *bytes.Buffer) (err error) {
	bt := new(bytes.Buffer)
	err = util.SetString(m.TopicName, bt)
	if m.FixedHeader.QosLevel != util.QosAtMostOnce {
		err = util.SetUint16(m.PacketIdentifier, bt)
	}

	var cp Properties
	if m.PublishProperties != nil {
		cp = m.PublishProperties
	} else {
		cp = new(PublishProperties)
	}
	err = Encode(&cp, bt)

	err = util.SetBytesNoLen(m.Payload, bt)

	m.FixedHeader.RemainingLength = uint32(bt.Len())
	err = m.FixedHeader.Encode(buf)
	buf.Write(bt.Bytes())
	return
}

func (m *Publish) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header

	l1 := p

	m.TopicName = util.GetString(b, &p)
	if m.FixedHeader.QosLevel != util.QosAtMostOnce {
		m.PacketIdentifier = util.GetUint16(b, &p)
	}

	var properties Properties
	properties = &PublishProperties{}
	Decode(&properties, b, &p)
	m.PublishProperties = properties.(*PublishProperties)

	l2 := p

	m.Payload = util.GetBytesNoLen(b, &p, int(header.RemainingLength)-l2+l1)
	return

}
