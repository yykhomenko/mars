package smpp

import (
	"log"
	"strings"
	"time"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"

	"github.com/yykhomenko/mars/internal/entity"
	"github.com/yykhomenko/mars/internal/service/router"
)

type SMPPConnector struct {
	addr     string
	user     string
	password string
	tx       *smpp.Transceiver
	router   router.Router
}

func NewSMPPConnector(addr, user, password string, router router.Router) *SMPPConnector {
	return &SMPPConnector{
		addr:     addr,
		user:     user,
		password: password,
		router:   router,
	}
}

func (c *SMPPConnector) Start() {
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

	log.Println("SMPPConnector listen:", c.addr)
}

func receiverHandler(router router.Router) smpp.HandlerFunc {
	return func(p pdu.Body) {
		switch p.Header().ID {
		case pdu.DeliverSMID:
			start := time.Now()
			from := p.Fields()[pdufield.SourceAddr].String()
			to := p.Fields()[pdufield.DestinationAddr].String()
			text := p.Fields()[pdufield.ShortMessage].String()

			router.Route(&entity.Message{
				From: from,
				To:   to,
				Text: text,
			})

			log.Printf("smpp: rx: duration: %s", time.Since(start))
		}
	}
}

func parseTLVStatus(text string) map[string]string {
	var m map[string]string
	ss := strings.Split(text, " ")
	m = make(map[string]string)

	for _, pair := range ss {
		z := strings.Split(pair, ":")
		if len(z) == 2 {
			m[z[0]] = z[1]
		}
	}

	return m
}
