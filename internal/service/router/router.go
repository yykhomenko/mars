package router

import (
	"log"

	"github.com/yykhomenko/mars/internal/entity"
)

type Router interface {
	Route(m *entity.Message)
}

type router struct {
	messages map[string][]*entity.Message
}

func NewRouter() Router {
	messages := make(map[string][]*entity.Message)
	return &router{messages: messages}
}

func (r *router) Route(m *entity.Message) {
	r.messages[m.From] = append(r.messages[m.From], m)
	log.Printf("router: message routed: %v\n", m)
}
