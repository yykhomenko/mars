package router

import (
	"fmt"

	"github.com/fiorix/go-smpp/smpp"
)

type Router struct {
}

func New() *Router {
	return &Router{}
}

func (r *Router) Route(m *smpp.ShortMessage) {
	fmt.Println("message routed:", m)
}
