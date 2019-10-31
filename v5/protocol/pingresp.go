package protocol

import "bytes"

// 下面是 具体协议编码解码
type Pingresp struct {
	FixedHeader *FixedHeader
}

func (m *Pingresp) Encode(buf *bytes.Buffer) (err error) {
	err = m.FixedHeader.Encode(buf)
	return
}

func (m *Pingresp) Decode(b []byte) {
	p := 0
	header := &FixedHeader{}
	header.Decode(b, &p)
	m.FixedHeader = header

	return

}
