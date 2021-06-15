package util

import (
	"strings"

	logs "github.com/souliot/siot-log"
)

func InitLog(path string, name string) {
	logs.SetLogFuncCall(true)
	logs.EnableFullFilePath(false)
	logs.SetLogFuncCallDepth(3)
	logs.SetLevel(logs.LevelInfo)
	// logs.Async()
	logs.WithPrefix(name)
	logs.WithPrefix(GetIPStr())
	filepath := strings.TrimRight(path, "/") + "/" + name + ".log"
	logs.SetLogger("file", `{"filename":"`+filepath+`","daily":true,"maxdays":10,"color":false}`)
	logs.SetLogger("console")
}
