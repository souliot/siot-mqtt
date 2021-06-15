package main

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/souliot/siot-mqtt/db"
	"github.com/souliot/siot-mqtt/server"
	"github.com/souliot/siot-mqtt/util"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	util.InitLog("logs", "server")
	ms := &db.MongoSetting{
		Hosts:    []string{"192.168.0.4:27017"},
		Username: "",
		Password: "",
		DbName:   "llz-mqtt",
	}
	db.InitMongo(ms)
	srv := server.NewServer()
	defer func() {
		srv.Stop()
	}()
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	_ = <-chSig
}
