package proxy

import (
	"github.com/KalbiProject/Kalbi"
	"github.com/KalbiProject/Kalbi/sip/message"
	"github.com/KalbiProject/Kalbi/sip/method"
	"github.com/KalbiProject/Kalbi/sip/status"
	"github.com/KalbiProject/Kalbi/sip/transaction"
)

type Proxy struct {
	stack            *kalbi.SipStack
	requestschannel  chan transaction.Transaction
	responseschannel chan transaction.Transaction
	RegisteredUsers  map[string]string
}

func (p *Proxy) HandleRequest(tx transaction.Transaction) {

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


func (p *Proxy) HandleResponse(response transaction.Transaction) {

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

func (p *Proxy) ServeRequests() {

	for {
		tx := <-p.requestschannel
		p.HandleRequest(tx)
	}

}

func (p *Proxy) ServeResponses() {
	for {
		tx := <-p.responseschannel
		p.HandleResponse(tx)
	}
}

func (p *Proxy) Start(host string, port int) {
	p.RegisteredUsers = make(map[string]string)
	p.stack = kalbi.NewSipStack("Basic")
	p.stack.CreateListenPoint("udp", host, port)
	p.requestschannel = p.stack.CreateRequestsChannel()
	p.responseschannel = p.stack.CreateResponseChannel()
	go p.stack.Start()
	go p.ServeRequests()
	
	
	p.ServeResponses()

}
