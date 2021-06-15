package server

import (
	"sync"

	"github.com/souliot/fetcp"
	logs "github.com/souliot/siot-log"
	util "github.com/souliot/siot-mqtt/util"
	v5 "github.com/souliot/siot-mqtt/v5"
)

type HandlerV5 struct {
	mutex *sync.Mutex
}

var _ iHandler = new(HandlerV5)

// 设备登录
func (m *HandlerV5) Connect(p *Packet, c *fetcp.Conn, srv *Server) {
	msg := p.Message.(*v5.Connect)
	clientid := msg.ClientId
	resCode := 0
	err := updateClientState(clientid, 1)
	if err != nil {
		logs.Error(err)
		resCode = 2
	}
	fixedHeader := msg.FixedHeader
	fixedHeader.MsgType = util.MsgConnAck
	fixedHeader.RemainingLength = uint32(2)

	// 登录应答
	res := &Packet{
		MsgType: util.MsgConnAck,
		Message: &v5.ConnAck{
			FixedHeader: fixedHeader,
			ConnectAcknowledgeFlags: &v5.ConnectAcknowledgeFlags{
				SessionPresentFlag: false,
			},
			ReasonCode: util.ReasonCode(resCode),
		},
	}

	DownCommand(c, res)

	if resCode != 0 {
		c.Close()
		return
	}

	var extraData = &ExtraData{
		ClientId:          clientid,
		PacketIdentifiers: make(map[uint16]struct{}),
		ProtocolLevel:     msg.ProtocolLevel,
		SubscribePayload:  &v5.SubscribePayload{},
	}
	c.PutExtraData(extraData)

	// 如果当前设备已登录 强制下线
	srv.AddClient(clientid, c)
	logs.Info("设备登录:", clientid)
}

// 发布响应
func (m *HandlerV5) Publish(p *Packet, c *fetcp.Conn, srv *Server) {
	msg := p.Message.(*v5.Publish)
	logs.Info(msg.Payload)
	logs.Info(string(msg.Payload))
	// clientid := getClientId(c)
	extraData := c.GetExtraData().(*ExtraData)

	fixedHeader := *(msg.FixedHeader)
	qos := fixedHeader.QosLevel
	// retain := fixedHeader.Retain
	// TODO:
	// 消息处理
	// storageMsg := &db.Msg{
	// 	Topic:   msg.TopicName,
	// 	Sender:  clientid,
	// 	Qos:     uint8(qos),
	// 	Retain:  retain,
	// 	Payload: string(msg.Payload),
	// }
	// err := storageMsg.Insert()
	// if err != nil {
	// 	logs.Error("Storage Message Error:", err)
	// }

	// 回复客户端
	if qos == 0 {
		publishMessageV5(p, srv)
		return
	}

	if qos == 1 {
		publishMessageV5(p, srv)

		fixedHeader.MsgType = util.MsgPubAck
		fixedHeader.RemainingLength = uint32(2)
		res := &Packet{
			Message: &v5.PubAck{
				FixedHeader:      &fixedHeader,
				PacketIdentifier: msg.PacketIdentifier,
			},
		}
		DownCommand(c, res)
	}

	if qos == 2 {
		if _, ok := extraData.PacketIdentifiers[msg.PacketIdentifier]; ok {
			p.Message.(*v5.Publish).FixedHeader.DupFlag = true
		}
		publishMessageV5(p, srv)

		m.mutex.Lock()
		extraData.PacketIdentifiers[msg.PacketIdentifier] = struct{}{}
		m.mutex.Unlock()

		fixedHeader.MsgType = util.MsgPubRec
		fixedHeader.RemainingLength = uint32(2)
		res := &Packet{
			Message: &v5.PubRec{
				FixedHeader:      &fixedHeader,
				PacketIdentifier: msg.PacketIdentifier,
			},
		}
		DownCommand(c, res)
	}

}

func publishMessageV5(p *Packet, srv *Server) {
	// msg := p.Message.(*v5.Publish)
	// for _, c := range srv.GetClientList() {
	// 	if extraData, ok := c.GetExtraData().(*ExtraData); ok {
	// 		subscribePayload := extraData.SubscribePayload
	// 		if subscribePayload.HasPublish(msg) {
	// 			// logs.Info("分发消息：", ClientId)
	// 			// logs.Info(p.Message.(*v4.Publish).PacketIdentifier)
	// 			DownCommand(c, p)
	// 		}
	// 	}
	// }
}
