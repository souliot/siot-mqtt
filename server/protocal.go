package server

import (
	"bytes"
	"sync"

	"github.com/souliot/fetcp"
	logs "github.com/souliot/siot-log"
	util "github.com/souliot/siot-mqtt/util"
	v5 "github.com/souliot/siot-mqtt/v5"
)

var (
	ClientList = make(map[string]*fetcp.Conn)
	Mutex      sync.Mutex
)

type ExtraData struct {
	ClientId          string
	PacketIdentifiers map[uint16]struct{}
	ProtocolLevel     uint8
	SubscribePayload  *v5.SubscribePayload
}

type Protocol struct {
}

func (m *Protocol) ReadPacket(c *fetcp.Conn) (fetcp.Packet, error) {
	// logs.Info("ReadPacket")
	conn := c.GetRawConn()

	fullBuf := bytes.NewBuffer([]byte{})

	data := make([]byte, 1024)

	readLengh, err := conn.Read(data)

	if err != nil { //EOF, or worse
		return nil, err
	}

	if readLengh == 0 { // Connection maybe closed by the client
		return nil, fetcp.ErrConnClosing
	} else {
		fullBuf.Write(data[:readLengh])
		// logs.Info("接收数据:", fullBuf.String())
		return NewPacket(fullBuf.Bytes(), c)
	}
}

type Callback struct {
	Srv *Server
}

func newCallback(srv *Server) (cb *Callback) {
	return &Callback{srv}
}

func (m *Callback) OnConnect(c *fetcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	logs.Info("地址:", addr, "——连接成功！")
	return true
}

func (m *Callback) OnMessage(c *fetcp.Conn, p fetcp.Packet) bool {
	packet := p.(*Packet)
	if packet != nil {
		var ihandler iHandler
		switch packet.ProtocolLevel {
		case 4:
			ihandler = new(HandlerV3)
			goto NEXT
		case 5:
			ihandler = new(HandlerV5)
			goto NEXT
		default:
			logs.Error("Not Support Protocol Level...")
			c.Close()
		}
	NEXT:
		go m.handleMessage(ihandler, packet, c)
	} else {
		logs.Error("数据异常")
		// c.Close()
	}
	return true
}

func (m *Callback) OnClose(c *fetcp.Conn) {
	addr := c.GetRawConn().RemoteAddr()
	if c.GetExtraData() != nil {
		client := c.GetExtraData().(*ExtraData)
		m.Srv.DelClient(client.ClientId)

		logs.Info("登出:", client.ClientId)
		logs.Info("地址:", addr, "——断开连接！")
	} else {
		logs.Info("地址:", addr, "——断开连接！")
	}
}

func (m *Callback) handleMessage(h iHandler, packet *Packet, c *fetcp.Conn) {
	switch packet.MsgType {
	case util.MsgConnect:
		logs.Info("MsgConnect")
		go h.Connect(packet, c, m.Srv)
	case util.MsgPublish:
		logs.Info("MsgPublish")
		go h.Publish(packet, c, m.Srv)
	// case util.MsgPubAck:
	// 	logs.Info("MsgPubAck")
	// 	go processPubAckV5(packet, c)
	// case util.MsgPubRec:
	// 	logs.Info("MsgPubRec")
	// 	go processPubRecV5(packet, c)
	// case util.MsgPubRel:
	// 	logs.Info("MsgPubRel")
	// 	go processPubRelV5(packet, c)
	// case util.MsgPubComp:
	// 	logs.Info("MsgPubComp")
	// 	go processPubCompV5(packet, c)
	// case util.MsgSubscribe:
	// 	logs.Info("MsgSubscribe")
	// 	go processSubscribeV5(packet, c)
	// case util.MsgUnsubscribe:
	// 	logs.Info("MsgUnsubscribe")
	// 	go processUnsubscribeV5(packet, c)
	// case util.MsgPingReq:
	// 	logs.Info("MsgPingReq")
	// 	go processPingReqV5(packet, c)
	// case util.MsgDisconnect:
	// 	logs.Info("MsgDisconnect")
	// 	go processDisconnectV5(packet, c)
	// case util.MsgAuth:
	// 	logs.Info("MsgAuth")
	// 	go processAuthV5(packet, c)
	default:
		logs.Error("Unkown Message Type...")
		// c.Close()
	}
	return
}
