package v3

import (
	"bytes"

	util "github.com/souliot/siot-mqtt/util"
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
	val := util.BoolToByte(m.UsernameFlag) << 7
	val |= util.BoolToByte(m.PasswordFlag) << 6
	val |= util.BoolToByte(m.WillRetain) << 5
	val |= byte(m.WillQos) << 3
	val |= util.BoolToByte(m.WillFlag) << 2
	val |= util.BoolToByte(m.CleanSession) << 1
	val |= util.BoolToByte(m.Reserved)

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
	err = util.SetString(m.ProtocolName, buf)
	err = util.SetUint8(m.ProtocolLevel, buf)
	err = m.ConnectFlags.Encode(buf)
	err = util.SetUint16(m.KeepAlive, buf)

	err = util.SetString(m.ClientId, buf)

	if m.ConnectFlags.WillFlag {
		err = util.SetString(m.WillTopic, buf)
		err = util.SetBytes(m.WillMessage, buf)
	}
	if m.ConnectFlags.UsernameFlag {
		err = util.SetString(m.Usename, buf)
	}
	if m.ConnectFlags.PasswordFlag {
		err = util.SetString(m.Password, buf)
	}

	return
}

func (m *Connect) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header

	m.ProtocolName = util.GetString(b, &p)
	m.ProtocolLevel = util.GetUint8(b, &p)

	connectFlags := &ConnectFlags{}
	connectFlags.Decode(b, &p)
	m.ConnectFlags = connectFlags
	m.KeepAlive = util.GetUint16(b, &p)

	m.ClientId = util.GetString(b, &p)

	if m.ConnectFlags.WillFlag {
		m.WillTopic = util.GetString(b, &p)
		m.WillMessage = util.GetBytes(b, &p)
	}
	if m.ConnectFlags.UsernameFlag {
		m.Usename = util.GetString(b, &p)
	}
	if m.ConnectFlags.PasswordFlag {
		m.Password = util.GetString(b, &p)
	}
	return

}
