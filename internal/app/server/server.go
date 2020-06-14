package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	http.HandleFunc("/messages", root(tx))

	log.Println("up")
	log.Fatal(http.ListenAndServe(":8027", nil))
}

func root(tx *smpp.Transceiver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		src := r.FormValue("src")
		dst := r.FormValue("dst")
		text := r.FormValue("text")

		fmt.Println(src, dst, text)

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
