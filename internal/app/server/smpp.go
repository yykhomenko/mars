package server

import (
	"log"
	"strconv"
	"strings"
	"unicode"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"

	"github.com/cbi-sh/messages/internal/app/router"
)

type SMPPConnector struct {
	addr     string
	user     string
	password string
	tx       *smpp.Transceiver
	router   *router.Router
}

func NewSMPPConnector(addr, user, password string, router *router.Router) *SMPPConnector {
	return &SMPPConnector{
		addr:     addr,
		user:     user,
		password: password,
		router:   router,
	}
}

func (c *SMPPConnector) Start() error {
	c.tx = &smpp.Transceiver{
		Addr:    c.addr,
		User:    c.user,
		Passwd:  c.password,
		Handler: receiverHandler(c.router),
	}

	statuses := c.tx.Bind()
	go func() {
		for s := range statuses {
			switch s.Status() {
			case smpp.Connected:
				log.Printf("smpp: tx: connected to %s", c.tx.Addr)
			case smpp.Disconnected:
				log.Printf("smpp: tx: disconnected from %s", c.tx.Addr)
				if err := c.tx.Close(); err != nil {
					log.Printf("smpp: tx: close: %v\n", err)
				}
			case smpp.ConnectionFailed:
				log.Printf("smpp: tx: unable to connect: %v\n", s.Error())
			default:
				log.Printf("smpp: tx: %s, err: %v\n", s.Status(), s.Error())
			}
		}
	}()

	log.Println("SMPP server listen:", c.addr)

	return nil
}

func receiverHandler(r *router.Router) smpp.HandlerFunc {
	return func(p pdu.Body) {
		switch p.Header().ID {
		case pdu.DeliverSMID:
			f := p.Fields()
			src := f[pdufield.SourceAddr].String()
			dst := f[pdufield.DestinationAddr].String()
			txt := f[pdufield.ShortMessage].String()
			log.Printf("Short message from=%q to=%q: %q", src, dst, txt)

			params := parseTLVStatus(txt)

			mid, e := strconv.ParseUint(params["id"], 10, 0)
			if e != nil {
				log.Fatalf("mid: %v\n", e.Error())
			}

			log.Println(mid)

			m := &smpp.ShortMessage{
				Src:           src,
				SourceAddrTON: getTON(src),
				Dst:           dst,
				Text:          pdutext.Raw(txt),
				Register:      pdufield.NoDeliveryReceipt,
			}

			r.Route(m)
		}
	}
}

func parseTLVStatus(text string) map[string]string {
	var m map[string]string
	var ss []string

	ss = strings.Split(text, " ")
	m = make(map[string]string)

	for _, pair := range ss {
		z := strings.Split(pair, ":")
		if len(z) == 2 {
			m[z[0]] = z[1]
		}
	}

	return m
}

func getTON(s string) uint8 {
	if isLetter(s) {
		return 0x05
	}
	return 0
}

func isLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// message := &smpp.ShortMessage{
// 	Src:           src,
// 	SourceAddrTON: getTON(src),
// 	Dst:           dst,
// 	Text:          pdutext.Raw(text),
// 	Register:      pdufield.FinalDeliveryReceipt,
// }
//
// resp, e := tx.Submit(message)
//
// if e == smpp.ErrNotConnected {
// 	http.Error(w, "Oops.", http.StatusServiceUnavailable)
// 	return
// }
//
// if e != nil {
// 	http.Error(w, e.Error(), http.StatusBadRequest)
// 	return
// }
//
// midStr := resp.RespID()
//
// mid, e := strconv.ParseUint(midStr, 16, 0)
// if e != nil {
// 	log.Println("parse MID error: ", e.Error())
// }
//
// log.Println("mid:", mid)
//
// _, _ = io.WriteString(w, midStr)
