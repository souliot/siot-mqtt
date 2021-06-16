package server

import (
	"sync"

	"github.com/souliot/fetcp"
	logs "github.com/souliot/siot-log"
	"github.com/souliot/siot-mqtt/db"
	util "github.com/souliot/siot-mqtt/util"
	v5 "github.com/souliot/siot-mqtt/v5"
)

type HandlerV5 struct {
	mutex *sync.Mutex
}

func NewHandlerV5() (h *HandlerV5) {
	return &HandlerV5{new(sync.Mutex)}
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
	srv.AddClient(clientid, c, msg.ConnectFlags.CleanStart)
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

// 发布响应 QOS1
func (m *HandlerV5) PubAck(p *Packet, c *fetcp.Conn, srv *Server) {
	// Do Nothing
}

// 发布响应 第1步 QOS2
func (m *HandlerV5) PubRec(p *Packet, c *fetcp.Conn, srv *Server) {
	msg := p.Message.(*v5.PubRec)
	extraData := c.GetExtraData().(*ExtraData)
	fixedHeader := msg.FixedHeader
	// TODO:
	m.mutex.Lock()
	extraData.PacketIdentifiers[msg.PacketIdentifier] = struct{}{}
	m.mutex.Unlock()

	fixedHeader.MsgType = util.MsgPubRel
	fixedHeader.RemainingLength = uint32(2)
	res := &Packet{
		Message: &v5.PubRel{
			FixedHeader:      fixedHeader,
			PacketIdentifier: msg.PacketIdentifier,
		},
	}

	DownCommand(c, res)
}

// 发布响应 第2步 QOS2
func (m *HandlerV5) PubRel(p *Packet, c *fetcp.Conn, srv *Server) {
	msg := p.Message.(*v5.PubRel)
	extraData := c.GetExtraData().(*ExtraData)
	fixedHeader := msg.FixedHeader
	// TODO:
	m.mutex.Lock()
	delete(extraData.PacketIdentifiers, msg.PacketIdentifier)
	m.mutex.Unlock()

	fixedHeader.MsgType = util.MsgPubComp
	fixedHeader.RemainingLength = uint32(2)
	res := &Packet{
		Message: &v5.PubComp{
			FixedHeader:      fixedHeader,
			PacketIdentifier: msg.PacketIdentifier,
		},
	}

	DownCommand(c, res)
}

// 发布响应 第3步 QOS2
func (m *HandlerV5) PubComp(p *Packet, c *fetcp.Conn, srv *Server) {
	msg := p.Message.(*v5.PubComp)
	extraData := c.GetExtraData().(*ExtraData)
	// 删除保存的报文标识符
	m.mutex.Lock()
	delete(extraData.PacketIdentifiers, msg.PacketIdentifier)
	m.mutex.Unlock()
}

// 订阅响应
func (m *HandlerV5) Subscribe(p *Packet, c *fetcp.Conn, srv *Server) {
	msg := p.Message.(*v5.Subscribe)
	clientid := getClientId(c)

	fixedHeader := msg.FixedHeader
	ReasonCodes := []util.ReasonCode{}
	extraData := c.GetExtraData().(*ExtraData)
	extraData.SubscribePayload.(*v5.SubscribePayload).Merger(msg.SubscribePayload)

	l := 2
	for _, v := range msg.SubscribePayload.SubscribeTopics {
		resCode := util.ReasonCode(v.SubscriptionOptions.QosLevel)
		sub := &db.Sub{
			ClientId: clientid,
			Topic:    v.TopicFilter,
			Qos:      uint8(v.SubscriptionOptions.QosLevel),
		}
		err := sub.Insert()
		if err != nil {
			resCode = util.ReasonCode(80)
		}
		ReasonCodes = append(ReasonCodes, resCode)
		l++
	}

	fixedHeader.MsgType = util.MsgSubAck
	fixedHeader.RemainingLength = uint32(l)
	res := &Packet{
		Message: &v5.SubAck{
			FixedHeader:      fixedHeader,
			PacketIdentifier: msg.PacketIdentifier,
			SubAckPayload: &v5.SubAckPayload{
				ReasonCodes: ReasonCodes,
			},
		},
	}

	DownCommand(c, res)
}

// 取消订阅响应
func (m *HandlerV5) Unsubscribe(p *Packet, c *fetcp.Conn, srv *Server) {
	msg := p.Message.(*v5.Unsubscribe)
	clientid := getClientId(c)

	fixedHeader := msg.FixedHeader
	extraData := c.GetExtraData().(*ExtraData)
	extraData.SubscribePayload.(*v5.SubscribePayload).Remove(msg.UnsubscribePayload)

	for _, v := range msg.UnsubscribePayload.UnsubscribeTopics {
		unsub := &db.Sub{
			ClientId: clientid,
			Topic:    v.TopicFilter,
		}
		unsub.Delete()
	}

	fixedHeader.MsgType = util.MsgUnsubAck
	fixedHeader.RemainingLength = uint32(2)
	res := &Packet{
		Message: &v5.UnsubAck{
			FixedHeader:      fixedHeader,
			PacketIdentifier: msg.PacketIdentifier,
		},
	}

	DownCommand(c, res)
}

// 心跳处理
func (m *HandlerV5) PingReq(p *Packet, c *fetcp.Conn, srv *Server) {
	msg := p.Message.(*v5.PingReq)
	fixedHeader := msg.FixedHeader

	fixedHeader.MsgType = util.MsgPingResp
	fixedHeader.RemainingLength = uint32(0)
	res := &Packet{
		Message: &v5.PingReq{
			FixedHeader: fixedHeader,
		},
	}

	DownCommand(c, res)
}

// 认证
func (m *HandlerV5) Auth(p *Packet, c *fetcp.Conn, srv *Server) {
	msg := p.Message.(*v5.PingReq)
	fixedHeader := msg.FixedHeader

	fixedHeader.MsgType = util.MsgPingResp
	fixedHeader.RemainingLength = uint32(0)
	res := &Packet{
		Message: &v5.PingReq{
			FixedHeader: fixedHeader,
		},
	}

	DownCommand(c, res)
}

// 断开连接操作
func (m *HandlerV5) Disconnect(p *Packet, c *fetcp.Conn, srv *Server) {
	c.Close()
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
