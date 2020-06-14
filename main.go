package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
	"golang.org/x/time/rate"
)

var messages map[uint64]*smpp.ShortMessage

func main() {

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

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
			log.Println(messages)
		}
	}

	lm := rate.NewLimiter(rate.Limit(1000), 1)

	tx := &smpp.Transceiver{
		Addr:        "192.168.0.2:3736",
		User:        "user",
		Passwd:      "password",
		Handler:     f,
		RateLimiter: lm,
	}

	conn := tx.Bind()
	go func() {
		for c := range conn {
			log.Println("SMPP connection status:", c.Status(), c.Error())
		}
	}()

	http.HandleFunc("/messages", root(tx))

	log.Println("up")
	log.Fatal(http.ListenAndServe(":8027", nil))
}

func root(tx *smpp.Transceiver) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Read body
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		log.Println(string(b))

		start := time.Now()

		// src := r.FormValue("src")
		// dst := r.FormValue("dst")
		// text := r.FormValue("text")

		src := "777"
		dst := "380671112222"
		text := "hello"

		message := &smpp.ShortMessage{
			Src:           src,
			SourceAddrTON: getTON(src),
			Dst:           dst,
			Text:          pdutext.Raw(text),
			Register:      pdufield.FinalDeliveryReceipt,
		}

		resp, e := tx.Submit(message)

		if e == smpp.ErrNotConnected {
			http.Error(w, "Oops.", http.StatusServiceUnavailable)
			return
		}

		if e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}

		midStr := resp.RespID()

		mid, e := strconv.ParseUint(midStr, 16, 0)

		if e != nil {
			log.Println("parse MID error: ", e.Error())
		}

		log.Println("mid:", mid)

		_, _ = io.WriteString(w, midStr)

		log.Printf("duration: %s", time.Now().Sub(start))
	}
}

func toMap(text string) map[string]string {

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
