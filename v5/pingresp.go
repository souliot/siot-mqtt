package v5

import (
	"bytes"
)

// 下面是 具体协议编码解码
type PingResp struct {
	FixedHeader *FixedHeader
}

func (m *PingResp) Encode(buf *bytes.Buffer) (err error) {
	err = m.FixedHeader.Encode(buf)
	return
}

func (m *PingResp) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header

	return

}
