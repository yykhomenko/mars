package server

import (
	"log"
	"strconv"
	"strings"
	"unicode"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
)

func NewSMPPConnector(addr, user, password string) {

	log.Println("SMPP server listen:", addr)

	tx := &smpp.Transceiver{
		Addr:    addr,
		User:    user,
		Passwd:  password,
		Handler: receiverHandler,
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

func receiverHandler(p pdu.Body) {
	switch p.Header().ID {
	case pdu.DeliverSMID:
		f := p.Fields()
		src := f[pdufield.SourceAddr]
		dst := f[pdufield.DestinationAddr]
		txt := f[pdufield.ShortMessage]
		log.Printf("Short message from=%q to=%q: %q", src, dst, txt)

		params := ParseTLVStatus(txt.String())

		mid, e := strconv.ParseUint(params["id"], 10, 0)
		if e != nil {
			log.Println("parse MID error: ", e.Error())
		}

		log.Println(mid)
	}
}

func ParseTLVStatus(text string) map[string]string {

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
