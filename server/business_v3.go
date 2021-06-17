package server

import (
	"sync"

	"github.com/souliot/fetcp"
	logs "github.com/souliot/siot-log"
	"github.com/souliot/siot-mqtt/db"
	util "github.com/souliot/siot-mqtt/util"
	v3 "github.com/souliot/siot-mqtt/v3"
)

type HandlerV3 struct {
	mutex *sync.Mutex
}

func NewHandlerV3() (h *HandlerV3) {
	return &HandlerV3{new(sync.Mutex)}
}

var _ iHandler = new(HandlerV3)

// 设备登录
func (m *HandlerV3) Connect(p *Packet, c *fetcp.Conn, srv *Server) {
	msg := p.Message.(*v3.Connect)
	if msg.KeepAlive == 0 {
		c.SetHeartBeatStatus(false)
	} else {
		c.SetKeepAlive(int64(msg.KeepAlive) * 3 / 2)
	}
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
		Message: &v3.ConnAck{
			FixedHeader: fixedHeader,
			ConnectAcknowledgeFlags: &v3.ConnectAcknowledgeFlags{
				SessionPresentFlag: false,
			},
			ReasonCode: util.ReasonCode(resCode),
		},
	}
	srv.AddClient(clientid, c, msg.ConnectFlags.CleanSession)
	DownCommand(c, res)

	if resCode != 0 {
		c.Close()
		return
	}

	var extraData = &ExtraData{
		ClientId:          clientid,
		PacketIdentifiers: make(map[uint16]struct{}),
		ProtocolLevel:     msg.ProtocolLevel,
		SubscribePayload:  &v3.SubscribePayload{},
	}
	c.PutExtraData(extraData)

	logs.Info("设备登录:", clientid)
}

// 发布响应
func (m *HandlerV3) Publish(p *Packet, c *fetcp.Conn, srv *Server) {
	msg := p.Message.(*v3.Publish)
	logs.Info(msg.Payload)
	logs.Info(string(msg.Payload))
	clientid := getClientId(c)
	extraData := c.GetExtraData().(*ExtraData)

	fixedHeader := *(msg.FixedHeader)
	qos := fixedHeader.QosLevel
	retain := fixedHeader.Retain
	// TODO:
	// 消息处理
	storageMsg := &db.Message{
		Topic:   msg.TopicName,
		Sender:  clientid,
		Qos:     uint8(qos),
		Retain:  retain,
		Payload: string(msg.Payload),
	}
	err := storageMsg.Insert()
	if err != nil {
		logs.Error("Storage Message Error:", err)
	}

	// 回复客户端
	if qos == 0 {
		publishMessageV3(p, srv)
		return
	}

	if qos == 1 {
		publishMessageV3(p, srv)

		fixedHeader.MsgType = util.MsgPubAck
		fixedHeader.RemainingLength = uint32(2)
		res := &Packet{
			Message: &v3.PubAck{
				FixedHeader:      &fixedHeader,
				PacketIdentifier: msg.PacketIdentifier,
			},
		}
		DownCommand(c, res)
	}

	if qos == 2 {
		if _, ok := extraData.PacketIdentifiers[msg.PacketIdentifier]; ok {
			p.Message.(*v3.Publish).FixedHeader.DupFlag = true
		}
		publishMessageV3(p, srv)

		m.mutex.Lock()
		extraData.PacketIdentifiers[msg.PacketIdentifier] = struct{}{}
		m.mutex.Unlock()

		fixedHeader.MsgType = util.MsgPubRec
		fixedHeader.RemainingLength = uint32(2)
		res := &Packet{
			Message: &v3.PubRec{
				FixedHeader:      &fixedHeader,
				PacketIdentifier: msg.PacketIdentifier,
			},
		}
		DownCommand(c, res)
	}
}

// 发布响应 QOS1
func (m *HandlerV3) PubAck(p *Packet, c *fetcp.Conn, srv *Server) {
	// Do Nothing
}

// 发布响应 第1步 QOS2
func (m *HandlerV3) PubRec(p *Packet, c *fetcp.Conn, srv *Server) {
	msg := p.Message.(*v3.PubRec)
	extraData := c.GetExtraData().(*ExtraData)
	fixedHeader := msg.FixedHeader
	// TODO:
	m.mutex.Lock()
	extraData.PacketIdentifiers[msg.PacketIdentifier] = struct{}{}
	m.mutex.Unlock()

	fixedHeader.MsgType = util.MsgPubRel
	fixedHeader.RemainingLength = uint32(2)
	res := &Packet{
		Message: &v3.PubRel{
			FixedHeader:      fixedHeader,
			PacketIdentifier: msg.PacketIdentifier,
		},
	}

	DownCommand(c, res)
}

// 发布响应 第2步 QOS2
func (m *HandlerV3) PubRel(p *Packet, c *fetcp.Conn, srv *Server) {
	msg := p.Message.(*v3.PubRel)
	extraData := c.GetExtraData().(*ExtraData)
	fixedHeader := msg.FixedHeader
	// TODO:
	m.mutex.Lock()
	delete(extraData.PacketIdentifiers, msg.PacketIdentifier)
	m.mutex.Unlock()

	fixedHeader.MsgType = util.MsgPubComp
	fixedHeader.RemainingLength = uint32(2)
	res := &Packet{
		Message: &v3.PubComp{
			FixedHeader:      fixedHeader,
			PacketIdentifier: msg.PacketIdentifier,
		},
	}

	DownCommand(c, res)
}

// 发布响应 第3步 QOS2
func (m *HandlerV3) PubComp(p *Packet, c *fetcp.Conn, srv *Server) {
	msg := p.Message.(*v3.PubComp)
	extraData := c.GetExtraData().(*ExtraData)
	// 删除保存的报文标识符
	m.mutex.Lock()
	delete(extraData.PacketIdentifiers, msg.PacketIdentifier)
	m.mutex.Unlock()
}

// 订阅响应
func (m *HandlerV3) Subscribe(p *Packet, c *fetcp.Conn, srv *Server) {
	msg := p.Message.(*v3.Subscribe)
	clientid := getClientId(c)

	fixedHeader := msg.FixedHeader
	ReasonCodes := []util.ReasonCode{}
	extraData := c.GetExtraData().(*ExtraData)
	extraData.SubscribePayload.(*v3.SubscribePayload).Merger(msg.SubscribePayload)

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
		Message: &v3.SubAck{
			FixedHeader:      fixedHeader,
			PacketIdentifier: msg.PacketIdentifier,
			SubAckPayload: &v3.SubAckPayload{
				ReasonCodes: ReasonCodes,
			},
		},
	}

	DownCommand(c, res)
}

// 取消订阅响应
func (m *HandlerV3) Unsubscribe(p *Packet, c *fetcp.Conn, srv *Server) {
	msg := p.Message.(*v3.Unsubscribe)
	clientid := getClientId(c)

	fixedHeader := msg.FixedHeader
	extraData := c.GetExtraData().(*ExtraData)
	extraData.SubscribePayload.(*v3.SubscribePayload).Remove(msg.UnsubscribePayload)

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
		Message: &v3.UnsubAck{
			FixedHeader:      fixedHeader,
			PacketIdentifier: msg.PacketIdentifier,
		},
	}

	DownCommand(c, res)
}

// 心跳处理
func (m *HandlerV3) PingReq(p *Packet, c *fetcp.Conn, srv *Server) {
	msg := p.Message.(*v3.PingReq)
	fixedHeader := msg.FixedHeader

	fixedHeader.MsgType = util.MsgPingResp
	fixedHeader.RemainingLength = uint32(0)
	res := &Packet{
		Message: &v3.PingReq{
			FixedHeader: fixedHeader,
		},
	}

	DownCommand(c, res)
}

// 心跳处理
func (m *HandlerV3) Auth(p *Packet, c *fetcp.Conn, srv *Server) {
	msg := p.Message.(*v3.PingReq)
	fixedHeader := msg.FixedHeader

	fixedHeader.MsgType = util.MsgPingResp
	fixedHeader.RemainingLength = uint32(0)
	res := &Packet{
		Message: &v3.PingReq{
			FixedHeader: fixedHeader,
		},
	}

	DownCommand(c, res)
}

// 断开连接操作
func (m *HandlerV3) Disconnect(p *Packet, c *fetcp.Conn, srv *Server) {
	clientid := getClientId(c)
	err := updateClientState(clientid, 0)
	if err != nil {
		logs.Error("Update Client State Error:", err)
	}
	c.Close()
}

func publishMessageV3(p *Packet, srv *Server) {
	msg := p.Message.(*v3.Publish)
	for _, c := range srv.GetClientList() {
		if extraData, ok := c.GetExtraData().(*ExtraData); ok {
			subscribePayload := extraData.SubscribePayload.(*v3.SubscribePayload)
			if subscribePayload.HasPublish(msg) {
				// logs.Info("分发消息：", ClientId)
				// logs.Info(p.Message.(*v3.Publish).PacketIdentifier)
				DownCommand(c, p)
			}
		}
	}
}
