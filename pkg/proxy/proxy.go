package proxy

import (
	"fmt"
	"github.com/KalbiProject/Kalbi"
	"github.com/KalbiProject/Kalbi/sip/message"
	"github.com/KalbiProject/Kalbi/sip/method"
	"github.com/KalbiProject/Kalbi/sip/status"
	"github.com/KalbiProject/Kalbi/sip/transaction"
	"net/http"
	_ "net/http/pprof"
	"strings"
)

type Proxy struct {
	stack            *kalbi.SipStack
	requestschannel  chan transaction.Transaction
	responseschannel chan transaction.Transaction
	RegisteredUsers  map[string]string
}

func (p *Proxy) HandleRequest(tx transaction.Transaction) {
	if string(tx.GetOrigin().Req.Method) == method.INVITE {

		msg := message.NewResponse(status.Trying, string(tx.GetOrigin().Contact.Host)+"@"+string(tx.GetOrigin().Contact.Host), "@")
		msg.CopyHeaders(tx.GetOrigin())
		msg.ContLen.SetValue("0")
		tx.Send(msg, string(tx.GetOrigin().Contact.Host), string(tx.GetOrigin().Contact.Port))

		user, exists := p.RegisteredUsers[string(tx.GetOrigin().To.User)]
		if exists == false {
			msg := message.NewResponse(status.NotFound, "@", "@")
			msg.CopyHeaders(tx.GetOrigin())
			tx.Send(msg, string(tx.GetOrigin().Contact.Host), string(tx.GetOrigin().Contact.Port))

		} else {
			msg2 := message.NewRequest(method.INVITE, string(tx.GetOrigin().To.User)+"@"+string(tx.GetOrigin().To.Host), string(tx.GetOrigin().From.User)+"@"+string(tx.GetOrigin().From.Host))
			msg2.CopyHeaders(tx.GetOrigin())
			msg2.CopySdp(tx.GetOrigin())
			TxMng := p.stack.GetTransactionManager()
			ctx := TxMng.NewClientTransaction(msg)
			ctx.SetServerTransaction(tx)
			user := strings.Split(user, ":")
			ctx.Send(msg2, user[0], user[1])

		}

	} else if string(tx.GetOrigin().Req.Method) == method.REGISTER {
		p.RegisteredUsers[string(tx.GetOrigin().Contact.User)] = string(tx.GetOrigin().Contact.Host) + ":" + string(tx.GetOrigin().Contact.Port)
		msg := message.NewResponse(status.OK, "@", "@")
		msg.CopyHeaders(tx.GetOrigin())
		msg.ContLen.SetValue("0")
		tx.Send(msg, string(tx.GetOrigin().Contact.Host), string(tx.GetOrigin().Contact.Port))

	} else if string(tx.GetOrigin().Req.Method) == method.BYE {
		msg := message.NewResponse(status.OK, "@", "@")
		msg.CopyHeaders(tx.GetOrigin())
		msg.ContLen.SetValue("0")
		tx.Send(msg, string(tx.GetOrigin().Contact.Host), string(tx.GetOrigin().Contact.Port))

	} else {
		msg := message.NewResponse(status.OK, "@", "@")
		msg.CopyHeaders(tx.GetOrigin())
		msg.ContLen.SetValue("0")
		tx.Send(msg, string(tx.GetOrigin().Contact.Host), string(tx.GetOrigin().Contact.Port))

	}

}


func (p *Proxy) HandleResponse(response transaction.Transaction) {
	 if response.GetLastMessage().GetStatusCode() == 100 {
         return
	 } else {
		  fmt.Println("I GET HERE ")
		  tx := response.GetServerTransaction()
		  tx.Send(response.GetLastMessage(), string(tx.GetOrigin().Contact.Host), string(tx.GetOrigin().Contact.Port) )
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

func (p *Proxy) Start() {
	p.RegisteredUsers = make(map[string]string)
	p.stack = kalbi.NewSipStack("Basic")
	p.stack.CreateListenPoint("udp", "0.0.0.0", 5060)
	p.requestschannel = p.stack.CreateRequestsChannel()
	p.responseschannel = p.stack.CreateResponseChannel()
	go p.stack.Start()
	go p.ServeRequests()
	go http.ListenAndServe("localhost:6060", nil)
	
	p.ServeResponses()

}
