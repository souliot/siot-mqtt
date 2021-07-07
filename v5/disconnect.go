package v5

import (
	"bytes"

	util "github.com/souliot/siot-mqtt/util"
)

type DisconnectProperties struct {
	SessionExpiryInterval uint32
	ReasonString          string
	UserProperty          map[string][]interface{}
	ServerReference       string
}

// 下面是 具体协议编码解码
type Disconnect struct {
	FixedHeader          *FixedHeader
	ReasonCode           util.ReasonCode
	DisconnectProperties *DisconnectProperties
}

func (m *Disconnect) Encode(buf *bytes.Buffer) (err error) {
	bt := new(bytes.Buffer)

	err = util.SetUint8(uint8(m.ReasonCode), bt)
	var cp Properties
	if m.DisconnectProperties != nil {
		cp = m.DisconnectProperties
	} else {
		cp = new(DisconnectProperties)
	}
	err = Encode(&cp, bt)

	m.FixedHeader.RemainingLength = uint32(bt.Len())
	err = m.FixedHeader.Encode(buf)
	buf.Write(bt.Bytes())
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

	m.ReasonCode = util.ReasonCode(util.GetUint8(b, &p))
	if header.RemainingLength < 2 {
		return
	}

	var properties Properties
	properties = &DisconnectProperties{}
	Decode(&properties, b, &p)
	m.DisconnectProperties = properties.(*DisconnectProperties)

	return

}
