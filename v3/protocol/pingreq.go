package protocol

import "bytes"

// 下面是 具体协议编码解码
type Pingreq struct {
	FixedHeader *FixedHeader
}

func (m *Pingreq) Encode(buf *bytes.Buffer) (err error) {
	err = m.FixedHeader.Encode(buf)
	return
}

func (m *Pingreq) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header

	return

}
