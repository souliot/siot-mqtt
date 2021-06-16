package db

import (
	"time"

	"github.com/souliot/siot-orm/orm"
	sutil "github.com/souliot/siot-util"
)

type Message struct {
	Id        string `bson:"_id"`
	Topic     string `bson:"Topic"`
	Msgid     string `bson:"Msgid"`
	Sender    string `bson:"Sender`
	Qos       uint8  `bson:"Qos"`
	Retain    bool   `bson:"Retain"`
	Payload   string `bson:"Payload"`
	ArrivedAt int64  `bson:"ArrivedAt"`
}

func init() {
	orm.RegisterModel(new(Message))
}

func (m *Message) Insert() (err error) {
	now := time.Now()
	m.Msgid = sutil.To_md5(m.Sender + ":" + m.Topic + ":" + now.Format("2006-01-02 15:04:05"))
	m.ArrivedAt = now.Unix()
	o := orm.NewOrm()
	o.Using("default")
	_, err = o.Insert(m)
	return
}

func (m *Message) Delete() (err error) {
	o := orm.NewOrm()
	o.Using("default")
	_, err = o.QueryTable("Message").Filter("Sender", m.Sender).Filter("Topic", m.Topic).Delete()
	return
}
