package db

import (
	"errors"
	"strings"

	"github.com/souliot/siot-orm/orm"

	logs "github.com/souliot/siot-log"
)

var (
	ClientInvalid    = errors.New("Client is Invalid")
	UnsubscribeError = errors.New("Unsubscribe Topic is Invalid")
)

var (
	mongodb = "default"
)

type MongoSetting struct {
	Hosts    []string `json:"Hosts"`
	Username string   `json:"Username"`
	Password string   `json:"Password"`
	DbName   string   `json:"DbName"`
}

func InitMongo(s *MongoSetting) {
	logs.Info("初始化 mongodb 配置信息...")
	if orm.HasDefaultDataBase() {
		mongodb = "mongodb"
	}
	mu := ""
	if s.Username != "" {
		mu = s.Username + ":" + s.Password
	}
	mongo_address := "mongodb://" + mu + "@" + strings.Join(s.Hosts, ",") + "/" + s.DbName
	if mu != "" {
		mongo_address += "?authSource=admin"
	}

	orm.RegisterDriver("mongo", orm.DRMongo)
	err := orm.RegisterDataBase(mongodb, "mongo", mongo_address, true)
	if err != nil {
		logs.Error("初始化mongodb错误：", mongo_address, err)
		return
	}
	logs.Info("初始化 mongodb 配置信息完成：", mongo_address)
}
