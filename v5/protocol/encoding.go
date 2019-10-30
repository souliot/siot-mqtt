package protocol

import (
	"bytes"
	"encoding/binary"
	"strconv"
)

func getBytes(b []byte, p *int) []byte {
	l := int(getUint16(b, p))
	*p += l
	return b[*p-l : *p]
}
func getUint8(b []byte, p *int) uint8 {
	*p += 1
	return uint8(b[*p-1])
}

func getUint16(b []byte, p *int) uint16 {
	*p += 2
	return binary.BigEndian.Uint16(b[*p-2 : *p])
}

func getUint32(b []byte, p *int) uint32 {
	*p += 4
	return binary.BigEndian.Uint32(b[*p-4 : *p])
}

func getString(b []byte, p *int) string {
	l := int(getUint16(b, p))
	*p += l
	return string(b[*p-l : *p])
}

func getBytesNoLen(b []byte, p *int) []byte {
	*p = len(b) - 1
	return b[*p:]
}

func setBytes(val []byte, buf *bytes.Buffer) (err error) {
	length := uint16(len(val))
	err = setUint16(length, buf)
	if err != nil {
		return
	}
	_, err = buf.Write(val)
	return
}

func setUint8(val uint8, buf *bytes.Buffer) (err error) {
	err = buf.WriteByte(byte(val))
	return
}

func setUint16(val uint16, buf *bytes.Buffer) (err error) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, val)

	_, err = buf.Write(b)
	return
}

func setUint32(val uint32, buf *bytes.Buffer) (err error) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, val)

	_, err = buf.Write(b)
	return
}

func setString(val string, buf *bytes.Buffer) (err error) {
	length := uint16(len(val))
	err = setUint16(length, buf)
	if err != nil {
		return
	}
	_, err = buf.WriteString(val)
	return
}

func setBytesNoLen(val []byte, buf *bytes.Buffer) (err error) {
	_, err = buf.Write(val)
	return
}

func boolToByte(val bool) byte {
	if val {
		return byte(1)
	}
	return byte(0)
}

func decodeLength(b []byte, p *int) uint32 {
	m := uint32(1)
	v := uint32(b[*p] & 0x7f)
	*p += 1
	for b[*p-1]&0x80 > 0 {
		m *= 128
		v += uint32(b[*p]&0x7f) * m
		*p += 1
	}
	return v
}

func encodeLength(length uint32, buf *bytes.Buffer) (err error) {
	if length == 0 {
		err = buf.WriteByte(byte(0))
		return
	}
	for length > 0 {
		digit := length & 0x7f
		length = length >> 7
		if length > 0 {
			digit = digit | 0x80
		}
		err = buf.WriteByte(byte(digit))
		if err != nil {
			return
		}
	}
	return
}

func interfaceToString(v interface{}) (str string) {
	switch v.(type) {
	case string:
		str = v.(string)
	case int:
		str = strconv.Itoa(v.(int))
	case int64:
		str = strconv.FormatInt(v.(int64), 10)
	case float32:
		str = strconv.FormatFloat(float64(v.(float32)), 'f', -1, 32)
	case float64:
		str = strconv.FormatFloat(v.(float64), 'f', -1, 64)

	}
	return
}
