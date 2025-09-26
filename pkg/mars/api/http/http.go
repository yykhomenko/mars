package http

import (
	"log"
	"net/http"
	"time"

	"github.com/yykhomenko/mars/pkg/mars/config"
	"github.com/yykhomenko/mars/pkg/mars/entity"
	"github.com/yykhomenko/mars/pkg/mars/service/router"
)

type HTTPServer struct {
	conf   *config.Config
	router router.Router
}

func NewHTTPServer(conf *config.Config, router router.Router) *HTTPServer {
	s := &HTTPServer{
		conf:   conf,
		router: router,
	}

	http.HandleFunc("/messages", s.messages)

	return s
}

func (s *HTTPServer) Start() error {
	log.Println("HTTP server listen:", s.conf.MarsApiAddr)
	return http.ListenAndServe(s.conf.MarsApiAddr, nil)
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
