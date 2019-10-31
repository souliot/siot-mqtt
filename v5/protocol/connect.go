package protocol

import "bytes"

type ConnectFlags struct {
	UsernameFlag, PasswordFlag, WillRetain bool
	WillQos                                uint8
	WillFlag, CleanStart, Reserved         bool
}

func (m *ConnectFlags) Encode(buf *bytes.Buffer) (err error) {
	if m == nil {
		return
	}
	val := boolToByte(m.UsernameFlag) << 7
	val |= boolToByte(m.PasswordFlag) << 6
	val |= boolToByte(m.WillRetain) << 5
	val |= byte(m.WillQos) << 3
	val |= boolToByte(m.WillFlag) << 2
	val |= boolToByte(m.CleanStart) << 1
	val |= boolToByte(m.Reserved)

	err = buf.WriteByte(val)
	return
}

func (m *ConnectFlags) Decode(b []byte, p *int) {
	if len(b) <= *p {
		return
	}
	byte1 := b[*p]
	*p += 1
	*m = ConnectFlags{
		UsernameFlag: byte1&0x80 > 0,
		PasswordFlag: byte1&0x40 > 0,
		WillRetain:   byte1&0x20 > 0,
		WillQos:      uint8(byte1 & 0x18 >> 3),
		WillFlag:     byte1&0x04 > 0,
		CleanStart:   byte1&0x02 > 0,
		Reserved:     byte1&0x01 > 0,
	}
	return
}

type ConnectProperties struct {
	// If the Session Expiry Interval is absent the value 0 is used.
	// If it is set to 0, or is absent, the Session ends when the Network Connection is closed.
	// If the Session Expiry Interval is 0xFFFFFFFF (UINT_MAX), the Session does not expire.
	// 单位 秒(S)
	SessionExpiryInterval      uint32
	ReceiveMaximum             uint16
	MaximumPacketSize          uint32
	TopicAliasMaximum          uint16
	RequestResponseInformation uint8
	RequestProblemInformation  uint8
	UserProperty               map[string]interface{}
	AuthenticationMethod       string
	AuthenticationData         []byte
}

type WillProperties struct {
	WillDelayInterval      uint32
	PayloadFormatIndicator uint8
	MessageExpiryInterval  uint32
	ContentType            string
	ResponseTopic          string
	CorrelationData        []byte
	UserProperty           map[string]interface{}
}

// 下面是 具体协议编码解码
type Connect struct {
	FixedHeader       *FixedHeader
	ProtocolName      string
	ProtocolLevel     uint8
	ConnectFlags      *ConnectFlags
	KeepAlive         uint16
	ConnectProperties *ConnectProperties
	ClientId          string
	WillProperties    *WillProperties
	WillTopic         string
	WillPayload       []byte
	Usename           string
	Password          string
}

func (m *Connect) Encode(buf *bytes.Buffer) (err error) {
	err = m.FixedHeader.Encode(buf)
	err = setString(m.ProtocolName, buf)
	err = setUint8(m.ProtocolLevel, buf)
	err = m.ConnectFlags.Encode(buf)
	err = setUint16(m.KeepAlive, buf)
	var cp Properties = m.ConnectProperties
	err = Encode(&cp, buf)
	err = setString(m.ClientId, buf)
	var wp Properties = m.WillProperties
	err = Encode(&wp, buf)
	err = setString(m.WillTopic, buf)
	err = setBytes(m.WillPayload, buf)
	err = setString(m.Usename, buf)
	err = setString(m.Password, buf)

	return
}

func (m *Connect) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header

	m.ProtocolName = getString(b, &p)
	m.ProtocolLevel = getUint8(b, &p)

	connectFlags := &ConnectFlags{}
	connectFlags.Decode(b, &p)
	m.ConnectFlags = connectFlags
	m.KeepAlive = getUint16(b, &p)

	var connectProperties Properties
	connectProperties = &ConnectProperties{}
	Decode(&connectProperties, b, &p)
	m.ConnectProperties = connectProperties.(*ConnectProperties)

	m.ClientId = getString(b, &p)

	var willProperties Properties
	willProperties = &WillProperties{}
	Decode(&willProperties, b, &p)
	m.WillProperties = willProperties.(*WillProperties)

	m.WillTopic = getString(b, &p)
	m.WillPayload = getBytes(b, &p)
	m.Usename = getString(b, &p)
	m.Password = getString(b, &p)

	return

}
