package server

import (
	"log"
	"strconv"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
)

func NewSMPPConnector() {

	f := func(p pdu.Body) {
		switch p.Header().ID {
		case pdu.DeliverSMID:
			f := p.Fields()
			src := f[pdufield.SourceAddr]
			dst := f[pdufield.DestinationAddr]
			txt := f[pdufield.ShortMessage]
			log.Printf("Short message from=%q to=%q: %q", src, dst, txt)

			params := toMap(txt.String())

			mid, e := strconv.ParseUint(params["id"], 10, 0)
			if e != nil {
				log.Println("parse MID error: ", e.Error())
			}

			log.Println(mid)
		}
	}

	tx := &smpp.Transceiver{
		Addr:    "192.168.0.2:3736",
		User:    "user",
		Passwd:  "password",
		Handler: f,
	}

	conn := tx.Bind()
	go func() {
		for c := range conn {
			log.Println("SMPP connection status:", c.Status(), c.Error())
			if c.Status() == smpp.Disconnected {
				if err := tx.Close(); err != nil {
					log.Println("error close tx:", err)
				}
				// log.Println("try rebind")
				// conn = tx.Bind() // todo check for leaks
			}
		}
	}()

}
