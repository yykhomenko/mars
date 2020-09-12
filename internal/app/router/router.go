package router

import (
	"fmt"

	"github.com/fiorix/go-smpp/smpp"
)

type Router struct {
	query []*smpp.ShortMessage
}

func New() *Router {
	return &Router{}
}

func (r *Router) Route(m *smpp.ShortMessage) {
	r.query = append(r.query, m)
	fmt.Println("message routed:", m)
}
