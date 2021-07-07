package server

import (
	"github.com/souliot/fetcp"
	logs "github.com/souliot/siot-log"
	"github.com/souliot/siot-mqtt/db"
)

type iHandler interface {
	Connect(*Packet, *fetcp.Conn, *Server)
	Publish(*Packet, *fetcp.Conn, *Server)
	PubAck(*Packet, *fetcp.Conn, *Server)
	PubRec(*Packet, *fetcp.Conn, *Server)
	PubRel(*Packet, *fetcp.Conn, *Server)
	PubComp(*Packet, *fetcp.Conn, *Server)
	Subscribe(*Packet, *fetcp.Conn, *Server)
	Unsubscribe(*Packet, *fetcp.Conn, *Server)
	PingReq(*Packet, *fetcp.Conn, *Server)
	Disconnect(*Packet, *fetcp.Conn, *Server)
	Auth(*Packet, *fetcp.Conn, *Server)
}

// **************************************************下发指令********************************************** //
// 指令下发同义操作
func DownCommand(c *fetcp.Conn, p *Packet) {
	err := c.AsyncWritePacket(p, 0)
	if err != nil {
		logs.Error(err)
	}
}

func getClientId(c *fetcp.Conn) (clientid string) {
	if c.GetExtraData() != nil {
		return c.GetExtraData().(*ExtraData).ClientId
	}
	return
}

func updateClientState(clientid string, state int8) (err error) {
	client := &db.Client{
		ClientId: clientid,
	}
	switch state {
	case 0:
		err = client.Disconnect()
	case 1:
		err = client.Connect()
	}

	return
}
