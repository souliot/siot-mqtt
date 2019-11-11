package protocol

import "bytes"

type DisconnectProperties struct {
	SessionExpiryInterval uint32
	ReasonString          string
	UserProperty          map[string][]interface{}
	ServerReference       string
}

// 下面是 具体协议编码解码
type Disconnect struct {
	FixedHeader          *FixedHeader
	ReasonCode           ReasonCode
	DisconnectProperties *DisconnectProperties
}

func (m *Disconnect) Encode(buf *bytes.Buffer) (err error) {
	err = m.FixedHeader.Encode(buf)

	if m.ReasonCode == 0 && m.DisconnectProperties == nil {
		return
	}

	err = setUint8(uint8(m.ReasonCode), buf)

	var cp Properties = m.DisconnectProperties
	err = Encode(&cp, buf)

	return
}

func (m *Disconnect) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header
	if header.RemainingLength < 1 {
		return
	}

	m.ReasonCode = ReasonCode(getUint8(b, &p))
	if header.RemainingLength < 2 {
		return
	}

	var properties Properties
	properties = &DisconnectProperties{}
	Decode(&properties, b, &p)
	m.DisconnectProperties = properties.(*DisconnectProperties)

	return

}
