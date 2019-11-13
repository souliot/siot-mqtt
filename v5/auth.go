package v5

import (
	"bytes"

	util "github.com/souliot/siot-mqtt/util"
)

type AuthProperties struct {
	AuthenticationMethod string
	AuthenticationData   []byte
	ReasonString         string
	UserProperty         map[string][]interface{}
}

// 下面是 具体协议编码解码
type Auth struct {
	FixedHeader    *FixedHeader
	ReasonCode     util.ReasonCode
	AuthProperties *AuthProperties
}

func (m *Auth) Encode(buf *bytes.Buffer) (err error) {
	err = m.FixedHeader.Encode(buf)

	if m.ReasonCode == 0 && m.AuthProperties == nil {
		return
	}

	err = util.SetUint8(uint8(m.ReasonCode), buf)

	var cp Properties = m.AuthProperties
	err = Encode(&cp, buf)

	return
}

func (m *Auth) Decode(b []byte) {
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
	properties = &AuthProperties{}
	Decode(&properties, b, &p)
	m.AuthProperties = properties.(*AuthProperties)

	return

}
