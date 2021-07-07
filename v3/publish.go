package v3

import (
	"bytes"

	util "github.com/souliot/siot-mqtt/util"
)

type Publish struct {
	FixedHeader      *FixedHeader
	TopicName        string
	PacketIdentifier uint16
	Payload          []byte
}

func (m *Publish) Encode(buf *bytes.Buffer) (err error) {
	bt := new(bytes.Buffer)
	err = util.SetString(m.TopicName, bt)

	if m.FixedHeader.QosLevel != util.QosAtMostOnce {
		err = util.SetUint16(m.PacketIdentifier, bt)
	}

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

	l2 := p

	m.Payload = util.GetBytesNoLen(b, &p, int(header.RemainingLength)-l2+l1)

	return
}
