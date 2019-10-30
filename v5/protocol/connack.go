package protocol

import "bytes"

type ConnectAcknowledgeFlags struct {
	SessionPresentFlag bool
}

func (m *ConnectAcknowledgeFlags) Encode(buf *bytes.Buffer) (err error) {
	if m == nil {
		return
	}
	val := boolToByte(m.SessionPresentFlag)
	err = buf.WriteByte(val)
	return
}

func (m *ConnectAcknowledgeFlags) Decode(b []byte, p *int) {
	if len(b) == 0 {
		return
	}
	byte1 := b[*p]
	*p += 1
	*m = ConnectAcknowledgeFlags{
		SessionPresentFlag: byte1&0x01 > 0,
	}
	return
}

type ConnackProperties struct {
	SessionExpiryInterval           uint32
	ReceiveMaximum                  uint16
	MaximumQoS                      uint8
	RetainAvailable                 uint8
	MaximumPacketSize               uint32
	AssignedClientIdentifier        string
	TopicAliasMaximum               uint16
	ReasonString                    string
	UserProperty                    map[string][]interface{}
	WildcardSubscriptionAvailable   uint8
	SubscriptionIdentifierAvailable uint8
	SharedSubscriptionAvailable     uint8
	ServerKeepAlive                 uint16
	ResponseInformation             string
	ServerReference                 string
	AuthenticationMethod            string
	AuthenticationData              []byte
}

// 下面是 具体协议编码解码
type Connack struct {
	FixedHeader             *FixedHeader
	ConnectAcknowledgeFlags *ConnectAcknowledgeFlags
	ReasonCode              ReasonCode
	ConnackProperties       *ConnackProperties
}

func (m *Connack) Encode(buf *bytes.Buffer) (err error) {
	err = m.FixedHeader.Encode(buf)
	err = m.ConnectAcknowledgeFlags.Encode(buf)
	err = setUint8(uint8(m.ReasonCode), buf)

	var cp Properties = m.ConnackProperties
	err = Encode(&cp, buf)

	return
}

func (m *Connack) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header

	connectAcknowledgeFlags := &ConnectAcknowledgeFlags{}
	connectAcknowledgeFlags.Decode(b, &p)
	m.ConnectAcknowledgeFlags = connectAcknowledgeFlags

	m.ReasonCode = ReasonCode(getUint8(b, &p))

	var properties Properties
	properties = &ConnackProperties{}
	Decode(&properties, b, &p)
	m.ConnackProperties = properties.(*ConnackProperties)

	return

}
