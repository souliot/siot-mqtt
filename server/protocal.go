package server

import (
	"bytes"
	"sync"

	"github.com/souliot/fetcp"
	logs "github.com/souliot/siot-log"
	util "github.com/souliot/siot-mqtt/util"
)

var (
	ClientList = make(map[string]*fetcp.Conn)
	Mutex      sync.Mutex
)

type ExtraData struct {
	ClientId          string
	PacketIdentifiers map[uint16]struct{}
	ProtocolLevel     uint8
	SubscribePayload  interface{}
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
			ihandler = NewHandlerV3()
			goto NEXT
		case 5:
			ihandler = NewHandlerV5()
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
	// logs.Info("-----------------", packet.MsgType)
	switch packet.MsgType {
	case util.MsgConnect:
		logs.Info("MsgConnect")
		go h.Connect(packet, c, m.Srv)
	case util.MsgPublish:
		logs.Info("MsgPublish")
		go h.Publish(packet, c, m.Srv)
	case util.MsgPubAck:
		logs.Info("MsgPubAck")
		go h.PubAck(packet, c, m.Srv)
	case util.MsgPubRec:
		logs.Info("MsgPubRec")
		go h.PubRec(packet, c, m.Srv)
	case util.MsgPubRel:
		logs.Info("MsgPubRel")
		go h.PubRel(packet, c, m.Srv)
	case util.MsgPubComp:
		logs.Info("MsgPubComp")
		go h.PubComp(packet, c, m.Srv)
	case util.MsgSubscribe:
		logs.Info("MsgSubscribe")
		go h.Subscribe(packet, c, m.Srv)
	case util.MsgUnsubscribe:
		logs.Info("MsgUnsubscribe")
		go h.Unsubscribe(packet, c, m.Srv)
	case util.MsgPingReq:
		logs.Info("MsgPingReq")
		go h.PingReq(packet, c, m.Srv)
	case util.MsgDisconnect:
		logs.Info("MsgDisconnect")
		go h.Disconnect(packet, c, m.Srv)
	case util.MsgAuth:
		logs.Info("MsgAuth")
		go h.Auth(packet, c, m.Srv)
	default:
		logs.Error("Unkown Message Type:", packet.MsgType)
		// c.Close()
	}
	return
}
