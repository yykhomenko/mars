package router

import (
	"sync"

	"github.com/yykhomenko/mars/pkg/mars/config"
	"github.com/yykhomenko/mars/pkg/mars/entity"
)

type Router interface {
	Route(m *entity.Message)
	GetNum() int
}

type router struct {
	config   *config.Config
	messages map[string][]*entity.Message
	Num      int
	mu       sync.Mutex
}

func NewRouter(config *config.Config) Router {
	messages := make(map[string][]*entity.Message)
	return &router{config: config, messages: messages}
}

func (r *router) Route(m *entity.Message) {
	//r.messages[m.From] = append(r.messages[m.From], m)
	r.mu.Lock()
	r.Num++
	r.mu.Unlock()
	//r.config.Log.Printf("router: message %d routed: %v\n", r.Num, m)
}

func (r *router) GetNum() int {
	return r.Num
}
