// curl "http://localhost:8080/messages?src=777&dst=380671234567&txt=Hello"
package server

import (
	"log"
	"net/http"
	"time"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"

	"github.com/yykhomenko/mars/internal/app/router"
)

type HTTPConnector struct {
	addr   string
	router *router.Router
}

func NewHTTPConnector(addr string, router *router.Router) *HTTPConnector {
	s := &HTTPConnector{
		addr:   addr,
		router: router,
	}

	s.configureRouter()

	return s
}

func (c *HTTPConnector) Start() error {
	log.Println("HTTP server listen:", c.addr)
	return http.ListenAndServe(c.addr, nil)
}

func (c *HTTPConnector) configureRouter() {
	http.HandleFunc("/messages", messages(c.router))
}

func messages(router *router.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		r.ParseForm()
		src := r.FormValue("src")
		dst := r.FormValue("dst")
		txt := r.FormValue("txt")

		router.Route(&smpp.ShortMessage{
			Src:           src,
			SourceAddrTON: getTON(src),
			Dst:           dst,
			Text:          pdutext.Raw(txt),
			Register:      pdufield.FinalDeliveryReceipt,
		})

		log.Printf("http: rx: duration: %s", time.Now().Sub(start))
	}
}
