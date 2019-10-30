package protocol

import "bytes"

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
	err = m.FixedHeader.Encode(buf)
	err = setString(m.TopicName, buf)
	if m.FixedHeader.QosLevel != QosAtMostOnce {
		err = setUint16(m.PacketIdentifier, buf)
	}

	var pp Properties = m.PublishProperties
	err = Encode(&pp, buf)

	err = setBytesNoLen(m.Payload, buf)

	return
}

func (m *Publish) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header

	m.TopicName = getString(b, &p)
	m.PacketIdentifier = getUint16(b, &p)

	var properties Properties
	properties = &PublishProperties{}
	Decode(&properties, b, &p)
	m.PublishProperties = properties.(*PublishProperties)

	m.Payload = getBytesNoLen(b, &p)
	return

}
