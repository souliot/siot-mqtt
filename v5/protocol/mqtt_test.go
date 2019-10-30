package protocol

import (
	"bytes"
	"fmt"
	"siot-mqtt/v5/protocol"
	"testing"
)

type Properties1 struct {
	PayloadFormatIndicator          uint8
	MessageExpiryInterval           uint32
	ContentType                     string
	ResponseTopic                   string
	CorrelationData                 []byte
	SubscriptionIdentifier          uint32
	SessionExpiryInterval           uint32
	AssignedClientIdentifier        string
	ServerKeepAlive                 uint16
	AuthenticationMethod            string
	AuthenticationData              []byte
	RequestProblemInformation       uint8
	WillDelayInterval               uint32
	RequestResponseInformation      uint8
	ResponseInformation             string
	ServerReference                 string
	ReasonString                    string
	ReceiveMaximum                  uint16
	TopicAliasMaximum               uint16
	TopicAlias                      uint16
	MaximumQoS                      uint8
	RetainAvailable                 uint8
	UserProperty                    map[string][]interface{}
	MaximumPacketSize               uint32
	WildcardSubscriptionAvailable   uint8
	SubscriptionIdentifierAvailable uint8
	SharedSubscriptionAvailable     uint8
}

var header = &protocol.FixedHeader{
	DupFlag:         true,
	Retain:          false,
	QosLevel:        2,
	MsgType:         3,
	RemainingLength: 123131,
}

var connectFlags = &protocol.ConnectFlags{
	true, true, false, 2, false, false, false,
}

var connectProperties = &protocol.ConnectProperties{
	SessionExpiryInterval:      3600,
	ReceiveMaximum:             1024,
	MaximumPacketSize:          65535,
	TopicAliasMaximum:          256,
	RequestResponseInformation: 16,
	RequestProblemInformation:  16,
	UserProperty: map[string]interface{}{
		"name": []interface{}{"llz", "xy", 1234},
		"age":  28,
		"sex":  "男",
	},
	AuthenticationMethod: "POST",
	// AuthenticationData         :
}

var WillProperties = &protocol.WillProperties{
	WillDelayInterval:      3600,
	PayloadFormatIndicator: 36,
	MessageExpiryInterval:  3600,
	ContentType:            "application/json",
	ResponseTopic:          "/device",
	// CorrelationData        :
	UserProperty: map[string]interface{}{
		"name": []interface{}{"llz", "xy", 1234},
		"age":  28,
		"sex":  "男",
	},
}

var connect = &protocol.Connect{
	FixedHeader:       header,
	ProtocolName:      "MQTT",
	ProtocolLevel:     5,
	ConnectFlags:      connectFlags,
	KeepAlive:         60,
	ConnectProperties: connectProperties,
	ClientId:          "feiciot-1001",
	WillProperties:    WillProperties,
	WillTopic:         "/user",
	WillPayload:       make([]byte, 2),
	Usename:           "souliot",
	Password:          "abcd1234",
}

func test_FixedHeader(t *testing.T) {

	buf := new(bytes.Buffer)
	header.Encode(buf)
	fmt.Println(buf.Bytes())
	p := 0

	h := &protocol.FixedHeader{}
	h.Decode(buf.Bytes(), &p)
	fmt.Println(h)
}

func test_ConnectFlags(t *testing.T) {

	buf := new(bytes.Buffer)
	connectFlags.Encode(buf)
	// fmt.Println(buf.Bytes())

	p := 0
	cf := &protocol.ConnectFlags{}
	cf.Decode(buf.Bytes(), &p)
	fmt.Println(cf)
}

func test_t(t *testing.T) {
	buf := new(bytes.Buffer)
	connect.Encode(buf)
	// fmt.Println(buf.Bytes())

	c := &protocol.Connect{}

	c.Decode(buf.Bytes())

	fmt.Println(c)
	fmt.Println(c.FixedHeader)
	fmt.Println(c.ConnectProperties)
	fmt.Println(c.WillProperties)
}

func TestAll(t *testing.T) {
	// t.Run("FixedHeader", test_FixedHeader)
	// t.Run("ConnectFlags", test_ConnectFlags)
	t.Run("Test", test_t)
}
