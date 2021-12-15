// curl "http://localhost:8080/messages?src=777&dst=380671234567&txt=Hello"
package http

import (
	"log"
	"net/http"
	"time"

	"github.com/yykhomenko/mars/internal/entity"
	"github.com/yykhomenko/mars/internal/service/router"
)

type HTTPServer struct {
	addr   string
	router router.Router
}

func NewHTTPServer(addr string, router router.Router) *HTTPServer {
	s := &HTTPServer{
		addr:   addr,
		router: router,
	}

	http.HandleFunc("/messages", s.messages)

	return s
}

func (s *HTTPServer) Start() error {
	log.Println("HTTP server listen:", s.addr)
	return http.ListenAndServe(s.addr, nil)
}

func (s *HTTPServer) messages(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	r.ParseForm()
	from := r.FormValue("from")
	to := r.FormValue("to")
	text := r.FormValue("text")

	s.router.Route(&entity.Message{
		From: from,
		To:   to,
		Text: text,
	})

	log.Printf("http: rx: duration: %s", time.Since(start))
}
