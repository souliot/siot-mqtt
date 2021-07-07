package v5

import (
	"bytes"

	util "github.com/souliot/siot-mqtt/util"
)

type ConnectAcknowledgeFlags struct {
	SessionPresentFlag bool
}

func (m *ConnectAcknowledgeFlags) Encode(buf *bytes.Buffer) (err error) {
	if m == nil {
		return
	}
	val := util.BoolToByte(m.SessionPresentFlag)
	err = buf.WriteByte(val)
	return
}

func (m *ConnectAcknowledgeFlags) Decode(b []byte, p *int) {
	if len(b) <= *p {
		return
	}
	byte1 := b[*p]
	*p += 1
	*m = ConnectAcknowledgeFlags{
		SessionPresentFlag: byte1&0x01 > 0,
	}
	return
}

type ConnAckProperties struct {
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
type ConnAck struct {
	FixedHeader             *FixedHeader
	ConnectAcknowledgeFlags *ConnectAcknowledgeFlags
	ReasonCode              util.ReasonCode
	ConnAckProperties       *ConnAckProperties
}

func (m *ConnAck) Encode(buf *bytes.Buffer) (err error) {
	bt := new(bytes.Buffer)
	err = m.ConnectAcknowledgeFlags.Encode(bt)
	err = util.SetUint8(uint8(m.ReasonCode), bt)

	var cp Properties
	if m.ConnAckProperties != nil {
		cp = m.ConnAckProperties
	} else {
		cp = new(ConnAckProperties)
	}

	err = Encode(&cp, bt)

	m.FixedHeader.RemainingLength = uint32(bt.Len())
	err = m.FixedHeader.Encode(buf)
	buf.Write(bt.Bytes())
	return
}

func (m *ConnAck) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header

	connectAcknowledgeFlags := &ConnectAcknowledgeFlags{}
	connectAcknowledgeFlags.Decode(b, &p)
	m.ConnectAcknowledgeFlags = connectAcknowledgeFlags

	m.ReasonCode = util.ReasonCode(util.GetUint8(b, &p))

	var properties Properties
	properties = &ConnAckProperties{}
	Decode(&properties, b, &p)
	m.ConnAckProperties = properties.(*ConnAckProperties)

	return

}
