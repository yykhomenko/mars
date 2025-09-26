package router

import (
	"github.com/yykhomenko/mars/pkg/mars/config"
	"github.com/yykhomenko/mars/pkg/mars/entity"
)

type Router interface {
	Route(m *entity.Message)
}

type router struct {
	config   *config.Config
	num      int
	messages map[string][]*entity.Message
}

func NewRouter(config *config.Config) Router {
	messages := make(map[string][]*entity.Message)
	return &router{config: config, messages: messages}
}

func (r *router) Route(m *entity.Message) {
	//r.messages[m.From] = append(r.messages[m.From], m)
	r.num++
	r.config.Log.Printf("router: message %d routed: %v\n", r.num, m)
}
