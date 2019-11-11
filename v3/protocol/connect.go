package protocol

import (
	"bytes"
)

type ConnectFlags struct {
	UsernameFlag, PasswordFlag, WillRetain bool
	WillQos                                uint8
	WillFlag, CleanSession, Reserved       bool
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
	val |= boolToByte(m.CleanSession) << 1
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
		CleanSession: byte1&0x02 > 0,
		Reserved:     byte1&0x01 > 0,
	}
	return
}

// 下面是 具体协议编码解码
type Connect struct {
	FixedHeader   *FixedHeader
	ProtocolName  string
	ProtocolLevel uint8
	ConnectFlags  *ConnectFlags
	KeepAlive     uint16
	ClientId      string
	WillTopic     string
	WillMessage   []byte
	Usename       string
	Password      string
}

func (m *Connect) Encode(buf *bytes.Buffer) (err error) {
	err = m.FixedHeader.Encode(buf)
	err = setString(m.ProtocolName, buf)
	err = setUint8(m.ProtocolLevel, buf)
	err = m.ConnectFlags.Encode(buf)
	err = setUint16(m.KeepAlive, buf)

	err = setString(m.ClientId, buf)

	if m.ConnectFlags.WillFlag {
		err = setString(m.WillTopic, buf)
		err = setBytes(m.WillMessage, buf)
	}
	if m.ConnectFlags.UsernameFlag {
		err = setString(m.Usename, buf)
	}
	if m.ConnectFlags.PasswordFlag {
		err = setString(m.Password, buf)
	}

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

	m.ClientId = getString(b, &p)

	if m.ConnectFlags.WillFlag {
		m.WillTopic = getString(b, &p)
		m.WillMessage = getBytes(b, &p)
	}
	if m.ConnectFlags.UsernameFlag {
		m.Usename = getString(b, &p)
	}
	if m.ConnectFlags.PasswordFlag {
		m.Password = getString(b, &p)
	}
	return

}
