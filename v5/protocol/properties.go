package protocol

import (
	"bytes"
	"reflect"
)

var Properties_type = map[string]reflect.Kind{
	"PayloadFormatIndicator":          reflect.Uint8,
	"MessageExpiryInterval":           reflect.Uint32,
	"ContentType":                     reflect.String,
	"ResponseTopic":                   reflect.String,
	"CorrelationData":                 reflect.Slice,
	"SubscriptionIdentifier":          reflect.Uint32,
	"SessionExpiryInterval":           reflect.Uint32,
	"AssignedClientIdentifier":        reflect.String,
	"ServerKeepAlive":                 reflect.Uint16,
	"AuthenticationMethod":            reflect.String,
	"AuthenticationData":              reflect.Slice,
	"RequestProblemInformation":       reflect.Uint8,
	"WillDelayInterval":               reflect.Uint32,
	"RequestResponseInformation":      reflect.Uint8,
	"ResponseInformation":             reflect.String,
	"ServerReference":                 reflect.String,
	"ReasonString":                    reflect.String,
	"ReceiveMaximum":                  reflect.Uint16,
	"TopicAliasMaximum":               reflect.Uint16,
	"TopicAlias":                      reflect.Uint16,
	"MaximumQoS":                      reflect.Uint8,
	"RetainAvailable":                 reflect.Uint8,
	"UserProperty":                    reflect.Map,
	"MaximumPacketSize":               reflect.Uint32,
	"WildcardSubscriptionAvailable":   reflect.Uint8,
	"SubscriptionIdentifierAvailable": reflect.Uint8,
	"SharedSubscriptionAvailable":     reflect.Uint8,
}
var Properties_value = map[string]uint8{
	"PayloadFormatIndicator":          1,
	"MessageExpiryInterval":           2,
	"ContentType":                     3,
	"ResponseTopic":                   8,
	"CorrelationData":                 9,
	"SubscriptionIdentifier":          11,
	"SessionExpiryInterval":           17,
	"AssignedClientIdentifier":        18,
	"ServerKeepAlive":                 19,
	"AuthenticationMethod":            21,
	"AuthenticationData":              22,
	"RequestProblemInformation":       23,
	"WillDelayInterval":               24,
	"RequestResponseInformation":      25,
	"ResponseInformation":             26,
	"ServerReference":                 28,
	"ReasonString":                    31,
	"ReceiveMaximum":                  33,
	"TopicAliasMaximum":               34,
	"TopicAlias":                      35,
	"MaximumQoS":                      36,
	"RetainAvailable":                 37,
	"UserProperty":                    38,
	"MaximumPacketSize":               39,
	"WildcardSubscriptionAvailable":   40,
	"SubscriptionIdentifierAvailable": 41,
	"SharedSubscriptionAvailable":     42,
}
var Properties_name = map[uint8]string{
	1:  "PayloadFormatIndicator",
	2:  "MessageExpiryInterval",
	3:  "ContentType",
	8:  "ResponseTopic",
	9:  "CorrelationData",
	11: "SubscriptionIdentifier",
	17: "SessionExpiryInterval",
	18: "AssignedClientIdentifier",
	19: "ServerKeepAlive",
	21: "AuthenticationMethod",
	22: "AuthenticationData",
	23: "RequestProblemInformation",
	24: "WillDelayInterval",
	25: "RequestResponseInformation",
	26: "ResponseInformation",
	28: "ServerReference",
	31: "ReasonString",
	33: "ReceiveMaximum",
	34: "TopicAliasMaximum",
	35: "TopicAlias",
	36: "MaximumQoS",
	37: "RetainAvailable",
	38: "UserProperty",
	39: "MaximumPacketSize",
	40: "WildcardSubscriptionAvailable",
	41: "SubscriptionIdentifierAvailable",
	42: "SharedSubscriptionAvailable",
}

type Properties interface {
}

func Encode(m *Properties, buf *bytes.Buffer) (err error) {
	if reflect.ValueOf(*m).IsNil() {
		return
	}
	bt := new(bytes.Buffer)
	t := reflect.TypeOf(*m).Elem()
	vs := reflect.ValueOf(*m).Elem()
	count := vs.NumField()

	for i := 0; i < count; i++ {
		f := t.Field(i)
		v := vs.Field(i)

		switch f.Type.Kind() {
		case reflect.Uint8:
			if v.Uint() == 0 {
				continue
			}
			err = setUint8(Properties_value[f.Name], bt)
			if err != nil {
				return
			}
			err = setUint8(uint8(v.Uint()), bt)
			if err != nil {
				return
			}
		case reflect.Uint16:
			if v.Uint() == 0 {
				continue
			}
			err = setUint8(Properties_value[f.Name], bt)
			if err != nil {
				return
			}
			err = setUint16(uint16(v.Uint()), bt)
			if err != nil {
				return
			}
		case reflect.Uint32:
			if v.Uint() == 0 {
				continue
			}
			err = setUint8(Properties_value[f.Name], bt)
			if err != nil {
				return
			}
			err = setUint32(uint32(v.Uint()), bt)
			if err != nil {
				return
			}
		case reflect.String:
			if v.Len() == 0 {
				continue
			}
			err = setUint8(Properties_value[f.Name], bt)
			if err != nil {
				return
			}
			err = setString(v.String(), bt)
			if err != nil {
				return
			}
		case reflect.Slice:
			if v.Len() == 0 {
				continue
			}
			err = setUint8(Properties_value[f.Name], bt)
			if err != nil {
				return
			}
			err = setBytes(v.Bytes(), bt)
			if err != nil {
				return
			}
		case reflect.Map:
			if v.Len() == 0 {
				continue
			}
			for _, val := range v.MapKeys() {
				if reflect.TypeOf(v.MapIndex(val).Interface()).Kind() == reflect.Slice {
					in := reflect.ValueOf(v.MapIndex(val).Interface())
					for j := 0; j < in.Len(); j++ {
						err = setUint8(Properties_value[f.Name], bt)
						if err != nil {
							return
						}
						err = setString(val.String(), bt)
						if err != nil {
							return
						}

						err = setString(interfaceToString(in.Index(j).Interface()), bt)
						if err != nil {
							return
						}
					}
				} else {
					err = setUint8(Properties_value[f.Name], bt)
					if err != nil {
						return
					}
					err = setString(val.String(), bt)
					if err != nil {
						return
					}

					err = setString(interfaceToString(v.MapIndex(val).Interface()), bt)
					if err != nil {
						return
					}
				}

			}
		}
	}

	l := uint32(bt.Len())
	err = encodeLength(l, buf)
	if err != nil {
		return
	}
	_, err = buf.Write(bt.Bytes())
	return
}

func Decode(m *Properties, b []byte, p *int) {
	if len(b) <= *p {
		return
	}
	t := reflect.ValueOf(*m).Elem()
	l := int(decodeLength(b, p))
	lt := l + *p
	for *p < lt {
		k := Properties_name[getUint8(b, p)]
		field := t.FieldByName(k)
		// fmt.Println(k, b[*p:])
		switch Properties_type[k] {
		case reflect.Uint8:
			field.Set(reflect.ValueOf(getUint8(b, p)))
		case reflect.Uint16:
			field.Set(reflect.ValueOf(getUint16(b, p)))
		case reflect.Uint32:
			field.Set(reflect.ValueOf(getUint32(b, p)))
		case reflect.String:
			field.SetString(getString(b, p))
		case reflect.Slice:
			field.SetBytes(getBytes(b, p))
		case reflect.Map:
			key := getString(b, p)
			val := getString(b, p)
			if field.IsNil() {
				mapReflect := reflect.MakeMap(reflect.TypeOf(make(map[string]interface{})))
				mapReflect.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
				field.Set(mapReflect)
			} else {
				st := field.MapIndex(reflect.ValueOf(key))
				if st.IsValid() {
					if reflect.ValueOf(st.Interface()).Kind() == reflect.Slice {
						field.SetMapIndex(reflect.ValueOf(key), reflect.Append(reflect.ValueOf(st.Interface()), reflect.ValueOf(val)))
					} else {
						sliceReflect := reflect.MakeSlice(reflect.TypeOf([]interface{}{reflect.ValueOf(val)}), 0, 10)
						sliceReflect = reflect.Append(sliceReflect, reflect.ValueOf(st.Interface()))
						sliceReflect = reflect.Append(sliceReflect, reflect.ValueOf(val))
						field.SetMapIndex(reflect.ValueOf(key), sliceReflect)
					}
				} else {
					field.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
				}
			}
		}
	}
	return
}
