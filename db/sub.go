package db

import (
	"github.com/souliot/siot-orm/orm"
)

type Sub struct {
	Id       string `bson:"_id"`
	ClientId string `bson:"ClientId"`
	Topic    string `bson:"Topic"`
	Qos      uint8  `bson:"Qos"`
}

func init() {
	orm.RegisterModel(new(Sub))
}

func (m *Sub) Insert() (err error) {
	o := orm.NewOrm()
	o.Using("default")
	exist := o.QueryTable("Sub").Filter("ClientId", m.ClientId).Filter("Topic", m.Topic).Exist()
	if exist {
		_, err = o.QueryTable("Sub").Filter("ClientId", m.ClientId).Filter("Topic", m.Topic).Update(orm.MgoSet, orm.Params{
			"Qos": m.Qos,
		})
		return
	}

	_, err = o.Insert(m)
	return
}

func (m *Sub) Delete() (err error) {
	o := orm.NewOrm()
	o.Using("default")
	_, err = o.QueryTable("Sub").Filter("ClientId", m.ClientId).Filter("Topic", m.Topic).Delete()
	return
}
