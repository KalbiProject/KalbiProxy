package proxy

import (
	"github.com/KalbiProject/Kalbi"
	"github.com/KalbiProject/Kalbi/interfaces"
	"github.com/KalbiProject/Kalbi/sip/message"
	"github.com/KalbiProject/Kalbi/sip/method"
	"github.com/KalbiProject/Kalbi/sip/status"
)

type Proxy struct {
	stack            *kalbi.SipStack
	RegisteredUsers  map[string]string
}

func (p *Proxy) HandleRequests(event interfaces.SipEventObject) {

	tx := event.GetTransaction()
	switch string(tx.GetLastMessage().Req.Method) {
	case method.CANCEL:
		go p.HandleCancel(tx)
	case method.INVITE:
		go p.HandleInvite(tx)
	case method.REGISTER:
		go p.HandleRegister(tx)
	case method.BYE:
		go p.HandleBye(tx)
	case method.ACK:
		
	default:
		msg := message.NewResponse(status.OK, "@", "@")
		msg.CopyHeaders(tx.GetOrigin())
		msg.ContLen.SetValue("0")
		tx.Send(msg, string(tx.GetOrigin().Contact.Host), string(tx.GetOrigin().Contact.Port))
	}

}


func (p *Proxy) HandleResponses(event interfaces.SipEventObject) {

	response := event.GetTransaction()

	switch response.GetLastMessage().GetStatusCode() {
	
	case 100:
		return
	case 180:
		return
	default:
		go p.Handle200(response)
	}

}

func (p *Proxy) AddToRegister(key string, contact string) {
	p.RegisteredUsers[key] = contact
}


func (p *Proxy) Start(host string, port int) {
	p.RegisteredUsers = make(map[string]string)
	p.stack = kalbi.NewSipStack("Basic")
	p.stack.CreateListenPoint("udp", host, port)
	p.stack.SetSipListener(p)
	go p.stack.Start()
	select{}//blocking action
}
