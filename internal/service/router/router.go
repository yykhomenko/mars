package router

import (
	"log"

	"github.com/yykhomenko/mars/internal/entity"
)

type Router struct {
	messages map[string][]*entity.Message
}

func New() *Router {
	messages := make(map[string][]*entity.Message)
	return &Router{messages: messages}
}

func (r *Router) Route(m *entity.Message) {
	r.messages[m.From] = append(r.messages[m.From], m)
	log.Printf("router: message routed: %v\n", m)
}
