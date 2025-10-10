package smpp

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"

	"github.com/yykhomenko/mars/pkg/mars/entity"
	"github.com/yykhomenko/mars/pkg/mars/service/router"
)

type SMPPConnector struct {
	tx     *smpp.Transceiver
	router router.Router
}

func NewSMPPConnector(addr, user, password string, router router.Router) *SMPPConnector {
	return &SMPPConnector{
		tx: &smpp.Transceiver{
			Addr:        addr,
			User:        user,
			Passwd:      password,
			Handler:     receiverHandler(router),
			RespTimeout: 5 * time.Second,
		},
		router: router,
	}
}

func (c *SMPPConnector) Start() {
	go func() {
		statuses := c.tx.Bind()
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

	log.Println("SMPPConnector listen:", c.tx.Addr)
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

func (c *SMPPConnector) Send(msg *entity.Message) (*smpp.ShortMessage, error) {
	if c.tx == nil {
		return nil, fmt.Errorf("smpp: transceiver not connected")
	}

	var enc = "latin1"
	var codec pdutext.Codec
	switch enc {
	case "ucs2", "ucs-2":
		codec = pdutext.UCS2(msg.Text)
	case "latin1", "latin-1":
		codec = pdutext.Latin1(msg.Text)
	default:
		codec = pdutext.Raw(msg.Text)
	}

	sm := &smpp.ShortMessage{
		Src:      msg.From,
		Dst:      msg.To,
		Text:     codec,
		Register: pdufield.NoDeliveryReceipt,
	}

	resp, err := c.tx.Submit(sm)
	if err != nil {
		return nil, fmt.Errorf("smpp: submit: %w", err)
	}

	return resp, nil
}
