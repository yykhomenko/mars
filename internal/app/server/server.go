package server

import (
	"log"
	"net/http"
	"time"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"

	"github.com/cbi-sh/messages/internal/app/router"
)

type HttpConnector struct {
	addr   string
	router *router.Router
}

func New(addr string, router *router.Router) *HttpConnector {
	s := &HttpConnector{
		addr:   addr,
		router: router,
	}

	s.configureRouter()

	return s
}

func (c *HttpConnector) Start() error {
	log.Println("HTTP server listen:", c.addr)
	return http.ListenAndServe(c.addr, nil)
}

func (c *HttpConnector) configureRouter() {
	http.HandleFunc("/messages", messages(c.router))
}

func messages(router *router.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		src := r.FormValue("src")
		dst := r.FormValue("dst")
		text := r.FormValue("text")

		m := &smpp.ShortMessage{
			Src:           src,
			SourceAddrTON: getTON(src),
			Dst:           dst,
			Text:          pdutext.Raw(text),
			Register:      pdufield.FinalDeliveryReceipt,
		}

		router.Route(m)

		log.Printf("duration: %s", time.Now().Sub(start))
	}
}
