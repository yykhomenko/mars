package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	http.HandleFunc("/messages", messages())

	log.Println("up")
	log.Fatal(http.ListenAndServe(":8027", nil))
}

func messages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		src := r.FormValue("src")
		dst := r.FormValue("dst")
		text := r.FormValue("text")

		fmt.Println(src, dst, text)

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

		log.Printf("duration: %s", time.Now().Sub(start))
	}
}
