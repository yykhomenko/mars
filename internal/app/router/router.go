package router

import (
	"log"

	"github.com/fiorix/go-smpp/smpp"
)

type Router struct {
	queries map[string]*smpp.ShortMessage
}

func New() *Router {
	return &Router{}
}

func (r *Router) Route(m *smpp.ShortMessage) {
	r.queries[m.Src] = m
	log.Println("message routed:", m)
}
