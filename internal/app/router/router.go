package router

import (
	"log"

	"github.com/fiorix/go-smpp/smpp"
)

type Router struct {
	messages map[string][]*smpp.ShortMessage
}

func New() *Router {
	messages := make(map[string][]*smpp.ShortMessage)
	return &Router{messages: messages}
}

func (r *Router) Route(m *smpp.ShortMessage) {
	r.messages[m.Src] = append(r.messages[m.Src], m)
	log.Printf("router: message routed: %v\n", m)
	log.Printf("router: text: %s\n", m.Text.Decode())
	log.Printf("router: messages to sent: %v\n", r.messages)
}
