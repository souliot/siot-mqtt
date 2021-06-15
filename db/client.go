package db

import (
	"time"

	logs "github.com/souliot/siot-log"
	"github.com/souliot/siot-orm/orm"
)

type Client struct {
	Id        string `bson:"_id"`
	ClientId  string `bson:"ClientId"`
	ServerId  string `bson:"ServerId"`
	OnlineAt  int64  `bson:"OnlineAt"`
	OfflineAt int64  `bson:"OfflineAt"`
	State     uint8  `bson:"State"` //在线状态 0 离线 1 在线
}

func init() {
	orm.RegisterModel(new(Client))
}

func (m *Client) Connect() (err error) {
	o := orm.NewOrm()
	o.Using("default")
	cnt, err := o.QueryTable("Client").Filter("ClientId", m.ClientId).Update(orm.MgoSetOnInsert, orm.Params{
		"OnlineAt": time.Now().Unix(),
		"State":    1,
	})
	if err != nil {
		logs.Error(err)
		return
	}
	if cnt == 0 {
		m.OnlineAt = time.Now().Unix()
		m.State = 1
		_, err = o.Insert(m)
	}
	return
}

func (m *Client) Disconnect() (err error) {
	o := orm.NewOrm()
	o.Using("default")
	cnt, err := o.QueryTable("Client").Filter("ClientId", m.ClientId).Update(orm.MgoSetOnInsert, orm.Params{
		"OfflineAt": time.Now().Unix(),
		"State":     0,
	})
	if err != nil {
		return
	}
	logs.Info(cnt)
	return
}
